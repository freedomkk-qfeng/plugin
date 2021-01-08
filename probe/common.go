package main

import (
	"strings"
)

type MetricValue struct {
	Metric      string  `json:"metric"`
	Timestamp   int64   `json:"timestamp"`
	Value       float64 `json:"value"`
	CounterType string  `json:"counterType"`
	Tags        string  `json:"tags"`
}

func NewMetricValue(ts int64, metric string, val float64, dataType string, tags ...string) *MetricValue {
	mv := MetricValue{
		Metric:      metric,
		Timestamp:   ts,
		Value:       val,
		CounterType: dataType,
	}

	size := len(tags)

	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return &mv
}

func GaugeValue(ts int64, metric string, val float64, tags ...string) *MetricValue {
	return NewMetricValue(ts, metric, val, "GAUGE", tags...)
}

func CounterValue(ts int64, metric string, val float64, tags ...string) *MetricValue {
	return NewMetricValue(ts, metric, val, "COUNTER", tags...)
}
