// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"os"

	"github.com/catenocrypt/nano-work-cache/restapi"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
	"github.com/catenocrypt/nano-work-cache/workcache"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// first optional paramter is config file name
	if len(os.Args) > 1 {
		workcache.SetConfigFile(os.Args[1])
	}

	fmt.Printf("Config: \n")
	rpcUrl := workcache.ConfigNodeRpc()
	if len(rpcUrl) == 0 {
		panic("No value configured for Main.NodeRpc")
	}
	rpcWorkUrl := workcache.ConfigNodeRpcWork()
	if len(rpcWorkUrl) == 0 {
		rpcWorkUrl = rpcUrl
	}
	fmt.Printf("  NodeRpc          %v \n", rpcUrl)
	fmt.Printf("  NodeRpcWork      %v \n", rpcWorkUrl)
	fmt.Printf("  ListenIpPort     %v \n", workcache.ConfigListenIpPort())
	fmt.Printf("  RestMaxActiveRequests  %v \n", workcache.ConfigRestMaxActiveRequests())
	fmt.Printf("  BackgroundWorkerCount  %v \n", workcache.ConfigBackgroundWorkerCount())
	fmt.Printf("  MaxOutRequests   %v \n", workcache.ConfigMaxOutRequests())
	fmt.Printf("  EnablePregeneration  %v \n", workcache.ConfigEnablePregeneration())
	fmt.Printf("  PregenerationQueueSize  %v \n", workcache.ConfigPregenerationQueueSize())
	fmt.Printf("  MaxCacheAgeDays  %v \n", workcache.ConfigMaxCacheAgeDays())

	rpcclient.Init(rpcUrl, rpcWorkUrl)
	workcache.Start()
	restapi.Start()
}
