denki_kakeibo_go
====

でんき家計簿の30分ごと電気使用量を取得

## Description

東京電力のでんき家計簿（自由化前の料金プラン向け利用状況確認ページ）から30分ごとの電気使用量を取得、CSV化して出力します。

## Demo

## Usage

### ダウンロード

```
$ go get github.com/kangaechu/denki_kakeibo_go
$ cd $GOPATH/github.com/kangaechu/denki_kakeibo_go
```

### 設定ファイルの作成

```
$ cp settings.yaml.sample settings.yaml
$ vi settings.yaml

userID: "USERID"     # ログインのユーザ名
password: "PASSWORD" # パスワード
days: 10             # 取得する日数
```
### 実行
```
$ go run main.go
```
実行が終わるとoutput.csvに30分ごとの電力量が出力されます。

## Install

クロスビルドも可能です。
```
# Linux
$ GOOS=linux GOARCH=amd64 go build -o denki_kakeibo main.go
# OSX
$ GOOS=darwin GOARCH=amd64 go build -o denki_kakeibo main.go
# Windows
$ GOOS=windows GOARCH=amd64 go build -o denki_kakeibo.exe main.go
```

## Licence

[MIT](https://github.com/tcnksm/tool/blob/master/LICENCE)

## Author

[kangaechu](https://github.com/kangaechu)