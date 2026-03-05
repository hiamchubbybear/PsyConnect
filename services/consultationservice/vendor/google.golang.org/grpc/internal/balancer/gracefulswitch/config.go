

package gracefulswitch

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/serviceconfig"
)

type lbConfig struct {
	serviceconfig.LoadBalancingConfig

	childBuilder balancer.Builder
	childConfig  serviceconfig.LoadBalancingConfig
}



func ChildName(l serviceconfig.LoadBalancingConfig) string {
	return l.(*lbConfig).childBuilder.Name()
}








func ParseConfig(cfg json.RawMessage) (serviceconfig.LoadBalancingConfig, error) {
	var lbCfg []map[string]json.RawMessage
	if err := json.Unmarshal(cfg, &lbCfg); err != nil {
		return nil, err
	}
	for i, e := range lbCfg {
		if len(e) != 1 {
			return nil, fmt.Errorf("expected a JSON struct with one entry; received entry %v at index %d", e, i)
		}

		var name string
		var jsonCfg json.RawMessage
		for name, jsonCfg = range e {
		}

		builder := balancer.Get(name)
		if builder == nil {
			
			continue
		}

		parser, ok := builder.(balancer.ConfigParser)
		if !ok {
			
			return &lbConfig{childBuilder: builder}, nil
		}

		cfg, err := parser.ParseConfig(jsonCfg)
		if err != nil {
			return nil, fmt.Errorf("error parsing config for policy %q: %v", name, err)
		}
		return &lbConfig{childBuilder: builder, childConfig: cfg}, nil
	}

	return nil, fmt.Errorf("no supported policies found in config: %v", string(cfg))
}
