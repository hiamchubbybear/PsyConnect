



package obj

import "github.com/twitchyliquid64/golang-asm/src"




































type InlTree struct {
	nodes []InlinedCall
}


type InlinedCall struct {
	Parent   int      
	Pos      src.XPos 
	Func     *LSym    
	ParentPC int32    
}


func (tree *InlTree) Add(parent int, pos src.XPos, func_ *LSym) int {
	r := len(tree.nodes)
	call := InlinedCall{
		Parent: parent,
		Pos:    pos,
		Func:   func_,
	}
	tree.nodes = append(tree.nodes, call)
	return r
}

func (tree *InlTree) Parent(inlIndex int) int {
	return tree.nodes[inlIndex].Parent
}

func (tree *InlTree) InlinedFunction(inlIndex int) *LSym {
	return tree.nodes[inlIndex].Func
}

func (tree *InlTree) CallPos(inlIndex int) src.XPos {
	return tree.nodes[inlIndex].Pos
}

func (tree *InlTree) setParentPC(inlIndex int, pc int32) {
	tree.nodes[inlIndex].ParentPC = pc
}





func (ctxt *Link) OutermostPos(xpos src.XPos) src.Pos {
	pos := ctxt.InnermostPos(xpos)

	outerxpos := xpos
	for ix := pos.Base().InliningIndex(); ix >= 0; {
		call := ctxt.InlTree.nodes[ix]
		ix = call.Parent
		outerxpos = call.Pos
	}
	return ctxt.PosTable.Pos(outerxpos)
}








func (ctxt *Link) InnermostPos(xpos src.XPos) src.Pos {
	return ctxt.PosTable.Pos(xpos)
}




func (ctxt *Link) AllPos(xpos src.XPos, result []src.Pos) []src.Pos {
	pos := ctxt.InnermostPos(xpos)
	result = result[:0]
	result = append(result, ctxt.PosTable.Pos(xpos))
	for ix := pos.Base().InliningIndex(); ix >= 0; {
		call := ctxt.InlTree.nodes[ix]
		ix = call.Parent
		result = append(result, ctxt.PosTable.Pos(call.Pos))
	}
	return result
}

func dumpInlTree(ctxt *Link, tree InlTree) {
	for i, call := range tree.nodes {
		pos := ctxt.PosTable.Pos(call.Pos)
		ctxt.Logf("%0d | %0d | %s (%s) pc=%d\n", i, call.Parent, call.Func, pos, call.ParentPC)
	}
}
