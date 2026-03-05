












package dns

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
)


type Resolver struct {
	
	LookupSRV func(string, string, string) (string, []*net.SRV, error)
	LookupTXT func(string) ([]string, error)
}


var DefaultResolver = &Resolver{net.LookupSRV, net.LookupTXT}


func (r *Resolver) ParseHosts(host string, srvName string, stopOnErr bool) ([]string, error) {
	parsedHosts := strings.Split(host, ",")

	if len(parsedHosts) != 1 {
		return nil, fmt.Errorf("URI with SRV must include one and only one hostname")
	}
	return r.fetchSeedlistFromSRV(parsedHosts[0], srvName, stopOnErr)
}


func (r *Resolver) GetConnectionArgsFromTXT(host string) ([]string, error) {
	var connectionArgsFromTXT []string

	
	
	recordsFromTXT, _ := r.LookupTXT(host)

	
	
	
	if runtime.GOOS == "windows" {
		recordsFromTXT = []string{strings.Join(recordsFromTXT, "")}
	}

	if len(recordsFromTXT) > 1 {
		return nil, errors.New("multiple records from TXT not supported")
	}
	if len(recordsFromTXT) > 0 {
		connectionArgsFromTXT = strings.FieldsFunc(recordsFromTXT[0], func(r rune) bool { return r == ';' || r == '&' })

		err := validateTXTResult(connectionArgsFromTXT)
		if err != nil {
			return nil, err
		}
	}

	return connectionArgsFromTXT, nil
}

func (r *Resolver) fetchSeedlistFromSRV(host string, srvName string, stopOnErr bool) ([]string, error) {
	var err error

	_, _, err = net.SplitHostPort(host)

	if err == nil {
		
		
		return nil, fmt.Errorf("URI with srv must not include a port number")
	}

	
	if srvName == "" {
		srvName = "mongodb"
	}
	_, addresses, err := r.LookupSRV(srvName, "tcp", host)
	if err != nil && strings.Contains(err.Error(), "cannot unmarshal DNS message") {
		return nil, fmt.Errorf("see https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#hdr-Potential_DNS_Issues: %w", err)
	} else if err != nil {
		return nil, err
	}

	trimmedHost := strings.TrimSuffix(host, ".")

	parsedHosts := make([]string, 0, len(addresses))
	for _, address := range addresses {
		trimmedAddressTarget := strings.TrimSuffix(address.Target, ".")
		err := validateSRVResult(trimmedAddressTarget, trimmedHost)
		if err != nil {
			if stopOnErr {
				return nil, err
			}
			continue
		}
		parsedHosts = append(parsedHosts, fmt.Sprintf("%s:%d", trimmedAddressTarget, address.Port))
	}
	return parsedHosts, nil
}

func validateSRVResult(recordFromSRV, inputHostName string) error {
	separatedInputDomain := strings.Split(strings.ToLower(inputHostName), ".")
	separatedRecord := strings.Split(strings.ToLower(recordFromSRV), ".")
	if len(separatedRecord) < 2 {
		return errors.New("DNS name must contain at least 2 labels")
	}
	if len(separatedRecord) < len(separatedInputDomain) {
		return errors.New("Domain suffix from SRV record not matched input domain")
	}

	inputDomainSuffix := separatedInputDomain[1:]
	domainSuffixOffset := len(separatedRecord) - (len(separatedInputDomain) - 1)

	recordDomainSuffix := separatedRecord[domainSuffixOffset:]
	for ix, label := range inputDomainSuffix {
		if label != recordDomainSuffix[ix] {
			return errors.New("Domain suffix from SRV record not matched input domain")
		}
	}
	return nil
}

var allowedTXTOptions = map[string]struct{}{
	"authsource":   {},
	"replicaset":   {},
	"loadbalanced": {},
}

func validateTXTResult(paramsFromTXT []string) error {
	for _, param := range paramsFromTXT {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) != 2 {
			return errors.New("Invalid TXT record")
		}
		key := strings.ToLower(kv[0])
		if _, ok := allowedTXTOptions[key]; !ok {
			return fmt.Errorf("Cannot specify option '%s' in TXT record", kv[0])
		}
	}
	return nil
}
