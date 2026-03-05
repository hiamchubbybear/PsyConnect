


package serviceconfig

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	externalserviceconfig "google.golang.org/grpc/serviceconfig"
)

var logger = grpclog.Component("core")








type BalancerConfig struct {
	Name   string
	Config externalserviceconfig.LoadBalancingConfig
}

type intermediateBalancerConfig []map[string]json.RawMessage





func (bc *BalancerConfig) MarshalJSON() ([]byte, error) {
	if bc.Config == nil {
		
		return []byte(fmt.Sprintf(`[{%q: %v}]`, bc.Name, "{}")), nil
	}
	c, err := json.Marshal(bc.Config)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`[{%q: %s}]`, bc.Name, c)), nil
}










func (bc *BalancerConfig) UnmarshalJSON(b []byte) error {
	var ir intermediateBalancerConfig
	err := json.Unmarshal(b, &ir)
	if err != nil {
		return err
	}

	var names []string
	for i, lbcfg := range ir {
		if len(lbcfg) != 1 {
			return fmt.Errorf("invalid loadBalancingConfig: entry %v does not contain exactly 1 policy/config pair: %q", i, lbcfg)
		}

		var (
			name    string
			jsonCfg json.RawMessage
		)
		
		
		for name, jsonCfg = range lbcfg {
		}

		names = append(names, name)
		builder := balancer.Get(name)
		if builder == nil {
			
			
			continue
		}
		bc.Name = name

		parser, ok := builder.(balancer.ConfigParser)
		if !ok {
			if string(jsonCfg) != "{}" {
				logger.Warningf("non-empty balancer configuration %q, but balancer does not implement ParseConfig", string(jsonCfg))
			}
			
			return nil
		}

		cfg, err := parser.ParseConfig(jsonCfg)
		if err != nil {
			return fmt.Errorf("error parsing loadBalancingConfig for policy %q: %v", name, err)
		}
		bc.Config = cfg
		return nil
	}
	
	
	
	
	return fmt.Errorf("invalid loadBalancingConfig: no supported policies found in %v", names)
}



type MethodConfig struct {
	
	
	
	WaitForReady *bool
	
	
	
	
	Timeout *time.Duration
	
	
	
	
	
	
	MaxReqSize *int
	
	
	MaxRespSize *int
	
	RetryPolicy *RetryPolicy
}




type RetryPolicy struct {
	
	
	
	MaxAttempts int

	
	
	
	
	
	
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64

	
	
	
	
	
	
	RetryableStatusCodes map[codes.Code]bool
}
