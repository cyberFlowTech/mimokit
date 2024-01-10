package metric

import (
	"github.com/zeromicro/go-zero/core/metric"
)

// 定义Counter
var zapryApiCounter = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "zapry",
	Subsystem: "api_counter",
	Name:      "ret_total",
	Help:      "zapry apis stat.",
	Labels:    []string{"PATH", "Result"},
})

var zapryApiHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
	Namespace: "zapry",
	Subsystem: "api_histogram",
	Name:      "duration_ms",
	Help:      "api execution time cost stat.",
	Labels:    []string{"PATH", "Result"},
	Buckets:   []float64{5, 10, 30, 60, 100, 200, 300, 400, 500, 1000},
})

// 执行结果和耗时上报
func ZapryApiStatReport(v int64, labels ...string) {
	// 用户每次访问资源通过如下方法增加计数
	zapryApiCounter.Inc(labels...)

	zapryApiHistogram.Observe(v, labels...)
}
