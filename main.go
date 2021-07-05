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
		fmt.Println("缺少参数id\n")
		fmt.Println("使用-help查看\n")
		return
	}
	coloringBookId := strconv.Itoa(*id)
	libs.ColoringBook(coloringBookId)
}
