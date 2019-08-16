# kustomize-parser

To run `kustomize-parser` run `go build` and then
```bash
./kustomize-parser --input {path to kustomize file}
```

You should see a list of files required for that specific file. Right now it only handles resources, patches and configMapGenerator.