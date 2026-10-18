package common

import "sync/atomic"

// activeConnections counts requests in flight on this instance. StatsMiddleware
// updates it, and /api/status/test and the system instance heartbeat read it.
var activeConnections atomic.Int64

func AddActiveConnections(delta int64) {
	activeConnections.Add(delta)
}

func GetActiveConnections() int64 {
	return activeConnections.Load()
}
