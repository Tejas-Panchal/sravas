package job

import "log"

// AggregateAnalytics rolls up raw events into daily video and channel stats
func AggregateAnalytics() {
	log.Println("[scheduler] aggregate analytics: rolling up yesterday's events")
	// TODO: query analytics events from yesterday
	// TODO: aggregate per-video: total views, watch time, likes, comments
	// TODO: aggregate per-channel: subscriber change, total views, top videos
	// TODO: write daily snapshot to analytics_daily table
}
