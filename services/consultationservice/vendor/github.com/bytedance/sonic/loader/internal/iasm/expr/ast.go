















package expr

import (
	"fmt"
)


type Type int

const (
	
	CONST Type = iota

	
	TERM

	
	EXPR
)

var typeNames = map[Type]string{
	EXPR:  "Expr",
	TERM:  "Term",
	CONST: "Const",
}


func (self Type) String() string {
	if v, ok := typeNames[self]; ok {
		return v
	} else {
		return fmt.Sprintf("expr.Type(%d)", self)
	}
}


type Operator uint8

const (
	
	ADD Operator = iota

	
	SUB

	
	MUL

	
	DIV

	
	MOD

	
	AND

	
	OR

	
	XOR

	
	SHL

	
	SHR

	
	POW

	
	NOT

	
	NEG
)

var operatorNames = map[Operator]string{
	ADD: "Add",
	SUB: "Subtract",
	MUL: "Multiply",
	DIV: "Divide",
	MOD: "Modulo",
	AND: "And",
	OR:  "Or",
	XOR: "ExclusiveOr",
	SHL: "ShiftLeft",
	SHR: "ShiftRight",
	POW: "Power",
	NOT: "Invert",
	NEG: "Negate",
}


func (self Operator) String() string {
	if v, ok := operatorNames[self]; ok {
		return v
	} else {
		return fmt.Sprintf("expr.Operator(%d)", self)
	}
}


type Expr struct {
	Type  Type
	Term  Term
	Op    Operator
	Left  *Expr
	Right *Expr
	Const int64
}


func Ref(t Term) (p *Expr) {
	p = newExpression()
	p.Term = t
	p.Type = TERM
	return
}


func Int(v int64) (p *Expr) {
	p = newExpression()
	p.Type = CONST
	p.Const = v
	return
}

func (self *Expr) clear() {
	if self.Term != nil {
		self.Term.Free()
	}
	if self.Left != nil {
		self.Left.Free()
	}
	if self.Right != nil {
		self.Right.Free()
	}
}



func (self *Expr) Free() {
	self.clear()
	freeExpression(self)
}



func (self *Expr) Evaluate() (int64, error) {
	switch self.Type {
	case EXPR:
		return self.eval()
	case TERM:
		return self.Term.Evaluate()
	case CONST:
		return self.Const, nil
	default:
		panic("invalid expression type: " + self.Type.String())
	}
}



func combine(a *Expr, op Operator, b *Expr) (r *Expr) {
	r = newExpression()
	r.Op = op
	r.Type = EXPR
	r.Left = a
	r.Right = b
	return
}

func (self *Expr) Add(v *Expr) *Expr { return combine(self, ADD, v) }
func (self *Expr) Sub(v *Expr) *Expr { return combine(self, SUB, v) }
func (self *Expr) Mul(v *Expr) *Expr { return combine(self, MUL, v) }
func (self *Expr) Div(v *Expr) *Expr { return combine(self, DIV, v) }
func (self *Expr) Mod(v *Expr) *Expr { return combine(self, MOD, v) }
func (self *Expr) And(v *Expr) *Expr { return combine(self, AND, v) }
func (self *Expr) Or(v *Expr) *Expr  { return combine(self, OR, v) }
func (self *Expr) Xor(v *Expr) *Expr { return combine(self, XOR, v) }
func (self *Expr) Shl(v *Expr) *Expr { return combine(self, SHL, v) }
func (self *Expr) Shr(v *Expr) *Expr { return combine(self, SHR, v) }
func (self *Expr) Pow(v *Expr) *Expr { return combine(self, POW, v) }
func (self *Expr) Not() *Expr        { return combine(self, NOT, nil) }
func (self *Expr) Neg() *Expr        { return combine(self, NEG, nil) }



var binaryEvaluators = [256]func(int64, int64) (int64, error){
	ADD: func(a, b int64) (int64, error) { return a + b, nil },
	SUB: func(a, b int64) (int64, error) { return a - b, nil },
	MUL: func(a, b int64) (int64, error) { return a * b, nil },
	DIV: idiv,
	MOD: imod,
	AND: func(a, b int64) (int64, error) { return a & b, nil },
	OR:  func(a, b int64) (int64, error) { return a | b, nil },
	XOR: func(a, b int64) (int64, error) { return a ^ b, nil },
	SHL: func(a, b int64) (int64, error) { return a << b, nil },
	SHR: func(a, b int64) (int64, error) { return a >> b, nil },
	POW: ipow,
}

func (self *Expr) eval() (int64, error) {
	var lhs int64
	var rhs int64
	var err error
	var vfn func(int64, int64) (int64, error)

	
	if lhs, err = self.Left.Evaluate(); err != nil {
		return 0, err
	}

	
	switch self.Op {
	case NOT:
		return self.unaryNot(lhs)
	case NEG:
		return self.unaryNeg(lhs)
	}

	
	if vfn = binaryEvaluators[self.Op]; vfn == nil {
		panic("invalid operator: " + self.Op.String())
	}

	
	if self.Right == nil {
		panic("operator " + self.Op.String() + " is a binary operator")
	}

	
	if rhs, err = self.Right.Evaluate(); err != nil {
		return 0, err
	} else {
		return vfn(lhs, rhs)
	}
}

func (self *Expr) unaryNot(v int64) (int64, error) {
	if self.Right == nil {
		return ^v, nil
	} else {
		panic("operator Invert is an unary operator")
	}
}

func (self *Expr) unaryNeg(v int64) (int64, error) {
	if self.Right == nil {
		return -v, nil
	} else {
		panic("operator Negate is an unary operator")
	}
}
