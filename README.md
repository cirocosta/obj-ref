# obj-ref

figures out a fully qualified reference to a Kubernetes object.

```console
$ obj-ref pipeline test
apiVersion: tekton.dev/v1beta1
kind: Pipeline
name: test
```

or

```console
$ obj-ref -o line pipeline test
pipelines.tekton.dev test
```

## why

sometimes, you need to provide the equivalent of a `corev1.ObjectReference` to fields in Kubernetes custom resources.

having a little cli that can do that for me makes things easier, e.g., by using
vim's `:read` command, I can, in the middle of a piece of YAML that I'm working
on:

```
:r! obj-ref pipeline test
```

and then boom! you have the object reference right there.


## LICENSE

MIT
