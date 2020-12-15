module github.com/cirocosta/obj-ref

go 1.15

require (
	github.com/go-logr/logr v0.2.1 // indirect
	github.com/jessevdk/go-flags v1.4.0
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v10.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.1-0.20200902144306-f2d4ad78c7ab

replace k8s.io/client-go => k8s.io/client-go v0.19.0
