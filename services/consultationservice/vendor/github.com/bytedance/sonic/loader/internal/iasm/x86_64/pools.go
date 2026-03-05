















package x86_64


func CreateLabel(name string) *Label {
	p := new(Label)

	
	p.refs = 1
	p.Name = name
	return p
}

func newProgram(arch *Arch) *Program {
	p := new(Program)

	
	p.arch = arch
	return p
}

func newInstruction(name string, argc int, argv Operands) *Instruction {
	p := new(Instruction)

	
	p.name = name
	p.argc = argc
	p.argv = argv
	return p
}


func CreateMemoryOperand() *MemoryOperand {
	p := new(MemoryOperand)

	
	p.refs = 1
	return p
}
