package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9/internal/hashtag"
)

type HashCmdable interface {
	HDel(ctx context.Context, key string, fields ...string) *IntCmd
	HExists(ctx context.Context, key, field string) *BoolCmd
	HGet(ctx context.Context, key, field string) *StringCmd
	HGetAll(ctx context.Context, key string) *MapStringStringCmd
	HGetDel(ctx context.Context, key string, fields ...string) *StringSliceCmd
	HGetEX(ctx context.Context, key string, fields ...string) *StringSliceCmd
	HGetEXWithArgs(ctx context.Context, key string, options *HGetEXOptions, fields ...string) *StringSliceCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd
	HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd
	HKeys(ctx context.Context, key string) *StringSliceCmd
	HLen(ctx context.Context, key string) *IntCmd
	HMGet(ctx context.Context, key string, fields ...string) *SliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmd
	HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmd
	HSetEX(ctx context.Context, key string, fieldsAndValues ...string) *IntCmd
	HSetEXWithArgs(ctx context.Context, key string, options *HSetEXOptions, fieldsAndValues ...string) *IntCmd
	HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
	HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
	HVals(ctx context.Context, key string) *StringSliceCmd
	HRandField(ctx context.Context, key string, count int) *StringSliceCmd
	HRandFieldWithValues(ctx context.Context, key string, count int) *KeyValueSliceCmd
	HStrLen(ctx context.Context, key, field string) *IntCmd
	HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *IntSliceCmd
	HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd
	HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *IntSliceCmd
	HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd
	HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *IntSliceCmd
	HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd
	HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *IntSliceCmd
	HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd
	HPersist(ctx context.Context, key string, fields ...string) *IntSliceCmd
	HExpireTime(ctx context.Context, key string, fields ...string) *IntSliceCmd
	HPExpireTime(ctx context.Context, key string, fields ...string) *IntSliceCmd
	HTTL(ctx context.Context, key string, fields ...string) *IntSliceCmd
	HPTTL(ctx context.Context, key string, fields ...string) *IntSliceCmd
}

func (c cmdable) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	args := make([]interface{}, 2+len(fields))
	args[0] = "hdel"
	args[1] = key
	for i, field := range fields {
		args[2+i] = field
	}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HExists(ctx context.Context, key, field string) *BoolCmd {
	cmd := NewBoolCmd(ctx, "hexists", key, field)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HGet(ctx context.Context, key, field string) *StringCmd {
	cmd := NewStringCmd(ctx, "hget", key, field)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HGetAll(ctx context.Context, key string) *MapStringStringCmd {
	cmd := NewMapStringStringCmd(ctx, "hgetall", key)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd {
	cmd := NewIntCmd(ctx, "hincrby", key, field, incr)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd {
	cmd := NewFloatCmd(ctx, "hincrbyfloat", key, field, incr)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HKeys(ctx context.Context, key string) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "hkeys", key)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HLen(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(ctx, "hlen", key)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) HMGet(ctx context.Context, key string, fields ...string) *SliceCmd {
	args := make([]interface{}, 2+len(fields))
	args[0] = "hmget"
	args[1] = key
	for i, field := range fields {
		args[2+i] = field
	}
	cmd := NewSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

























func (c cmdable) HSet(ctx context.Context, key string, values ...interface{}) *IntCmd {
	args := make([]interface{}, 2, 2+len(values))
	args[0] = "hset"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmd {
	args := make([]interface{}, 2, 2+len(values))
	args[0] = "hmset"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd {
	cmd := NewBoolCmd(ctx, "hsetnx", key, field, value)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HVals(ctx context.Context, key string) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "hvals", key)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) HRandField(ctx context.Context, key string, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "hrandfield", key, count)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) HRandFieldWithValues(ctx context.Context, key string, count int) *KeyValueSliceCmd {
	cmd := NewKeyValueSliceCmd(ctx, "hrandfield", key, count, "withvalues")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	args := []interface{}{"hscan", key, cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	cmd := NewScanCmd(ctx, c, args...)
	if hashtag.Present(match) {
		cmd.SetFirstKeyPos(4)
	}
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HStrLen(ctx context.Context, key, field string) *IntCmd {
	cmd := NewIntCmd(ctx, "hstrlen", key, field)
	_ = c(ctx, cmd)
	return cmd
}
func (c cmdable) HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	args := []interface{}{"hscan", key, cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	args = append(args, "novalues")
	cmd := NewScanCmd(ctx, c, args...)
	if hashtag.Present(match) {
		cmd.SetFirstKeyPos(4)
	}
	_ = c(ctx, cmd)
	return cmd
}

type HExpireArgs struct {
	NX bool
	XX bool
	GT bool
	LT bool
}







func (c cmdable) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *IntSliceCmd {
	args := []interface{}{"HEXPIRE", key, formatSec(ctx, expiration), "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}








func (c cmdable) HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd {
	args := []interface{}{"HEXPIRE", key, formatSec(ctx, expiration)}

	
	
	if expirationArgs.NX {
		args = append(args, "NX")
	} else if expirationArgs.XX {
		args = append(args, "XX")
	} else if expirationArgs.GT {
		args = append(args, "GT")
	} else if expirationArgs.LT {
		args = append(args, "LT")
	}

	args = append(args, "FIELDS", len(fields))

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}








func (c cmdable) HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPEXPIRE", key, formatMs(ctx, expiration), "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}








func (c cmdable) HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPEXPIRE", key, formatMs(ctx, expiration)}

	
	
	if expirationArgs.NX {
		args = append(args, "NX")
	} else if expirationArgs.XX {
		args = append(args, "XX")
	} else if expirationArgs.GT {
		args = append(args, "GT")
	} else if expirationArgs.LT {
		args = append(args, "LT")
	}

	args = append(args, "FIELDS", len(fields))

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}








func (c cmdable) HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *IntSliceCmd {

	args := []interface{}{"HEXPIREAT", key, tm.Unix(), "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd {
	args := []interface{}{"HEXPIREAT", key, tm.Unix()}

	
	
	if expirationArgs.NX {
		args = append(args, "NX")
	} else if expirationArgs.XX {
		args = append(args, "XX")
	} else if expirationArgs.GT {
		args = append(args, "GT")
	} else if expirationArgs.LT {
		args = append(args, "LT")
	}

	args = append(args, "FIELDS", len(fields))

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}







func (c cmdable) HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPEXPIREAT", key, tm.UnixNano() / int64(time.Millisecond), "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs HExpireArgs, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPEXPIREAT", key, tm.UnixNano() / int64(time.Millisecond)}

	
	
	if expirationArgs.NX {
		args = append(args, "NX")
	} else if expirationArgs.XX {
		args = append(args, "XX")
	} else if expirationArgs.GT {
		args = append(args, "GT")
	} else if expirationArgs.LT {
		args = append(args, "LT")
	}

	args = append(args, "FIELDS", len(fields))

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}








func (c cmdable) HPersist(ctx context.Context, key string, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPERSIST", key, "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}









func (c cmdable) HExpireTime(ctx context.Context, key string, fields ...string) *IntSliceCmd {
	args := []interface{}{"HEXPIRETIME", key, "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}









func (c cmdable) HPExpireTime(ctx context.Context, key string, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPEXPIRETIME", key, "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}








func (c cmdable) HTTL(ctx context.Context, key string, fields ...string) *IntSliceCmd {
	args := []interface{}{"HTTL", key, "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}









func (c cmdable) HPTTL(ctx context.Context, key string, fields ...string) *IntSliceCmd {
	args := []interface{}{"HPTTL", key, "FIELDS", len(fields)}

	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HGetDel(ctx context.Context, key string, fields ...string) *StringSliceCmd {
	args := []interface{}{"HGETDEL", key, "FIELDS", len(fields)}
	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HGetEX(ctx context.Context, key string, fields ...string) *StringSliceCmd {
	args := []interface{}{"HGETEX", key, "FIELDS", len(fields)}
	for _, field := range fields {
		args = append(args, field)
	}
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}


type HGetEXExpirationType string

const (
	HGetEXExpirationEX      HGetEXExpirationType = "EX"
	HGetEXExpirationPX      HGetEXExpirationType = "PX"
	HGetEXExpirationEXAT    HGetEXExpirationType = "EXAT"
	HGetEXExpirationPXAT    HGetEXExpirationType = "PXAT"
	HGetEXExpirationPERSIST HGetEXExpirationType = "PERSIST"
)

type HGetEXOptions struct {
	ExpirationType HGetEXExpirationType
	ExpirationVal  int64
}

func (c cmdable) HGetEXWithArgs(ctx context.Context, key string, options *HGetEXOptions, fields ...string) *StringSliceCmd {
	args := []interface{}{"HGETEX", key}
	if options.ExpirationType != "" {
		args = append(args, string(options.ExpirationType))
		if options.ExpirationType != HGetEXExpirationPERSIST {
			args = append(args, options.ExpirationVal)
		}
	}

	args = append(args, "FIELDS", len(fields))
	for _, field := range fields {
		args = append(args, field)
	}

	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

type HSetEXCondition string

const (
	HSetEXFNX HSetEXCondition = "FNX" 
	HSetEXFXX HSetEXCondition = "FXX" 
)

type HSetEXExpirationType string

const (
	HSetEXExpirationEX      HSetEXExpirationType = "EX"
	HSetEXExpirationPX      HSetEXExpirationType = "PX"
	HSetEXExpirationEXAT    HSetEXExpirationType = "EXAT"
	HSetEXExpirationPXAT    HSetEXExpirationType = "PXAT"
	HSetEXExpirationKEEPTTL HSetEXExpirationType = "KEEPTTL"
)

type HSetEXOptions struct {
	Condition      HSetEXCondition
	ExpirationType HSetEXExpirationType
	ExpirationVal  int64
}

func (c cmdable) HSetEX(ctx context.Context, key string, fieldsAndValues ...string) *IntCmd {
	args := []interface{}{"HSETEX", key, "FIELDS", len(fieldsAndValues) / 2}
	for _, field := range fieldsAndValues {
		args = append(args, field)
	}

	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HSetEXWithArgs(ctx context.Context, key string, options *HSetEXOptions, fieldsAndValues ...string) *IntCmd {
	args := []interface{}{"HSETEX", key}
	if options.Condition != "" {
		args = append(args, string(options.Condition))
	}
	if options.ExpirationType != "" {
		args = append(args, string(options.ExpirationType))
		if options.ExpirationType != HSetEXExpirationKEEPTTL {
			args = append(args, options.ExpirationVal)
		}
	}
	args = append(args, "FIELDS", len(fieldsAndValues)/2)
	for _, field := range fieldsAndValues {
		args = append(args, field)
	}

	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}
