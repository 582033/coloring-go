package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Mkdir(dir string) bool {
	isExist := false
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsExist(err) {
			isExist = true
		}
	}
	if isExist {
		//fmt.Println("已存在")
		return false
	}
	if os.Mkdir(dir, os.ModePerm) == nil {
		//fmt.Println("创建失败")
		return false
	}
	return true
}

func Download(url string, savePath string, callBack func()) {
	fileInfo := strings.Split(url, "/")
	fileName := fileInfo[len(fileInfo)-1]
	//创建目录
	Mkdir(savePath)
	//fmt.Println(fileName)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("下载失败:%s", resp.Request.URL)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取数据失败")
	}

	//fmt.Printf(string(data))
	saveTo := savePath + fileName
	fmt.Printf("文件:" + saveTo + "下载成功\n")
	ioutil.WriteFile(saveTo, data, 0644)
	defer resp.Body.Close()
	callBack()
}
