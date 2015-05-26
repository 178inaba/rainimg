package rainimg

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"time"
	"github.com/golang/glog"
	"flag"
)

const (
	baseUrl    = "http://tokyo-ame.jwa.or.jp/"
	mapPath    = "map/map000.jpg"
	mskPath    = "map/msk000.png"
	meshFmt    = "mesh/000/%s.gif"
	encDir     = "encode_img/"
	encNameFmt = "%s.png"
	srcDir     = "source_img/"
	layout     = "200601021504"
)

func init() {
	flag.Parse()
	glog.Info("flag parse")
}

func GetImgPath() string {
	tcTime := time.Now().Truncate(5 * time.Minute)
	encFilePath := fmt.Sprintf(encDir+encNameFmt, tcTime.Format(layout))
	_, err := os.Stat(encFilePath)
	if err == nil {
		glog.Info("it is already")
		return encFilePath
	}

	mapFPath, mskFPath := getBaseFile()
	mapSrc := getImgSrc(mapFPath)
	mskSrc := getImgSrc(mskFPath)

	// mesh
	var meshFPath string
	for meshFPath == "" {
		meshFPath = dlImg(baseUrl + fmt.Sprintf(meshFmt, tcTime.Format(layout)))
		if meshFPath == "" {
			tcTime = tcTime.Add(-5 * time.Minute)
			encFilePath = fmt.Sprintf(encDir+encNameFmt, tcTime.Format(layout))
		}
	}
	meshSrc := getImgSrc(meshFPath)

	rgba := image.NewRGBA(image.Rect(0, 0, 770, 480))
	b := rgba.Bounds()
	p := image.Point{0, 0}
	draw.Draw(rgba, b, mapSrc, p, draw.Src)
	draw.Draw(rgba, b, mskSrc, p, draw.Over)
	draw.Draw(rgba, b, meshSrc, p, draw.Over)

	os.Mkdir(encDir, os.ModePerm)
	encFile, _ := os.Create(encFilePath)
	defer encFile.Close()

	png.Encode(encFile, rgba)

	return encFilePath
}

func dlImg(url string) string {
	// req
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Info(url, ": status code is not 200")
		return ""
	}

	// save path
	_, fName := path.Split(url)
	sPath := srcDir + fName

	// save
	f, _ := os.Create(sPath)
	defer f.Close()
	io.Copy(f, resp.Body)

	return sPath
}

func getImgSrc(path string) image.Image {
	f, _ := os.Open(path)
	defer f.Close()
	src, _, _ := image.Decode(f)

	return src
}

func getBaseFile() (string, string) {
	_, mapFName := path.Split(mapPath)
	_, err := os.Stat(srcDir + mapFName)
	var mapFPath, mskFPath string
	if err != nil {
		glog.Info("get base path")
		os.Mkdir(srcDir, os.ModePerm)
		mapFPath = dlImg(baseUrl + mapPath)
		mskFPath = dlImg(baseUrl + mskPath)
	} else {
		mapFPath = srcDir + mapFName

		_, mskFName := path.Split(mskPath)
		mskFPath = srcDir + mskFName
	}

	return mapFPath, mskFPath
}
