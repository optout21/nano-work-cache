// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"github.com/catenocrypt/nano-work-cache/workcache"
	"fmt"
)

// Return the inner status info of the service, in Json string
func getStatus() string {
	cacheSize := workcache.StatusCacheSize()
	workOutReqCount := workcache.StatusWorkOutReqCount()
	workOutRespCount := workcache.StatusWorkOutRespCount()
	workInReqCount := workcache.StatusWorkInReqCount()
	workInReqFromCache := workcache.StatusWorkInReqFromCache()
	var workInReqCacheRatio float32 = 0
	if workInReqCount > 0 { workInReqCacheRatio = float32(workInReqFromCache) / float32(workInReqCount) }
	activeHandlerCount := ActiveHandlerCount()
	activeWorkOutReqCount := workcache.StatusActiveWorkOutReqCount()
	pregenerQueSize := workcache.StatusPregenerQueueSize()
	return fmt.Sprintf(`{"cache_size": %v, "work_in_req_count": %v, "work_in_req_from_cache": %v, "work_in_req_cache_ratio": %v, "work_out_req_count": %v, "work_out_resp_count": %v, "active_handler_count": %v, "active_work_out_req_count": %v, "pregenr_que_size": %v}`,
		cacheSize, workInReqCount, workInReqFromCache, workInReqCacheRatio, workOutReqCount, workOutRespCount, activeHandlerCount, activeWorkOutReqCount, pregenerQueSize)
}
