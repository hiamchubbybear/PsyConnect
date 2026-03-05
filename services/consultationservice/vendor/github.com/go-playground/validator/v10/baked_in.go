package validator

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/mail"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/sha3"
	"golang.org/x/text/language"

	"github.com/gabriel-vasile/mimetype"
	urn "github.com/leodido/go-urn"
)



type Func func(fl FieldLevel) bool



type FuncCtx func(ctx context.Context, fl FieldLevel) bool


func wrapFunc(fn Func) FuncCtx {
	if fn == nil {
		return nil 
	}
	return func(ctx context.Context, fl FieldLevel) bool {
		return fn(fl)
	}
}

var (
	restrictedTags = map[string]struct{}{
		diveTag:           {},
		keysTag:           {},
		endKeysTag:        {},
		structOnlyTag:     {},
		omitzero:          {},
		omitempty:         {},
		omitnil:           {},
		skipValidationTag: {},
		utf8HexComma:      {},
		utf8Pipe:          {},
		noStructLevelTag:  {},
		requiredTag:       {},
		isdefault:         {},
	}

	
	
	
	bakedInAliases = map[string]string{
		"iscolor":         "hexcolor|rgb|rgba|hsl|hsla",
		"country_code":    "iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric",
		"eu_country_code": "iso3166_1_alpha2_eu|iso3166_1_alpha3_eu|iso3166_1_alpha_numeric_eu",
	}

	
	
	
	bakedInValidators = map[string]Func{
		"required":                      hasValue,
		"required_if":                   requiredIf,
		"required_unless":               requiredUnless,
		"skip_unless":                   skipUnless,
		"required_with":                 requiredWith,
		"required_with_all":             requiredWithAll,
		"required_without":              requiredWithout,
		"required_without_all":          requiredWithoutAll,
		"excluded_if":                   excludedIf,
		"excluded_unless":               excludedUnless,
		"excluded_with":                 excludedWith,
		"excluded_with_all":             excludedWithAll,
		"excluded_without":              excludedWithout,
		"excluded_without_all":          excludedWithoutAll,
		"isdefault":                     isDefault,
		"len":                           hasLengthOf,
		"min":                           hasMinOf,
		"max":                           hasMaxOf,
		"eq":                            isEq,
		"eq_ignore_case":                isEqIgnoreCase,
		"ne":                            isNe,
		"ne_ignore_case":                isNeIgnoreCase,
		"lt":                            isLt,
		"lte":                           isLte,
		"gt":                            isGt,
		"gte":                           isGte,
		"eqfield":                       isEqField,
		"eqcsfield":                     isEqCrossStructField,
		"necsfield":                     isNeCrossStructField,
		"gtcsfield":                     isGtCrossStructField,
		"gtecsfield":                    isGteCrossStructField,
		"ltcsfield":                     isLtCrossStructField,
		"ltecsfield":                    isLteCrossStructField,
		"nefield":                       isNeField,
		"gtefield":                      isGteField,
		"gtfield":                       isGtField,
		"ltefield":                      isLteField,
		"ltfield":                       isLtField,
		"fieldcontains":                 fieldContains,
		"fieldexcludes":                 fieldExcludes,
		"alpha":                         isAlpha,
		"alphanum":                      isAlphanum,
		"alphaunicode":                  isAlphaUnicode,
		"alphanumunicode":               isAlphanumUnicode,
		"boolean":                       isBoolean,
		"numeric":                       isNumeric,
		"number":                        isNumber,
		"hexadecimal":                   isHexadecimal,
		"hexcolor":                      isHEXColor,
		"rgb":                           isRGB,
		"rgba":                          isRGBA,
		"hsl":                           isHSL,
		"hsla":                          isHSLA,
		"e164":                          isE164,
		"email":                         isEmail,
		"url":                           isURL,
		"http_url":                      isHttpURL,
		"uri":                           isURI,
		"urn_rfc2141":                   isUrnRFC2141, 
		"file":                          isFile,
		"filepath":                      isFilePath,
		"base32":                        isBase32,
		"base64":                        isBase64,
		"base64url":                     isBase64URL,
		"base64rawurl":                  isBase64RawURL,
		"contains":                      contains,
		"containsany":                   containsAny,
		"containsrune":                  containsRune,
		"excludes":                      excludes,
		"excludesall":                   excludesAll,
		"excludesrune":                  excludesRune,
		"startswith":                    startsWith,
		"endswith":                      endsWith,
		"startsnotwith":                 startsNotWith,
		"endsnotwith":                   endsNotWith,
		"image":                         isImage,
		"isbn":                          isISBN,
		"isbn10":                        isISBN10,
		"isbn13":                        isISBN13,
		"issn":                          isISSN,
		"eth_addr":                      isEthereumAddress,
		"eth_addr_checksum":             isEthereumAddressChecksum,
		"btc_addr":                      isBitcoinAddress,
		"btc_addr_bech32":               isBitcoinBech32Address,
		"uuid":                          isUUID,
		"uuid3":                         isUUID3,
		"uuid4":                         isUUID4,
		"uuid5":                         isUUID5,
		"uuid_rfc4122":                  isUUIDRFC4122,
		"uuid3_rfc4122":                 isUUID3RFC4122,
		"uuid4_rfc4122":                 isUUID4RFC4122,
		"uuid5_rfc4122":                 isUUID5RFC4122,
		"ulid":                          isULID,
		"md4":                           isMD4,
		"md5":                           isMD5,
		"sha256":                        isSHA256,
		"sha384":                        isSHA384,
		"sha512":                        isSHA512,
		"ripemd128":                     isRIPEMD128,
		"ripemd160":                     isRIPEMD160,
		"tiger128":                      isTIGER128,
		"tiger160":                      isTIGER160,
		"tiger192":                      isTIGER192,
		"ascii":                         isASCII,
		"printascii":                    isPrintableASCII,
		"multibyte":                     hasMultiByteCharacter,
		"datauri":                       isDataURI,
		"latitude":                      isLatitude,
		"longitude":                     isLongitude,
		"ssn":                           isSSN,
		"ipv4":                          isIPv4,
		"ipv6":                          isIPv6,
		"ip":                            isIP,
		"cidrv4":                        isCIDRv4,
		"cidrv6":                        isCIDRv6,
		"cidr":                          isCIDR,
		"tcp4_addr":                     isTCP4AddrResolvable,
		"tcp6_addr":                     isTCP6AddrResolvable,
		"tcp_addr":                      isTCPAddrResolvable,
		"udp4_addr":                     isUDP4AddrResolvable,
		"udp6_addr":                     isUDP6AddrResolvable,
		"udp_addr":                      isUDPAddrResolvable,
		"ip4_addr":                      isIP4AddrResolvable,
		"ip6_addr":                      isIP6AddrResolvable,
		"ip_addr":                       isIPAddrResolvable,
		"unix_addr":                     isUnixAddrResolvable,
		"mac":                           isMAC,
		"hostname":                      isHostnameRFC952,  
		"hostname_rfc1123":              isHostnameRFC1123, 
		"fqdn":                          isFQDN,
		"unique":                        isUnique,
		"oneof":                         isOneOf,
		"oneofci":                       isOneOfCI,
		"html":                          isHTML,
		"html_encoded":                  isHTMLEncoded,
		"url_encoded":                   isURLEncoded,
		"dir":                           isDir,
		"dirpath":                       isDirPath,
		"json":                          isJSON,
		"jwt":                           isJWT,
		"hostname_port":                 isHostnamePort,
		"port":                          isPort,
		"lowercase":                     isLowercase,
		"uppercase":                     isUppercase,
		"datetime":                      isDatetime,
		"timezone":                      isTimeZone,
		"iso3166_1_alpha2":              isIso3166Alpha2,
		"iso3166_1_alpha2_eu":           isIso3166Alpha2EU,
		"iso3166_1_alpha3":              isIso3166Alpha3,
		"iso3166_1_alpha3_eu":           isIso3166Alpha3EU,
		"iso3166_1_alpha_numeric":       isIso3166AlphaNumeric,
		"iso3166_1_alpha_numeric_eu":    isIso3166AlphaNumericEU,
		"iso3166_2":                     isIso31662,
		"iso4217":                       isIso4217,
		"iso4217_numeric":               isIso4217Numeric,
		"bcp47_language_tag":            isBCP47LanguageTag,
		"postcode_iso3166_alpha2":       isPostcodeByIso3166Alpha2,
		"postcode_iso3166_alpha2_field": isPostcodeByIso3166Alpha2Field,
		"bic":                           isIsoBicFormat,
		"semver":                        isSemverFormat,
		"dns_rfc1035_label":             isDnsRFC1035LabelFormat,
		"credit_card":                   isCreditCard,
		"cve":                           isCveFormat,
		"luhn_checksum":                 hasLuhnChecksum,
		"mongodb":                       isMongoDBObjectId,
		"mongodb_connection_string":     isMongoDBConnectionString,
		"cron":                          isCron,
		"spicedb":                       isSpiceDB,
		"ein":                           isEIN,
	}
)

var (
	oneofValsCache       = map[string][]string{}
	oneofValsCacheRWLock = sync.RWMutex{}
)

func parseOneOfParam2(s string) []string {
	oneofValsCacheRWLock.RLock()
	vals, ok := oneofValsCache[s]
	oneofValsCacheRWLock.RUnlock()
	if !ok {
		oneofValsCacheRWLock.Lock()
		vals = splitParamsRegex().FindAllString(s, -1)
		for i := 0; i < len(vals); i++ {
			vals[i] = strings.ReplaceAll(vals[i], "'", "")
		}
		oneofValsCache[s] = vals
		oneofValsCacheRWLock.Unlock()
	}
	return vals
}

func isURLEncoded(fl FieldLevel) bool {
	return uRLEncodedRegex().MatchString(fl.Field().String())
}

func isHTMLEncoded(fl FieldLevel) bool {
	return hTMLEncodedRegex().MatchString(fl.Field().String())
}

func isHTML(fl FieldLevel) bool {
	return hTMLRegex().MatchString(fl.Field().String())
}

func isOneOf(fl FieldLevel) bool {
	vals := parseOneOfParam2(fl.Param())

	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	for i := 0; i < len(vals); i++ {
		if vals[i] == v {
			return true
		}
	}
	return false
}


func isOneOfCI(fl FieldLevel) bool {
	vals := parseOneOfParam2(fl.Param())
	field := fl.Field()

	if field.Kind() != reflect.String {
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	v := field.String()
	for _, val := range vals {
		if strings.EqualFold(val, v) {
			return true
		}
	}
	return false
}


func isUnique(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()
	v := reflect.ValueOf(struct{}{})

	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		elem := field.Type().Elem()
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		if param == "" {
			m := reflect.MakeMap(reflect.MapOf(elem, v.Type()))

			for i := 0; i < field.Len(); i++ {
				m.SetMapIndex(reflect.Indirect(field.Index(i)), v)
			}
			return field.Len() == m.Len()
		}

		sf, ok := elem.FieldByName(param)
		if !ok {
			panic(fmt.Sprintf("Bad field name %s", param))
		}

		sfTyp := sf.Type
		if sfTyp.Kind() == reflect.Ptr {
			sfTyp = sfTyp.Elem()
		}

		m := reflect.MakeMap(reflect.MapOf(sfTyp, v.Type()))
		var fieldlen int
		for i := 0; i < field.Len(); i++ {
			key := reflect.Indirect(reflect.Indirect(field.Index(i)).FieldByName(param))
			if key.IsValid() {
				fieldlen++
				m.SetMapIndex(key, v)
			}
		}
		return fieldlen == m.Len()
	case reflect.Map:
		var m reflect.Value
		if field.Type().Elem().Kind() == reflect.Ptr {
			m = reflect.MakeMap(reflect.MapOf(field.Type().Elem().Elem(), v.Type()))
		} else {
			m = reflect.MakeMap(reflect.MapOf(field.Type().Elem(), v.Type()))
		}

		for _, k := range field.MapKeys() {
			m.SetMapIndex(reflect.Indirect(field.MapIndex(k)), v)
		}

		return field.Len() == m.Len()
	default:
		if parent := fl.Parent(); parent.Kind() == reflect.Struct {
			uniqueField := parent.FieldByName(param)
			if uniqueField == reflect.ValueOf(nil) {
				panic(fmt.Sprintf("Bad field name provided %s", param))
			}

			if uniqueField.Kind() != field.Kind() {
				panic(fmt.Sprintf("Bad field type %T:%T", field.Interface(), uniqueField.Interface()))
			}

			return field.Interface() != uniqueField.Interface()
		}

		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
}


func isMAC(fl FieldLevel) bool {
	_, err := net.ParseMAC(fl.Field().String())

	return err == nil
}


func isCIDRv4(fl FieldLevel) bool {
	ip, net, err := net.ParseCIDR(fl.Field().String())

	return err == nil && ip.To4() != nil && net.IP.Equal(ip)
}


func isCIDRv6(fl FieldLevel) bool {
	ip, _, err := net.ParseCIDR(fl.Field().String())

	return err == nil && ip.To4() == nil
}


func isCIDR(fl FieldLevel) bool {
	_, _, err := net.ParseCIDR(fl.Field().String())

	return err == nil
}


func isIPv4(fl FieldLevel) bool {
	ip := net.ParseIP(fl.Field().String())

	return ip != nil && ip.To4() != nil
}


func isIPv6(fl FieldLevel) bool {
	ip := net.ParseIP(fl.Field().String())

	return ip != nil && ip.To4() == nil
}


func isIP(fl FieldLevel) bool {
	ip := net.ParseIP(fl.Field().String())

	return ip != nil
}


func isSSN(fl FieldLevel) bool {
	field := fl.Field()

	if field.Len() != 11 {
		return false
	}

	return sSNRegex().MatchString(field.String())
}


func isLongitude(fl FieldLevel) bool {
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	return longitudeRegex().MatchString(v)
}


func isLatitude(fl FieldLevel) bool {
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	return latitudeRegex().MatchString(v)
}


func isDataURI(fl FieldLevel) bool {
	uri := strings.SplitN(fl.Field().String(), ",", 2)

	if len(uri) != 2 {
		return false
	}

	if !dataURIRegex().MatchString(uri[0]) {
		return false
	}

	return base64Regex().MatchString(uri[1])
}


func hasMultiByteCharacter(fl FieldLevel) bool {
	field := fl.Field()

	if field.Len() == 0 {
		return true
	}

	return multibyteRegex().MatchString(field.String())
}


func isPrintableASCII(fl FieldLevel) bool {
	return printableASCIIRegex().MatchString(fl.Field().String())
}


func isASCII(fl FieldLevel) bool {
	return aSCIIRegex().MatchString(fl.Field().String())
}


func isUUID5(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUID5Regex, fl)
}


func isUUID4(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUID4Regex, fl)
}


func isUUID3(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUID3Regex, fl)
}


func isUUID(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUIDRegex, fl)
}


func isUUID5RFC4122(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUID5RFC4122Regex, fl)
}


func isUUID4RFC4122(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUID4RFC4122Regex, fl)
}


func isUUID3RFC4122(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUID3RFC4122Regex, fl)
}


func isUUIDRFC4122(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uUIDRFC4122Regex, fl)
}


func isULID(fl FieldLevel) bool {
	return fieldMatchesRegexByStringerValOrString(uLIDRegex, fl)
}


func isMD4(fl FieldLevel) bool {
	return md4Regex().MatchString(fl.Field().String())
}


func isMD5(fl FieldLevel) bool {
	return md5Regex().MatchString(fl.Field().String())
}


func isSHA256(fl FieldLevel) bool {
	return sha256Regex().MatchString(fl.Field().String())
}


func isSHA384(fl FieldLevel) bool {
	return sha384Regex().MatchString(fl.Field().String())
}


func isSHA512(fl FieldLevel) bool {
	return sha512Regex().MatchString(fl.Field().String())
}


func isRIPEMD128(fl FieldLevel) bool {
	return ripemd128Regex().MatchString(fl.Field().String())
}


func isRIPEMD160(fl FieldLevel) bool {
	return ripemd160Regex().MatchString(fl.Field().String())
}


func isTIGER128(fl FieldLevel) bool {
	return tiger128Regex().MatchString(fl.Field().String())
}


func isTIGER160(fl FieldLevel) bool {
	return tiger160Regex().MatchString(fl.Field().String())
}


func isTIGER192(fl FieldLevel) bool {
	return tiger192Regex().MatchString(fl.Field().String())
}


func isISBN(fl FieldLevel) bool {
	return isISBN10(fl) || isISBN13(fl)
}


func isISBN13(fl FieldLevel) bool {
	s := strings.Replace(strings.Replace(fl.Field().String(), "-", "", 4), " ", "", 4)

	if !iSBN13Regex().MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	factor := []int32{1, 3}

	for i = 0; i < 12; i++ {
		checksum += factor[i%2] * int32(s[i]-'0')
	}

	return (int32(s[12]-'0'))-((10-(checksum%10))%10) == 0
}


func isISBN10(fl FieldLevel) bool {
	s := strings.Replace(strings.Replace(fl.Field().String(), "-", "", 3), " ", "", 3)

	if !iSBN10Regex().MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	for i = 0; i < 9; i++ {
		checksum += (i + 1) * int32(s[i]-'0')
	}

	if s[9] == 'X' {
		checksum += 10 * 10
	} else {
		checksum += 10 * int32(s[9]-'0')
	}

	return checksum%11 == 0
}


func isISSN(fl FieldLevel) bool {
	s := fl.Field().String()

	if !iSSNRegex().MatchString(s) {
		return false
	}
	s = strings.ReplaceAll(s, "-", "")

	pos := 8
	checksum := 0

	for i := 0; i < 7; i++ {
		checksum += pos * int(s[i]-'0')
		pos--
	}

	if s[7] == 'X' {
		checksum += 10
	} else {
		checksum += int(s[7] - '0')
	}

	return checksum%11 == 0
}


func isEthereumAddress(fl FieldLevel) bool {
	address := fl.Field().String()

	return ethAddressRegex().MatchString(address)
}


func isEthereumAddressChecksum(fl FieldLevel) bool {
	address := fl.Field().String()

	if !ethAddressRegex().MatchString(address) {
		return false
	}
	
	address = address[2:] 
	h := sha3.NewLegacyKeccak256()
	
	_, _ = h.Write([]byte(strings.ToLower(address)))
	hash := hex.EncodeToString(h.Sum(nil))

	for i := 0; i < len(address); i++ {
		if address[i] <= '9' { 
			continue
		}
		if hash[i] > '7' && address[i] >= 'a' || hash[i] <= '7' && address[i] <= 'F' {
			return false
		}
	}

	return true
}


func isBitcoinAddress(fl FieldLevel) bool {
	address := fl.Field().String()

	if !btcAddressRegex().MatchString(address) {
		return false
	}

	alphabet := []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	decode := [25]byte{}

	for _, n := range []byte(address) {
		d := bytes.IndexByte(alphabet, n)

		for i := 24; i >= 0; i-- {
			d += 58 * int(decode[i])
			decode[i] = byte(d % 256)
			d /= 256
		}
	}

	h := sha256.New()
	_, _ = h.Write(decode[:21])
	d := h.Sum([]byte{})
	h = sha256.New()
	_, _ = h.Write(d)

	validchecksum := [4]byte{}
	computedchecksum := [4]byte{}

	copy(computedchecksum[:], h.Sum(d[:0]))
	copy(validchecksum[:], decode[21:])

	return validchecksum == computedchecksum
}


func isBitcoinBech32Address(fl FieldLevel) bool {
	address := fl.Field().String()

	if !btcLowerAddressRegexBech32().MatchString(address) && !btcUpperAddressRegexBech32().MatchString(address) {
		return false
	}

	am := len(address) % 8

	if am == 0 || am == 3 || am == 5 {
		return false
	}

	address = strings.ToLower(address)

	alphabet := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

	hr := []int{3, 3, 0, 2, 3} 
	addr := address[3:]
	dp := make([]int, 0, len(addr))

	for _, c := range addr {
		dp = append(dp, strings.IndexRune(alphabet, c))
	}

	ver := dp[0]

	if ver < 0 || ver > 16 {
		return false
	}

	if ver == 0 {
		if len(address) != 42 && len(address) != 62 {
			return false
		}
	}

	values := append(hr, dp...)

	GEN := []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

	p := 1

	for _, v := range values {
		b := p >> 25
		p = (p&0x1ffffff)<<5 ^ v

		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				p ^= GEN[i]
			}
		}
	}

	if p != 1 {
		return false
	}

	b := uint(0)
	acc := 0
	mv := (1 << 5) - 1
	var sw []int

	for _, v := range dp[1 : len(dp)-6] {
		acc = (acc << 5) | v
		b += 5
		for b >= 8 {
			b -= 8
			sw = append(sw, (acc>>b)&mv)
		}
	}

	if len(sw) < 2 || len(sw) > 40 {
		return false
	}

	return true
}


func excludesRune(fl FieldLevel) bool {
	return !containsRune(fl)
}


func excludesAll(fl FieldLevel) bool {
	return !containsAny(fl)
}


func excludes(fl FieldLevel) bool {
	return !contains(fl)
}


func containsRune(fl FieldLevel) bool {
	r, _ := utf8.DecodeRuneInString(fl.Param())

	return strings.ContainsRune(fl.Field().String(), r)
}


func containsAny(fl FieldLevel) bool {
	return strings.ContainsAny(fl.Field().String(), fl.Param())
}


func contains(fl FieldLevel) bool {
	return strings.Contains(fl.Field().String(), fl.Param())
}


func startsWith(fl FieldLevel) bool {
	return strings.HasPrefix(fl.Field().String(), fl.Param())
}


func endsWith(fl FieldLevel) bool {
	return strings.HasSuffix(fl.Field().String(), fl.Param())
}


func startsNotWith(fl FieldLevel) bool {
	return !startsWith(fl)
}


func endsNotWith(fl FieldLevel) bool {
	return !endsWith(fl)
}


func fieldContains(fl FieldLevel) bool {
	field := fl.Field()

	currentField, _, ok := fl.GetStructFieldOK()

	if !ok {
		return false
	}

	return strings.Contains(field.String(), currentField.String())
}


func fieldExcludes(fl FieldLevel) bool {
	field := fl.Field()

	currentField, _, ok := fl.GetStructFieldOK()
	if !ok {
		return true
	}

	return !strings.Contains(field.String(), currentField.String())
}


func isNeField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()

	if !ok || currentKind != kind {
		return true
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() != currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() != currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) != int64(currentField.Len())

	case reflect.Bool:
		return field.Bool() != currentField.Bool()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && currentField.Type().ConvertibleTo(timeType) {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return !fieldTime.Equal(t)
		}

		
		if fieldType != currentField.Type() {
			return true
		}
	}

	
	return field.String() != currentField.String()
}


func isNe(fl FieldLevel) bool {
	return !isEq(fl)
}



func isNeIgnoreCase(fl FieldLevel) bool {
	return !isEqIgnoreCase(fl)
}


func isLteCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() <= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() <= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() <= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) <= int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && topField.Type().ConvertibleTo(timeType) {

			fieldTime := field.Convert(timeType).Interface().(time.Time)
			topTime := topField.Convert(timeType).Interface().(time.Time)

			return fieldTime.Before(topTime) || fieldTime.Equal(topTime)
		}

		
		if fieldType != topField.Type() {
			return false
		}
	}

	
	return field.String() <= topField.String()
}



func isLtCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() < topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() < topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) < int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && topField.Type().ConvertibleTo(timeType) {

			fieldTime := field.Convert(timeType).Interface().(time.Time)
			topTime := topField.Convert(timeType).Interface().(time.Time)

			return fieldTime.Before(topTime)
		}

		
		if fieldType != topField.Type() {
			return false
		}
	}

	
	return field.String() < topField.String()
}


func isGteCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() >= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() >= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() >= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) >= int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && topField.Type().ConvertibleTo(timeType) {

			fieldTime := field.Convert(timeType).Interface().(time.Time)
			topTime := topField.Convert(timeType).Interface().(time.Time)

			return fieldTime.After(topTime) || fieldTime.Equal(topTime)
		}

		
		if fieldType != topField.Type() {
			return false
		}
	}

	
	return field.String() >= topField.String()
}


func isGtCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() > topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() > topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) > int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && topField.Type().ConvertibleTo(timeType) {

			fieldTime := field.Convert(timeType).Interface().(time.Time)
			topTime := topField.Convert(timeType).Interface().(time.Time)

			return fieldTime.After(topTime)
		}

		
		if fieldType != topField.Type() {
			return false
		}
	}

	
	return field.String() > topField.String()
}


func isNeCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return true
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() != field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() != field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() != field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) != int64(field.Len())

	case reflect.Bool:
		return topField.Bool() != field.Bool()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && topField.Type().ConvertibleTo(timeType) {

			t := field.Convert(timeType).Interface().(time.Time)
			fieldTime := topField.Convert(timeType).Interface().(time.Time)

			return !fieldTime.Equal(t)
		}

		
		if fieldType != topField.Type() {
			return true
		}
	}

	
	return topField.String() != field.String()
}


func isEqCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() == field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() == field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() == field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) == int64(field.Len())

	case reflect.Bool:
		return topField.Bool() == field.Bool()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && topField.Type().ConvertibleTo(timeType) {

			t := field.Convert(timeType).Interface().(time.Time)
			fieldTime := topField.Convert(timeType).Interface().(time.Time)

			return fieldTime.Equal(t)
		}

		
		if fieldType != topField.Type() {
			return false
		}
	}

	
	return topField.String() == field.String()
}


func isEqField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() == currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == int64(currentField.Len())

	case reflect.Bool:
		return field.Bool() == currentField.Bool()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && currentField.Type().ConvertibleTo(timeType) {

			t := currentField.Convert(timeType).Interface().(time.Time)
			fieldTime := field.Convert(timeType).Interface().(time.Time)

			return fieldTime.Equal(t)
		}

		
		if fieldType != currentField.Type() {
			return false
		}
	}

	
	return field.String() == currentField.String()
}


func isEq(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		return field.String() == param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() == p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() == p

	case reflect.Bool:
		p := asBool(param)

		return field.Bool() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}




func isEqIgnoreCase(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return strings.EqualFold(field.String(), param)
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}



func isPostcodeByIso3166Alpha2(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	postcodeRegexInit.Do(initPostcodes)
	reg, found := postCodeRegexDict[param]
	if !found {
		return false
	}

	return reg.MatchString(field.String())
}



func isPostcodeByIso3166Alpha2Field(fl FieldLevel) bool {
	field := fl.Field()
	params := parseOneOfParam2(fl.Param())

	if len(params) != 1 {
		return false
	}

	currentField, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), params[0])
	if !found {
		return false
	}

	if kind != reflect.String {
		panic(fmt.Sprintf("Bad field type %T", currentField.Interface()))
	}

	postcodeRegexInit.Do(initPostcodes)
	reg, found := postCodeRegexDict[currentField.String()]
	if !found {
		return false
	}

	return reg.MatchString(field.String())
}


func isBase32(fl FieldLevel) bool {
	return base32Regex().MatchString(fl.Field().String())
}


func isBase64(fl FieldLevel) bool {
	return base64Regex().MatchString(fl.Field().String())
}


func isBase64URL(fl FieldLevel) bool {
	return base64URLRegex().MatchString(fl.Field().String())
}


func isBase64RawURL(fl FieldLevel) bool {
	return base64RawURLRegex().MatchString(fl.Field().String())
}


func isURI(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:

		s := field.String()

		
		
		if i := strings.Index(s, "#"); i > -1 {
			s = s[:i]
		}

		if len(s) == 0 {
			return false
		}

		_, err := url.ParseRequestURI(s)

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isFileURL(path string) bool {
	if !strings.HasPrefix(path, "file:/") {
		return false
	}
	_, err := url.ParseRequestURI(path)
	return err == nil
}


func isURL(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:

		s := strings.ToLower(field.String())

		if len(s) == 0 {
			return false
		}

		if isFileURL(s) {
			return true
		}

		url, err := url.Parse(s)
		if err != nil || url.Scheme == "" {
			return false
		}

		if url.Host == "" && url.Fragment == "" && url.Opaque == "" {
			return false
		}

		return true
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isHttpURL(fl FieldLevel) bool {
	if !isURL(fl) {
		return false
	}

	field := fl.Field()
	switch field.Kind() {
	case reflect.String:

		s := strings.ToLower(field.String())

		url, err := url.Parse(s)
		if err != nil || url.Host == "" {
			return false
		}

		return url.Scheme == "http" || url.Scheme == "https"
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isUrnRFC2141(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:

		str := field.String()

		_, match := urn.Parse([]byte(str))

		return match
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isFile(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		fileInfo, err := os.Stat(field.String())
		if err != nil {
			return false
		}

		return !fileInfo.IsDir()
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isImage(fl FieldLevel) bool {
	mimetypes := map[string]bool{
		"image/bmp":                true,
		"image/cis-cod":            true,
		"image/gif":                true,
		"image/ief":                true,
		"image/jpeg":               true,
		"image/jp2":                true,
		"image/jpx":                true,
		"image/jpm":                true,
		"image/pipeg":              true,
		"image/png":                true,
		"image/svg+xml":            true,
		"image/tiff":               true,
		"image/webp":               true,
		"image/x-cmu-raster":       true,
		"image/x-cmx":              true,
		"image/x-icon":             true,
		"image/x-portable-anymap":  true,
		"image/x-portable-bitmap":  true,
		"image/x-portable-graymap": true,
		"image/x-portable-pixmap":  true,
		"image/x-rgb":              true,
		"image/x-xbitmap":          true,
		"image/x-xpixmap":          true,
		"image/x-xwindowdump":      true,
	}
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		filePath := field.String()
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return false
		}

		if fileInfo.IsDir() {
			return false
		}

		file, err := os.Open(filePath)
		if err != nil {
			return false
		}
		defer func() {
			_ = file.Close()
		}()

		mime, err := mimetype.DetectReader(file)
		if err != nil {
			return false
		}

		if _, ok := mimetypes[mime.String()]; ok {
			return true
		}
	}
	return false
}


func isFilePath(fl FieldLevel) bool {
	var exists bool
	var err error

	field := fl.Field()

	
	if isDir(fl) {
		return false
	}
	
	
	if exists = isFile(fl); exists {
		return true
	}

	
	switch field.Kind() {
	case reflect.String:
		
		
		
		if strings.TrimSpace(field.String()) == "" {
			return false
		}
		
		if strings.HasSuffix(field.String(), string(os.PathSeparator)) {
			return false
		}
		if _, err = os.Stat(field.String()); err != nil {
			switch t := err.(type) {
			case *fs.PathError:
				if t.Err == syscall.EINVAL {
					
					return false
				}
				
				
				return true
			default:
				
				
				panic(err)
			}
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isE164(fl FieldLevel) bool {
	return e164Regex().MatchString(fl.Field().String())
}


func isEmail(fl FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	if err != nil {
		return false
	}
	return emailRegex().MatchString(fl.Field().String())
}


func isHSLA(fl FieldLevel) bool {
	return hslaRegex().MatchString(fl.Field().String())
}


func isHSL(fl FieldLevel) bool {
	return hslRegex().MatchString(fl.Field().String())
}


func isRGBA(fl FieldLevel) bool {
	return rgbaRegex().MatchString(fl.Field().String())
}


func isRGB(fl FieldLevel) bool {
	return rgbRegex().MatchString(fl.Field().String())
}


func isHEXColor(fl FieldLevel) bool {
	return hexColorRegex().MatchString(fl.Field().String())
}


func isHexadecimal(fl FieldLevel) bool {
	return hexadecimalRegex().MatchString(fl.Field().String())
}


func isNumber(fl FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return numberRegex().MatchString(fl.Field().String())
	}
}


func isNumeric(fl FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return numericRegex().MatchString(fl.Field().String())
	}
}


func isAlphanum(fl FieldLevel) bool {
	return alphaNumericRegex().MatchString(fl.Field().String())
}


func isAlpha(fl FieldLevel) bool {
	return alphaRegex().MatchString(fl.Field().String())
}


func isAlphanumUnicode(fl FieldLevel) bool {
	return alphaUnicodeNumericRegex().MatchString(fl.Field().String())
}


func isAlphaUnicode(fl FieldLevel) bool {
	return alphaUnicodeRegex().MatchString(fl.Field().String())
}


func isBoolean(fl FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Bool:
		return true
	default:
		_, err := strconv.ParseBool(fl.Field().String())
		return err == nil
	}
}


func isDefault(fl FieldLevel) bool {
	return !hasValue(fl)
}


func hasValue(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		if fl.(*validate).fldIsPointer && field.Interface() != nil {
			return true
		}
		return field.IsValid() && !field.IsZero()
	}
}


func hasNotZeroValue(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		if fl.(*validate).fldIsPointer && field.Interface() != nil {
			return !field.IsZero()
		}
		return field.IsValid() && !field.IsZero()
	}
}


func requireCheckFieldKind(fl FieldLevel, param string, defaultNotFoundValue bool) bool {
	field := fl.Field()
	kind := field.Kind()
	var nullable, found bool
	if len(param) > 0 {
		field, kind, nullable, found = fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
		if !found {
			return defaultNotFoundValue
		}
	}
	switch kind {
	case reflect.Invalid:
		return defaultNotFoundValue
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		if nullable && field.Interface() != nil {
			return false
		}
		return field.IsValid() && field.IsZero()
	}
}


func requireCheckFieldValue(
	fl FieldLevel, param string, value string, defaultNotFoundValue bool,
) bool {
	field, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
	if !found {
		return defaultNotFoundValue
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == asInt(value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == asUint(value)

	case reflect.Float32:
		return field.Float() == asFloat32(value)

	case reflect.Float64:
		return field.Float() == asFloat64(value)

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == asInt(value)

	case reflect.Bool:
		return field.Bool() == (value == "true")

	case reflect.Ptr:
		if field.IsNil() {
			return value == "nil"
		}
		
		return requireCheckFieldValue(fl, param, value, defaultNotFoundValue)
	}

	
	return field.String() == value
}



func requiredIf(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for required_if %s", fl.FieldName()))
	}
	for i := 0; i < len(params); i += 2 {
		if !requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return true
		}
	}
	return hasValue(fl)
}



func excludedIf(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for excluded_if %s", fl.FieldName()))
	}

	for i := 0; i < len(params); i += 2 {
		if !requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return true
		}
	}
	return !hasValue(fl)
}



func requiredUnless(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for required_unless %s", fl.FieldName()))
	}

	for i := 0; i < len(params); i += 2 {
		if requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return true
		}
	}
	return hasValue(fl)
}



func skipUnless(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for skip_unless %s", fl.FieldName()))
	}
	for i := 0; i < len(params); i += 2 {
		if !requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return true
		}
	}
	return hasValue(fl)
}



func excludedUnless(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	if len(params)%2 != 0 {
		panic(fmt.Sprintf("Bad param number for excluded_unless %s", fl.FieldName()))
	}
	for i := 0; i < len(params); i += 2 {
		if !requireCheckFieldValue(fl, params[i], params[i+1], false) {
			return !hasValue(fl)
		}
	}
	return true
}



func excludedWith(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return !hasValue(fl)
		}
	}
	return true
}



func requiredWith(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return hasValue(fl)
		}
	}
	return true
}



func excludedWithAll(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	for _, param := range params {
		if requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return !hasValue(fl)
}



func requiredWithAll(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	for _, param := range params {
		if requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return hasValue(fl)
}



func excludedWithout(fl FieldLevel) bool {
	if requireCheckFieldKind(fl, strings.TrimSpace(fl.Param()), true) {
		return !hasValue(fl)
	}
	return true
}



func requiredWithout(fl FieldLevel) bool {
	if requireCheckFieldKind(fl, strings.TrimSpace(fl.Param()), true) {
		return hasValue(fl)
	}
	return true
}



func excludedWithoutAll(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return !hasValue(fl)
}



func requiredWithoutAll(fl FieldLevel) bool {
	params := parseOneOfParam2(fl.Param())
	for _, param := range params {
		if !requireCheckFieldKind(fl, param, true) {
			return true
		}
	}
	return hasValue(fl)
}


func isGteField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() >= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() >= currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() >= currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && currentField.Type().ConvertibleTo(timeType) {

			t := currentField.Convert(timeType).Interface().(time.Time)
			fieldTime := field.Convert(timeType).Interface().(time.Time)

			return fieldTime.After(t) || fieldTime.Equal(t)
		}

		
		if fieldType != currentField.Type() {
			return false
		}
	}

	
	return len(field.String()) >= len(currentField.String())
}


func isGtField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() > currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() > currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() > currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && currentField.Type().ConvertibleTo(timeType) {

			t := currentField.Convert(timeType).Interface().(time.Time)
			fieldTime := field.Convert(timeType).Interface().(time.Time)

			return fieldTime.After(t)
		}

		
		if fieldType != currentField.Type() {
			return false
		}
	}

	
	return len(field.String()) > len(currentField.String())
}


func isGte(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() >= p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() >= p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() >= p

	case reflect.Struct:

		if field.Type().ConvertibleTo(timeType) {

			now := time.Now().UTC()
			t := field.Convert(timeType).Interface().(time.Time)

			return t.After(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isGt(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() > p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() > p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() > p

	case reflect.Struct:

		if field.Type().ConvertibleTo(timeType) {
			return field.Convert(timeType).Interface().(time.Time).After(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func hasLengthOf(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() == p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func hasMinOf(fl FieldLevel) bool {
	return isGte(fl)
}


func isLteField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() <= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() <= currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() <= currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && currentField.Type().ConvertibleTo(timeType) {

			t := currentField.Convert(timeType).Interface().(time.Time)
			fieldTime := field.Convert(timeType).Interface().(time.Time)

			return fieldTime.Before(t) || fieldTime.Equal(t)
		}

		
		if fieldType != currentField.Type() {
			return false
		}
	}

	
	return len(field.String()) <= len(currentField.String())
}


func isLtField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() < currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() < currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() < currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		if fieldType.ConvertibleTo(timeType) && currentField.Type().ConvertibleTo(timeType) {

			t := currentField.Convert(timeType).Interface().(time.Time)
			fieldTime := field.Convert(timeType).Interface().(time.Time)

			return fieldTime.Before(t)
		}

		
		if fieldType != currentField.Type() {
			return false
		}
	}

	
	return len(field.String()) < len(currentField.String())
}


func isLte(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) <= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() <= p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() <= p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() <= p

	case reflect.Struct:

		if field.Type().ConvertibleTo(timeType) {

			now := time.Now().UTC()
			t := field.Convert(timeType).Interface().(time.Time)

			return t.Before(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isLt(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) < p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() < p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() < p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() < p

	case reflect.Struct:

		if field.Type().ConvertibleTo(timeType) {
			return field.Convert(timeType).Interface().(time.Time).Before(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func hasMaxOf(fl FieldLevel) bool {
	return isLte(fl)
}


func isTCP4AddrResolvable(fl FieldLevel) bool {
	if !isIP4Addr(fl) {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp4", fl.Field().String())
	return err == nil
}


func isTCP6AddrResolvable(fl FieldLevel) bool {
	if !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp6", fl.Field().String())

	return err == nil
}


func isTCPAddrResolvable(fl FieldLevel) bool {
	if !isIP4Addr(fl) && !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp", fl.Field().String())

	return err == nil
}


func isUDP4AddrResolvable(fl FieldLevel) bool {
	if !isIP4Addr(fl) {
		return false
	}

	_, err := net.ResolveUDPAddr("udp4", fl.Field().String())

	return err == nil
}


func isUDP6AddrResolvable(fl FieldLevel) bool {
	if !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveUDPAddr("udp6", fl.Field().String())

	return err == nil
}


func isUDPAddrResolvable(fl FieldLevel) bool {
	if !isIP4Addr(fl) && !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveUDPAddr("udp", fl.Field().String())

	return err == nil
}


func isIP4AddrResolvable(fl FieldLevel) bool {
	if !isIPv4(fl) {
		return false
	}

	_, err := net.ResolveIPAddr("ip4", fl.Field().String())

	return err == nil
}


func isIP6AddrResolvable(fl FieldLevel) bool {
	if !isIPv6(fl) {
		return false
	}

	_, err := net.ResolveIPAddr("ip6", fl.Field().String())

	return err == nil
}


func isIPAddrResolvable(fl FieldLevel) bool {
	if !isIP(fl) {
		return false
	}

	_, err := net.ResolveIPAddr("ip", fl.Field().String())

	return err == nil
}


func isUnixAddrResolvable(fl FieldLevel) bool {
	_, err := net.ResolveUnixAddr("unix", fl.Field().String())

	return err == nil
}

func isIP4Addr(fl FieldLevel) bool {
	val := fl.Field().String()

	if idx := strings.LastIndex(val, ":"); idx != -1 {
		val = val[0:idx]
	}

	ip := net.ParseIP(val)

	return ip != nil && ip.To4() != nil
}

func isIP6Addr(fl FieldLevel) bool {
	val := fl.Field().String()

	if idx := strings.LastIndex(val, ":"); idx != -1 {
		if idx != 0 && val[idx-1:idx] == "]" {
			val = val[1 : idx-1]
		}
	}

	ip := net.ParseIP(val)

	return ip != nil && ip.To4() == nil
}

func isHostnameRFC952(fl FieldLevel) bool {
	return hostnameRegexRFC952().MatchString(fl.Field().String())
}

func isHostnameRFC1123(fl FieldLevel) bool {
	return hostnameRegexRFC1123().MatchString(fl.Field().String())
}

func isFQDN(fl FieldLevel) bool {
	val := fl.Field().String()

	if val == "" {
		return false
	}

	return fqdnRegexRFC1123().MatchString(val)
}


func isDir(fl FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.String {
		fileInfo, err := os.Stat(field.String())
		if err != nil {
			return false
		}

		return fileInfo.IsDir()
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isDirPath(fl FieldLevel) bool {
	var exists bool
	var err error

	field := fl.Field()

	
	
	if exists = isDir(fl); exists {
		return true
	}

	
	switch field.Kind() {
	case reflect.String:
		
		
		
		if strings.TrimSpace(field.String()) == "" {
			return false
		}
		if _, err = os.Stat(field.String()); err != nil {
			switch t := err.(type) {
			case *fs.PathError:
				if t.Err == syscall.EINVAL {
					
					return false
				}
				
				
				
				if strings.HasSuffix(field.String(), string(os.PathSeparator)) {
					return true
				} else {
					return false
				}
			default:
				
				
				panic(err)
			}
		}
		
		if strings.HasSuffix(field.String(), string(os.PathSeparator)) {
			return true
		} else {
			return false
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isJSON(fl FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		val := field.String()
		return json.Valid([]byte(val))
	case reflect.Slice:
		fieldType := field.Type()

		if fieldType.ConvertibleTo(byteSliceType) {
			b := field.Convert(byteSliceType).Interface().([]byte)
			return json.Valid(b)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isJWT(fl FieldLevel) bool {
	return jWTRegex().MatchString(fl.Field().String())
}


func isHostnamePort(fl FieldLevel) bool {
	val := fl.Field().String()
	host, port, err := net.SplitHostPort(val)
	if err != nil {
		return false
	}
	
	if portNum, err := strconv.ParseInt(
		port, 10, 32,
	); err != nil || portNum > 65535 || portNum < 1 {
		return false
	}

	
	if host != "" {
		return hostnameRegexRFC1123().MatchString(host)
	}
	return true
}


func isPort(fl FieldLevel) bool {
	val := fl.Field().Uint()

	return val >= 1 && val <= 65535
}


func isLowercase(fl FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.String {
		if field.String() == "" {
			return false
		}
		return field.String() == strings.ToLower(field.String())
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isUppercase(fl FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.String {
		if field.String() == "" {
			return false
		}
		return field.String() == strings.ToUpper(field.String())
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isDatetime(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	if field.Kind() == reflect.String {
		_, err := time.Parse(param, field.String())

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isTimeZone(fl FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.String {
		
		if field.String() == "" {
			return false
		}

		
		if strings.ToLower(field.String()) == "local" {
			return false
		}

		_, err := time.LoadLocation(field.String())
		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isIso3166Alpha2(fl FieldLevel) bool {
	_, ok := iso3166_1_alpha2[fl.Field().String()]
	return ok
}


func isIso3166Alpha2EU(fl FieldLevel) bool {
	_, ok := iso3166_1_alpha2_eu[fl.Field().String()]
	return ok
}


func isIso3166Alpha3(fl FieldLevel) bool {
	_, ok := iso3166_1_alpha3[fl.Field().String()]
	return ok
}


func isIso3166Alpha3EU(fl FieldLevel) bool {
	_, ok := iso3166_1_alpha3_eu[fl.Field().String()]
	return ok
}


func isIso3166AlphaNumeric(fl FieldLevel) bool {
	field := fl.Field()

	var code int
	switch field.Kind() {
	case reflect.String:
		i, err := strconv.Atoi(field.String())
		if err != nil {
			return false
		}
		code = i % 1000
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		code = int(field.Int() % 1000)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		code = int(field.Uint() % 1000)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	_, ok := iso3166_1_alpha_numeric[code]
	return ok
}


func isIso3166AlphaNumericEU(fl FieldLevel) bool {
	field := fl.Field()

	var code int
	switch field.Kind() {
	case reflect.String:
		i, err := strconv.Atoi(field.String())
		if err != nil {
			return false
		}
		code = i % 1000
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		code = int(field.Int() % 1000)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		code = int(field.Uint() % 1000)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	_, ok := iso3166_1_alpha_numeric_eu[code]
	return ok
}


func isIso31662(fl FieldLevel) bool {
	_, ok := iso3166_2[fl.Field().String()]
	return ok
}


func isIso4217(fl FieldLevel) bool {
	_, ok := iso4217[fl.Field().String()]
	return ok
}


func isIso4217Numeric(fl FieldLevel) bool {
	field := fl.Field()

	var code int
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		code = int(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		code = int(field.Uint())
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	_, ok := iso4217_numeric[code]
	return ok
}


func isBCP47LanguageTag(fl FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.String {
		_, err := language.Parse(field.String())
		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}


func isIsoBicFormat(fl FieldLevel) bool {
	bicString := fl.Field().String()

	return bicRegex().MatchString(bicString)
}


func isSemverFormat(fl FieldLevel) bool {
	semverString := fl.Field().String()

	return semverRegex().MatchString(semverString)
}


func isCveFormat(fl FieldLevel) bool {
	cveString := fl.Field().String()

	return cveRegex().MatchString(cveString)
}




func isDnsRFC1035LabelFormat(fl FieldLevel) bool {
	val := fl.Field().String()

	size := len(val)
	if size > 63 {
		return false
	}

	return dnsRegexRFC1035Label().MatchString(val)
}


func digitsHaveLuhnChecksum(digits []string) bool {
	size := len(digits)
	sum := 0
	for i, digit := range digits {
		value, err := strconv.Atoi(digit)
		if err != nil {
			return false
		}
		if size%2 == 0 && i%2 == 0 || size%2 == 1 && i%2 == 1 {
			v := value * 2
			if v >= 10 {
				sum += 1 + (v % 10)
			} else {
				sum += v
			}
		} else {
			sum += value
		}
	}
	return (sum % 10) == 0
}


func isMongoDBObjectId(fl FieldLevel) bool {
	val := fl.Field().String()
	return mongodbIdRegex().MatchString(val)
}


func isMongoDBConnectionString(fl FieldLevel) bool {
	val := fl.Field().String()
	return mongodbConnectionRegex().MatchString(val)
}


func isSpiceDB(fl FieldLevel) bool {
	val := fl.Field().String()
	param := fl.Param()

	switch param {
	case "permission":
		return spicedbPermissionRegex().MatchString(val)
	case "type":
		return spicedbTypeRegex().MatchString(val)
	case "id", "":
		return spicedbIDRegex().MatchString(val)
	}

	panic("Unrecognized parameter: " + param)
}


func isCreditCard(fl FieldLevel) bool {
	val := fl.Field().String()
	var creditCard bytes.Buffer
	segments := strings.Split(val, " ")
	for _, segment := range segments {
		if len(segment) < 3 {
			return false
		}
		creditCard.WriteString(segment)
	}

	ccDigits := strings.Split(creditCard.String(), "")
	size := len(ccDigits)
	if size < 12 || size > 19 {
		return false
	}

	return digitsHaveLuhnChecksum(ccDigits)
}


func hasLuhnChecksum(fl FieldLevel) bool {
	field := fl.Field()
	var str string 
	switch field.Kind() {
	case reflect.String:
		str = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	size := len(str)
	if size < 2 { 
		return false
	}
	digits := strings.Split(str, "")
	return digitsHaveLuhnChecksum(digits)
}


func isCron(fl FieldLevel) bool {
	cronString := fl.Field().String()
	return cronRegex().MatchString(cronString)
}


func isEIN(fl FieldLevel) bool {
	field := fl.Field()

	if field.Len() != 10 {
		return false
	}

	return einRegex().MatchString(field.String())
}
