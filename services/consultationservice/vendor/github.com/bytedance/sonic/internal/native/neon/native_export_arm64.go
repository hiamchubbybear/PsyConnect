



package neon

var (
    S_f64toa = _subr__f64toa
    S_f32toa = _subr__f32toa
    S_i64toa = _subr__i64toa
    S_u64toa = _subr__u64toa
    S_lspace = _subr__lspace
)

var (
    S_quote   = _subr__quote
    S_unquote = _subr__unquote
)

var (
    S_value     = _subr__value
    S_vstring   = _subr__vstring
    S_vnumber   = _subr__vnumber
    S_vsigned   = _subr__vsigned
    S_vunsigned = _subr__vunsigned
)

var (
    S_skip_one    = _subr__skip_one
    S_skip_one_fast = _subr__skip_one_fast
    S_skip_array  = _subr__skip_array
    S_skip_object = _subr__skip_object
    S_skip_number = _subr__skip_number
    S_get_by_path = _subr__get_by_path
	S_lookup_small_key = _subr__lookup_small_key
	S_parse_with_padding = _subr__parse_with_padding
)
