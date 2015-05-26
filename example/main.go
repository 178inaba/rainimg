package main

import (
	"flag"
	"github.com/178inaba/rainimg"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	glog.Info("img path:", rainimg.GetImgPath())
}
