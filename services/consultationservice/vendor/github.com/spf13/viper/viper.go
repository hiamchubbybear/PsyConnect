


















package viper

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"

	"github.com/spf13/viper/internal/features"
)


type ConfigMarshalError struct {
	err error
}


func (e ConfigMarshalError) Error() string {
	return fmt.Sprintf("While marshaling config: %s", e.err.Error())
}

var v *Viper

func init() {
	v = New()
}



type UnsupportedConfigError string


func (str UnsupportedConfigError) Error() string {
	return fmt.Sprintf("Unsupported Config Type %q", string(str))
}


type ConfigFileNotFoundError struct {
	name, locations string
}


func (fnfe ConfigFileNotFoundError) Error() string {
	return fmt.Sprintf("Config File %q Not Found in %q", fnfe.name, fnfe.locations)
}


type ConfigFileAlreadyExistsError string


func (faee ConfigFileAlreadyExistsError) Error() string {
	return fmt.Sprintf("Config File %q Already Exists", string(faee))
}



type DecoderConfigOption func(*mapstructure.DecoderConfig)








func DecodeHook(hook mapstructure.DecodeHookFunc) DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = hook
	}
}





































type Viper struct {
	
	
	keyDelim string

	
	configPaths []string

	
	fs afero.Fs

	finder Finder

	
	remoteProviders []*defaultRemoteProvider

	
	configName        string
	configFile        string
	configType        string
	configPermissions os.FileMode
	envPrefix         string

	automaticEnvApplied bool
	envKeyReplacer      StringReplacer
	allowEmptyEnv       bool

	parents        []string
	config         map[string]any
	override       map[string]any
	defaults       map[string]any
	kvstore        map[string]any
	pflags         map[string]FlagValue
	env            map[string][]string
	aliases        map[string]string
	typeByDefValue bool

	onConfigChange func(fsnotify.Event)

	logger *slog.Logger

	encoderRegistry EncoderRegistry
	decoderRegistry DecoderRegistry

	decodeHook mapstructure.DecodeHookFunc

	experimentalFinder     bool
	experimentalBindStruct bool
}


func New() *Viper {
	v := new(Viper)
	v.keyDelim = "."
	v.configName = "config"
	v.configPermissions = os.FileMode(0o644)
	v.fs = afero.NewOsFs()
	v.config = make(map[string]any)
	v.parents = []string{}
	v.override = make(map[string]any)
	v.defaults = make(map[string]any)
	v.kvstore = make(map[string]any)
	v.pflags = make(map[string]FlagValue)
	v.env = make(map[string][]string)
	v.aliases = make(map[string]string)
	v.typeByDefValue = false
	v.logger = slog.New(&discardHandler{})

	codecRegistry := NewCodecRegistry()

	v.encoderRegistry = codecRegistry
	v.decoderRegistry = codecRegistry

	v.experimentalFinder = features.Finder
	v.experimentalBindStruct = features.BindStruct

	return v
}





type Option interface {
	apply(v *Viper)
}

type optionFunc func(v *Viper)

func (fn optionFunc) apply(v *Viper) {
	fn(v)
}



func KeyDelimiter(d string) Option {
	return optionFunc(func(v *Viper) {
		v.keyDelim = d
	})
}


type StringReplacer interface {
	
	Replace(s string) string
}


func EnvKeyReplacer(r StringReplacer) Option {
	return optionFunc(func(v *Viper) {
		if r == nil {
			return
		}

		v.envKeyReplacer = r
	})
}


func WithDecodeHook(h mapstructure.DecodeHookFunc) Option {
	return optionFunc(func(v *Viper) {
		if h == nil {
			return
		}

		v.decodeHook = h
	})
}


func NewWithOptions(opts ...Option) *Viper {
	v := New()

	for _, opt := range opts {
		opt.apply(v)
	}

	return v
}





func SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt.apply(v)
	}
}




func Reset() {
	v = New()
	SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}

	resetRemote()
}


var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}


func OnConfigChange(run func(in fsnotify.Event)) { v.OnConfigChange(run) }


func (v *Viper) OnConfigChange(run func(in fsnotify.Event)) {
	v.onConfigChange = run
}


func WatchConfig() { v.WatchConfig() }


func (v *Viper) WatchConfig() {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			v.logger.Error(fmt.Sprintf("failed to create watcher: %s", err))
			os.Exit(1)
		}
		defer watcher.Close()
		
		filename, err := v.getConfigFile()
		if err != nil {
			v.logger.Error(fmt.Sprintf("get config file: %s", err))
			initWG.Done()
			return
		}

		configFile := filepath.Clean(filename)
		configDir, _ := filepath.Split(configFile)
		realConfigFile, _ := filepath.EvalSymlinks(filename)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { 
						eventsWG.Done()
						return
					}
					currentConfigFile, _ := filepath.EvalSymlinks(filename)
					
					
					
					if (filepath.Clean(event.Name) == configFile &&
						(event.Has(fsnotify.Write) || event.Has(fsnotify.Create))) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {
						realConfigFile = currentConfigFile
						err := v.ReadInConfig()
						if err != nil {
							v.logger.Error(fmt.Sprintf("read config file: %s", err))
						}
						if v.onConfigChange != nil {
							v.onConfigChange(event)
						}
					} else if filepath.Clean(event.Name) == configFile && event.Has(fsnotify.Remove) {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok { 
						v.logger.Error(fmt.Sprintf("watcher error: %s", err))
					}
					eventsWG.Done()
					return
				}
			}
		}()
		watcher.Add(configDir)
		initWG.Done()   
		eventsWG.Wait() 
	}()
	initWG.Wait() 
}



func SetConfigFile(in string) { v.SetConfigFile(in) }

func (v *Viper) SetConfigFile(in string) {
	if in != "" {
		v.configFile = in
	}
}




func SetEnvPrefix(in string) { v.SetEnvPrefix(in) }

func (v *Viper) SetEnvPrefix(in string) {
	if in != "" {
		v.envPrefix = in
	}
}

func GetEnvPrefix() string { return v.GetEnvPrefix() }

func (v *Viper) GetEnvPrefix() string {
	return v.envPrefix
}

func (v *Viper) mergeWithEnvPrefix(in string) string {
	if v.envPrefix != "" {
		return strings.ToUpper(v.envPrefix + "_" + in)
	}

	return strings.ToUpper(in)
}




func AllowEmptyEnv(allowEmptyEnv bool) { v.AllowEmptyEnv(allowEmptyEnv) }

func (v *Viper) AllowEmptyEnv(allowEmptyEnv bool) {
	v.allowEmptyEnv = allowEmptyEnv
}








func (v *Viper) getEnv(key string) (string, bool) {
	if v.envKeyReplacer != nil {
		key = v.envKeyReplacer.Replace(key)
	}

	val, ok := os.LookupEnv(key)

	return val, ok && (v.allowEmptyEnv || val != "")
}


func ConfigFileUsed() string            { return v.ConfigFileUsed() }
func (v *Viper) ConfigFileUsed() string { return v.configFile }



func AddConfigPath(in string) { v.AddConfigPath(in) }

func (v *Viper) AddConfigPath(in string) {
	if v.finder != nil {
		v.logger.Warn("ineffective call to function: custom finder takes precedence", slog.String("function", "AddConfigPath"))
	}

	if in != "" {
		absin := absPathify(v.logger, in)

		v.logger.Info("adding path to search paths", "path", absin)
		if !slices.Contains(v.configPaths, absin) {
			v.configPaths = append(v.configPaths, absin)
		}
	}
}




func (v *Viper) searchMap(source map[string]any, path []string) any {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		
		if len(path) == 1 {
			return next
		}

		
		switch next := next.(type) {
		case map[any]any:
			return v.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]any:
			
			
			return v.searchMap(next, path[1:])
		default:
			
			return nil
		}
	}
	return nil
}












func (v *Viper) searchIndexableWithPathPrefixes(source any, path []string) any {
	if len(path) == 0 {
		return source
	}

	
	for i := len(path); i > 0; i-- {
		prefixKey := strings.ToLower(strings.Join(path[0:i], v.keyDelim))

		var val any
		switch sourceIndexable := source.(type) {
		case []any:
			val = v.searchSliceWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		case map[string]any:
			val = v.searchMapWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		}
		if val != nil {
			return val
		}
	}

	
	return nil
}





func (v *Viper) searchSliceWithPathPrefixes(
	sourceSlice []any,
	prefixKey string,
	pathIndex int,
	path []string,
) any {
	
	index, err := strconv.Atoi(prefixKey)
	if err != nil || len(sourceSlice) <= index {
		return nil
	}

	next := sourceSlice[index]

	
	if pathIndex == len(path) {
		return next
	}

	switch n := next.(type) {
	case map[any]any:
		return v.searchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:])
	case map[string]any, []any:
		return v.searchIndexableWithPathPrefixes(n, path[pathIndex:])
	default:
		
	}

	
	return nil
}





func (v *Viper) searchMapWithPathPrefixes(
	sourceMap map[string]any,
	prefixKey string,
	pathIndex int,
	path []string,
) any {
	next, ok := sourceMap[prefixKey]
	if !ok {
		return nil
	}

	
	if pathIndex == len(path) {
		return next
	}

	
	switch n := next.(type) {
	case map[any]any:
		return v.searchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:])
	case map[string]any, []any:
		return v.searchIndexableWithPathPrefixes(n, path[pathIndex:])
	default:
		
	}

	
	return nil
}






func (v *Viper) isPathShadowedInDeepMap(path []string, m map[string]any) string {
	var parentVal any
	for i := 1; i < len(path); i++ {
		parentVal = v.searchMap(m, path[0:i])
		if parentVal == nil {
			
			return ""
		}
		switch parentVal.(type) {
		case map[any]any:
			continue
		case map[string]any:
			continue
		default:
			
			return strings.Join(path[0:i], v.keyDelim)
		}
	}
	return ""
}






func (v *Viper) isPathShadowedInFlatMap(path []string, mi any) string {
	
	var m map[string]interface{}
	switch miv := mi.(type) {
	case map[string]string:
		m = castMapStringToMapInterface(miv)
	case map[string]FlagValue:
		m = castMapFlagToMapInterface(miv)
	default:
		return ""
	}

	
	var parentKey string
	for i := 1; i < len(path); i++ {
		parentKey = strings.Join(path[0:i], v.keyDelim)
		if _, ok := m[parentKey]; ok {
			return parentKey
		}
	}
	return ""
}






func (v *Viper) isPathShadowedInAutoEnv(path []string) string {
	var parentKey string
	for i := 1; i < len(path); i++ {
		parentKey = strings.Join(path[0:i], v.keyDelim)
		if _, ok := v.getEnv(v.mergeWithEnvPrefix(parentKey)); ok {
			return parentKey
		}
	}
	return ""
}















func SetTypeByDefaultValue(enable bool) { v.SetTypeByDefaultValue(enable) }

func (v *Viper) SetTypeByDefaultValue(enable bool) {
	v.typeByDefValue = enable
}


func GetViper() *Viper {
	return v
}








func Get(key string) any { return v.Get(key) }

func (v *Viper) Get(key string) any {
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey, true)
	if val == nil {
		return nil
	}

	if v.typeByDefValue {
		
		valType := val
		path := strings.Split(lcaseKey, v.keyDelim)
		defVal := v.searchMap(v.defaults, path)
		if defVal != nil {
			valType = defVal
		}

		switch valType.(type) {
		case bool:
			return cast.ToBool(val)
		case string:
			return cast.ToString(val)
		case int32, int16, int8, int:
			return cast.ToInt(val)
		case uint:
			return cast.ToUint(val)
		case uint32:
			return cast.ToUint32(val)
		case uint64:
			return cast.ToUint64(val)
		case int64:
			return cast.ToInt64(val)
		case float64, float32:
			return cast.ToFloat64(val)
		case time.Time:
			return cast.ToTime(val)
		case time.Duration:
			return cast.ToDuration(val)
		case []string:
			return cast.ToStringSlice(val)
		case []int:
			return cast.ToIntSlice(val)
		case []time.Duration:
			return cast.ToDurationSlice(val)
		}
	}

	return val
}



func Sub(key string) *Viper { return v.Sub(key) }

func (v *Viper) Sub(key string) *Viper {
	subv := New()
	data := v.Get(key)
	if data == nil {
		return nil
	}

	if reflect.TypeOf(data).Kind() == reflect.Map {
		subv.parents = append([]string(nil), v.parents...)
		subv.parents = append(subv.parents, strings.ToLower(key))
		subv.automaticEnvApplied = v.automaticEnvApplied
		subv.envPrefix = v.envPrefix
		subv.envKeyReplacer = v.envKeyReplacer
		subv.keyDelim = v.keyDelim
		subv.config = cast.ToStringMap(data)
		return subv
	}
	return nil
}


func GetString(key string) string { return v.GetString(key) }

func (v *Viper) GetString(key string) string {
	return cast.ToString(v.Get(key))
}


func GetBool(key string) bool { return v.GetBool(key) }

func (v *Viper) GetBool(key string) bool {
	return cast.ToBool(v.Get(key))
}


func GetInt(key string) int { return v.GetInt(key) }

func (v *Viper) GetInt(key string) int {
	return cast.ToInt(v.Get(key))
}


func GetInt32(key string) int32 { return v.GetInt32(key) }

func (v *Viper) GetInt32(key string) int32 {
	return cast.ToInt32(v.Get(key))
}


func GetInt64(key string) int64 { return v.GetInt64(key) }

func (v *Viper) GetInt64(key string) int64 {
	return cast.ToInt64(v.Get(key))
}


func GetUint8(key string) uint8 { return v.GetUint8(key) }

func (v *Viper) GetUint8(key string) uint8 {
	return cast.ToUint8(v.Get(key))
}


func GetUint(key string) uint { return v.GetUint(key) }

func (v *Viper) GetUint(key string) uint {
	return cast.ToUint(v.Get(key))
}


func GetUint16(key string) uint16 { return v.GetUint16(key) }

func (v *Viper) GetUint16(key string) uint16 {
	return cast.ToUint16(v.Get(key))
}


func GetUint32(key string) uint32 { return v.GetUint32(key) }

func (v *Viper) GetUint32(key string) uint32 {
	return cast.ToUint32(v.Get(key))
}


func GetUint64(key string) uint64 { return v.GetUint64(key) }

func (v *Viper) GetUint64(key string) uint64 {
	return cast.ToUint64(v.Get(key))
}


func GetFloat64(key string) float64 { return v.GetFloat64(key) }

func (v *Viper) GetFloat64(key string) float64 {
	return cast.ToFloat64(v.Get(key))
}


func GetTime(key string) time.Time { return v.GetTime(key) }

func (v *Viper) GetTime(key string) time.Time {
	return cast.ToTime(v.Get(key))
}


func GetDuration(key string) time.Duration { return v.GetDuration(key) }

func (v *Viper) GetDuration(key string) time.Duration {
	return cast.ToDuration(v.Get(key))
}


func GetIntSlice(key string) []int { return v.GetIntSlice(key) }

func (v *Viper) GetIntSlice(key string) []int {
	return cast.ToIntSlice(v.Get(key))
}


func GetStringSlice(key string) []string { return v.GetStringSlice(key) }

func (v *Viper) GetStringSlice(key string) []string {
	return cast.ToStringSlice(v.Get(key))
}


func GetStringMap(key string) map[string]any { return v.GetStringMap(key) }

func (v *Viper) GetStringMap(key string) map[string]any {
	return cast.ToStringMap(v.Get(key))
}


func GetStringMapString(key string) map[string]string { return v.GetStringMapString(key) }

func (v *Viper) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(v.Get(key))
}


func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }

func (v *Viper) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(v.Get(key))
}



func GetSizeInBytes(key string) uint { return v.GetSizeInBytes(key) }

func (v *Viper) GetSizeInBytes(key string) uint {
	sizeStr := cast.ToString(v.Get(key))
	return parseSizeInBytes(sizeStr)
}


func UnmarshalKey(key string, rawVal any, opts ...DecoderConfigOption) error {
	return v.UnmarshalKey(key, rawVal, opts...)
}

func (v *Viper) UnmarshalKey(key string, rawVal any, opts ...DecoderConfigOption) error {
	return decode(v.Get(key), v.defaultDecoderConfig(rawVal, opts...))
}



func Unmarshal(rawVal any, opts ...DecoderConfigOption) error {
	return v.Unmarshal(rawVal, opts...)
}

func (v *Viper) Unmarshal(rawVal any, opts ...DecoderConfigOption) error {
	keys := v.AllKeys()

	if v.experimentalBindStruct {
		
		structKeys, err := v.decodeStructKeys(rawVal, opts...)
		if err != nil {
			return err
		}

		keys = append(keys, structKeys...)
	}

	
	return decode(v.getSettings(keys), v.defaultDecoderConfig(rawVal, opts...))
}

func (v *Viper) decodeStructKeys(input any, opts ...DecoderConfigOption) ([]string, error) {
	var structKeyMap map[string]any

	err := decode(input, v.defaultDecoderConfig(&structKeyMap, opts...))
	if err != nil {
		return nil, err
	}

	flattenedStructKeyMap := v.flattenAndMergeMap(map[string]bool{}, structKeyMap, "")

	r := make([]string, 0, len(flattenedStructKeyMap))
	for v := range flattenedStructKeyMap {
		r = append(r, v)
	}

	return r, nil
}



func (v *Viper) defaultDecoderConfig(output any, opts ...DecoderConfigOption) *mapstructure.DecoderConfig {
	decodeHook := v.decodeHook
	if decodeHook == nil {
		decodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			
			stringToWeakSliceHookFunc(","),
		)
	}

	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
	}

	for _, opt := range opts {
		opt(c)
	}

	
	c.Result = output

	return c
}




func stringToWeakSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Slice {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		return strings.Split(raw, sep), nil
	}
}


func decode(input any, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}



func UnmarshalExact(rawVal any, opts ...DecoderConfigOption) error {
	return v.UnmarshalExact(rawVal, opts...)
}

func (v *Viper) UnmarshalExact(rawVal any, opts ...DecoderConfigOption) error {
	config := v.defaultDecoderConfig(rawVal, opts...)
	config.ErrorUnused = true

	keys := v.AllKeys()

	if v.experimentalBindStruct {
		
		structKeys, err := v.decodeStructKeys(rawVal, opts...)
		if err != nil {
			return err
		}

		keys = append(keys, structKeys...)
	}

	
	return decode(v.getSettings(keys), config)
}



func BindPFlags(flags *pflag.FlagSet) error { return v.BindPFlags(flags) }

func (v *Viper) BindPFlags(flags *pflag.FlagSet) error {
	return v.BindFlagValues(pflagValueSet{flags})
}






func BindPFlag(key string, flag *pflag.Flag) error { return v.BindPFlag(key, flag) }

func (v *Viper) BindPFlag(key string, flag *pflag.Flag) error {
	if flag == nil {
		return fmt.Errorf("flag for %q is nil", key)
	}
	return v.BindFlagValue(key, pflagValue{flag})
}



func BindFlagValues(flags FlagValueSet) error { return v.BindFlagValues(flags) }

func (v *Viper) BindFlagValues(flags FlagValueSet) (err error) {
	flags.VisitAll(func(flag FlagValue) {
		if err = v.BindFlagValue(flag.Name(), flag); err != nil {
			return
		}
	})
	return nil
}


func BindFlagValue(key string, flag FlagValue) error { return v.BindFlagValue(key, flag) }

func (v *Viper) BindFlagValue(key string, flag FlagValue) error {
	if flag == nil {
		return fmt.Errorf("flag for %q is nil", key)
	}
	v.pflags[strings.ToLower(key)] = flag
	return nil
}







func BindEnv(input ...string) error { return v.BindEnv(input...) }

func (v *Viper) BindEnv(input ...string) error {
	if len(input) == 0 {
		return fmt.Errorf("missing key to bind to")
	}

	key := strings.ToLower(input[0])

	if len(input) == 1 {
		v.env[key] = append(v.env[key], v.mergeWithEnvPrefix(key))
	} else {
		v.env[key] = append(v.env[key], input[1:]...)
	}

	return nil
}




func MustBindEnv(input ...string) { v.MustBindEnv(input...) }

func (v *Viper) MustBindEnv(input ...string) {
	if err := v.BindEnv(input...); err != nil {
		panic(fmt.Sprintf("error while binding environment variable: %v", err))
	}
}










func (v *Viper) find(lcaseKey string, flagDefault bool) any {
	var (
		val    any
		exists bool
		path   = strings.Split(lcaseKey, v.keyDelim)
		nested = len(path) > 1
	)

	
	if nested && v.isPathShadowedInDeepMap(path, castMapStringToMapInterface(v.aliases)) != "" {
		return nil
	}

	
	lcaseKey = v.realKey(lcaseKey)
	path = strings.Split(lcaseKey, v.keyDelim)
	nested = len(path) > 1

	
	val = v.searchMap(v.override, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.override) != "" {
		return nil
	}

	
	flag, exists := v.pflags[lcaseKey]
	if exists && flag.HasChanged() {
		switch flag.ValueType() {
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		case "stringSlice", "stringArray":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := readAsCSV(s)
			return res
		case "intSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := readAsCSV(s)
			return cast.ToIntSlice(res)
		case "durationSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			slice := strings.Split(s, ",")
			return cast.ToDurationSlice(slice)
		case "stringToString":
			return stringToStringConv(flag.ValueString())
		case "stringToInt":
			return stringToIntConv(flag.ValueString())
		default:
			return flag.ValueString()
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.pflags) != "" {
		return nil
	}

	
	if v.automaticEnvApplied {
		envKey := strings.Join(append(v.parents, lcaseKey), ".")
		
		
		if val, ok := v.getEnv(v.mergeWithEnvPrefix(envKey)); ok {
			return val
		}
		if nested && v.isPathShadowedInAutoEnv(path) != "" {
			return nil
		}
	}
	envkeys, exists := v.env[lcaseKey]
	if exists {
		for _, envkey := range envkeys {
			if val, ok := v.getEnv(envkey); ok {
				return val
			}
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.env) != "" {
		return nil
	}

	
	val = v.searchIndexableWithPathPrefixes(v.config, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.config) != "" {
		return nil
	}

	
	val = v.searchMap(v.kvstore, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.kvstore) != "" {
		return nil
	}

	
	val = v.searchMap(v.defaults, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.defaults) != "" {
		return nil
	}

	if flagDefault {
		
		
		if flag, exists := v.pflags[lcaseKey]; exists {
			switch flag.ValueType() {
			case "int", "int8", "int16", "int32", "int64":
				return cast.ToInt(flag.ValueString())
			case "bool":
				return cast.ToBool(flag.ValueString())
			case "stringSlice", "stringArray":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return res
			case "intSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return cast.ToIntSlice(res)
			case "stringToString":
				return stringToStringConv(flag.ValueString())
			case "stringToInt":
				return stringToIntConv(flag.ValueString())
			case "durationSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				slice := strings.Split(s, ",")
				return cast.ToDurationSlice(slice)
			default:
				return flag.ValueString()
			}
		}
		
	}

	return nil
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}



func stringToStringConv(val string) any {
	val = strings.Trim(val, "[]")
	
	if val == "" {
		return map[string]any{}
	}
	r := csv.NewReader(strings.NewReader(val))
	ss, err := r.Read()
	if err != nil {
		return nil
	}
	out := make(map[string]any, len(ss))
	for _, pair := range ss {
		k, vv, found := strings.Cut(pair, "=")
		if !found {
			return nil
		}
		out[k] = vv
	}
	return out
}



func stringToIntConv(val string) any {
	val = strings.Trim(val, "[]")
	
	if val == "" {
		return map[string]any{}
	}
	ss := strings.Split(val, ",")
	out := make(map[string]any, len(ss))
	for _, pair := range ss {
		k, vv, found := strings.Cut(pair, "=")
		if !found {
			return nil
		}
		var err error
		out[k], err = strconv.Atoi(vv)
		if err != nil {
			return nil
		}
	}
	return out
}



func IsSet(key string) bool { return v.IsSet(key) }

func (v *Viper) IsSet(key string) bool {
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey, false)
	return val != nil
}



func AutomaticEnv() { v.AutomaticEnv() }

func (v *Viper) AutomaticEnv() {
	v.automaticEnvApplied = true
}




func SetEnvKeyReplacer(r *strings.Replacer) { v.SetEnvKeyReplacer(r) }

func (v *Viper) SetEnvKeyReplacer(r *strings.Replacer) {
	v.envKeyReplacer = r
}



func RegisterAlias(alias, key string) { v.RegisterAlias(alias, key) }

func (v *Viper) RegisterAlias(alias, key string) {
	v.registerAlias(alias, strings.ToLower(key))
}

func (v *Viper) registerAlias(alias, key string) {
	alias = strings.ToLower(alias)
	if alias != key && alias != v.realKey(key) {
		_, exists := v.aliases[alias]

		if !exists {
			
			
			
			if val, ok := v.config[alias]; ok {
				delete(v.config, alias)
				v.config[key] = val
			}
			if val, ok := v.kvstore[alias]; ok {
				delete(v.kvstore, alias)
				v.kvstore[key] = val
			}
			if val, ok := v.defaults[alias]; ok {
				delete(v.defaults, alias)
				v.defaults[key] = val
			}
			if val, ok := v.override[alias]; ok {
				delete(v.override, alias)
				v.override[key] = val
			}
			v.aliases[alias] = key
		}
	} else {
		v.logger.Warn("creating circular reference alias", "alias", alias, "key", key, "real_key", v.realKey(key))
	}
}

func (v *Viper) realKey(key string) string {
	newkey, exists := v.aliases[key]
	if exists {
		v.logger.Debug("key is an alias", "alias", key, "to", newkey)

		return v.realKey(newkey)
	}
	return key
}


func InConfig(key string) bool { return v.InConfig(key) }

func (v *Viper) InConfig(key string) bool {
	lcaseKey := strings.ToLower(key)

	
	lcaseKey = v.realKey(lcaseKey)
	path := strings.Split(lcaseKey, v.keyDelim)

	return v.searchIndexableWithPathPrefixes(v.config, path) != nil
}




func SetDefault(key string, value any) { v.SetDefault(key, value) }

func (v *Viper) SetDefault(key string, value any) {
	
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.defaults, path[0:len(path)-1])

	
	deepestMap[lastKey] = value
}





func Set(key string, value any) { v.Set(key, value) }

func (v *Viper) Set(key string, value any) {
	
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.override, path[0:len(path)-1])

	
	deepestMap[lastKey] = value
}



func ReadInConfig() error { return v.ReadInConfig() }

func (v *Viper) ReadInConfig() error {
	v.logger.Info("attempting to read in config file")
	filename, err := v.getConfigFile()
	if err != nil {
		return err
	}

	if !slices.Contains(SupportedExts, v.getConfigType()) {
		return UnsupportedConfigError(v.getConfigType())
	}

	v.logger.Debug("reading file", "file", filename)
	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	config := make(map[string]any)

	err = v.unmarshalReader(bytes.NewReader(file), config)
	if err != nil {
		return err
	}

	v.config = config
	return nil
}


func MergeInConfig() error { return v.MergeInConfig() }

func (v *Viper) MergeInConfig() error {
	v.logger.Info("attempting to merge in config file")
	filename, err := v.getConfigFile()
	if err != nil {
		return err
	}

	if !slices.Contains(SupportedExts, v.getConfigType()) {
		return UnsupportedConfigError(v.getConfigType())
	}

	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	return v.MergeConfig(bytes.NewReader(file))
}



func ReadConfig(in io.Reader) error { return v.ReadConfig(in) }

func (v *Viper) ReadConfig(in io.Reader) error {
	if v.configType == "" {
		return errors.New("cannot decode configuration: config type is not set")
	}

	v.config = make(map[string]any)
	return v.unmarshalReader(in, v.config)
}


func MergeConfig(in io.Reader) error { return v.MergeConfig(in) }

func (v *Viper) MergeConfig(in io.Reader) error {
	if v.configType == "" {
		return errors.New("cannot decode configuration: config type is not set")
	}

	cfg := make(map[string]any)
	if err := v.unmarshalReader(in, cfg); err != nil {
		return err
	}
	return v.MergeConfigMap(cfg)
}



func MergeConfigMap(cfg map[string]any) error { return v.MergeConfigMap(cfg) }

func (v *Viper) MergeConfigMap(cfg map[string]any) error {
	if v.config == nil {
		v.config = make(map[string]any)
	}
	insensitiviseMap(cfg)
	mergeMaps(cfg, v.config, nil)
	return nil
}


func WriteConfig() error { return v.WriteConfig() }

func (v *Viper) WriteConfig() error {
	filename, err := v.getConfigFile()
	if err != nil {
		return err
	}
	return v.writeConfig(filename, true)
}


func SafeWriteConfig() error { return v.SafeWriteConfig() }

func (v *Viper) SafeWriteConfig() error {
	if len(v.configPaths) < 1 {
		return errors.New("missing configuration for 'configPath'")
	}
	return v.SafeWriteConfigAs(filepath.Join(v.configPaths[0], v.configName+"."+v.configType))
}


func WriteConfigAs(filename string) error { return v.WriteConfigAs(filename) }

func (v *Viper) WriteConfigAs(filename string) error {
	return v.writeConfig(filename, true)
}


func WriteConfigTo(w io.Writer) error { return v.WriteConfigTo(w) }

func (v *Viper) WriteConfigTo(w io.Writer) error {
	format := strings.ToLower(v.getConfigType())

	if !slices.Contains(SupportedExts, format) {
		return UnsupportedConfigError(format)
	}

	return v.marshalWriter(w, format)
}


func SafeWriteConfigAs(filename string) error { return v.SafeWriteConfigAs(filename) }

func (v *Viper) SafeWriteConfigAs(filename string) error {
	alreadyExists, err := afero.Exists(v.fs, filename)
	if alreadyExists && err == nil {
		return ConfigFileAlreadyExistsError(filename)
	}
	return v.writeConfig(filename, false)
}

func (v *Viper) writeConfig(filename string, force bool) error {
	v.logger.Info("attempting to write configuration to file")

	var configType string

	ext := filepath.Ext(filename)
	if ext != "" && ext != filepath.Base(filename) {
		configType = ext[1:]
	} else {
		configType = v.configType
	}
	if configType == "" {
		return fmt.Errorf("config type could not be determined for %s", filename)
	}

	if !slices.Contains(SupportedExts, configType) {
		return UnsupportedConfigError(configType)
	}
	if v.config == nil {
		v.config = make(map[string]any)
	}
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	if !force {
		flags |= os.O_EXCL
	}
	f, err := v.fs.OpenFile(filename, flags, v.configPermissions)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := v.marshalWriter(f, configType); err != nil {
		return err
	}

	return f.Sync()
}

func (v *Viper) unmarshalReader(in io.Reader, c map[string]any) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	format := strings.ToLower(v.getConfigType())

	if !slices.Contains(SupportedExts, format) {
		return UnsupportedConfigError(format)
	}

	decoder, err := v.decoderRegistry.Decoder(format)
	if err != nil {
		return ConfigParseError{err}
	}

	err = decoder.Decode(buf.Bytes(), c)
	if err != nil {
		return ConfigParseError{err}
	}

	insensitiviseMap(c)
	return nil
}


func (v *Viper) marshalWriter(w io.Writer, configType string) error {
	c := v.AllSettings()

	encoder, err := v.encoderRegistry.Encoder(configType)
	if err != nil {
		return ConfigMarshalError{err}
	}

	b, err := encoder.Encode(c)
	if err != nil {
		return ConfigMarshalError{err}
	}

	_, err = w.Write(b)
	if err != nil {
		return ConfigMarshalError{err}
	}

	return nil
}

func keyExists(k string, m map[string]any) string {
	lk := strings.ToLower(k)
	for mk := range m {
		lmk := strings.ToLower(mk)
		if lmk == lk {
			return mk
		}
	}
	return ""
}

func castToMapStringInterface(
	src map[any]any,
) map[string]any {
	tgt := map[string]any{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

func castMapStringSliceToMapInterface(src map[string][]string) map[string]any {
	tgt := map[string]any{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}

func castMapStringToMapInterface(src map[string]string) map[string]any {
	tgt := map[string]any{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}

func castMapFlagToMapInterface(src map[string]FlagValue) map[string]any {
	tgt := map[string]any{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}






func mergeMaps(src, tgt map[string]any, itgt map[any]any) {
	for sk, sv := range src {
		tk := keyExists(sk, tgt)
		if tk == "" {
			v.logger.Debug("", "tk", "\"\"", fmt.Sprintf("tgt[%s]", sk), sv)
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		tv, ok := tgt[tk]
		if !ok {
			v.logger.Debug("", fmt.Sprintf("ok[%s]", tk), false, fmt.Sprintf("tgt[%s]", sk), sv)
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)

		v.logger.Debug(
			"processing",
			"key", sk,
			"st", svType,
			"tt", tvType,
			"sv", sv,
			"tv", tv,
		)

		switch ttv := tv.(type) {
		case map[any]any:
			v.logger.Debug("merging maps (must convert)")
			tsv, ok := sv.(map[any]any)
			if !ok {
				v.logger.Error(
					"Could not cast sv to map[any]any",
					"key", sk,
					"st", svType,
					"tt", tvType,
					"sv", sv,
					"tv", tv,
				)
				continue
			}

			ssv := castToMapStringInterface(tsv)
			stv := castToMapStringInterface(ttv)
			mergeMaps(ssv, stv, ttv)
		case map[string]any:
			v.logger.Debug("merging maps")
			tsv, ok := sv.(map[string]any)
			if !ok {
				v.logger.Error(
					"Could not cast sv to map[string]any",
					"key", sk,
					"st", svType,
					"tt", tvType,
					"sv", sv,
					"tv", tv,
				)
				continue
			}
			mergeMaps(tsv, ttv, nil)
		default:
			v.logger.Debug("setting value")
			tgt[tk] = sv
			if itgt != nil {
				itgt[tk] = sv
			}
		}
	}
}



func AllKeys() []string { return v.AllKeys() }

func (v *Viper) AllKeys() []string {
	m := map[string]bool{}
	
	m = v.flattenAndMergeMap(m, castMapStringToMapInterface(v.aliases), "")
	m = v.flattenAndMergeMap(m, v.override, "")
	m = v.mergeFlatMap(m, castMapFlagToMapInterface(v.pflags))
	m = v.mergeFlatMap(m, castMapStringSliceToMapInterface(v.env))
	m = v.flattenAndMergeMap(m, v.config, "")
	m = v.flattenAndMergeMap(m, v.kvstore, "")
	m = v.flattenAndMergeMap(m, v.defaults, "")

	
	a := make([]string, 0, len(m))
	for x := range m {
		a = append(a, x)
	}
	return a
}








func (v *Viper) flattenAndMergeMap(shadow map[string]bool, m map[string]any, prefix string) map[string]bool {
	if shadow != nil && prefix != "" && shadow[prefix] {
		
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]bool)
	}

	var m2 map[string]any
	if prefix != "" {
		prefix += v.keyDelim
	}
	for k, val := range m {
		fullKey := prefix + k
		switch val := val.(type) {
		case map[string]any:
			m2 = val
		case map[any]any:
			m2 = cast.ToStringMap(val)
		default:
			
			shadow[strings.ToLower(fullKey)] = true
			continue
		}
		
		shadow = v.flattenAndMergeMap(shadow, m2, fullKey)
	}
	return shadow
}



func (v *Viper) mergeFlatMap(shadow map[string]bool, m map[string]any) map[string]bool {
	
outer:
	for k := range m {
		path := strings.Split(k, v.keyDelim)
		
		var parentKey string
		for i := 1; i < len(path); i++ {
			parentKey = strings.Join(path[0:i], v.keyDelim)
			if shadow[parentKey] {
				
				continue outer
			}
		}
		
		shadow[strings.ToLower(k)] = true
	}
	return shadow
}


func AllSettings() map[string]any { return v.AllSettings() }

func (v *Viper) AllSettings() map[string]any {
	return v.getSettings(v.AllKeys())
}

func (v *Viper) getSettings(keys []string) map[string]any {
	m := map[string]any{}
	
	for _, k := range keys {
		value := v.Get(k)
		if value == nil {
			
			
			continue
		}
		path := strings.Split(k, v.keyDelim)
		lastKey := strings.ToLower(path[len(path)-1])
		deepestMap := deepSearch(m, path[0:len(path)-1])
		
		deepestMap[lastKey] = value
	}
	return m
}


func SetFs(fs afero.Fs) { v.SetFs(fs) }

func (v *Viper) SetFs(fs afero.Fs) {
	v.fs = fs
}



func SetConfigName(in string) { v.SetConfigName(in) }

func (v *Viper) SetConfigName(in string) {
	if v.finder != nil {
		v.logger.Warn("ineffective call to function: custom finder takes precedence", slog.String("function", "SetConfigName"))
	}

	if in != "" {
		v.configName = in
		v.configFile = ""
	}
}



func SetConfigType(in string) { v.SetConfigType(in) }

func (v *Viper) SetConfigType(in string) {
	if in != "" {
		v.configType = in
	}
}


func SetConfigPermissions(perm os.FileMode) { v.SetConfigPermissions(perm) }

func (v *Viper) SetConfigPermissions(perm os.FileMode) {
	v.configPermissions = perm.Perm()
}

func (v *Viper) getConfigType() string {
	if v.configType != "" {
		return v.configType
	}

	cf, err := v.getConfigFile()
	if err != nil {
		return ""
	}

	ext := filepath.Ext(cf)

	if len(ext) > 1 {
		return ext[1:]
	}

	return ""
}

func (v *Viper) getConfigFile() (string, error) {
	if v.configFile == "" {
		cf, err := v.findConfigFile()
		if err != nil {
			return "", err
		}
		v.configFile = cf
	}
	return v.configFile, nil
}



func Debug()              { v.Debug() }
func DebugTo(w io.Writer) { v.DebugTo(w) }

func (v *Viper) Debug() { v.DebugTo(os.Stdout) }

func (v *Viper) DebugTo(w io.Writer) {
	fmt.Fprintf(w, "Aliases:\n%#v\n", v.aliases)
	fmt.Fprintf(w, "Override:\n%#v\n", v.override)
	fmt.Fprintf(w, "PFlags:\n%#v\n", v.pflags)
	fmt.Fprintf(w, "Env:\n%#v\n", v.env)
	fmt.Fprintf(w, "Key/Value Store:\n%#v\n", v.kvstore)
	fmt.Fprintf(w, "Config:\n%#v\n", v.config)
	fmt.Fprintf(w, "Defaults:\n%#v\n", v.defaults)
}
