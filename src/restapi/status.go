// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/catenocrypt/nano-work-cache/rpcclient"
	"github.com/catenocrypt/nano-work-cache/workcache"
)

var startTime time.Time = time.Now()

// Return the inner status info of the service, in Json string
func getStatus() string {
	cacheSize := workcache.StatusCacheSize()
	workOutReqCount := workcache.StatusWorkOutReqCount()
	workOutRespCount := workcache.StatusWorkOutRespCount()
	workOutDurAvg := workcache.StatusWorkOutDurationAvg()
	workInReqCount := workcache.StatusWorkInReqCount()
	workInReqFromCache := workcache.StatusWorkInReqFromCache()
	workInReqError := workcache.StatusWorkInReqError()
	var workInReqCacheRatio float32 = 0
	if workInReqCount > 0 {
		workInReqCacheRatio = float32(workInReqFromCache) / float32(workInReqCount)
	}
	activeHandlerCount := ActiveHandlerCount()
	activeWorkOutReqCount := workcache.StatusActiveWorkOutReqCount()
	pregenerQueSize := workcache.StatusPregenerQueueSize()
	uptime := time.Now().Sub(startTime)
	return fmt.Sprintf(`{"cache_size": %v, "work_in_req_count": %v, "work_in_req_from_cache": %v, "work_in_req_error": %v, "work_in_req_cache_ratio": %v, "work_out_req_count": %v, "work_out_resp_count": %v, "work_out_dur_avg": %v, "active_handler_count": %v, "active_work_out_req_count": %v, "pregenr_que_size": %v, "diff": "%v", "hrs", %v}`,
		cacheSize, workInReqCount, workInReqFromCache, workInReqError, workInReqCacheRatio, workOutReqCount, workOutRespCount, workOutDurAvg, activeHandlerCount, activeWorkOutReqCount, pregenerQueSize,
		strconv.FormatUint(rpcclient.GetDifficultyCached(nanoNodeUrl), 16), uptime.Hours())
}
