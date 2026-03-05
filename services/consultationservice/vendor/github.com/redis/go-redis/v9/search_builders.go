package redis

import (
	"context"
)








type SearchBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	query   string
	options *FTSearchOptions
}



func (c *Client) NewSearchBuilder(ctx context.Context, index, query string) *SearchBuilder {
	b := &SearchBuilder{c: c, ctx: ctx, index: index, query: query, options: &FTSearchOptions{LimitOffset: -1}}
	return b
}


func (b *SearchBuilder) WithScores() *SearchBuilder {
	b.options.WithScores = true
	return b
}


func (b *SearchBuilder) NoContent() *SearchBuilder { b.options.NoContent = true; return b }


func (b *SearchBuilder) Verbatim() *SearchBuilder { b.options.Verbatim = true; return b }


func (b *SearchBuilder) NoStopWords() *SearchBuilder { b.options.NoStopWords = true; return b }


func (b *SearchBuilder) WithPayloads() *SearchBuilder {
	b.options.WithPayloads = true
	return b
}


func (b *SearchBuilder) WithSortKeys() *SearchBuilder {
	b.options.WithSortKeys = true
	return b
}


func (b *SearchBuilder) Filter(field string, min, max interface{}) *SearchBuilder {
	b.options.Filters = append(b.options.Filters, FTSearchFilter{
		FieldName: field,
		Min:       min,
		Max:       max,
	})
	return b
}


func (b *SearchBuilder) GeoFilter(field string, lon, lat, radius float64, unit string) *SearchBuilder {
	b.options.GeoFilter = append(b.options.GeoFilter, FTSearchGeoFilter{
		FieldName: field,
		Longitude: lon,
		Latitude:  lat,
		Radius:    radius,
		Unit:      unit,
	})
	return b
}


func (b *SearchBuilder) InKeys(keys ...interface{}) *SearchBuilder {
	b.options.InKeys = append(b.options.InKeys, keys...)
	return b
}


func (b *SearchBuilder) InFields(fields ...interface{}) *SearchBuilder {
	b.options.InFields = append(b.options.InFields, fields...)
	return b
}


func (b *SearchBuilder) ReturnFields(fields ...string) *SearchBuilder {
	for _, f := range fields {
		b.options.Return = append(b.options.Return, FTSearchReturn{FieldName: f})
	}
	return b
}


func (b *SearchBuilder) ReturnAs(field, alias string) *SearchBuilder {
	b.options.Return = append(b.options.Return, FTSearchReturn{FieldName: field, As: alias})
	return b
}


func (b *SearchBuilder) Slop(slop int) *SearchBuilder {
	b.options.Slop = slop
	return b
}


func (b *SearchBuilder) Timeout(timeout int) *SearchBuilder {
	b.options.Timeout = timeout
	return b
}


func (b *SearchBuilder) InOrder() *SearchBuilder {
	b.options.InOrder = true
	return b
}


func (b *SearchBuilder) Language(lang string) *SearchBuilder {
	b.options.Language = lang
	return b
}


func (b *SearchBuilder) Expander(expander string) *SearchBuilder {
	b.options.Expander = expander
	return b
}


func (b *SearchBuilder) Scorer(scorer string) *SearchBuilder {
	b.options.Scorer = scorer
	return b
}


func (b *SearchBuilder) ExplainScore() *SearchBuilder {
	b.options.ExplainScore = true
	return b
}


func (b *SearchBuilder) Payload(payload string) *SearchBuilder {
	b.options.Payload = payload
	return b
}


func (b *SearchBuilder) SortBy(field string, asc bool) *SearchBuilder {
	b.options.SortBy = append(b.options.SortBy, FTSearchSortBy{
		FieldName: field,
		Asc:       asc,
		Desc:      !asc,
	})
	return b
}


func (b *SearchBuilder) WithSortByCount() *SearchBuilder {
	b.options.SortByWithCount = true
	return b
}


func (b *SearchBuilder) Param(key string, value interface{}) *SearchBuilder {
	if b.options.Params == nil {
		b.options.Params = make(map[string]interface{}, 1)
	}
	b.options.Params[key] = value
	return b
}


func (b *SearchBuilder) ParamsMap(p map[string]interface{}) *SearchBuilder {
	if b.options.Params == nil {
		b.options.Params = make(map[string]interface{}, len(p))
	}
	for k, v := range p {
		b.options.Params[k] = v
	}
	return b
}


func (b *SearchBuilder) Dialect(version int) *SearchBuilder {
	b.options.DialectVersion = version
	return b
}


func (b *SearchBuilder) Limit(offset, count int) *SearchBuilder {
	b.options.LimitOffset = offset
	b.options.Limit = count
	return b
}
func (b *SearchBuilder) CountOnly() *SearchBuilder { b.options.CountOnly = true; return b }


func (b *SearchBuilder) Run() (FTSearchResult, error) {
	cmd := b.c.FTSearchWithArgs(b.ctx, b.index, b.query, b.options)
	return cmd.Result()
}





type AggregateBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	query   string
	options *FTAggregateOptions
}



func (c *Client) NewAggregateBuilder(ctx context.Context, index, query string) *AggregateBuilder {
	return &AggregateBuilder{c: c, ctx: ctx, index: index, query: query, options: &FTAggregateOptions{LimitOffset: -1}}
}


func (b *AggregateBuilder) Verbatim() *AggregateBuilder { b.options.Verbatim = true; return b }


func (b *AggregateBuilder) AddScores() *AggregateBuilder { b.options.AddScores = true; return b }


func (b *AggregateBuilder) Scorer(s string) *AggregateBuilder {
	b.options.Scorer = s
	return b
}


func (b *AggregateBuilder) LoadAll() *AggregateBuilder {
	b.options.LoadAll = true
	return b
}



func (b *AggregateBuilder) Load(field string, alias ...string) *AggregateBuilder {
	
	l := FTAggregateLoad{Field: field}
	if len(alias) > 0 {
		l.As = alias[0]
	}
	b.options.Load = append(b.options.Load, l)
	return b
}


func (b *AggregateBuilder) Timeout(ms int) *AggregateBuilder {
	b.options.Timeout = ms
	return b
}


func (b *AggregateBuilder) Apply(field string, alias ...string) *AggregateBuilder {
	a := FTAggregateApply{Field: field}
	if len(alias) > 0 {
		a.As = alias[0]
	}
	b.options.Apply = append(b.options.Apply, a)
	return b
}


func (b *AggregateBuilder) GroupBy(fields ...interface{}) *AggregateBuilder {
	b.options.GroupBy = append(b.options.GroupBy, FTAggregateGroupBy{
		Fields: fields,
	})
	return b
}


func (b *AggregateBuilder) Reduce(fn SearchAggregator, args ...interface{}) *AggregateBuilder {
	if len(b.options.GroupBy) == 0 {
		
		return b
	}
	idx := len(b.options.GroupBy) - 1
	b.options.GroupBy[idx].Reduce = append(b.options.GroupBy[idx].Reduce, FTAggregateReducer{
		Reducer: fn,
		Args:    args,
	})
	return b
}


func (b *AggregateBuilder) ReduceAs(fn SearchAggregator, alias string, args ...interface{}) *AggregateBuilder {
	if len(b.options.GroupBy) == 0 {
		return b
	}
	idx := len(b.options.GroupBy) - 1
	b.options.GroupBy[idx].Reduce = append(b.options.GroupBy[idx].Reduce, FTAggregateReducer{
		Reducer: fn,
		Args:    args,
		As:      alias,
	})
	return b
}


func (b *AggregateBuilder) SortBy(field string, asc bool) *AggregateBuilder {
	sb := FTAggregateSortBy{FieldName: field, Asc: asc, Desc: !asc}
	b.options.SortBy = append(b.options.SortBy, sb)
	return b
}


func (b *AggregateBuilder) SortByMax(max int) *AggregateBuilder {
	b.options.SortByMax = max
	return b
}


func (b *AggregateBuilder) Filter(expr string) *AggregateBuilder {
	b.options.Filter = expr
	return b
}


func (b *AggregateBuilder) WithCursor(count, maxIdle int) *AggregateBuilder {
	b.options.WithCursor = true
	if b.options.WithCursorOptions == nil {
		b.options.WithCursorOptions = &FTAggregateWithCursor{}
	}
	b.options.WithCursorOptions.Count = count
	b.options.WithCursorOptions.MaxIdle = maxIdle
	return b
}


func (b *AggregateBuilder) Params(p map[string]interface{}) *AggregateBuilder {
	if b.options.Params == nil {
		b.options.Params = make(map[string]interface{}, len(p))
	}
	for k, v := range p {
		b.options.Params[k] = v
	}
	return b
}


func (b *AggregateBuilder) Dialect(version int) *AggregateBuilder {
	b.options.DialectVersion = version
	return b
}


func (b *AggregateBuilder) Run() (*FTAggregateResult, error) {
	cmd := b.c.FTAggregateWithArgs(b.ctx, b.index, b.query, b.options)
	return cmd.Result()
}






type CreateIndexBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	options *FTCreateOptions
	schema  []*FieldSchema
}



func (c *Client) NewCreateIndexBuilder(ctx context.Context, index string) *CreateIndexBuilder {
	return &CreateIndexBuilder{c: c, ctx: ctx, index: index, options: &FTCreateOptions{}}
}


func (b *CreateIndexBuilder) OnHash() *CreateIndexBuilder { b.options.OnHash = true; return b }


func (b *CreateIndexBuilder) OnJSON() *CreateIndexBuilder { b.options.OnJSON = true; return b }


func (b *CreateIndexBuilder) Prefix(prefixes ...interface{}) *CreateIndexBuilder {
	b.options.Prefix = prefixes
	return b
}


func (b *CreateIndexBuilder) Filter(filter string) *CreateIndexBuilder {
	b.options.Filter = filter
	return b
}


func (b *CreateIndexBuilder) DefaultLanguage(lang string) *CreateIndexBuilder {
	b.options.DefaultLanguage = lang
	return b
}


func (b *CreateIndexBuilder) LanguageField(field string) *CreateIndexBuilder {
	b.options.LanguageField = field
	return b
}


func (b *CreateIndexBuilder) Score(score float64) *CreateIndexBuilder {
	b.options.Score = score
	return b
}


func (b *CreateIndexBuilder) ScoreField(field string) *CreateIndexBuilder {
	b.options.ScoreField = field
	return b
}


func (b *CreateIndexBuilder) PayloadField(field string) *CreateIndexBuilder {
	b.options.PayloadField = field
	return b
}


func (b *CreateIndexBuilder) NoOffsets() *CreateIndexBuilder { b.options.NoOffsets = true; return b }


func (b *CreateIndexBuilder) Temporary(sec int) *CreateIndexBuilder {
	b.options.Temporary = sec
	return b
}


func (b *CreateIndexBuilder) NoHL() *CreateIndexBuilder { b.options.NoHL = true; return b }


func (b *CreateIndexBuilder) NoFields() *CreateIndexBuilder { b.options.NoFields = true; return b }


func (b *CreateIndexBuilder) NoFreqs() *CreateIndexBuilder { b.options.NoFreqs = true; return b }


func (b *CreateIndexBuilder) StopWords(words ...interface{}) *CreateIndexBuilder {
	b.options.StopWords = words
	return b
}


func (b *CreateIndexBuilder) SkipInitialScan() *CreateIndexBuilder {
	b.options.SkipInitialScan = true
	return b
}


func (b *CreateIndexBuilder) Schema(field *FieldSchema) *CreateIndexBuilder {
	b.schema = append(b.schema, field)
	return b
}


func (b *CreateIndexBuilder) Run() (string, error) {
	cmd := b.c.FTCreate(b.ctx, b.index, b.options, b.schema...)
	return cmd.Result()
}






type DropIndexBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	options *FTDropIndexOptions
}



func (c *Client) NewDropIndexBuilder(ctx context.Context, index string) *DropIndexBuilder {
	return &DropIndexBuilder{c: c, ctx: ctx, index: index}
}


func (b *DropIndexBuilder) DeleteDocs() *DropIndexBuilder { b.options.DeleteDocs = true; return b }


func (b *DropIndexBuilder) Run() (string, error) {
	cmd := b.c.FTDropIndexWithArgs(b.ctx, b.index, b.options)
	return cmd.Result()
}






type AliasBuilder struct {
	c      *Client
	ctx    context.Context
	alias  string
	index  string
	action string 
}



func (c *Client) NewAliasBuilder(ctx context.Context, alias string) *AliasBuilder {
	return &AliasBuilder{c: c, ctx: ctx, alias: alias}
}


func (b *AliasBuilder) Action(action string) *AliasBuilder {
	b.action = action
	return b
}


func (b *AliasBuilder) Add(index string) *AliasBuilder {
	b.action = "add"
	b.index = index
	return b
}


func (b *AliasBuilder) Del() *AliasBuilder {
	b.action = "del"
	return b
}


func (b *AliasBuilder) Update(index string) *AliasBuilder {
	b.action = "update"
	b.index = index
	return b
}


func (b *AliasBuilder) Run() (string, error) {
	switch b.action {
	case "add":
		cmd := b.c.FTAliasAdd(b.ctx, b.index, b.alias)
		return cmd.Result()
	case "del":
		cmd := b.c.FTAliasDel(b.ctx, b.alias)
		return cmd.Result()
	case "update":
		cmd := b.c.FTAliasUpdate(b.ctx, b.index, b.alias)
		return cmd.Result()
	}
	return "", nil
}






type ExplainBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	query   string
	options *FTExplainOptions
}



func (c *Client) NewExplainBuilder(ctx context.Context, index, query string) *ExplainBuilder {
	return &ExplainBuilder{c: c, ctx: ctx, index: index, query: query, options: &FTExplainOptions{}}
}


func (b *ExplainBuilder) Dialect(d string) *ExplainBuilder { b.options.Dialect = d; return b }


func (b *ExplainBuilder) Run() (string, error) {
	cmd := b.c.FTExplainWithArgs(b.ctx, b.index, b.query, b.options)
	return cmd.Result()
}





type FTInfoBuilder struct {
	c     *Client
	ctx   context.Context
	index string
}


func (c *Client) NewSearchInfoBuilder(ctx context.Context, index string) *FTInfoBuilder {
	return &FTInfoBuilder{c: c, ctx: ctx, index: index}
}


func (b *FTInfoBuilder) Run() (FTInfoResult, error) {
	cmd := b.c.FTInfo(b.ctx, b.index)
	return cmd.Result()
}






type SpellCheckBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	query   string
	options *FTSpellCheckOptions
}



func (c *Client) NewSpellCheckBuilder(ctx context.Context, index, query string) *SpellCheckBuilder {
	return &SpellCheckBuilder{c: c, ctx: ctx, index: index, query: query, options: &FTSpellCheckOptions{}}
}


func (b *SpellCheckBuilder) Distance(d int) *SpellCheckBuilder { b.options.Distance = d; return b }


func (b *SpellCheckBuilder) Terms(include bool, dictionary string, terms ...interface{}) *SpellCheckBuilder {
	if b.options.Terms == nil {
		b.options.Terms = &FTSpellCheckTerms{}
	}
	if include {
		b.options.Terms.Inclusion = "INCLUDE"
	} else {
		b.options.Terms.Inclusion = "EXCLUDE"
	}
	b.options.Terms.Dictionary = dictionary
	b.options.Terms.Terms = terms
	return b
}


func (b *SpellCheckBuilder) Dialect(d int) *SpellCheckBuilder { b.options.Dialect = d; return b }


func (b *SpellCheckBuilder) Run() ([]SpellCheckResult, error) {
	cmd := b.c.FTSpellCheckWithArgs(b.ctx, b.index, b.query, b.options)
	return cmd.Result()
}






type DictBuilder struct {
	c      *Client
	ctx    context.Context
	dict   string
	terms  []interface{}
	action string 
}



func (c *Client) NewDictBuilder(ctx context.Context, dict string) *DictBuilder {
	return &DictBuilder{c: c, ctx: ctx, dict: dict}
}


func (b *DictBuilder) Action(action string) *DictBuilder {
	b.action = action
	return b
}


func (b *DictBuilder) Add(terms ...interface{}) *DictBuilder {
	b.action = "add"
	b.terms = terms
	return b
}


func (b *DictBuilder) Del(terms ...interface{}) *DictBuilder {
	b.action = "del"
	b.terms = terms
	return b
}


func (b *DictBuilder) Dump() *DictBuilder {
	b.action = "dump"
	return b
}


func (b *DictBuilder) Run() (interface{}, error) {
	switch b.action {
	case "add":
		cmd := b.c.FTDictAdd(b.ctx, b.dict, b.terms...)
		return cmd.Result()
	case "del":
		cmd := b.c.FTDictDel(b.ctx, b.dict, b.terms...)
		return cmd.Result()
	case "dump":
		cmd := b.c.FTDictDump(b.ctx, b.dict)
		return cmd.Result()
	}
	return nil, nil
}






type TagValsBuilder struct {
	c     *Client
	ctx   context.Context
	index string
	field string
}



func (c *Client) NewTagValsBuilder(ctx context.Context, index, field string) *TagValsBuilder {
	return &TagValsBuilder{c: c, ctx: ctx, index: index, field: field}
}


func (b *TagValsBuilder) Run() ([]string, error) {
	cmd := b.c.FTTagVals(b.ctx, b.index, b.field)
	return cmd.Result()
}






type CursorBuilder struct {
	c        *Client
	ctx      context.Context
	index    string
	cursorId int64
	count    int
	action   string 
}



func (c *Client) NewCursorBuilder(ctx context.Context, index string, cursorId int64) *CursorBuilder {
	return &CursorBuilder{c: c, ctx: ctx, index: index, cursorId: cursorId}
}


func (b *CursorBuilder) Action(action string) *CursorBuilder {
	b.action = action
	return b
}


func (b *CursorBuilder) Read() *CursorBuilder {
	b.action = "read"
	return b
}


func (b *CursorBuilder) Del() *CursorBuilder {
	b.action = "del"
	return b
}


func (b *CursorBuilder) Count(count int) *CursorBuilder { b.count = count; return b }


func (b *CursorBuilder) Run() (interface{}, error) {
	switch b.action {
	case "read":
		cmd := b.c.FTCursorRead(b.ctx, b.index, int(b.cursorId), b.count)
		return cmd.Result()
	case "del":
		cmd := b.c.FTCursorDel(b.ctx, b.index, int(b.cursorId))
		return cmd.Result()
	}
	return nil, nil
}






type SynUpdateBuilder struct {
	c       *Client
	ctx     context.Context
	index   string
	groupId interface{}
	options *FTSynUpdateOptions
	terms   []interface{}
}



func (c *Client) NewSynUpdateBuilder(ctx context.Context, index string, groupId interface{}) *SynUpdateBuilder {
	return &SynUpdateBuilder{c: c, ctx: ctx, index: index, groupId: groupId, options: &FTSynUpdateOptions{}}
}


func (b *SynUpdateBuilder) SkipInitialScan() *SynUpdateBuilder {
	b.options.SkipInitialScan = true
	return b
}


func (b *SynUpdateBuilder) Terms(terms ...interface{}) *SynUpdateBuilder { b.terms = terms; return b }


func (b *SynUpdateBuilder) Run() (string, error) {
	cmd := b.c.FTSynUpdateWithArgs(b.ctx, b.index, b.groupId, b.options, b.terms)
	return cmd.Result()
}
