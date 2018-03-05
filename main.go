package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("settings")  // 設定ファイル名を拡張子抜きで指定する
	viper.AddConfigPath(".")         // 現在のワーキングディレクトリを探索することもできる
	vipererr := viper.ReadInConfig() // 設定ファイルを探索して読み取る
	if vipererr != nil {             // 設定ファイルの読み取りエラー対応
		panic(fmt.Errorf("設定ファイル読み込みエラー: %s", vipererr))
	}
	var userID = viper.GetString("userID")
	var password = viper.GetString("password")

	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	err = c.Run(ctxt, getDenkiKakeibo30MinData(userID, password, &site, &res))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("saved screenshot from search result listing `%s` (%s)", res, site)
}

func getDenkiKakeibo30MinData(userID string, password string, site, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// ログインページ
		chromedp.Navigate(`https://www.kakeibo.tepco.co.jp/dk/aut/login/`),
		chromedp.WaitVisible(`#idLogin`, chromedp.ByID),
		chromedp.SendKeys(`input[name="id"]`, userID, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="password"]`, password, chromedp.ByQuery),
		chromedp.Click(`#idLogin`),
		// トップページ
		chromedp.WaitVisible(`.login_info`, chromedp.ByQuery),
		chromedp.Click(`.box01 a`, chromedp.ByQuery), // 使用量と料金をグラフで見る
		// 日毎の使用量のページ
		chromedp.WaitVisible(`.graph_head a`, chromedp.ByQuery),
		chromedp.Click(`.graph_head a`, chromedp.ByQuery), // 時間別グラフはこちら
		chromedp.WaitVisible(`.hogehoge`, chromedp.ByQuery),
		// 30分ごとの使用量のページ
		// chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
		// 	return ioutil.WriteFile("screenshot.png", buf, 0644)
		// }),
	}
}
