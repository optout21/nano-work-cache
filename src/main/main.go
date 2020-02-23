// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"github.com/catenocrypt/nano-work-cache/restapi"
	"github.com/catenocrypt/nano-work-cache/workcache"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// first optional paramter is config file name
	if (len(os.Args) > 1) {
		workcache.SetConfigFile(os.Args[1])
	}

	rpcUrl := workcache.ConfigGetString("Main.NodeRpc")
	if (len(rpcUrl) == 0) {
		panic("No value configured for Main.NodeRpc")
	}
	listenIpPort := workcache.ConfigGetString("Main.ListenIpPort")

	workcache.Start()

	restapi.Start(rpcUrl, listenIpPort)
}
