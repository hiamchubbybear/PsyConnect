


package metric 

import "go.opentelemetry.io/otel/attribute"


type MeterConfig struct {
	instrumentationVersion string
	schemaURL              string
	attrs                  attribute.Set

	
	noCmp [0]func() 
}



func (cfg MeterConfig) InstrumentationVersion() string {
	return cfg.instrumentationVersion
}



func (cfg MeterConfig) InstrumentationAttributes() attribute.Set {
	return cfg.attrs
}


func (cfg MeterConfig) SchemaURL() string {
	return cfg.schemaURL
}


type MeterOption interface {
	
	applyMeter(MeterConfig) MeterConfig
}



func NewMeterConfig(opts ...MeterOption) MeterConfig {
	var config MeterConfig
	for _, o := range opts {
		config = o.applyMeter(config)
	}
	return config
}

type meterOptionFunc func(MeterConfig) MeterConfig

func (fn meterOptionFunc) applyMeter(cfg MeterConfig) MeterConfig {
	return fn(cfg)
}


func WithInstrumentationVersion(version string) MeterOption {
	return meterOptionFunc(func(config MeterConfig) MeterConfig {
		config.instrumentationVersion = version
		return config
	})
}




func WithInstrumentationAttributes(attr ...attribute.KeyValue) MeterOption {
	return meterOptionFunc(func(config MeterConfig) MeterConfig {
		config.attrs = attribute.NewSet(attr...)
		return config
	})
}


func WithSchemaURL(schemaURL string) MeterOption {
	return meterOptionFunc(func(config MeterConfig) MeterConfig {
		config.schemaURL = schemaURL
		return config
	})
}
