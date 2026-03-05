








package serviceconfig


type Config interface {
	isServiceConfig()
}



type LoadBalancingConfig interface {
	isLoadBalancingConfig()
}



type ParseResult struct {
	Config Config
	Err    error
}
