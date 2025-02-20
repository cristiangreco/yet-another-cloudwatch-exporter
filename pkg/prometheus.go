package exporter

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	cloudwatchAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_requests_total",
		Help: "Help is not implemented yet.",
	})
	cloudwatchAPIErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_request_errors",
		Help: "Help is not implemented yet.",
	})
	cloudwatchGetMetricDataAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_getmetricdata_requests_total",
		Help: "Help is not implemented yet.",
	})
	cloudwatchGetMetricStatisticsAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_getmetricstatistics_requests_total",
		Help: "Help is not implemented yet.",
	})
	resourceGroupTaggingAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_resourcegrouptaggingapi_requests_total",
		Help: "Help is not implemented yet.",
	})
	autoScalingAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_autoscalingapi_requests_total",
		Help: "Help is not implemented yet.",
	})
	targetGroupsAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_targetgroupapi_requests_total",
		Help: "Help is not implemented yet.",
	})
	apiGatewayAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_apigatewayapi_requests_total",
	})
	ec2APICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_ec2api_requests_total",
		Help: "Help is not implemented yet.",
	})
	dmsAPICounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "yace_cloudwatch_dmsapi_requests_total",
		Help: "Help is not implemented yet.",
	})
)

var replacer = strings.NewReplacer(
	" ", "_",
	",", "_",
	"\t", "_",
	"/", "_",
	"\\", "_",
	".", "_",
	"-", "_",
	":", "_",
	"=", "_",
	"“", "_",
	"@", "_",
	"<", "_",
	">", "_",
	"%", "_percent",
)
var splitRegexp = regexp.MustCompile(`([a-z0-9])([A-Z])`)

type PrometheusMetric struct {
	name             *string
	labels           map[string]string
	value            *float64
	includeTimestamp bool
	timestamp        time.Time
}

type PrometheusCollector struct {
	metrics []*PrometheusMetric
}

func NewPrometheusCollector(metrics []*PrometheusMetric) *PrometheusCollector {
	return &PrometheusCollector{
		metrics: removeDuplicatedMetrics(metrics),
	}
}

func (p *PrometheusCollector) Describe(descs chan<- *prometheus.Desc) {
	for _, metric := range p.metrics {
		descs <- createDesc(metric)
	}
}

func (p *PrometheusCollector) Collect(metrics chan<- prometheus.Metric) {
	for _, metric := range p.metrics {
		metrics <- createMetric(metric)
	}
}

func createDesc(metric *PrometheusMetric) *prometheus.Desc {
	return prometheus.NewDesc(
		*metric.name,
		"Help is not implemented yet.",
		nil,
		metric.labels,
	)
}

func createMetric(metric *PrometheusMetric) prometheus.Metric {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        *metric.name,
		Help:        "Help is not implemented yet.",
		ConstLabels: metric.labels,
	})

	gauge.Set(*metric.value)

	if !metric.includeTimestamp {
		return gauge
	}

	return prometheus.NewMetricWithTimestamp(metric.timestamp, gauge)
}

func removeDuplicatedMetrics(metrics []*PrometheusMetric) []*PrometheusMetric {
	keys := make(map[string]bool)
	filteredMetrics := []*PrometheusMetric{}
	for _, metric := range metrics {
		check := *metric.name + combineLabels(metric.labels)
		if _, value := keys[check]; !value {
			keys[check] = true
			filteredMetrics = append(filteredMetrics, metric)
		}
	}
	return filteredMetrics
}

func combineLabels(labels map[string]string) string {
	var combinedLabels string
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		combinedLabels += promString(k) + promString(labels[k])
	}
	return combinedLabels
}

func promString(text string) string {
	text = splitString(text)
	return strings.ToLower(sanitize(text))
}

func promStringTag(text string, labelsSnakeCase bool) string {
	if labelsSnakeCase {
		return promString(text)
	}
	return sanitize(text)
}

func sanitize(text string) string {
	return replacer.Replace(text)
}

func splitString(text string) string {
	return splitRegexp.ReplaceAllString(text, `$1.$2`)
}
