



//go:build go1.10

package bidirule

func (t *Transformer) isFinal() bool {
	return t.state == ruleLTRFinal || t.state == ruleRTLFinal || t.state == ruleInitial
}
