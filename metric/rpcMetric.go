package metric

import (
	"github.com/zeromicro/go-zero/core/metric"
)

// 定义Counter
var zapryRpcCounter = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zapry",
	Subsystem: "rpc_counter",
	Name:      "ret_total",
	Help:      "zapry rpc stat.",
	Labels:    []string{"ServerName", "Path", "Result", "Client"},
})

var zapryRpcHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
	Namespace: "zapry",
	Subsystem: "rpc_histogram",
	Name:      "duration_ms",
	Help:      "rpc execution time cost stat.",
	Labels:    []string{"ServerName", "Path", "Result", "Client"},
	Buckets:   []float64{5, 10, 30, 60, 100, 200, 300, 400, 500, 1000},
})

// rpc结果和耗时上报
func ZapryRpcStatReport(v int64, labels ...string) {
	// 用户每次访问资源通过如下方法增加计数
	zapryRpcCounter.Inc(labels...)

	zapryRpcHistogram.Observe(v, labels...)
}
