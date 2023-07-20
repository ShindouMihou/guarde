package healthcheck

import "sync/atomic"

const (
	AverageRequestSize     AverageKey = "request_size"
	AverageResponseSize    AverageKey = "response_size"
	AverageRequestLatency  AverageKey = "request_latency"
	AverageForwardLatency  AverageKey = "forward_latency"
	AverageFallbackLatency AverageKey = "fallback_latency"
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
	AverageRequestLatency:  NewAverage[int](25),
	AverageForwardLatency:  NewAverage[int](25),
	AverageFallbackLatency: NewAverage[int](25),
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

func report() map[string]any {
	reports := make(map[string]any)
	setAverages(reports, AverageRequestSize, AverageResponseSize, AverageRequestLatency, AverageForwardLatency, AverageFallbackLatency)
	setCounters(reports, RequestsForwarded, RequestsRejected, RequestsReceived, RequestsErrored, WhoIsRequests, FallbackFailures)
	return reports
}
func setAverages(m map[string]any, keys ...AverageKey) {
	for _, key := range keys {
		avg := averages[key]
		m[string(key)] = map[string]any{
			"avg":    int64(avg.Stale()),
			"values": avg.Values(),
		}
	}
}

func setCounters(m map[string]any, keys ...CounterKey) {
	for _, key := range keys {
		m[string(key)] = counters[key].Load()
	}
}
