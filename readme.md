denki_kakeibo_go
====

でんき家計簿の30分ごと電気使用量を取得

## Description

東京電力のでんき家計簿（自由化前の料金プラン向け利用状況確認ページ）から30分ごとの電気使用量を取得、CSV化して出力します。

## Download

https://github.com/kangaechu/denki_kakeibo_go/releases/

## Requirements

Google Chrome

## Usage

### バイナリから実行(Windows)

#### ダウンロード

https://github.com/kangaechu/denki_kakeibo_go/releases/ から最新の denki_kakeibo_windows.zip をダウンロード

#### 設定ファイルの作成

settings.yaml.sampleをsettings.yamlにコピーする。
settings.yamlにログインのアカウント情報と繰り返し回数を指定する。
```
userID: "USERID"     # ログインのユーザ名
password: "PASSWORD" # パスワード
days: 10             # 取得する日数
```

#### 実行

denki_kakeibo.exeをダブルクリック。
1度目の実行時にはWindows Defender Smartscreenのダイアログが表示されるので、詳細情報→実行を選択する。
Google Chromeが立ち上がり、実行が開始される。
認証ダイアログは適宜入力、閉じること。
実行が終わるとoutput.csvに30分ごとの電力量が出力される。

### ソースから実行

#### ダウンロード

```
$ go get github.com/kangaechu/denki_kakeibo_go
$ cd $GOPATH/github.com/kangaechu/denki_kakeibo_go
```

#### 設定ファイルの作成

```
$ cp settings.yaml.sample settings.yaml
$ vi settings.yaml

userID: "USERID"     # ログインのユーザ名
password: "PASSWORD" # パスワード
days: 10             # 取得する日数
```

#### その他

実行時にGoogle Chromenについての画面が出力される場合は
https://github.com/rjeczalik/chromedp/commit/520a76514fd911e7544c72412d307fec5ae524ad を適用すると修正されます。

#### 実行
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