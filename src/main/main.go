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

	rpcUrl := workcache.GetNodeRpc()
	if len(rpcUrl) == 0 {
		panic("No value configured for Main.NodeRpc")
	}
	fmt.Printf("Config: NodeRpc  %v \n", rpcUrl)
	listenIpPort := workcache.GetListenIpPort()
	fmt.Printf("Config: ListenIpPort  %v \n", listenIpPort)
	restMaxActiveRequests := workcache.GetRestMaxActiveRequests()
	fmt.Printf("Config: RestMaxActiveRequests  %v \n", restMaxActiveRequests)
	fmt.Printf("Config: BackgroundWorkerCount  %v \n", workcache.GetBackgroundWorkerCount())
	fmt.Printf("Config: MaxOutRequests  %v \n", workcache.GetMaxOutRequests())
	fmt.Printf("Config: EnablePregeneration  %v \n", workcache.GetEnablePregeneration())
	fmt.Printf("Config: PregenerationQueueSize  %v \n", workcache.GetPregenerationQueueSize())
	fmt.Printf("Config: MaxCacheAgeDays  %v \n", workcache.GetMaxCacheAgeDays())

	workcache.Start()

	restapi.Start(rpcUrl, listenIpPort, restMaxActiveRequests)
}
