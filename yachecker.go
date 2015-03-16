package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type handlerError struct {
	Code int
}

func (e handlerError) Error() string {
	return http.StatusText(e.Code)
}

type slackPayload struct {
	Text string `json:"text"`
}

func main() {
	http.HandleFunc("/slackbot", func(w http.ResponseWriter, r *http.Request) {
		b, err := topHandler(r)
		switch t := err.(type) {
		case nil:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(b)
		case handlerError:
			http.Error(w, t.Error(), t.Code)
		default:
			log.Fatal("Unknown Error: ", err)
		}
	})

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func topHandler(r *http.Request) ([]byte, error) {

	// POST以外を排除
	if r.Method != "POST" {
		return nil, handlerError{http.StatusMethodNotAllowed}
	}

	// tokenチェック
	token := r.FormValue("token")
	if len(os.Getenv("SLACK_TOKEN")) < 1 || token != os.Getenv("SLACK_TOKEN") {
		return nil, handlerError{http.StatusBadRequest}
	}

	// 「ヤフオク」を含む？
	text := r.FormValue("text")
	if !strings.Contains(text, "ヤフオク") {
		return nil, handlerError{http.StatusNoContent}
	}

	values := url.Values{}
	values.Add("auccat", "0")
	values.Add("ei", "UTF-8")
	values.Add("n", "20")
	values.Add("p", "mixi アカウント")
	values.Add("tab_ex", "commerce")

	url := "http://closedsearch.auctions.yahoo.co.jp/jp/closedsearch" + "?" + values.Encode()

	// ヤフオクのページを取得
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Print("Get: ", err)
		return nil, handlerError{http.StatusInternalServerError}
	}

	num, sum, max, min := 0, 0, 0, math.MaxInt32

	rgxp, _ := regexp.Compile("[^0-9-]")

	// 落札価格を集計
	doc.Find("span.ePrice").Each(func(_ int, s *goquery.Selection) {
		price, _ := strconv.Atoi(rgxp.ReplaceAllString(s.Text(), ""))
		sum = sum + price
		num = num + 1
		if price < min {
			min = price
		}
		if price > max {
			max = price
		}
	})

	if num == 0 {
		return nil, handlerError{http.StatusNoContent}
	}

	message := fmt.Sprintf("最近の平均落札価格は%.2f円(%d円～%d円)です。\n詳細は<%s|こちら>をクリック。",
		float64(sum)/float64(num), min, max, url)

	b, err := json.Marshal(slackPayload{message})
	if err != nil {
		log.Print("json.Marchal: ", err)
		return nil, handlerError{http.StatusInternalServerError}
	}

	return b, nil
}
