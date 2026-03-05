




package metadata 

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/internal"
)

func init() {
	internal.FromOutgoingContextRaw = fromOutgoingContextRaw
}




func DecodeKeyValue(k, v string) (string, string, error) {
	return k, v, nil
}



type MD map[string][]string













func New(m map[string]string) MD {
	md := make(MD, len(m))
	for k, val := range m {
		key := strings.ToLower(k)
		md[key] = append(md[key], val)
	}
	return md
}














func Pairs(kv ...string) MD {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := make(MD, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		md[key] = append(md[key], kv[i+1])
	}
	return md
}


func (md MD) Len() int {
	return len(md)
}


func (md MD) Copy() MD {
	out := make(MD, len(md))
	for k, v := range md {
		out[k] = copyOf(v)
	}
	return out
}




func (md MD) Get(k string) []string {
	k = strings.ToLower(k)
	return md[k]
}




func (md MD) Set(k string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	k = strings.ToLower(k)
	md[k] = vals
}





func (md MD) Append(k string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	k = strings.ToLower(k)
	md[k] = append(md[k], vals...)
}



func (md MD) Delete(k string) {
	k = strings.ToLower(k)
	delete(md, k)
}





func Join(mds ...MD) MD {
	out := MD{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return out
}

type mdIncomingKey struct{}
type mdOutgoingKey struct{}



func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}





func NewOutgoingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, rawMD{md: md})
}




func AppendToOutgoingContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := ctx.Value(mdOutgoingKey{}).(rawMD)
	added := make([][]string, len(md.added)+1)
	copy(added, md.added)
	kvCopy := make([]string, 0, len(kv))
	for i := 0; i < len(kv); i += 2 {
		kvCopy = append(kvCopy, strings.ToLower(kv[i]), kv[i+1])
	}
	added[len(added)-1] = kvCopy
	return context.WithValue(ctx, mdOutgoingKey{}, rawMD{md: md.md, added: added})
}




func FromIncomingContext(ctx context.Context) (MD, bool) {
	md, ok := ctx.Value(mdIncomingKey{}).(MD)
	if !ok {
		return nil, false
	}
	out := make(MD, len(md))
	for k, v := range md {
		
		
		
		key := strings.ToLower(k)
		out[key] = copyOf(v)
	}
	return out, true
}




func ValueFromIncomingContext(ctx context.Context, key string) []string {
	md, ok := ctx.Value(mdIncomingKey{}).(MD)
	if !ok {
		return nil
	}

	if v, ok := md[key]; ok {
		return copyOf(v)
	}
	for k, v := range md {
		
		
		
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}

func copyOf(v []string) []string {
	vals := make([]string, len(v))
	copy(vals, v)
	return vals
}







func fromOutgoingContextRaw(ctx context.Context) (MD, [][]string, bool) {
	raw, ok := ctx.Value(mdOutgoingKey{}).(rawMD)
	if !ok {
		return nil, nil, false
	}

	return raw.md, raw.added, true
}




func FromOutgoingContext(ctx context.Context) (MD, bool) {
	raw, ok := ctx.Value(mdOutgoingKey{}).(rawMD)
	if !ok {
		return nil, false
	}

	mdSize := len(raw.md)
	for i := range raw.added {
		mdSize += len(raw.added[i]) / 2
	}

	out := make(MD, mdSize)
	for k, v := range raw.md {
		
		
		
		key := strings.ToLower(k)
		out[key] = copyOf(v)
	}
	for _, added := range raw.added {
		if len(added)%2 == 1 {
			panic(fmt.Sprintf("metadata: FromOutgoingContext got an odd number of input pairs for metadata: %d", len(added)))
		}

		for i := 0; i < len(added); i += 2 {
			key := strings.ToLower(added[i])
			out[key] = append(out[key], added[i+1])
		}
	}
	return out, ok
}

type rawMD struct {
	md    MD
	added [][]string
}
