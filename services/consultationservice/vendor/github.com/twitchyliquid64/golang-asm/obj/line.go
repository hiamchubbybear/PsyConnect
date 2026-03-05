



package obj

import (
	"github.com/twitchyliquid64/golang-asm/goobj"
	"github.com/twitchyliquid64/golang-asm/src"
)


func (ctxt *Link) AddImport(pkg string, fingerprint goobj.FingerprintType) {
	ctxt.Imports = append(ctxt.Imports, goobj.ImportedPkg{Pkg: pkg, Fingerprint: fingerprint})
}

func linkgetlineFromPos(ctxt *Link, xpos src.XPos) (f string, l int32) {
	pos := ctxt.PosTable.Pos(xpos)
	if !pos.IsKnown() {
		pos = src.Pos{}
	}
	
	return pos.SymFilename(), int32(pos.RelLine())
}


func getFileIndexAndLine(ctxt *Link, xpos src.XPos) (int, int32) {
	f, l := linkgetlineFromPos(ctxt, xpos)
	return ctxt.PosTable.FileIndex(f), l
}
