# mackernews

## Build

```
go build -o mackernews.app -ldflags="-s -w \
        -X 'main.icon=<icon-file-location>' \
        -X 'main.BundleID=<bundle-identifier>' \
        -X 'main.Version=<version>'" .
```
