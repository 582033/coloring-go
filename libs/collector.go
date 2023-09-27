package libs

import (
	"coloring/common"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/panjf2000/ants/v2"
)

type urlParams struct {
	uri string
	url string
}

func ColoringBook(id string) {
	htmlUri := "https://downloads.khinsider.com/"
	htmlUrl := htmlUri + "game-soundtracks/album/the-legend-of-zelda-tears-of-the-kingdom-o.s.t-switch-gamerip-2023"

	res, err := http.Get(htmlUrl)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	//从列表页缩略图中查出所有原图链接地址
	var mp3List []string
	//doc.Find("td.clickable-row>a").Each(func(i int, s *goquery.Selection) {
	doc.Find("div.playTrack").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Parent().Parent().Find("a").Attr("href")
		//src, _ := s.Attr("href")
		mp3List = append(mp3List, htmlUri+src)
	})

	if debugBytes, _ := json.Marshal(mp3List); len(debugBytes) > 0 {
		fmt.Printf("RequestID:%v DebugMessage:%s Value:%s", nil, "mp3List", string(debugBytes))
	}

	defer ants.Release()
	p, _ := ants.NewPool(200)

	//用于等待协程完成
	var wg sync.WaitGroup

	for _, url := range mp3List {
		wg.Add(1)
		p.Submit(getColoringBookOriginPicUrlAndDownload(htmlUri, url, func() {
			defer wg.Done()
		}))
	}

	wg.Wait()
}

func getColoringBookOriginPicUrlAndDownload(curUri string, url string, callBack func()) func() {
	return func() {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()
		if res.StatusCode != 200 {
			fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		href, _ := doc.Find("audio#audio").Attr("src")
		picUrl := href

		//获取当前目录下的`downloads`作为存储目录
		savePath, _ := os.Getwd()
		savePath = savePath + "/downloads/"

		common.Download(picUrl, savePath, callBack)
		//fmt.Println(picUrl)
	}
}
