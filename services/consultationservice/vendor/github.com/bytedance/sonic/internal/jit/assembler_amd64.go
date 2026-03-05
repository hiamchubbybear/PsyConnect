

package jit

import (
    `encoding/binary`
    `strconv`
    `strings`
    `sync`

    `github.com/bytedance/sonic/loader`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/twitchyliquid64/golang-asm/obj`
    `github.com/twitchyliquid64/golang-asm/obj/x86`
)

const (
    _LB_jump_pc = "_jump_pc_"
)

type BaseAssembler struct {
    i        int
    f        func()
    c        []byte
    o        sync.Once
    pb       *Backend
    xrefs    map[string][]*obj.Prog
    labels   map[string]*obj.Prog
    pendings map[string][]*obj.Prog
}



var _NOPS = [][16]byte {
    {0x90},                                                     
    {0x66, 0x90},                                               
    {0x0f, 0x1f, 0x00},                                         
    {0x0f, 0x1f, 0x40, 0x00},                                   
    {0x0f, 0x1f, 0x44, 0x00, 0x00},                             
    {0x66, 0x0f, 0x1f, 0x44, 0x00, 0x00},                       
    {0x0f, 0x1f, 0x80, 0x00, 0x00, 0x00, 0x00},                 
    {0x0f, 0x1f, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},           
    {0x66, 0x0f, 0x1f, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},     
}

func (self *BaseAssembler) NOP() *obj.Prog {
    p := self.pb.New()
    p.As = obj.ANOP
    self.pb.Append(p)
    return p
}

func (self *BaseAssembler) NOPn(n int) {
    for i := len(_NOPS); i > 0 && n > 0; i-- {
        for ; n >= i; n -= i {
            self.Byte(_NOPS[i - 1][:i]...)
        }
    }
}

func (self *BaseAssembler) Byte(v ...byte) {
    for ; len(v) >= 8; v = v[8:] { self.From("QUAD", Imm(rt.Get64(v))) }
    for ; len(v) >= 4; v = v[4:] { self.From("LONG", Imm(int64(rt.Get32(v)))) }
    for ; len(v) >= 2; v = v[2:] { self.From("WORD", Imm(int64(rt.Get16(v)))) }
    for ; len(v) >= 1; v = v[1:] { self.From("BYTE", Imm(int64(v[0]))) }
}

func (self *BaseAssembler) Mark(pc int) {
    self.i++
    self.Link(_LB_jump_pc + strconv.Itoa(pc))
}

func (self *BaseAssembler) Link(to string) {
    var p *obj.Prog
    var v []*obj.Prog

    
    if strings.Contains(to, "{n}") {
        to = strings.ReplaceAll(to, "{n}", strconv.Itoa(self.i))
    }

    
    if _, ok := self.labels[to]; ok {
        panic("label " + to + " has already been linked")
    }

    
    p = self.NOP()
    v = self.pendings[to]

    
    for _, q := range v {
        q.To.Val = p
    }

    
    self.labels[to] = p
    delete(self.pendings, to)
}

func (self *BaseAssembler) Xref(pc int, d int64) {
    self.Sref(_LB_jump_pc + strconv.Itoa(pc), d)
}

func (self *BaseAssembler) Sref(to string, d int64) {
    p := self.pb.New()
    p.As = x86.ALONG
    p.From = Imm(-d)

    
    if strings.Contains(to, "{n}") {
        to = strings.ReplaceAll(to, "{n}", strconv.Itoa(self.i))
    }

    
    self.pb.Append(p)
    self.xrefs[to] = append(self.xrefs[to], p)
}

func (self *BaseAssembler) Xjmp(op string, to int) {
    self.Sjmp(op, _LB_jump_pc + strconv.Itoa(to))
}

func (self *BaseAssembler) Sjmp(op string, to string) {
    p := self.pb.New()
    p.As = As(op)

    
    if strings.Contains(to, "{n}") {
        to = strings.ReplaceAll(to, "{n}", strconv.Itoa(self.i))
    }

    
    if v, ok := self.labels[to]; ok {
        p.To.Val = v
    } else {
        self.pendings[to] = append(self.pendings[to], p)
    }

    
    p.To.Type = obj.TYPE_BRANCH
    self.pb.Append(p)
}

func (self *BaseAssembler) Rjmp(op string, to obj.Addr) {
    p := self.pb.New()
    p.To = to
    p.As = As(op)
    self.pb.Append(p)
}

func (self *BaseAssembler) From(op string, val obj.Addr) {
    p := self.pb.New()
    p.As = As(op)
    p.From = val
    self.pb.Append(p)
}

func (self *BaseAssembler) Emit(op string, args ...obj.Addr) {
    p := self.pb.New()
    p.As = As(op)
    self.assignOperands(p, args)
    self.pb.Append(p)
}

func (self *BaseAssembler) assignOperands(p *obj.Prog, args []obj.Addr) {
    switch len(args) {
        case 0  :
        case 1  : p.To                     = args[0]
        case 2  : p.To, p.From             = args[1], args[0]
        case 3  : p.To, p.From, p.RestArgs = args[2], args[0], args[1:2]
        case 4  : p.To, p.From, p.RestArgs = args[2], args[3], args[:2]
        default : panic("invalid operands")
    }
}



func (self *BaseAssembler) Size() int {
    self.build()
    return len(self.c)
}

func (self *BaseAssembler) Init(f func()) {
    self.i = 0
    self.f = f
    self.c = nil
    self.o = sync.Once{}
}

var jitLoader = loader.Loader{
    Name: "sonic.jit.",
    File: "github.com/bytedance/sonic/jit.go",
    Options: loader.Options{
        NoPreempt: true,
    },
}

func (self *BaseAssembler) Load(name string, frameSize int, argSize int, argStackmap []bool, localStackmap []bool) loader.Function {
    self.build()
    return jitLoader.LoadOne(self.c, name, frameSize, argSize, argStackmap, localStackmap)
}



func (self *BaseAssembler) init() {
    self.pb       = newBackend("amd64")
    self.xrefs    = map[string][]*obj.Prog{}
    self.labels   = map[string]*obj.Prog{}
    self.pendings = map[string][]*obj.Prog{}
}

func (self *BaseAssembler) build() {
    self.o.Do(func() {
        self.init()
        self.f()
        self.validate()
        self.assemble()
        self.resolve()
        self.release()
    })
}

func (self *BaseAssembler) release() {
    self.pb.Release()
    self.pb = nil
    self.xrefs = nil
    self.labels = nil
    self.pendings = nil
}

func (self *BaseAssembler) resolve() {
    for s, v := range self.xrefs {
        for _, prog := range v {
            if prog.As != x86.ALONG {
                panic("invalid RIP relative reference")
            } else if p, ok := self.labels[s]; !ok {
                panic("links are not fully resolved: " + s)
            } else {
                off := prog.From.Offset + p.Pc - prog.Pc
                binary.LittleEndian.PutUint32(self.c[prog.Pc:], uint32(off))
            }
        }
    }
}

func (self *BaseAssembler) validate() {
    for key := range self.pendings {
        panic("links are not fully resolved: " + key)
    }
}

func (self *BaseAssembler) assemble() {
    self.c = self.pb.Assemble()
}
