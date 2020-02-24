// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"github.com/catenocrypt/nano-work-cache/workcache"
	"fmt"
)

// Return the inner status info of the service, in Json string
func getStatus() string {
	cacheSize := workcache.StatusCacheSize()
	workReqCount := workcache.StatusWorkReqCount()
	workRespCount := workcache.StatusWorkRespCount()
	activeHandlerCount := ActiveHandlerCount()
	activeWorkOutReqCount := workcache.StatusActiveWorkOutReqCount()
	pregenerQueSize := workcache.StatusPregenerQueueSize()
	return fmt.Sprintf(`{"cache_size": %v, "work_req_count": %v, "work_resp_count": %v, "active_handler_count": %v, "active_work_out_req_count": %v, "pregenr_que_size": %v}`,
		cacheSize, workReqCount, workRespCount, activeHandlerCount, activeWorkOutReqCount, pregenerQueSize)
}
