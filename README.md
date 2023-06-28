# Mackernews

A a tiny Mac menu bar app that enables you to quickly access latest & most popular Hacker News stories

https://github.com/mburakerman/mackernews/assets/17620102/f70bf1c7-6d8d-4b23-9353-759c36d54796

<br />
<hr />

### 📦 Bundle

```
go build -o mackernews.app -ldflags="-s -w \
        -X 'main.icon=<icon-file-location>' \
        -X 'main.BundleID=<bundle-identifier>' \
        -X 'main.Version=<version>'" .
```
