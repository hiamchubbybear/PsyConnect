









package v4

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws"
	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)

const (
	authorizationHeader     = "Authorization"
	authHeaderSignatureElem = "Signature="

	authHeaderPrefix = "AWS4-HMAC-SHA256"
	timeFormat       = "20060102T150405Z"
	shortTimeFormat  = "20060102"
	awsV4Request     = "aws4_request"

	
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
)

var ignoredHeaders = rules{
	excludeList{
		mapRule{
			authorizationHeader: struct{}{},
			"User-Agent":        struct{}{},
			"X-Amzn-Trace-Id":   struct{}{},
		},
	},
}



type Signer struct {
	
	
	Credentials *credentials.Credentials
}


func NewSigner(credentials *credentials.Credentials) *Signer {
	v4 := &Signer{
		Credentials: credentials,
	}

	return v4
}

type signingCtx struct {
	ServiceName      string
	Region           string
	Request          *http.Request
	Body             io.ReadSeeker
	Query            url.Values
	Time             time.Time
	SignedHeaderVals http.Header

	credValues credentials.Value

	bodyDigest       string
	signedHeaders    string
	canonicalHeaders string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
}


























func (v4 Signer) Sign(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) {
	return v4.signWithBody(r, body, service, region, signTime)
}

func (v4 Signer) signWithBody(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) {
	ctx := &signingCtx{
		Request:     r,
		Body:        body,
		Query:       r.URL.Query(),
		Time:        signTime,
		ServiceName: service,
		Region:      region,
	}

	for key := range ctx.Query {
		sort.Strings(ctx.Query[key])
	}

	if ctx.isRequestSigned() {
		ctx.Time = time.Now()
	}

	var err error
	ctx.credValues, err = v4.Credentials.GetWithContext(r.Context())
	if err != nil {
		return http.Header{}, err
	}

	ctx.sanitizeHostForHeader()
	ctx.assignAmzQueryValues()
	if err := ctx.build(); err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	if body != nil {
		var ok bool
		if reader, ok = body.(io.ReadCloser); !ok {
			reader = ioutil.NopCloser(body)
		}
	}
	r.Body = reader

	return ctx.SignedHeaderVals, nil
}


func (ctx *signingCtx) sanitizeHostForHeader() {
	r := ctx.Request
	host := getHost(r)
	port := portOnly(host)
	if port != "" && isDefaultPort(r.URL.Scheme, port) {
		r.Host = stripPort(host)
	}
}

func (ctx *signingCtx) assignAmzQueryValues() {
	if ctx.credValues.SessionToken != "" {
		ctx.Request.Header.Set("X-Amz-Security-Token", ctx.credValues.SessionToken)
	}
}

func (ctx *signingCtx) build() error {
	ctx.buildTime()             
	ctx.buildCredentialString() 

	if err := ctx.buildBodyDigest(); err != nil {
		return err
	}

	unsignedHeaders := ctx.Request.Header

	ctx.buildCanonicalHeaders(ignoredHeaders, unsignedHeaders)
	ctx.buildCanonicalString() 
	ctx.buildStringToSign()    
	ctx.buildSignature()       

	parts := []string{
		authHeaderPrefix + " Credential=" + ctx.credValues.AccessKeyID + "/" + ctx.credentialString,
		"SignedHeaders=" + ctx.signedHeaders,
		authHeaderSignatureElem + ctx.signature,
	}
	ctx.Request.Header.Set(authorizationHeader, strings.Join(parts, ", "))

	return nil
}

func (ctx *signingCtx) buildTime() {
	ctx.Request.Header.Set("X-Amz-Date", formatTime(ctx.Time))
}

func (ctx *signingCtx) buildCredentialString() {
	ctx.credentialString = buildSigningScope(ctx.Region, ctx.ServiceName, ctx.Time)
}

func (ctx *signingCtx) buildCanonicalHeaders(r rule, header http.Header) {
	headers := make([]string, 0, len(header)+1)
	headers = append(headers, "host")
	for k, v := range header {
		if !r.IsValid(k) {
			continue 
		}
		if ctx.SignedHeaderVals == nil {
			ctx.SignedHeaderVals = make(http.Header)
		}

		lowerCaseKey := strings.ToLower(k)
		if _, ok := ctx.SignedHeaderVals[lowerCaseKey]; ok {
			
			ctx.SignedHeaderVals[lowerCaseKey] = append(ctx.SignedHeaderVals[lowerCaseKey], v...)
			continue
		}

		headers = append(headers, lowerCaseKey)
		ctx.SignedHeaderVals[lowerCaseKey] = v
	}
	sort.Strings(headers)

	ctx.signedHeaders = strings.Join(headers, ";")

	headerItems := make([]string, len(headers))
	for i, k := range headers {
		if k == "host" {
			if ctx.Request.Host != "" {
				headerItems[i] = "host:" + ctx.Request.Host
			} else {
				headerItems[i] = "host:" + ctx.Request.URL.Host
			}
		} else {
			headerValues := make([]string, len(ctx.SignedHeaderVals[k]))
			for i, v := range ctx.SignedHeaderVals[k] {
				headerValues[i] = strings.TrimSpace(v)
			}
			headerItems[i] = k + ":" +
				strings.Join(headerValues, ",")
		}
	}
	stripExcessSpaces(headerItems)
	ctx.canonicalHeaders = strings.Join(headerItems, "\n")
}

func (ctx *signingCtx) buildCanonicalString() {
	ctx.Request.URL.RawQuery = strings.Replace(ctx.Query.Encode(), "+", "%20", -1)

	uri := getURIPath(ctx.Request.URL)

	uri = EscapePath(uri, false)

	ctx.canonicalString = strings.Join([]string{
		ctx.Request.Method,
		uri,
		ctx.Request.URL.RawQuery,
		ctx.canonicalHeaders + "\n",
		ctx.signedHeaders,
		ctx.bodyDigest,
	}, "\n")
}

func (ctx *signingCtx) buildStringToSign() {
	ctx.stringToSign = strings.Join([]string{
		authHeaderPrefix,
		formatTime(ctx.Time),
		ctx.credentialString,
		hex.EncodeToString(hashSHA256([]byte(ctx.canonicalString))),
	}, "\n")
}

func (ctx *signingCtx) buildSignature() {
	creds := deriveSigningKey(ctx.Region, ctx.ServiceName, ctx.credValues.SecretAccessKey, ctx.Time)
	signature := hmacSHA256(creds, []byte(ctx.stringToSign))
	ctx.signature = hex.EncodeToString(signature)
}

func (ctx *signingCtx) buildBodyDigest() error {
	hash := ctx.Request.Header.Get("X-Amz-Content-Sha256")
	if hash == "" {
		if ctx.Body == nil {
			hash = emptyStringSHA256
		} else {
			if !aws.IsReaderSeekable(ctx.Body) {
				return fmt.Errorf("cannot use unseekable request body %T, for signed request with body", ctx.Body)
			}
			hashBytes, err := makeSha256Reader(ctx.Body)
			if err != nil {
				return err
			}
			hash = hex.EncodeToString(hashBytes)
		}
	}
	ctx.bodyDigest = hash

	return nil
}


func (ctx *signingCtx) isRequestSigned() bool {
	return ctx.Request.Header.Get("Authorization") != ""
}

func hmacSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func makeSha256Reader(reader io.ReadSeeker) (hashBytes []byte, err error) {
	hash := sha256.New()
	start, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	defer func() {
		
		_, err = reader.Seek(start, io.SeekStart)
	}()

	
	
	size, err := aws.SeekerLen(reader)
	if err != nil {
		_, _ = io.Copy(hash, reader)
	} else {
		_, _ = io.CopyN(hash, reader, size)
	}

	return hash.Sum(nil), nil
}

const doubleSpace = "  "



func stripExcessSpaces(vals []string) {
	var j, k, l, m, spaces int
	for i, str := range vals {
		

		
		for j = len(str) - 1; j >= 0 && str[j] == ' '; j-- {
		}

		
		for k = 0; k < j && str[k] == ' '; k++ {
		}

		

		str = str[k : j+1]

		
		j = strings.Index(str, doubleSpace)
		if j < 0 {
			vals[i] = str
			continue
		}

		buf := []byte(str)
		for k, m, l = j, j, len(buf); k < l; k++ {
			if buf[k] == ' ' {
				if spaces == 0 {
					
					buf[m] = buf[k]
					m++
				}
				spaces++
			} else {
				
				spaces = 0
				buf[m] = buf[k]
				m++
			}
		}

		vals[i] = string(buf[:m])
	}
}

func buildSigningScope(region, service string, dt time.Time) string {
	return strings.Join([]string{
		formatShortTime(dt),
		region,
		service,
		awsV4Request,
	}, "/")
}

func deriveSigningKey(region, service, secretKey string, dt time.Time) []byte {
	keyDate := hmacSHA256([]byte("AWS4"+secretKey), []byte(formatShortTime(dt)))
	keyRegion := hmacSHA256(keyDate, []byte(region))
	keyService := hmacSHA256(keyRegion, []byte(service))
	signingKey := hmacSHA256(keyService, []byte(awsV4Request))
	return signingKey
}

func formatShortTime(dt time.Time) string {
	return dt.UTC().Format(shortTimeFormat)
}

func formatTime(dt time.Time) string {
	return dt.UTC().Format(timeFormat)
}
