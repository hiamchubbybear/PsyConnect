package redis

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9/internal/hashtag"
)

type SortedSetCmdable interface {
	BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd
	BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd
	BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZSliceWithKeyCmd
	ZAdd(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddLT(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddGT(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddNX(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddXX(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddArgs(ctx context.Context, key string, args ZAddArgs) *IntCmd
	ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) *FloatCmd
	ZCard(ctx context.Context, key string) *IntCmd
	ZCount(ctx context.Context, key, min, max string) *IntCmd
	ZLexCount(ctx context.Context, key, min, max string) *IntCmd
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd
	ZInter(ctx context.Context, store *ZStore) *StringSliceCmd
	ZInterWithScores(ctx context.Context, store *ZStore) *ZSliceCmd
	ZInterCard(ctx context.Context, limit int64, keys ...string) *IntCmd
	ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd
	ZMPop(ctx context.Context, order string, count int64, keys ...string) *ZSliceWithKeyCmd
	ZMScore(ctx context.Context, key string, members ...string) *FloatSliceCmd
	ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd
	ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd
	ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRangeArgs(ctx context.Context, z ZRangeArgs) *StringSliceCmd
	ZRangeArgsWithScores(ctx context.Context, z ZRangeArgs) *ZSliceCmd
	ZRangeStore(ctx context.Context, dst string, z ZRangeArgs) *IntCmd
	ZRank(ctx context.Context, key, member string) *IntCmd
	ZRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd
	ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd
	ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRevRank(ctx context.Context, key, member string) *IntCmd
	ZRevRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd
	ZScore(ctx context.Context, key, member string) *FloatCmd
	ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd
	ZRandMember(ctx context.Context, key string, count int) *StringSliceCmd
	ZRandMemberWithScores(ctx context.Context, key string, count int) *ZSliceCmd
	ZUnion(ctx context.Context, store ZStore) *StringSliceCmd
	ZUnionWithScores(ctx context.Context, store ZStore) *ZSliceCmd
	ZDiff(ctx context.Context, keys ...string) *StringSliceCmd
	ZDiffWithScores(ctx context.Context, keys ...string) *ZSliceCmd
	ZDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
}


func (c cmdable) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd {
	args := make([]interface{}, 1+len(keys)+1)
	args[0] = "bzpopmax"
	for i, key := range keys {
		args[1+i] = key
	}
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewZWithKeyCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd {
	args := make([]interface{}, 1+len(keys)+1)
	args[0] = "bzpopmin"
	for i, key := range keys {
		args[1+i] = key
	}
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewZWithKeyCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}






func (c cmdable) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZSliceWithKeyCmd {
	args := make([]interface{}, 3+len(keys), 6+len(keys))
	args[0] = "bzmpop"
	args[1] = formatSec(ctx, timeout)
	args[2] = len(keys)
	for i, key := range keys {
		args[3+i] = key
	}
	args = append(args, strings.ToLower(order), "count", count)
	cmd := NewZSliceWithKeyCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}


type ZAddArgs struct {
	NX      bool
	XX      bool
	LT      bool
	GT      bool
	Ch      bool
	Members []Z
}

func (c cmdable) zAddArgs(key string, args ZAddArgs, incr bool) []interface{} {
	a := make([]interface{}, 0, 6+2*len(args.Members))
	a = append(a, "zadd", key)

	
	if args.NX {
		a = append(a, "nx")
	} else {
		if args.XX {
			a = append(a, "xx")
		}
		if args.GT {
			a = append(a, "gt")
		} else if args.LT {
			a = append(a, "lt")
		}
	}
	if args.Ch {
		a = append(a, "ch")
	}
	if incr {
		a = append(a, "incr")
	}
	for _, m := range args.Members {
		a = append(a, m.Score)
		a = append(a, m.Member)
	}
	return a
}

func (c cmdable) ZAddArgs(ctx context.Context, key string, args ZAddArgs) *IntCmd {
	cmd := NewIntCmd(ctx, c.zAddArgs(key, args, false)...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) *FloatCmd {
	cmd := NewFloatCmd(ctx, c.zAddArgs(key, args, true)...)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) ZAdd(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		Members: members,
	})
}


func (c cmdable) ZAddLT(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		LT:      true,
		Members: members,
	})
}


func (c cmdable) ZAddGT(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		GT:      true,
		Members: members,
	})
}


func (c cmdable) ZAddNX(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		NX:      true,
		Members: members,
	})
}


func (c cmdable) ZAddXX(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		XX:      true,
		Members: members,
	})
}

func (c cmdable) ZCard(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(ctx, "zcard", key)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZCount(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd(ctx, "zcount", key, min, max)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZLexCount(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd(ctx, "zlexcount", key, min, max)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd {
	cmd := NewFloatCmd(ctx, "zincrby", key, increment, member)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd {
	args := make([]interface{}, 0, 3+store.len())
	args = append(args, "zinterstore", destination, len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewIntCmd(ctx, args...)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInter(ctx context.Context, store *ZStore) *StringSliceCmd {
	args := make([]interface{}, 0, 2+store.len())
	args = append(args, "zinter", len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInterWithScores(ctx context.Context, store *ZStore) *ZSliceCmd {
	args := make([]interface{}, 0, 3+store.len())
	args = append(args, "zinter", len(store.Keys))
	args = store.appendArgs(args)
	args = append(args, "withscores")
	cmd := NewZSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInterCard(ctx context.Context, limit int64, keys ...string) *IntCmd {
	numKeys := len(keys)
	args := make([]interface{}, 4+numKeys)
	args[0] = "zintercard"
	args[1] = numKeys
	for i, key := range keys {
		args[2+i] = key
	}
	args[2+numKeys] = "limit"
	args[3+numKeys] = limit
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}




func (c cmdable) ZMPop(ctx context.Context, order string, count int64, keys ...string) *ZSliceWithKeyCmd {
	args := make([]interface{}, 2+len(keys), 5+len(keys))
	args[0] = "zmpop"
	args[1] = len(keys)
	for i, key := range keys {
		args[2+i] = key
	}
	args = append(args, strings.ToLower(order), "count", count)
	cmd := NewZSliceWithKeyCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZMScore(ctx context.Context, key string, members ...string) *FloatSliceCmd {
	args := make([]interface{}, 2+len(members))
	args[0] = "zmscore"
	args[1] = key
	for i, member := range members {
		args[2+i] = member
	}
	cmd := NewFloatSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd {
	args := []interface{}{
		"zpopmax",
		key,
	}

	switch len(count) {
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	}

	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd {
	args := []interface{}{
		"zpopmin",
		key,
	}

	switch len(count) {
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	}

	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}













type ZRangeArgs struct {
	Key string

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Start interface{}
	Stop  interface{}

	
	ByScore bool
	ByLex   bool

	Rev bool

	
	Offset int64
	Count  int64
}

func (z ZRangeArgs) appendArgs(args []interface{}) []interface{} {
	
	if z.Rev && (z.ByScore || z.ByLex) {
		args = append(args, z.Key, z.Stop, z.Start)
	} else {
		args = append(args, z.Key, z.Start, z.Stop)
	}

	if z.ByScore {
		args = append(args, "byscore")
	} else if z.ByLex {
		args = append(args, "bylex")
	}
	if z.Rev {
		args = append(args, "rev")
	}
	if z.Offset != 0 || z.Count != 0 {
		args = append(args, "limit", z.Offset, z.Count)
	}
	return args
}

func (c cmdable) ZRangeArgs(ctx context.Context, z ZRangeArgs) *StringSliceCmd {
	args := make([]interface{}, 0, 9)
	args = append(args, "zrange")
	args = z.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRangeArgsWithScores(ctx context.Context, z ZRangeArgs) *ZSliceCmd {
	args := make([]interface{}, 0, 10)
	args = append(args, "zrange")
	args = z.appendArgs(args)
	args = append(args, "withscores")
	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	return c.ZRangeArgs(ctx, ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  stop,
	})
}

func (c cmdable) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	return c.ZRangeArgsWithScores(ctx, ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  stop,
	})
}

type ZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

func (c cmdable) zRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy, withScores bool) *StringSliceCmd {
	args := []interface{}{zcmd, key, opt.Min, opt.Max}
	if withScores {
		args = append(args, "withscores")
	}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRangeBy(ctx, "zrangebyscore", key, opt, false)
}

func (c cmdable) ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRangeBy(ctx, "zrangebylex", key, opt, false)
}

func (c cmdable) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	args := []interface{}{"zrangebyscore", key, opt.Min, opt.Max, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRangeStore(ctx context.Context, dst string, z ZRangeArgs) *IntCmd {
	args := make([]interface{}, 0, 10)
	args = append(args, "zrangestore", dst)
	args = z.appendArgs(args)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRank(ctx context.Context, key, member string) *IntCmd {
	cmd := NewIntCmd(ctx, "zrank", key, member)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) ZRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd {
	cmd := NewRankWithScoreCmd(ctx, "zrank", key, member, "withscore")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	args := make([]interface{}, 2, 2+len(members))
	args[0] = "zrem"
	args[1] = key
	args = appendArgs(args, members)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd {
	cmd := NewIntCmd(
		ctx,
		"zremrangebyrank",
		key,
		start,
		stop,
	)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd(ctx, "zremrangebyscore", key, min, max)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd(ctx, "zremrangebylex", key, min, max)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "zrevrange", key, start, stop)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	cmd := NewZSliceCmd(ctx, "zrevrange", key, start, stop, "withscores")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) zRevRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy) *StringSliceCmd {
	args := []interface{}{zcmd, key, opt.Max, opt.Min}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRevRangeBy(ctx, "zrevrangebyscore", key, opt)
}

func (c cmdable) ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRevRangeBy(ctx, "zrevrangebylex", key, opt)
}

func (c cmdable) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	args := []interface{}{"zrevrangebyscore", key, opt.Max, opt.Min, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRank(ctx context.Context, key, member string) *IntCmd {
	cmd := NewIntCmd(ctx, "zrevrank", key, member)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd {
	cmd := NewRankWithScoreCmd(ctx, "zrevrank", key, member, "withscore")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZScore(ctx context.Context, key, member string) *FloatCmd {
	cmd := NewFloatCmd(ctx, "zscore", key, member)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZUnion(ctx context.Context, store ZStore) *StringSliceCmd {
	args := make([]interface{}, 0, 2+store.len())
	args = append(args, "zunion", len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZUnionWithScores(ctx context.Context, store ZStore) *ZSliceCmd {
	args := make([]interface{}, 0, 3+store.len())
	args = append(args, "zunion", len(store.Keys))
	args = store.appendArgs(args)
	args = append(args, "withscores")
	cmd := NewZSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd {
	args := make([]interface{}, 0, 3+store.len())
	args = append(args, "zunionstore", dest, len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewIntCmd(ctx, args...)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) ZRandMember(ctx context.Context, key string, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "zrandmember", key, count)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) ZRandMemberWithScores(ctx context.Context, key string, count int) *ZSliceCmd {
	cmd := NewZSliceCmd(ctx, "zrandmember", key, count, "withscores")
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) ZDiff(ctx context.Context, keys ...string) *StringSliceCmd {
	args := make([]interface{}, 2+len(keys))
	args[0] = "zdiff"
	args[1] = len(keys)
	for i, key := range keys {
		args[i+2] = key
	}

	cmd := NewStringSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) ZDiffWithScores(ctx context.Context, keys ...string) *ZSliceCmd {
	args := make([]interface{}, 3+len(keys))
	args[0] = "zdiff"
	args[1] = len(keys)
	for i, key := range keys {
		args[i+2] = key
	}
	args[len(keys)+2] = "withscores"

	cmd := NewZSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}


func (c cmdable) ZDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	args := make([]interface{}, 0, 3+len(keys))
	args = append(args, "zdiffstore", destination, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	args := []interface{}{"zscan", key, cursor}
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


type Z struct {
	Score  float64
	Member interface{}
}


type ZWithKey struct {
	Z
	Key string
}


type ZStore struct {
	Keys    []string
	Weights []float64
	
	Aggregate string
}

func (z ZStore) len() (n int) {
	n = len(z.Keys)
	if len(z.Weights) > 0 {
		n += 1 + len(z.Weights)
	}
	if z.Aggregate != "" {
		n += 2
	}
	return n
}

func (z ZStore) appendArgs(args []interface{}) []interface{} {
	for _, key := range z.Keys {
		args = append(args, key)
	}
	if len(z.Weights) > 0 {
		args = append(args, "weights")
		for _, weights := range z.Weights {
			args = append(args, weights)
		}
	}
	if z.Aggregate != "" {
		args = append(args, "aggregate", z.Aggregate)
	}
	return args
}
