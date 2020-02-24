// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"math"
	"os"
	"github.com/catenocrypt/nano-work-cache/restapi"
	"github.com/catenocrypt/nano-work-cache/workcache"
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
	fmt.Printf("Config: NodeRpc  %v \n", rpcUrl)
	listenIpPort := workcache.ConfigGetString("Main.ListenIpPort")
	fmt.Printf("Config: ListenIpPort  %v \n", listenIpPort)
	restMaxConcRequests := workcache.ConfigGetIntWithDefault("Main.RestMaxConcRequests", 200)
	restMaxConcRequests = int(math.Max(float64(restMaxConcRequests), float64(20)))
	fmt.Printf("Config: RestMaxConcRequests  %v \n", restMaxConcRequests)
	
	workcache.Start()

	restapi.Start(rpcUrl, listenIpPort, restMaxConcRequests)
}
