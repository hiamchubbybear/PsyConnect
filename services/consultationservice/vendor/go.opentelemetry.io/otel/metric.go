


package otel 

import (
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
)














func Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return GetMeterProvider().Meter(name, opts...)
}








func GetMeterProvider() metric.MeterProvider {
	return global.MeterProvider()
}


func SetMeterProvider(mp metric.MeterProvider) {
	global.SetMeterProvider(mp)
}
