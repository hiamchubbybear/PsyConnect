

package grpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/pickfirst"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/balancer/gracefulswitch"
	internalserviceconfig "google.golang.org/grpc/internal/serviceconfig"
	"google.golang.org/grpc/serviceconfig"
)

const maxInt = int(^uint(0) >> 1)







type MethodConfig = internalserviceconfig.MethodConfig







type ServiceConfig struct {
	serviceconfig.Config

	
	
	lbConfig serviceconfig.LoadBalancingConfig

	
	
	
	
	
	
	Methods map[string]MethodConfig

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	retryThrottling *retryThrottlingPolicy
	
	
	healthCheckConfig *healthCheckConfig
	
	
	rawJSONString string
}


type healthCheckConfig struct {
	
	ServiceName string
}

type jsonRetryPolicy struct {
	MaxAttempts          int
	InitialBackoff       internalserviceconfig.Duration
	MaxBackoff           internalserviceconfig.Duration
	BackoffMultiplier    float64
	RetryableStatusCodes []codes.Code
}




type retryThrottlingPolicy struct {
	
	
	
	
	MaxTokens float64
	
	
	
	
	
	TokenRatio float64
}

type jsonName struct {
	Service string
	Method  string
}

var (
	errDuplicatedName             = errors.New("duplicated name")
	errEmptyServiceNonEmptyMethod = errors.New("cannot combine empty 'service' and non-empty 'method'")
)

func (j jsonName) generatePath() (string, error) {
	if j.Service == "" {
		if j.Method != "" {
			return "", errEmptyServiceNonEmptyMethod
		}
		return "", nil
	}
	res := "/" + j.Service + "/"
	if j.Method != "" {
		res += j.Method
	}
	return res, nil
}


type jsonMC struct {
	Name                    *[]jsonName
	WaitForReady            *bool
	Timeout                 *internalserviceconfig.Duration
	MaxRequestMessageBytes  *int64
	MaxResponseMessageBytes *int64
	RetryPolicy             *jsonRetryPolicy
}


type jsonSC struct {
	LoadBalancingPolicy *string
	LoadBalancingConfig *json.RawMessage
	MethodConfig        *[]jsonMC
	RetryThrottling     *retryThrottlingPolicy
	HealthCheckConfig   *healthCheckConfig
}

func init() {
	internal.ParseServiceConfig = func(js string) *serviceconfig.ParseResult {
		return parseServiceConfig(js, defaultMaxCallAttempts)
	}
}

func parseServiceConfig(js string, maxAttempts int) *serviceconfig.ParseResult {
	if len(js) == 0 {
		return &serviceconfig.ParseResult{Err: fmt.Errorf("no JSON service config provided")}
	}
	var rsc jsonSC
	err := json.Unmarshal([]byte(js), &rsc)
	if err != nil {
		logger.Warningf("grpc: unmarshalling service config %s: %v", js, err)
		return &serviceconfig.ParseResult{Err: err}
	}
	sc := ServiceConfig{
		Methods:           make(map[string]MethodConfig),
		retryThrottling:   rsc.RetryThrottling,
		healthCheckConfig: rsc.HealthCheckConfig,
		rawJSONString:     js,
	}
	c := rsc.LoadBalancingConfig
	if c == nil {
		name := pickfirst.Name
		if rsc.LoadBalancingPolicy != nil {
			name = *rsc.LoadBalancingPolicy
		}
		if balancer.Get(name) == nil {
			name = pickfirst.Name
		}
		cfg := []map[string]any{{name: struct{}{}}}
		strCfg, err := json.Marshal(cfg)
		if err != nil {
			return &serviceconfig.ParseResult{Err: fmt.Errorf("unexpected error marshaling simple LB config: %w", err)}
		}
		r := json.RawMessage(strCfg)
		c = &r
	}
	cfg, err := gracefulswitch.ParseConfig(*c)
	if err != nil {
		return &serviceconfig.ParseResult{Err: err}
	}
	sc.lbConfig = cfg

	if rsc.MethodConfig == nil {
		return &serviceconfig.ParseResult{Config: &sc}
	}

	paths := map[string]struct{}{}
	for _, m := range *rsc.MethodConfig {
		if m.Name == nil {
			continue
		}

		mc := MethodConfig{
			WaitForReady: m.WaitForReady,
			Timeout:      (*time.Duration)(m.Timeout),
		}
		if mc.RetryPolicy, err = convertRetryPolicy(m.RetryPolicy, maxAttempts); err != nil {
			logger.Warningf("grpc: unmarshalling service config %s: %v", js, err)
			return &serviceconfig.ParseResult{Err: err}
		}
		if m.MaxRequestMessageBytes != nil {
			if *m.MaxRequestMessageBytes > int64(maxInt) {
				mc.MaxReqSize = newInt(maxInt)
			} else {
				mc.MaxReqSize = newInt(int(*m.MaxRequestMessageBytes))
			}
		}
		if m.MaxResponseMessageBytes != nil {
			if *m.MaxResponseMessageBytes > int64(maxInt) {
				mc.MaxRespSize = newInt(maxInt)
			} else {
				mc.MaxRespSize = newInt(int(*m.MaxResponseMessageBytes))
			}
		}
		for i, n := range *m.Name {
			path, err := n.generatePath()
			if err != nil {
				logger.Warningf("grpc: error unmarshalling service config %s due to methodConfig[%d]: %v", js, i, err)
				return &serviceconfig.ParseResult{Err: err}
			}

			if _, ok := paths[path]; ok {
				err = errDuplicatedName
				logger.Warningf("grpc: error unmarshalling service config %s due to methodConfig[%d]: %v", js, i, err)
				return &serviceconfig.ParseResult{Err: err}
			}
			paths[path] = struct{}{}
			sc.Methods[path] = mc
		}
	}

	if sc.retryThrottling != nil {
		if mt := sc.retryThrottling.MaxTokens; mt <= 0 || mt > 1000 {
			return &serviceconfig.ParseResult{Err: fmt.Errorf("invalid retry throttling config: maxTokens (%v) out of range (0, 1000]", mt)}
		}
		if tr := sc.retryThrottling.TokenRatio; tr <= 0 {
			return &serviceconfig.ParseResult{Err: fmt.Errorf("invalid retry throttling config: tokenRatio (%v) may not be negative", tr)}
		}
	}
	return &serviceconfig.ParseResult{Config: &sc}
}

func isValidRetryPolicy(jrp *jsonRetryPolicy) bool {
	return jrp.MaxAttempts > 1 &&
		jrp.InitialBackoff > 0 &&
		jrp.MaxBackoff > 0 &&
		jrp.BackoffMultiplier > 0 &&
		len(jrp.RetryableStatusCodes) > 0
}

func convertRetryPolicy(jrp *jsonRetryPolicy, maxAttempts int) (p *internalserviceconfig.RetryPolicy, err error) {
	if jrp == nil {
		return nil, nil
	}

	if !isValidRetryPolicy(jrp) {
		return nil, fmt.Errorf("invalid retry policy (%+v): ", jrp)
	}

	if jrp.MaxAttempts < maxAttempts {
		maxAttempts = jrp.MaxAttempts
	}
	rp := &internalserviceconfig.RetryPolicy{
		MaxAttempts:          maxAttempts,
		InitialBackoff:       time.Duration(jrp.InitialBackoff),
		MaxBackoff:           time.Duration(jrp.MaxBackoff),
		BackoffMultiplier:    jrp.BackoffMultiplier,
		RetryableStatusCodes: make(map[codes.Code]bool),
	}
	for _, code := range jrp.RetryableStatusCodes {
		rp.RetryableStatusCodes[code] = true
	}
	return rp, nil
}

func minPointers(a, b *int) *int {
	if *a < *b {
		return a
	}
	return b
}

func getMaxSize(mcMax, doptMax *int, defaultVal int) *int {
	if mcMax == nil && doptMax == nil {
		return &defaultVal
	}
	if mcMax != nil && doptMax != nil {
		return minPointers(mcMax, doptMax)
	}
	if mcMax != nil {
		return mcMax
	}
	return doptMax
}

func newInt(b int) *int {
	return &b
}

func init() {
	internal.EqualServiceConfigForTesting = equalServiceConfig
}





func equalServiceConfig(a, b serviceconfig.Config) bool {
	if a == nil && b == nil {
		return true
	}
	aa, ok := a.(*ServiceConfig)
	if !ok {
		return false
	}
	bb, ok := b.(*ServiceConfig)
	if !ok {
		return false
	}
	aaRaw := aa.rawJSONString
	aa.rawJSONString = ""
	bbRaw := bb.rawJSONString
	bb.rawJSONString = ""
	defer func() {
		aa.rawJSONString = aaRaw
		bb.rawJSONString = bbRaw
	}()
	
	
	
	return reflect.DeepEqual(aa, bb)
}
