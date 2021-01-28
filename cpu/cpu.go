package cpu

type reg int

const (
	r0 reg = iota
	r1
	r2
	r3
	r4
	r5
	r6
	r7
	r8
	r9
	r10
	r11
	r12
	r13
	r14
	rPc
	r8Fiq
	r9Fiq
	r10Fiq
	r11Fiq
	r12Fiq
	r13Fiq
	r14Fiq
	r13Irq
	r14Irq
	r13Svc
	r14Svc
	r13Abt
	r14Abt
	r13Und
	r14Und
	rCpsr
	rSpsrFiq
	rSpsrIrq
	rSpsrSvc
	rSpsrAbt
	rSpsrUnd
)

var bankedRegMap = map[cpuMode][]reg{
	user:       {r8, r9, r10, r11, r12, r13, r14},
	fiq:        {r8Fiq, r9Fiq, r10Fiq, r11Fiq, r12Fiq, r13Fiq, r14Fiq},
	irq:        {r8, r9, r10, r11, r12, r13Irq, r14Irq},
	supervisor: {r8, r9, r10, r11, r12, r13Svc, r14Svc},
	abort:      {r8, r9, r10, r11, r12, r13Abt, r14Abt},
	undefined:  {r8, r9, r10, r11, r12, r13Und, r14Und},
	system:     {r8, r9, r10, r11, r12, r13, r14},
}

var bankedSpsrMap = map[cpuMode]reg{
	fiq:        rSpsrFiq,
	irq:        rSpsrIrq,
	supervisor: rSpsrSvc,
	abort:      rSpsrAbt,
	undefined:  rSpsrUnd,
}

// Architecture is a typedef for the ARM CPU archs
type Architecture int

const (
	// V5 is version 5
	V5 Architecture = iota
	// V4 is version 4
	V4
)

// ArmCPU describes a generic Arm CPU
type ArmCPU interface {
	readReg(register reg) uint32
	setReg(reg, uint32)
	readPc() uint32
	setPc(value uint32)
	readCpsr() uint32
	setCpsr(value uint32)
	readSpsr() uint32
	setSpsr(value uint32)
	currentInstruction() uint32
	nextInstruction() uint32
	isFlag(fl flag) bool
	setFlag(fl flag, value bool)
	highIrqVectors() bool
	Architecture() Architecture
	mode() cpuMode
	setMode(mode cpuMode)
	WaitForTick()
	Tick()
	Run()
	Bus() DataBus
}
