





package driver




type LegacyOperationKind uint


const (
	LegacyNone LegacyOperationKind = iota
	LegacyFind
	LegacyGetMore
	LegacyKillCursors
	LegacyListCollections
	LegacyListIndexes
	LegacyHandshake
)
