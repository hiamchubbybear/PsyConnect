



package language

import "errors"

type scriptRegionFlags uint8

const (
	isList = 1 << iota
	scriptInFrom
	regionInFrom
)

func (t *Tag) setUndefinedLang(id Language) {
	if t.LangID == 0 {
		t.LangID = id
	}
}

func (t *Tag) setUndefinedScript(id Script) {
	if t.ScriptID == 0 {
		t.ScriptID = id
	}
}

func (t *Tag) setUndefinedRegion(id Region) {
	if t.RegionID == 0 || t.RegionID.Contains(id) {
		t.RegionID = id
	}
}



var ErrMissingLikelyTagsData = errors.New("missing likely tags data")





func (t Tag) addLikelySubtags() (Tag, error) {
	id, err := addTags(t)
	if err != nil {
		return t, err
	} else if id.equalTags(t) {
		return t, nil
	}
	id.RemakeString()
	return id, nil
}


func specializeRegion(t *Tag) bool {
	if i := regionInclusion[t.RegionID]; i < nRegionGroups {
		x := likelyRegionGroup[i]
		if Language(x.lang) == t.LangID && Script(x.script) == t.ScriptID {
			t.RegionID = Region(x.region)
		}
		return true
	}
	return false
}


func (t Tag) Maximize() (Tag, error) {
	return addTags(t)
}

func addTags(t Tag) (Tag, error) {
	
	if t.IsPrivateUse() {
		return t, nil
	}
	if t.ScriptID != 0 && t.RegionID != 0 {
		if t.LangID != 0 {
			
			specializeRegion(&t)
			return t, nil
		}
		
		
		list := likelyRegion[t.RegionID : t.RegionID+1]
		if x := list[0]; x.flags&isList != 0 {
			list = likelyRegionList[x.lang : x.lang+uint16(x.script)]
		}
		for _, x := range list {
			
			if Script(x.script) == t.ScriptID {
				t.setUndefinedLang(Language(x.lang))
				return t, nil
			}
		}
	}
	if t.LangID != 0 {
		
		if t.LangID < langNoIndexOffset {
			x := likelyLang[t.LangID]
			if x.flags&isList != 0 {
				list := likelyLangList[x.region : x.region+uint16(x.script)]
				if t.ScriptID != 0 {
					for _, x := range list {
						if Script(x.script) == t.ScriptID && x.flags&scriptInFrom != 0 {
							t.setUndefinedRegion(Region(x.region))
							return t, nil
						}
					}
				} else if t.RegionID != 0 {
					count := 0
					goodScript := true
					tt := t
					for _, x := range list {
						
						
						
						
						if x.flags&scriptInFrom == 0 && t.RegionID.Contains(Region(x.region)) {
							tt.RegionID = Region(x.region)
							tt.setUndefinedScript(Script(x.script))
							goodScript = goodScript && tt.ScriptID == Script(x.script)
							count++
						}
					}
					if count == 1 {
						return tt, nil
					}
					
					
					if goodScript {
						t.ScriptID = tt.ScriptID
					}
				}
			}
		}
	} else {
		
		if t.ScriptID != 0 {
			x := likelyScript[t.ScriptID]
			if x.region != 0 {
				t.setUndefinedRegion(Region(x.region))
				t.setUndefinedLang(Language(x.lang))
				return t, nil
			}
		}
		
		
		if t.RegionID != 0 {
			if i := regionInclusion[t.RegionID]; i < nRegionGroups {
				x := likelyRegionGroup[i]
				if x.region != 0 {
					t.setUndefinedLang(Language(x.lang))
					t.setUndefinedScript(Script(x.script))
					t.RegionID = Region(x.region)
				}
			} else {
				x := likelyRegion[t.RegionID]
				if x.flags&isList != 0 {
					x = likelyRegionList[x.lang]
				}
				if x.script != 0 && x.flags != scriptInFrom {
					t.setUndefinedLang(Language(x.lang))
					t.setUndefinedScript(Script(x.script))
					return t, nil
				}
			}
		}
	}

	
	if t.LangID < langNoIndexOffset {
		x := likelyLang[t.LangID]
		if x.flags&isList != 0 {
			x = likelyLangList[x.region]
		}
		if x.region != 0 {
			t.setUndefinedScript(Script(x.script))
			t.setUndefinedRegion(Region(x.region))
		}
		specializeRegion(&t)
		if t.LangID == 0 {
			t.LangID = _en 
		}
		return t, nil
	}
	return t, ErrMissingLikelyTagsData
}

func (t *Tag) setTagsFrom(id Tag) {
	t.LangID = id.LangID
	t.ScriptID = id.ScriptID
	t.RegionID = id.RegionID
}



func (t Tag) minimize() (Tag, error) {
	t, err := minimizeTags(t)
	if err != nil {
		return t, err
	}
	t.RemakeString()
	return t, nil
}


func minimizeTags(t Tag) (Tag, error) {
	if t.equalTags(Und) {
		return t, nil
	}
	max, err := addTags(t)
	if err != nil {
		return t, err
	}
	for _, id := range [...]Tag{
		{LangID: t.LangID},
		{LangID: t.LangID, RegionID: t.RegionID},
		{LangID: t.LangID, ScriptID: t.ScriptID},
	} {
		if x, err := addTags(id); err == nil && max.equalTags(x) {
			t.setTagsFrom(id)
			break
		}
	}
	return t, nil
}
