

package stats

import (
	"maps"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/stats"
)

func init() {
	internal.SnapshotMetricRegistryForTesting = snapshotMetricsRegistryForTesting
}

var logger = grpclog.Component("metrics-registry")




var DefaultMetrics = stats.NewMetricSet()


type MetricDescriptor struct {
	
	
	
	
	Name string
	
	Description string
	
	Unit string
	
	
	Labels []string
	
	
	OptionalLabels []string
	
	Default bool
	
	
	Type MetricType
	
	
	
	Bounds []float64
}


type MetricType int


const (
	MetricTypeIntCount MetricType = iota
	MetricTypeFloatCount
	MetricTypeIntHisto
	MetricTypeFloatHisto
	MetricTypeIntGauge
)




type Int64CountHandle MetricDescriptor



func (h *Int64CountHandle) Descriptor() *MetricDescriptor {
	return (*MetricDescriptor)(h)
}


func (h *Int64CountHandle) Record(recorder MetricsRecorder, incr int64, labels ...string) {
	recorder.RecordInt64Count(h, incr, labels...)
}



type Float64CountHandle MetricDescriptor



func (h *Float64CountHandle) Descriptor() *MetricDescriptor {
	return (*MetricDescriptor)(h)
}


func (h *Float64CountHandle) Record(recorder MetricsRecorder, incr float64, labels ...string) {
	recorder.RecordFloat64Count(h, incr, labels...)
}



type Int64HistoHandle MetricDescriptor



func (h *Int64HistoHandle) Descriptor() *MetricDescriptor {
	return (*MetricDescriptor)(h)
}


func (h *Int64HistoHandle) Record(recorder MetricsRecorder, incr int64, labels ...string) {
	recorder.RecordInt64Histo(h, incr, labels...)
}




type Float64HistoHandle MetricDescriptor



func (h *Float64HistoHandle) Descriptor() *MetricDescriptor {
	return (*MetricDescriptor)(h)
}


func (h *Float64HistoHandle) Record(recorder MetricsRecorder, incr float64, labels ...string) {
	recorder.RecordFloat64Histo(h, incr, labels...)
}



type Int64GaugeHandle MetricDescriptor



func (h *Int64GaugeHandle) Descriptor() *MetricDescriptor {
	return (*MetricDescriptor)(h)
}


func (h *Int64GaugeHandle) Record(recorder MetricsRecorder, incr int64, labels ...string) {
	recorder.RecordInt64Gauge(h, incr, labels...)
}


var registeredMetrics = make(map[string]bool)




var metricsRegistry = make(map[string]*MetricDescriptor)




func DescriptorForMetric(metricName string) *MetricDescriptor {
	return metricsRegistry[metricName]
}

func registerMetric(metricName string, def bool) {
	if registeredMetrics[metricName] {
		logger.Fatalf("metric %v already registered", metricName)
	}
	registeredMetrics[metricName] = true
	if def {
		DefaultMetrics = DefaultMetrics.Add(metricName)
	}
}







func RegisterInt64Count(descriptor MetricDescriptor) *Int64CountHandle {
	registerMetric(descriptor.Name, descriptor.Default)
	descriptor.Type = MetricTypeIntCount
	descPtr := &descriptor
	metricsRegistry[descriptor.Name] = descPtr
	return (*Int64CountHandle)(descPtr)
}







func RegisterFloat64Count(descriptor MetricDescriptor) *Float64CountHandle {
	registerMetric(descriptor.Name, descriptor.Default)
	descriptor.Type = MetricTypeFloatCount
	descPtr := &descriptor
	metricsRegistry[descriptor.Name] = descPtr
	return (*Float64CountHandle)(descPtr)
}







func RegisterInt64Histo(descriptor MetricDescriptor) *Int64HistoHandle {
	registerMetric(descriptor.Name, descriptor.Default)
	descriptor.Type = MetricTypeIntHisto
	descPtr := &descriptor
	metricsRegistry[descriptor.Name] = descPtr
	return (*Int64HistoHandle)(descPtr)
}







func RegisterFloat64Histo(descriptor MetricDescriptor) *Float64HistoHandle {
	registerMetric(descriptor.Name, descriptor.Default)
	descriptor.Type = MetricTypeFloatHisto
	descPtr := &descriptor
	metricsRegistry[descriptor.Name] = descPtr
	return (*Float64HistoHandle)(descPtr)
}







func RegisterInt64Gauge(descriptor MetricDescriptor) *Int64GaugeHandle {
	registerMetric(descriptor.Name, descriptor.Default)
	descriptor.Type = MetricTypeIntGauge
	descPtr := &descriptor
	metricsRegistry[descriptor.Name] = descPtr
	return (*Int64GaugeHandle)(descPtr)
}




func snapshotMetricsRegistryForTesting() func() {
	oldDefaultMetrics := DefaultMetrics
	oldRegisteredMetrics := registeredMetrics
	oldMetricsRegistry := metricsRegistry

	registeredMetrics = make(map[string]bool)
	metricsRegistry = make(map[string]*MetricDescriptor)
	maps.Copy(registeredMetrics, registeredMetrics)
	maps.Copy(metricsRegistry, metricsRegistry)

	return func() {
		DefaultMetrics = oldDefaultMetrics
		registeredMetrics = oldRegisteredMetrics
		metricsRegistry = oldMetricsRegistry
	}
}
