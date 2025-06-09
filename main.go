package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/godfather1103/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var f *os.File
	var rootPath string
	flag.StringVar(&rootPath, "prefix", ".", "下载的根路径")
	flag.Parse()
	who, _ := os.Hostname()
	var realPathPrefix = rootPath + "/" + who + "/Wallpaper/"
	log.Println("下载路径为：" + realPathPrefix)
	exists, _ := utils.PathExists(realPathPrefix)
	if !exists {
		utils.PathMkdir(realPathPrefix)
	}
	resp, err := http.Get("http://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=10")
	if err != nil {
		log.Printf("下载出错：%s", err)
		os.Exit(0)
	}
	defer resp.Body.Close()
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)
	if buf != nil {
		var imagesJson = map[string]interface{}{}
		json.Unmarshal(buf.Bytes(), &imagesJson)
		var x = imagesJson["images"].([]interface{})
		var imageUrls = make([]string, len(x))
		var endTimes = make([]string, len(x))
		for index, item := range x {
			var url = item.(map[string]interface{})["url"].(string)
			// url = strings.Replace(url, "1920x1080", "1366x768", 1)
			imageUrls[index] = "https://cn.bing.com" + url
			endTimes[index] = item.(map[string]interface{})["enddate"].(string)
		}

		for index, item := range imageUrls {
			fileName := realPathPrefix + "/" + endTimes[index] + ".jpg"
			exists, _ = utils.PathExists(fileName)
			if exists {
				f, _ = os.OpenFile(fileName, os.O_RDWR, 0666)
			} else {
				f, err = os.Create(fileName)
				if err != nil {
					log.Printf("创建文件出错：%s", err)
					os.Exit(0)
				}
			}
			log.Printf("开始下载第%d张图", index+1)
			resp, _ = http.Get(item)
			body, _ := ioutil.ReadAll(resp.Body)
			if f != nil && len(body) > 0 {
				io.Copy(f, bytes.NewReader(body))
			}
		}
	} else {
		log.Println("下载出错：请求结果为空!")
		os.Exit(0)
	}
}
