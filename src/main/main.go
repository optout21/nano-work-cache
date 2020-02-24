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
	restMaxActiveRequests := workcache.ConfigGetIntWithDefault("Main.RestMaxActiveRequests", 200)
	restMaxActiveRequests = int(math.Max(float64(restMaxActiveRequests), float64(20)))
	fmt.Printf("Config: RestMaxActiveRequests  %v \n", restMaxActiveRequests)
	backgroundWorkerCount := workcache.ConfigGetIntWithDefault("Main.BackgroundWorkerCount", 4)
	backgroundWorkerCount = int(math.Max(float64(backgroundWorkerCount), float64(2)))
	backgroundWorkerCount = int(math.Min(float64(backgroundWorkerCount), float64(20)))
	fmt.Printf("Config: BackgroundWorkerCount  %v \n", backgroundWorkerCount)
	maxOutRequests := workcache.ConfigGetIntWithDefault("Main.MaxOutRequests", 8)
	maxOutRequests = int(math.Max(float64(maxOutRequests), float64(3)))
	maxOutRequests = int(math.Min(float64(maxOutRequests), float64(30)))
	maxOutRequests = int(math.Max(float64(maxOutRequests), float64(backgroundWorkerCount + 1)))
	fmt.Printf("Config: MaxOutRequests  %v \n", maxOutRequests)
	
	workcache.Start(backgroundWorkerCount, maxOutRequests)

	restapi.Start(rpcUrl, listenIpPort, restMaxActiveRequests)
}
