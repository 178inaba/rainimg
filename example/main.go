package main

import (
	"github.com/178inaba/rainimg"
	"github.com/golang/glog"
	"flag"
)

func main() {
	flag.Parse()

	glog.Info("img path:", rainimg.GetImgPath())
}
