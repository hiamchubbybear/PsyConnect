

package stats

import "maps"






type MetricSet struct {
	
	metrics map[string]bool
}


func NewMetricSet(metricNames ...string) *MetricSet {
	newMetrics := make(map[string]bool)
	for _, metric := range metricNames {
		newMetrics[metric] = true
	}
	return &MetricSet{metrics: newMetrics}
}



func (m *MetricSet) Metrics() map[string]bool {
	return m.metrics
}



func (m *MetricSet) Add(metricNames ...string) *MetricSet {
	newMetrics := make(map[string]bool)
	for metric := range m.metrics {
		newMetrics[metric] = true
	}

	for _, metric := range metricNames {
		newMetrics[metric] = true
	}
	return &MetricSet{metrics: newMetrics}
}



func (m *MetricSet) Join(metrics *MetricSet) *MetricSet {
	newMetrics := make(map[string]bool)
	maps.Copy(newMetrics, m.metrics)
	maps.Copy(newMetrics, metrics.metrics)
	return &MetricSet{metrics: newMetrics}
}



func (m *MetricSet) Remove(metricNames ...string) *MetricSet {
	newMetrics := make(map[string]bool)
	for metric := range m.metrics {
		newMetrics[metric] = true
	}

	for _, metric := range metricNames {
		delete(newMetrics, metric)
	}
	return &MetricSet{metrics: newMetrics}
}
