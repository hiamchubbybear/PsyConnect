

//go:build !appengine && !noasm && gc && !noasm

package s2

func _dummy_()





//go:noescape
func encodeBlockAsm(dst []byte, src []byte) int





//go:noescape
func encodeBlockAsm4MB(dst []byte, src []byte) int





//go:noescape
func encodeBlockAsm12B(dst []byte, src []byte) int





//go:noescape
func encodeBlockAsm10B(dst []byte, src []byte) int





//go:noescape
func encodeBlockAsm8B(dst []byte, src []byte) int





//go:noescape
func encodeBetterBlockAsm(dst []byte, src []byte) int





//go:noescape
func encodeBetterBlockAsm4MB(dst []byte, src []byte) int





//go:noescape
func encodeBetterBlockAsm12B(dst []byte, src []byte) int





//go:noescape
func encodeBetterBlockAsm10B(dst []byte, src []byte) int





//go:noescape
func encodeBetterBlockAsm8B(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBlockAsm(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBlockAsm64K(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBlockAsm12B(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBlockAsm10B(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBlockAsm8B(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBetterBlockAsm(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBetterBlockAsm64K(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBetterBlockAsm12B(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBetterBlockAsm10B(dst []byte, src []byte) int





//go:noescape
func encodeSnappyBetterBlockAsm8B(dst []byte, src []byte) int





//go:noescape
func calcBlockSize(src []byte) int





//go:noescape
func calcBlockSizeSmall(src []byte) int








//go:noescape
func emitLiteral(dst []byte, lit []byte) int




//go:noescape
func emitRepeat(dst []byte, offset int, length int) int









//go:noescape
func emitCopy(dst []byte, offset int, length int) int









//go:noescape
func emitCopyNoRepeat(dst []byte, offset int, length int) int







//go:noescape
func matchLen(a []byte, b []byte) int



//go:noescape
func cvtLZ4BlockAsm(dst []byte, src []byte) (uncompressed int, dstUsed int)



//go:noescape
func cvtLZ4sBlockAsm(dst []byte, src []byte) (uncompressed int, dstUsed int)



//go:noescape
func cvtLZ4BlockSnappyAsm(dst []byte, src []byte) (uncompressed int, dstUsed int)



//go:noescape
func cvtLZ4sBlockSnappyAsm(dst []byte, src []byte) (uncompressed int, dstUsed int)
