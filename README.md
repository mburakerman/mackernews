# Mackernews

Mackernews is a tiny Mac menu bar app that enables you to quickly access latest & most popular Hacker News stories

https://github.com/mburakerman/mackernews/assets/17620102/07bd5f6f-bc17-4067-811f-4bf0b992e09f


### 📦 Bundle

```
go build -o mackernews.app -ldflags="-s -w \
        -X 'main.icon=<icon-file-location>' \
        -X 'main.BundleID=<bundle-identifier>' \
        -X 'main.Version=<version>'" .
```


#### 📌 Note
In order to create Go byte, first install `2goarray` package and run command below to get auto created `icon` package

```
go install github.com/cratonica/2goarray
```

```
$GOPATH/bin/2goarray Data icon < icon/icon.png > icon.go
```

<hr />


