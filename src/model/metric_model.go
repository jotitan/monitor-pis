package model

type MetricsRequest struct {
	Metrics map[string]MetricPoint
	Name string
}

type MetricPoint struct{
	Timestamp int64
	Value float32
}

func NewMetricPoint(timestamp int64, value float32)MetricPoint{
	return MetricPoint{Timestamp: timestamp,Value:value}
}