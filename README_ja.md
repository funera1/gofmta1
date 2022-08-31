# gofmtal
gofmtalはGoプログラムをフォーマットします。
gofmtalはgofmtの機能に加え、コメント内のコードにもフォーマットを掛けるコマンドです。

## インストール
Goのバージョンによってインストール方法が異なります

** Go 1.16未満 **
```
$ go get -u github.com/funera1/gofmtal
```

** Go 1.16以上 **
```
$ go install github.com/funera1/gofmtal@latest
```

## 使用方法
```
gofmtal [flags] [path ...]
```
