












package ocsp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"golang.org/x/crypto/ocsp"
	"golang.org/x/sync/errgroup"
)

var (
	tlsFeatureExtensionOID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 1, 24}
	mustStapleFeatureValue = big.NewInt(5)
)


type Error struct {
	wrapped error
}


func (e *Error) Error() string {
	return fmt.Sprintf("OCSP verification failed: %v", e.wrapped)
}


func (e *Error) Unwrap() error {
	return e.wrapped
}

func newOCSPError(wrapped error) error {
	return &Error{wrapped: wrapped}
}



type ResponseDetails struct {
	Status     int
	NextUpdate time.Time
}

func extractResponseDetails(res *ocsp.Response) *ResponseDetails {
	return &ResponseDetails{
		Status:     res.Status,
		NextUpdate: res.NextUpdate,
	}
}


func Verify(ctx context.Context, connState tls.ConnectionState, opts *VerifyOptions) error {
	if opts.Cache == nil {
		
		
		
		return newOCSPError(errors.New("no OCSP cache provided"))
	}
	if len(connState.VerifiedChains) == 0 {
		return newOCSPError(errors.New("no verified certificate chains reported after TLS handshake"))
	}

	certChain := connState.VerifiedChains[0]
	if numCerts := len(certChain); numCerts == 0 {
		return newOCSPError(errors.New("verified chain contained no certificates"))
	}

	ocspCfg, err := newConfig(certChain, opts)
	if err != nil {
		return newOCSPError(err)
	}

	res, err := getParsedResponse(ctx, ocspCfg, connState)
	if err != nil {
		return err
	}
	if res == nil {
		
		
		return nil
	}

	if res.Status == ocsp.Revoked {
		return newOCSPError(errors.New("certificate is revoked"))
	}
	return nil
}



func getParsedResponse(ctx context.Context, cfg config, connState tls.ConnectionState) (*ResponseDetails, error) {
	stapledResponse, err := processStaple(cfg, connState.OCSPResponse)
	if err != nil {
		return nil, err
	}

	if stapledResponse != nil {
		
		
		return cfg.cache.Update(cfg.ocspRequest, stapledResponse), nil
	}
	if cachedResponse := cfg.cache.Get(cfg.ocspRequest); cachedResponse != nil {
		return cachedResponse, nil
	}

	
	
	if cfg.disableEndpointChecking {
		return nil, nil
	}
	externalResponse := contactResponders(ctx, cfg)
	if externalResponse == nil {
		
		return nil, nil
	}

	
	
	return cfg.cache.Update(cfg.ocspRequest, externalResponse), nil
}








func processStaple(cfg config, staple []byte) (*ResponseDetails, error) {
	mustStaple, err := isMustStapleCertificate(cfg.serverCert)
	if err != nil {
		return nil, err
	}

	
	if mustStaple && len(staple) == 0 {
		return nil, errors.New("server provided a certificate with the Must-Staple extension but did not " +
			"provide a stapled OCSP response")
	}

	if len(staple) == 0 {
		return nil, nil
	}

	parsedResponse, err := ocsp.ParseResponseForCert(staple, cfg.serverCert, cfg.issuer)
	if err != nil {
		
		
		
		return nil, fmt.Errorf("error parsing stapled response: %w", err)
	}
	if err = verifyResponse(cfg, parsedResponse); err != nil {
		return nil, fmt.Errorf("error validating stapled response: %w", err)
	}

	return extractResponseDetails(parsedResponse), nil
}


func isMustStapleCertificate(cert *x509.Certificate) (bool, error) {
	var featureExtension pkix.Extension
	var foundExtension bool
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(tlsFeatureExtensionOID) {
			featureExtension = ext
			foundExtension = true
			break
		}
	}
	if !foundExtension {
		return false, nil
	}

	
	
	
	
	
	var featureValues []*big.Int
	if _, err := asn1.Unmarshal(featureExtension.Value, &featureValues); err != nil {
		return false, fmt.Errorf("error unmarshalling TLS feature extension values: %w", err)
	}

	for _, value := range featureValues {
		if value.Cmp(mustStapleFeatureValue) == 0 {
			return true, nil
		}
	}
	return false, nil
}





func contactResponders(ctx context.Context, cfg config) *ResponseDetails {
	if len(cfg.serverCert.OCSPServer) == 0 {
		return nil
	}

	
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	ocspResponses := make(chan *ocsp.Response, len(cfg.serverCert.OCSPServer))
	defer close(ocspResponses)

	for _, endpoint := range cfg.serverCert.OCSPServer {
		
		
		endpoint := endpoint

		
		
		
		
		group.Go(func() error {
			
			
			
			request, err := http.NewRequest("POST", endpoint, bytes.NewReader(cfg.ocspRequestBytes))
			if err != nil {
				return nil
			}
			request = request.WithContext(ctx)

			httpResponse, err := cfg.httpClient.Do(request)
			if err != nil {
				return nil
			}
			defer func() {
				_ = httpResponse.Body.Close()
			}()

			if httpResponse.StatusCode != 200 {
				return nil
			}

			httpBytes, err := ioutil.ReadAll(httpResponse.Body)
			if err != nil {
				return nil
			}

			ocspResponse, err := ocsp.ParseResponseForCert(httpBytes, cfg.serverCert, cfg.issuer)
			if err != nil || verifyResponse(cfg, ocspResponse) != nil || ocspResponse.Status == ocsp.Unknown {
				
				
				return nil
			}

			
			
			ocspResponses <- ocspResponse
			return errors.New("done")
		})
	}

	_ = group.Wait()
	select {
	case res := <-ocspResponses:
		return extractResponseDetails(res)
	default:
		
		
		return nil
	}
}


func verifyResponse(cfg config, res *ocsp.Response) error {
	if err := verifyExtendedKeyUsage(cfg, res); err != nil {
		return err
	}

	currTime := time.Now().UTC()
	if res.ThisUpdate.After(currTime) {
		return fmt.Errorf("reported thisUpdate time %s is after current time %s", res.ThisUpdate, currTime)
	}
	if !res.NextUpdate.IsZero() && res.NextUpdate.Before(currTime) {
		return fmt.Errorf("reported nextUpdate time %s is before current time %s", res.NextUpdate, currTime)
	}
	return nil
}

func verifyExtendedKeyUsage(cfg config, res *ocsp.Response) error {
	if res.Certificate == nil {
		return nil
	}

	namesMatch := res.RawResponderName != nil && bytes.Equal(res.RawResponderName, cfg.issuer.RawSubject)
	keyHashesMatch := res.ResponderKeyHash != nil && bytes.Equal(res.ResponderKeyHash, cfg.ocspRequest.IssuerKeyHash)
	if namesMatch || keyHashesMatch {
		
		return nil
	}

	
	for _, extKeyUsage := range res.Certificate.ExtKeyUsage {
		if extKeyUsage == x509.ExtKeyUsageOCSPSigning {
			return nil
		}
	}

	return errors.New("delegate responder certificate is missing the OCSP signing extended key usage")
}
