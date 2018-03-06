package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

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
	var days = viper.GetInt("days")

	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithErrorf(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// ログイン→30分ごとページへの遷移
	var site, res string
	err = c.Run(ctxt, login(userID, password, &site, &res))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < days; i++ {
		// データを取得
		var html string
		err = c.Run(ctxt, get30MinData(&site, &html))
		if err != nil {
			log.Fatal(err)
		}
		rDate := regexp.MustCompile(`(\d{4}\/\d{2}\/\d{2})　の電気使用量`)
		rData := regexp.MustCompile(`var items = \[\["日次", (.*?)\]`)

		fmt.Println(rDate.FindStringSubmatch(html)[1], rData.FindStringSubmatch(html)[1])

		err = c.Run(ctxt, navPrevDate(&site, &res))
		if err != nil {
			log.Fatal(err)
		}
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
}

func login(userID string, password string, site, res *string) chromedp.Tasks {
	tasks := chromedp.Tasks{
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
	}
	return tasks
}

func get30MinData(site, res *string) chromedp.Tasks {
	tasks := chromedp.Tasks{
		// 日毎の使用量のページ
		chromedp.Sleep(5 * time.Second),
		chromedp.OuterHTML(`html`, res, chromedp.ByQuery),
	}
	return tasks
}

func navPrevDate(site, res *string) chromedp.Tasks {
	tasks := chromedp.Tasks{
		chromedp.Click(`#doPrevious`),
	}
	return tasks
}
