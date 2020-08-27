// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
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
	if len(os.Args) > 1 {
		workcache.SetConfigFile(os.Args[1])
	}

	fmt.Printf("Config: \n")
	rpcUrl := workcache.ConfigNodeRpc()
	if len(rpcUrl) == 0 {
		panic("No value configured for Main.NodeRpc")
	}
	fmt.Printf("  NodeRpc  %v \n", rpcUrl)
	listenIpPort := workcache.ConfigListenIpPort()
	fmt.Printf("  RestMaxActiveRequests  %v \n", workcache.ConfigRestMaxActiveRequests())
	fmt.Printf("  BackgroundWorkerCount  %v \n", workcache.ConfigBackgroundWorkerCount())
	fmt.Printf("  MaxOutRequests  %v \n", workcache.ConfigMaxOutRequests())
	fmt.Printf("  EnablePregeneration  %v \n", workcache.ConfigEnablePregeneration())
	fmt.Printf("  PregenerationQueueSize  %v \n", workcache.ConfigPregenerationQueueSize())
	fmt.Printf("  MaxCacheAgeDays  %v \n", workcache.ConfigMaxCacheAgeDays())

	workcache.Start()

	restapi.Start(rpcUrl, listenIpPort)
}
