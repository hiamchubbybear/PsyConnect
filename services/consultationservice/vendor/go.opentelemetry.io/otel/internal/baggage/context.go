


package baggage 

import "context"

type baggageContextKeyType int

const baggageKey baggageContextKeyType = iota


type SetHookFunc func(context.Context, List) context.Context


type GetHookFunc func(context.Context, List) List

type baggageState struct {
	list List

	setHook SetHookFunc
	getHook GetHookFunc
}





func ContextWithSetHook(parent context.Context, hook SetHookFunc) context.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.setHook = hook
	return context.WithValue(parent, baggageKey, s)
}





func ContextWithGetHook(parent context.Context, hook GetHookFunc) context.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.getHook = hook
	return context.WithValue(parent, baggageKey, s)
}



func ContextWithList(parent context.Context, list List) context.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.list = list
	ctx := context.WithValue(parent, baggageKey, s)
	if s.setHook != nil {
		ctx = s.setHook(ctx, list)
	}

	return ctx
}


func ListFromContext(ctx context.Context) List {
	switch v := ctx.Value(baggageKey).(type) {
	case baggageState:
		if v.getHook != nil {
			return v.getHook(ctx, v.list)
		}
		return v.list
	default:
		return nil
	}
}
