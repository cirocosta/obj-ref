package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jessevdk/go-flags"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/yaml"
)

var cli = struct {
	Output string `long:"output" short:"o" choice:"yaml" choice:"line" default:"yaml"`

	Positional struct {
		Resource string
		Name     string
	} `positional-args:"yes" required:"true"`
}{}

type Mapper struct {
	rm meta.RESTMapper
}

func NewMapper(cfg *rest.Config) (*Mapper, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("new discovery client for config: %w", err)
	}

	rm := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return &Mapper{
		rm: rm,
	}, nil
}

func (m *Mapper) MappingFor(resource, name string) (*meta.RESTMapping, error) {
	gvk, err := m.GVKFor(resource, name)
	if err != nil {
		return nil, fmt.Errorf("gvk for '%s %s': %w", resource, name, err)
	}

	mapping, err := m.rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, fmt.Errorf("rest mapping: %w", err)
	}

	return mapping, nil
}

func (m *Mapper) GVKFor(resource, name string) (*schema.GroupVersionKind, error) {
	fullGVR, gr := schema.ParseResourceArg(resource)
	if fullGVR != nil {
		gvk, err := m.rm.KindFor(*fullGVR)
		if err != nil {
			return nil, fmt.Errorf("kind for %s: %w", fullGVR, err)
		}

		return &gvk, nil
	}

	gvk, err := m.rm.KindFor(gr.WithVersion(""))
	if err != nil {
		return nil, fmt.Errorf("kind for %s: %w", fullGVR, err)
	}

	return &gvk, nil
}

func showObjectReference(ref *corev1.ObjectReference) error {
	d, err := yaml.Marshal(ref)
	if err != nil {
		return fmt.Errorf("marhsal ref: %w", err)
	}

	fmt.Print(string(d))
	return nil
}

func run(ctx context.Context) error {
	_, err := flags.Parse(&cli)
	if err != nil {
		return fmt.Errorf("flags parse: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	mapper, err := NewMapper(cfg)
	if err != nil {
		return fmt.Errorf("new mapper: %w", err)
	}

	mapping, err := mapper.MappingFor(cli.Positional.Resource, cli.Positional.Name)
	if err != nil {
		return fmt.Errorf("mapping for: %w", err)
	}

	switch cli.Output {
	case "yaml":
		showObjectReference(&corev1.ObjectReference{
			Name:       cli.Positional.Name,
			Kind:       mapping.GroupVersionKind.Kind,
			APIVersion: mapping.GroupVersionKind.GroupVersion().String(),
		})
	case "line":
		fmt.Printf("%s %s",
			mapping.Resource.GroupResource().String(),
			cli.Positional.Name,
		)
	default:
		return fmt.Errorf("unknown output %s", cli.Output)
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := run(ctx); err != nil {
		panic(err)
	}

	return
}
