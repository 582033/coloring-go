package main

import (
	"coloring/libs"
	"flag"
	"fmt"
	"strconv"
)

var id = flag.Int("id", 0, "要下载的coloring id")

func main() {
	flag.Parse()
	if *id <= 0 {
		fmt.Println("缺少参数id")
		fmt.Println("使用-help查看")
		return
	}
	coloringBookId := strconv.Itoa(*id)
	libs.ColoringBook(coloringBookId)
}
