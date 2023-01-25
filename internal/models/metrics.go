package models

var (
	AllMetrics = [...]Metric{
		{"Alloc", Gauge, Value{}, RuntimeSource},
		{"BuckHashSys", Gauge, Value{}, RuntimeSource},
		{"Frees", Gauge, Value{}, RuntimeSource},
		{"GCCPUFraction", Gauge, Value{}, RuntimeSource},
		{"HeapAlloc", Gauge, Value{}, RuntimeSource},
		{"HeapIdle", Gauge, Value{}, RuntimeSource},
		{"HeapInuse", Gauge, Value{}, RuntimeSource},
		{"HeapObjects", Gauge, Value{}, RuntimeSource},
		{"HeapReleased", Gauge, Value{}, RuntimeSource},
		{"HeapSys", Gauge, Value{}, RuntimeSource},
		{"LastGC", Gauge, Value{}, RuntimeSource},
		{"Lookups", Gauge, Value{}, RuntimeSource},
		{"Lookups", Gauge, Value{}, RuntimeSource},
		{"MCacheSys", Gauge, Value{}, RuntimeSource},
		{"MSpanInuse", Gauge, Value{}, RuntimeSource},
		{"MSpanSys", Gauge, Value{}, RuntimeSource},
		{"Mallocs", Gauge, Value{}, RuntimeSource},
		{"NextGC", Gauge, Value{}, RuntimeSource},
		{"NumForcedGC", Gauge, Value{}, RuntimeSource},
		{"NumGC", Gauge, Value{}, RuntimeSource},
		{"OtherSys", Gauge, Value{}, RuntimeSource},
		{"PauseTotalNs", Gauge, Value{}, RuntimeSource},
		{"StackInuse", Gauge, Value{}, RuntimeSource},
		{"StackSys", Gauge, Value{}, RuntimeSource},
		{"Sys", Gauge, Value{}, RuntimeSource},
		{"TotalAlloc", Gauge, Value{}, RuntimeSource},
		{"PollCount", Counter, Value{}, CounterSource},
		{"RandomValue", Gauge, Value{}, RandomSource},
	}
)
