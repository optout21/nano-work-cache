// Copyright © 2019-2019 catenocrypt.  See LICENSE file for license information.

package main

import (
	"github.com/catenocrypt/nano-work-cache/restapi"
	"github.com/catenocrypt/nano-work-cache/workcache"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rpcUrl := workcache.ConfigGetString("Main.NodeRpc")
	if (len(rpcUrl) == 0) {
		panic("No value configured for Main.NodeRpc")
	}
	listenIpPort := workcache.ConfigGetString("Main.ListenIpPort")

	restapi.Start(rpcUrl, listenIpPort)
}
