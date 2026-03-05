



package language


type CompactCoreInfo uint32



func GetCompactCore(t Tag) (cci CompactCoreInfo, ok bool) {
	if t.LangID > langNoIndexOffset {
		return 0, false
	}
	cci |= CompactCoreInfo(t.LangID) << (8 + 12)
	cci |= CompactCoreInfo(t.ScriptID) << 12
	cci |= CompactCoreInfo(t.RegionID)
	return cci, true
}


func (c CompactCoreInfo) Tag() Tag {
	return Tag{
		LangID:   Language(c >> 20),
		RegionID: Region(c & 0x3ff),
		ScriptID: Script(c>>12) & 0xff,
	}
}
