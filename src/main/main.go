package main

import (
	//"fmt"
	"github.com/catenocrypt/nano-work-cache/restapi"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var url = "http://54.227.21.124:7176"
	restapi.Start(url)
}
