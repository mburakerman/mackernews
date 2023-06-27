# mackernews

## Build

to create Go byte, first install 2goarray package and run command below to get auto created 'icon' package

```
go install github.com/cratonica/2goarray
```

```
$GOPATH/bin/2goarray Data icon < icon/icon.png > icon.go
```

<hr />

bundle app

```
go build -o mackernews.app -ldflags="-s -w \
        -X 'main.icon=<icon-file-location>' \
        -X 'main.BundleID=<bundle-identifier>' \
        -X 'main.Version=<version>'" .
```
