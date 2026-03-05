




package sse

import (
    `github.com/bytedance/sonic/loader`
)

func Use() {
    loader.WrapGoC(_text_b64encode, _cfunc_b64encode, []loader.GoC{{"_b64encode", &S_b64encode, &F_b64encode}}, "sse", "sse/b64encode.c")
    loader.WrapGoC(_text_b64decode, _cfunc_b64decode, []loader.GoC{{"_b64decode", &S_b64decode, &F_b64decode}}, "sse", "sse/b64decode.c")
}