package healthcheck

import "sync/atomic"

const (
	AverageRequestSize     AverageKey = "avg_request_size"
	AverageResponseSize    AverageKey = "avg_response_size"
	AverageRequestLatency  AverageKey = "avg_request_latency"
	AverageForwardLatency  AverageKey = "avg_forward_latency"
	AverageFallbackLatency AverageKey = "avg_fallback_latency"
	RequestsForwarded      CounterKey = "requests_forwarded"
	RequestsRejected       CounterKey = "requests_rejected"
	RequestsReceived       CounterKey = "requests_received"
	RequestsErrored        CounterKey = "requests_errored"
	FallbackFailures       CounterKey = "fallback_fails"
	WhoIsRequests          CounterKey = "whois_requests"
)

type AverageKey string
type CounterKey string

var averages = map[AverageKey]*Average[int]{
	AverageRequestSize:     NewAverage[int](10),
	AverageResponseSize:    NewAverage[int](10),
	AverageRequestLatency:  NewAverage[int](50),
	AverageForwardLatency:  NewAverage[int](50),
	AverageFallbackLatency: NewAverage[int](50),
}

var counters = map[CounterKey]*atomic.Int64{
	RequestsForwarded: {},
	RequestsRejected:  {},
	RequestsReceived:  {},
	RequestsErrored:   {},
	WhoIsRequests:     {},
	FallbackFailures:  {},
}

func (key AverageKey) Add(value int) {
	go averages[key].Add(value)
}

func (key CounterKey) Increment() {
	counters[key].Add(1)
}

func (key AverageKey) Stale() int {
	return averages[key].stale
}

func report() map[string]int64 {
	reports := make(map[string]int64)
	setAverages(reports, AverageRequestSize, AverageResponseSize, AverageRequestLatency, AverageForwardLatency, AverageFallbackLatency)
	setCounters(reports, RequestsForwarded, RequestsRejected, RequestsReceived, RequestsErrored, WhoIsRequests, FallbackFailures)
	return reports
}
func setAverages(m map[string]int64, keys ...AverageKey) {
	for _, key := range keys {
		m[string(key)] = int64(averages[key].Stale())
	}
}

func setCounters(m map[string]int64, keys ...CounterKey) {
	for _, key := range keys {
		m[string(key)] = counters[key].Load()
	}
}
