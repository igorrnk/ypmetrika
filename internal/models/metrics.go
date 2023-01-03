package models

var (
	AllMetrics = [...]Metric{
		{"Alloc", GaugeType, Value{}, RuntimeSource},
		{"BuckHashSys", GaugeType, Value{}, RuntimeSource},
		{"Frees", GaugeType, Value{}, RuntimeSource},
		{"GCCPUFraction", GaugeType, Value{}, RuntimeSource},
		{"HeapAlloc", GaugeType, Value{}, RuntimeSource},
		{"HeapIdle", GaugeType, Value{}, RuntimeSource},
		{"HeapInuse", GaugeType, Value{}, RuntimeSource},
		{"HeapObjects", GaugeType, Value{}, RuntimeSource},
		{"HeapReleased", GaugeType, Value{}, RuntimeSource},
		{"HeapSys", GaugeType, Value{}, RuntimeSource},
		{"LastGC", GaugeType, Value{}, RuntimeSource},
		{"Lookups", GaugeType, Value{}, RuntimeSource},
		{"Lookups", GaugeType, Value{}, RuntimeSource},
		{"MCacheSys", GaugeType, Value{}, RuntimeSource},
		{"MSpanInuse", GaugeType, Value{}, RuntimeSource},
		{"MSpanSys", GaugeType, Value{}, RuntimeSource},
		{"Mallocs", GaugeType, Value{}, RuntimeSource},
		{"NextGC", GaugeType, Value{}, RuntimeSource},
		{"NumForcedGC", GaugeType, Value{}, RuntimeSource},
		{"NumGC", GaugeType, Value{}, RuntimeSource},
		{"OtherSys", GaugeType, Value{}, RuntimeSource},
		{"PauseTotalNs", GaugeType, Value{}, RuntimeSource},
		{"StackInuse", GaugeType, Value{}, RuntimeSource},
		{"StackSys", GaugeType, Value{}, RuntimeSource},
		{"Sys", GaugeType, Value{}, RuntimeSource},
		{"TotalAlloc", GaugeType, Value{}, RuntimeSource},
		{"PollCount", CounterType, Value{}, CounterSource},
		{"RandomValue", GaugeType, Value{}, RandomSource},
	}
)
