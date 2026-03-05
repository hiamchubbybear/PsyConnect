package redis

import (
	"context"
	"encoding/json"
	"strconv"
)


type VectorSetCmdable interface {
	VAdd(ctx context.Context, key, element string, val Vector) *BoolCmd
	VAddWithArgs(ctx context.Context, key, element string, val Vector, addArgs *VAddArgs) *BoolCmd
	VCard(ctx context.Context, key string) *IntCmd
	VDim(ctx context.Context, key string) *IntCmd
	VEmb(ctx context.Context, key, element string, raw bool) *SliceCmd
	VGetAttr(ctx context.Context, key, element string) *StringCmd
	VInfo(ctx context.Context, key string) *MapStringInterfaceCmd
	VLinks(ctx context.Context, key, element string) *StringSliceCmd
	VLinksWithScores(ctx context.Context, key, element string) *VectorScoreSliceCmd
	VRandMember(ctx context.Context, key string) *StringCmd
	VRandMemberCount(ctx context.Context, key string, count int) *StringSliceCmd
	VRem(ctx context.Context, key, element string) *BoolCmd
	VSetAttr(ctx context.Context, key, element string, attr interface{}) *BoolCmd
	VClearAttributes(ctx context.Context, key, element string) *BoolCmd
	VSim(ctx context.Context, key string, val Vector) *StringSliceCmd
	VSimWithScores(ctx context.Context, key string, val Vector) *VectorScoreSliceCmd
	VSimWithArgs(ctx context.Context, key string, val Vector, args *VSimArgs) *StringSliceCmd
	VSimWithArgsWithScores(ctx context.Context, key string, val Vector, args *VSimArgs) *VectorScoreSliceCmd
}

type Vector interface {
	Value() []any
}

const (
	vectorFormatFP32   string = "FP32"
	vectorFormatValues string = "Values"
)

type VectorFP32 struct {
	Val []byte
}

func (v *VectorFP32) Value() []any {
	return []any{vectorFormatFP32, v.Val}
}

var _ Vector = (*VectorFP32)(nil)

type VectorValues struct {
	Val []float64
}

func (v *VectorValues) Value() []any {
	res := make([]any, 2+len(v.Val))
	res[0] = vectorFormatValues
	res[1] = len(v.Val)
	for i, v := range v.Val {
		res[2+i] = v
	}
	return res
}

var _ Vector = (*VectorValues)(nil)

type VectorRef struct {
	Name string 
}

func (v *VectorRef) Value() []any {
	return []any{"ele", v.Name}
}

var _ Vector = (*VectorRef)(nil)

type VectorScore struct {
	Name  string
	Score float64
}



func (c cmdable) VAdd(ctx context.Context, key, element string, val Vector) *BoolCmd {
	return c.VAddWithArgs(ctx, key, element, val, &VAddArgs{})
}

type VAddArgs struct {
	
	Reduce int64
	Cas    bool

	
	NoQuant bool
	Q8      bool
	Bin     bool

	EF      int64
	SetAttr string
	M       int64
}

func (v VAddArgs) reduce() int64 {
	return v.Reduce
}

func (v VAddArgs) appendArgs(args []any) []any {
	if v.Cas {
		args = append(args, "cas")
	}

	if v.NoQuant {
		args = append(args, "noquant")
	} else if v.Q8 {
		args = append(args, "q8")
	} else if v.Bin {
		args = append(args, "bin")
	}

	if v.EF > 0 {
		args = append(args, "ef", strconv.FormatInt(v.EF, 10))
	}
	if len(v.SetAttr) > 0 {
		args = append(args, "setattr", v.SetAttr)
	}
	if v.M > 0 {
		args = append(args, "m", strconv.FormatInt(v.M, 10))
	}
	return args
}



func (c cmdable) VAddWithArgs(ctx context.Context, key, element string, val Vector, addArgs *VAddArgs) *BoolCmd {
	if addArgs == nil {
		addArgs = &VAddArgs{}
	}
	args := []any{"vadd", key}
	if addArgs.reduce() > 0 {
		args = append(args, "reduce", addArgs.reduce())
	}
	args = append(args, val.Value()...)
	args = append(args, element)
	args = addArgs.appendArgs(args)
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VCard(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(ctx, "vcard", key)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VDim(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(ctx, "vdim", key)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VEmb(ctx context.Context, key, element string, raw bool) *SliceCmd {
	args := []any{"vemb", key, element}
	if raw {
		args = append(args, "raw")
	}
	cmd := NewSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VGetAttr(ctx context.Context, key, element string) *StringCmd {
	cmd := NewStringCmd(ctx, "vgetattr", key, element)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VInfo(ctx context.Context, key string) *MapStringInterfaceCmd {
	cmd := NewMapStringInterfaceCmd(ctx, "vinfo", key)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VLinks(ctx context.Context, key, element string) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "vlinks", key, element)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VLinksWithScores(ctx context.Context, key, element string) *VectorScoreSliceCmd {
	cmd := NewVectorInfoSliceCmd(ctx, "vlinks", key, element, "withscores")
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VRandMember(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(ctx, "vrandmember", key)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VRandMemberCount(ctx context.Context, key string, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd(ctx, "vrandmember", key, count)
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VRem(ctx context.Context, key, element string) *BoolCmd {
	cmd := NewBoolCmd(ctx, "vrem", key, element)
	_ = c(ctx, cmd)
	return cmd
}






func (c cmdable) VSetAttr(ctx context.Context, key, element string, attr interface{}) *BoolCmd {
	var attrStr string
	var err error
	switch v := attr.(type) {
	case string:
		attrStr = v
	case []byte:
		attrStr = string(v)
	default:
		var bytes []byte
		bytes, err = json.Marshal(v)
		if err != nil {
			
			cmd := NewBoolCmd(ctx, "vsetattr", key, element, "")
			cmd.SetErr(err)
			return cmd
		}
		attrStr = string(bytes)
	}
	cmd := NewBoolCmd(ctx, "vsetattr", key, element, attrStr)
	_ = c(ctx, cmd)
	return cmd
}




func (c cmdable) VClearAttributes(ctx context.Context, key, element string) *BoolCmd {
	cmd := NewBoolCmd(ctx, "vsetattr", key, element, "")
	_ = c(ctx, cmd)
	return cmd
}



func (c cmdable) VSim(ctx context.Context, key string, val Vector) *StringSliceCmd {
	return c.VSimWithArgs(ctx, key, val, &VSimArgs{})
}



func (c cmdable) VSimWithScores(ctx context.Context, key string, val Vector) *VectorScoreSliceCmd {
	return c.VSimWithArgsWithScores(ctx, key, val, &VSimArgs{})
}

type VSimArgs struct {
	Count    int64
	EF       int64
	Filter   string
	FilterEF int64
	Truth    bool
	NoThread bool
	Epsilon  float64
}

func (v VSimArgs) appendArgs(args []any) []any {
	if v.Count > 0 {
		args = append(args, "count", v.Count)
	}
	if v.EF > 0 {
		args = append(args, "ef", v.EF)
	}
	if len(v.Filter) > 0 {
		args = append(args, "filter", v.Filter)
	}
	if v.FilterEF > 0 {
		args = append(args, "filter-ef", v.FilterEF)
	}
	if v.Truth {
		args = append(args, "truth")
	}
	if v.NoThread {
		args = append(args, "nothread")
	}
	if v.Epsilon > 0 {
		args = append(args, "Epsilon", v.Epsilon)
	}
	return args
}




func (c cmdable) VSimWithArgs(ctx context.Context, key string, val Vector, simArgs *VSimArgs) *StringSliceCmd {
	if simArgs == nil {
		simArgs = &VSimArgs{}
	}
	args := []any{"vsim", key}
	args = append(args, val.Value()...)
	args = simArgs.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}




func (c cmdable) VSimWithArgsWithScores(ctx context.Context, key string, val Vector, simArgs *VSimArgs) *VectorScoreSliceCmd {
	if simArgs == nil {
		simArgs = &VSimArgs{}
	}
	args := []any{"vsim", key}
	args = append(args, val.Value()...)
	args = append(args, "withscores")
	args = simArgs.appendArgs(args)
	cmd := NewVectorInfoSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}
