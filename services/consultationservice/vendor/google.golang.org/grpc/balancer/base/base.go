













package base

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)


type PickerBuilder interface {
	
	Build(info PickerBuildInfo) balancer.Picker
}



type PickerBuildInfo struct {
	
	
	ReadySCs map[balancer.SubConn]SubConnInfo
}



type SubConnInfo struct {
	Address resolver.Address 
}


type Config struct {
	
	HealthCheck bool
}


func NewBalancerBuilder(name string, pb PickerBuilder, config Config) balancer.Builder {
	return &baseBuilder{
		name:          name,
		pickerBuilder: pb,
		config:        config,
	}
}
