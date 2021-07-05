package libs

import (
	"coloring/common"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/panjf2000/ants/v2"
)

type urlParams struct {
	uri string
	url string
}

func ColoringBook(id string) {
	htmlUri := "https://www.coloring-book.info/coloring/"
	htmlUrl := htmlUri + "coloring_page.php?id=" + id

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
	var imgList []string
	doc.Find("a>img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		//fmt.Println(src)

		reg := regexp.MustCompile(`thumbs`)
		if reg.FindStringIndex(src) != nil {
			href, _ := s.Parent().Attr("href")
			imgList = append(imgList, htmlUri+href)
		}
	})

	defer ants.Release()
	p, _ := ants.NewPool(20)

	//用于等待协程完成
	var wg sync.WaitGroup

	for _, url := range imgList {
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
		href, _ := doc.Find("img.print").Attr("src")
		picUrl := curUri + href

		//获取当前目录下的`downloads`作为存储目录
		savePath, _ := os.Getwd()
		savePath = savePath + "/downloads/"

		common.Download(picUrl, savePath, callBack)
		//fmt.Println(picUrl)
	}
}
