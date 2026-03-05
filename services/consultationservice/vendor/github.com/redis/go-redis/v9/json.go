package redis

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/redis/go-redis/v9/internal/proto"
	"github.com/redis/go-redis/v9/internal/util"
)



type JSONCmdable interface {
	JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *IntSliceCmd
	JSONArrIndex(ctx context.Context, key, path string, value ...interface{}) *IntSliceCmd
	JSONArrIndexWithArgs(ctx context.Context, key, path string, options *JSONArrIndexArgs, value ...interface{}) *IntSliceCmd
	JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *IntSliceCmd
	JSONArrLen(ctx context.Context, key, path string) *IntSliceCmd
	JSONArrPop(ctx context.Context, key, path string, index int) *StringSliceCmd
	JSONArrTrim(ctx context.Context, key, path string) *IntSliceCmd
	JSONArrTrimWithArgs(ctx context.Context, key, path string, options *JSONArrTrimArgs) *IntSliceCmd
	JSONClear(ctx context.Context, key, path string) *IntCmd
	JSONDebugMemory(ctx context.Context, key, path string) *IntCmd
	JSONDel(ctx context.Context, key, path string) *IntCmd
	JSONForget(ctx context.Context, key, path string) *IntCmd
	JSONGet(ctx context.Context, key string, paths ...string) *JSONCmd
	JSONGetWithArgs(ctx context.Context, key string, options *JSONGetArgs, paths ...string) *JSONCmd
	JSONMerge(ctx context.Context, key, path string, value string) *StatusCmd
	JSONMSetArgs(ctx context.Context, docs []JSONSetArgs) *StatusCmd
	JSONMSet(ctx context.Context, params ...interface{}) *StatusCmd
	JSONMGet(ctx context.Context, path string, keys ...string) *JSONSliceCmd
	JSONNumIncrBy(ctx context.Context, key, path string, value float64) *JSONCmd
	JSONObjKeys(ctx context.Context, key, path string) *SliceCmd
	JSONObjLen(ctx context.Context, key, path string) *IntPointerSliceCmd
	JSONSet(ctx context.Context, key, path string, value interface{}) *StatusCmd
	JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *StatusCmd
	JSONStrAppend(ctx context.Context, key, path, value string) *IntPointerSliceCmd
	JSONStrLen(ctx context.Context, key, path string) *IntPointerSliceCmd
	JSONToggle(ctx context.Context, key, path string) *IntPointerSliceCmd
	JSONType(ctx context.Context, key, path string) *JSONSliceCmd
}

type JSONSetArgs struct {
	Key   string
	Path  string
	Value interface{}
}

type JSONArrIndexArgs struct {
	Start int
	Stop  *int
}

type JSONArrTrimArgs struct {
	Start int
	Stop  *int
}

type JSONCmd struct {
	baseCmd
	val      string
	expanded interface{}
}

var _ Cmder = (*JSONCmd)(nil)

func newJSONCmd(ctx context.Context, args ...interface{}) *JSONCmd {
	return &JSONCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *JSONCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *JSONCmd) SetVal(val string) {
	cmd.val = val
}

func (cmd *JSONCmd) Val() string {
	if len(cmd.val) == 0 && cmd.expanded != nil {
		val, err := json.Marshal(cmd.expanded)
		if err != nil {
			cmd.SetErr(err)
			return ""
		}
		return string(val)

	} else {
		return cmd.val
	}
}

func (cmd *JSONCmd) Result() (string, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *JSONCmd) Expanded() (interface{}, error) {
	if len(cmd.val) != 0 && cmd.expanded == nil {
		err := json.Unmarshal([]byte(cmd.val), &cmd.expanded)
		if err != nil {
			return nil, err
		}
	}

	return cmd.expanded, nil
}

func (cmd *JSONCmd) readReply(rd *proto.Reader) error {
	
	if cmd.baseCmd.Err() == Nil {
		cmd.val = ""
		return Nil
	}

	if readType, err := rd.PeekReplyType(); err != nil {
		return err
	} else if readType == proto.RespArray {

		size, err := rd.ReadArrayLen()
		if err != nil {
			return err
		}

		expanded := make([]interface{}, size)

		for i := 0; i < size; i++ {
			if expanded[i], err = rd.ReadReply(); err != nil {
				return err
			}
		}
		cmd.expanded = expanded

	} else {
		if str, err := rd.ReadString(); err != nil && err != Nil {
			return err
		} else if str == "" || err == Nil {
			cmd.val = ""
		} else {
			cmd.val = str
		}
	}

	return nil
}



type JSONSliceCmd struct {
	baseCmd
	val []interface{}
}

func NewJSONSliceCmd(ctx context.Context, args ...interface{}) *JSONSliceCmd {
	return &JSONSliceCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *JSONSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *JSONSliceCmd) SetVal(val []interface{}) {
	cmd.val = val
}

func (cmd *JSONSliceCmd) Val() []interface{} {
	return cmd.val
}

func (cmd *JSONSliceCmd) Result() ([]interface{}, error) {
	return cmd.val, cmd.err
}

func (cmd *JSONSliceCmd) readReply(rd *proto.Reader) error {
	if cmd.baseCmd.Err() == Nil {
		cmd.val = nil
		return Nil
	}

	if readType, err := rd.PeekReplyType(); err != nil {
		return err
	} else if readType == proto.RespArray {
		response, err := rd.ReadReply()
		if err != nil {
			return nil
		} else {
			cmd.val = response.([]interface{})
		}

	} else {
		n, err := rd.ReadArrayLen()
		if err != nil {
			return err
		}
		cmd.val = make([]interface{}, n)
		for i := 0; i < len(cmd.val); i++ {
			switch s, err := rd.ReadString(); {
			case err == Nil:
				cmd.val[i] = ""
			case err != nil:
				return err
			default:
				cmd.val[i] = s
			}
		}
	}
	return nil
}



type IntPointerSliceCmd struct {
	baseCmd
	val []*int64
}


func NewIntPointerSliceCmd(ctx context.Context, args ...interface{}) *IntPointerSliceCmd {
	return &IntPointerSliceCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *IntPointerSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *IntPointerSliceCmd) SetVal(val []*int64) {
	cmd.val = val
}

func (cmd *IntPointerSliceCmd) Val() []*int64 {
	return cmd.val
}

func (cmd *IntPointerSliceCmd) Result() ([]*int64, error) {
	return cmd.val, cmd.err
}

func (cmd *IntPointerSliceCmd) readReply(rd *proto.Reader) error {
	n, err := rd.ReadArrayLen()
	if err != nil {
		return err
	}
	cmd.val = make([]*int64, n)

	for i := 0; i < len(cmd.val); i++ {
		val, err := rd.ReadInt()
		if err != nil && err != Nil {
			return err
		} else if err != Nil {
			cmd.val[i] = &val
		}
	}

	return nil
}





func (c cmdable) JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *IntSliceCmd {
	args := []interface{}{"JSON.ARRAPPEND", key, path}
	args = append(args, values...)
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONArrIndex(ctx context.Context, key, path string, value ...interface{}) *IntSliceCmd {
	args := []interface{}{"JSON.ARRINDEX", key, path}
	args = append(args, value...)
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}




func (c cmdable) JSONArrIndexWithArgs(ctx context.Context, key, path string, options *JSONArrIndexArgs, value ...interface{}) *IntSliceCmd {
	args := []interface{}{"JSON.ARRINDEX", key, path}
	args = append(args, value...)

	if options != nil {
		args = append(args, options.Start)
		if options.Stop != nil {
			args = append(args, *options.Stop)
		}
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *IntSliceCmd {
	args := []interface{}{"JSON.ARRINSERT", key, path, index}
	args = append(args, values...)
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONArrLen(ctx context.Context, key, path string) *IntSliceCmd {
	args := []interface{}{"JSON.ARRLEN", key, path}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONArrPop(ctx context.Context, key, path string, index int) *StringSliceCmd {
	args := []interface{}{"JSON.ARRPOP", key, path, index}
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONArrTrim(ctx context.Context, key, path string) *IntSliceCmd {
	args := []interface{}{"JSON.ARRTRIM", key, path}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONArrTrimWithArgs(ctx context.Context, key, path string, options *JSONArrTrimArgs) *IntSliceCmd {
	args := []interface{}{"JSON.ARRTRIM", key, path}

	if options != nil {
		args = append(args, options.Start)

		if options.Stop != nil {
			args = append(args, *options.Stop)
		}
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONClear(ctx context.Context, key, path string) *IntCmd {
	args := []interface{}{"JSON.CLEAR", key, path}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONDebugMemory(ctx context.Context, key, path string) *IntCmd {
	panic("not implemented")
}



func (c cmdable) JSONDel(ctx context.Context, key, path string) *IntCmd {
	args := []interface{}{"JSON.DEL", key, path}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONForget(ctx context.Context, key, path string) *IntCmd {
	args := []interface{}{"JSON.FORGET", key, path}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}





func (c cmdable) JSONGet(ctx context.Context, key string, paths ...string) *JSONCmd {
	args := make([]interface{}, len(paths)+2)
	args[0] = "JSON.GET"
	args[1] = key
	for n, path := range paths {
		args[n+2] = path
	}
	cmd := newJSONCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

type JSONGetArgs struct {
	Indent  string
	Newline string
	Space   string
}





func (c cmdable) JSONGetWithArgs(ctx context.Context, key string, options *JSONGetArgs, paths ...string) *JSONCmd {
	args := []interface{}{"JSON.GET", key}
	if options != nil {
		if options.Indent != "" {
			args = append(args, "INDENT", options.Indent)
		}
		if options.Newline != "" {
			args = append(args, "NEWLINE", options.Newline)
		}
		if options.Space != "" {
			args = append(args, "SPACE", options.Space)
		}
		for _, path := range paths {
			args = append(args, path)
		}
	}
	cmd := newJSONCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONMerge(ctx context.Context, key, path string, value string) *StatusCmd {
	args := []interface{}{"JSON.MERGE", key, path, value}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}





func (c cmdable) JSONMGet(ctx context.Context, path string, keys ...string) *JSONSliceCmd {
	args := make([]interface{}, len(keys)+1)
	args[0] = "JSON.MGET"
	for n, key := range keys {
		args[n+1] = key
	}
	args = append(args, path)
	cmd := NewJSONSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONMSetArgs(ctx context.Context, docs []JSONSetArgs) *StatusCmd {
	args := []interface{}{"JSON.MSET"}
	for _, doc := range docs {
		args = append(args, doc.Key, doc.Path, doc.Value)
	}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) JSONMSet(ctx context.Context, params ...interface{}) *StatusCmd {
	args := []interface{}{"JSON.MSET"}
	args = append(args, params...)
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONNumIncrBy(ctx context.Context, key, path string, value float64) *JSONCmd {
	args := []interface{}{"JSON.NUMINCRBY", key, path, value}
	cmd := newJSONCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONObjKeys(ctx context.Context, key, path string) *SliceCmd {
	args := []interface{}{"JSON.OBJKEYS", key, path}
	cmd := NewSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONObjLen(ctx context.Context, key, path string) *IntPointerSliceCmd {
	args := []interface{}{"JSON.OBJLEN", key, path}
	cmd := NewIntPointerSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}





func (c cmdable) JSONSet(ctx context.Context, key, path string, value interface{}) *StatusCmd {
	return c.JSONSetMode(ctx, key, path, value, "")
}





func (c cmdable) JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *StatusCmd {
	var bytes []byte
	var err error
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		bytes, err = json.Marshal(v)
	}
	args := []interface{}{"JSON.SET", key, path, util.BytesToString(bytes)}
	if mode != "" {
		switch strings.ToUpper(mode) {
		case "XX", "NX":
			args = append(args, strings.ToUpper(mode))

		default:
			panic("redis: JSON.SET mode must be NX or XX")
		}
	}
	cmd := NewStatusCmd(ctx, args...)
	if err != nil {
		cmd.SetErr(err)
	} else {
		_ = c(ctx, cmd)
	}
	return cmd
}



func (c cmdable) JSONStrAppend(ctx context.Context, key, path, value string) *IntPointerSliceCmd {
	args := []interface{}{"JSON.STRAPPEND", key, path, value}
	cmd := NewIntPointerSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONStrLen(ctx context.Context, key, path string) *IntPointerSliceCmd {
	args := []interface{}{"JSON.STRLEN", key, path}
	cmd := NewIntPointerSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONToggle(ctx context.Context, key, path string) *IntPointerSliceCmd {
	args := []interface{}{"JSON.TOGGLE", key, path}
	cmd := NewIntPointerSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) JSONType(ctx context.Context, key, path string) *JSONSliceCmd {
	args := []interface{}{"JSON.TYPE", key, path}
	cmd := NewJSONSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}
