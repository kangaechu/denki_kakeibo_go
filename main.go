package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

func main() {
	// settings.yamlから設定を読み込み
	viper.SetConfigName("settings")  // 設定ファイル名を拡張子抜きで指定する
	viper.AddConfigPath(".")         // 現在のワーキングディレクトリを探索することもできる
	vipererr := viper.ReadInConfig() // 設定ファイルを探索して読み取る
	if vipererr != nil {             // 設定ファイルの読み取りエラー対応
		panic(fmt.Errorf("設定ファイル読み込みエラー: %s", vipererr))
	}
	var userID = viper.GetString("userID")
	var password = viper.GetString("password")
	var days = viper.GetInt("days")

	// 出力ファイルを開く
	file, err := os.Create(`output.csv`)
	if err != nil {
		panic(fmt.Errorf("出力ファイルオープンエラー: %s", err))
	}
	defer file.Close()

	header := "DATE, 0:00, 0:30, 1:00, 1:30, 2:00, 2:30, 3:00, 3:30, 4:00, 4:30, 5:00, 5:30, 6:00, 6:30, 7:00, 7:30, 8:00, 8:30, 9:00, 9:30, 10:00, 10:30, 11:00, 11:30, 12:00, 12:30, 13:00, 13:30, 14:00, 14:30, 15:00, 15:30, 16:00, 16:30, 17:00, 17:30, 18:00, 18:30, 19:00, 19:30, 20:00, 20:30, 21:00, 21:30, 22:00, 22:30, 23:00, 23:30, 0:00"
	file.Write(([]byte)(header))

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

		oneDayData := fmt.Sprintf("%s, %s\n", rDate.FindStringSubmatch(html)[1], rData.FindStringSubmatch(html)[1])
		file.Write(([]byte)(oneDayData))
		fmt.Println(oneDayData)

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
