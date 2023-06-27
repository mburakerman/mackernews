# Mackernews

Mackernews is a tiny Mac menu bar app that enables you to quickly access latest & most popular Hacker News stories


https://github.com/mburakerman/mackernews/assets/17620102/3b6a7356-1caf-4307-af9a-c1a46bb12f8c


<br />
<hr />

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


