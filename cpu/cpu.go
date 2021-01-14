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

var userSysModReg = [7]reg{r8, r9, r10, r11, r12, r13, r14}

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

type instructionHandler func(cpu *CPU)

// CPU represents an ARM CPU
type CPU struct {
	registers          [37]uint32
	bankedRegisters    []reg
	currentInstruction uint32
	currentSpsr        reg
	currentMode        cpuMode
	nextInstruction    uint32
	instructionMap     map[uint32]instructionHandler
}

// New constructs a new cpu object
func New() *CPU {
	return &CPU{
		registers:       [37]uint32{},
		bankedRegisters: bankedRegMap[user],
		currentMode:     user,
	}
}

func (cpu *CPU) readReg(register reg) uint32 {
	if register == rPc {
		return cpu.readPc()
	}
	if register < r8 {
		return cpu.registers[register]
	}

	return cpu.registers[cpu.bankedRegisters[register-8]]
}

func (cpu *CPU) setReg(register reg, value uint32) {
	if register == rPc {
		cpu.setPc(value)
		return
	}

	if register < r8 {
		cpu.registers[register] = value
		return
	}

	cpu.registers[cpu.bankedRegisters[register-8]] = value
}

func (cpu *CPU) readPc() uint32 {
	return cpu.registers[rPc]
}

func (cpu *CPU) setPc(value uint32) {
	cpu.registers[rPc] = value
}

func (cpu *CPU) readCpsr() uint32 {
	return cpu.registers[rCpsr]
}

func (cpu *CPU) setCpsr(value uint32) {
	cpu.registers[rCpsr] = value
}

func (cpu *CPU) readSpsr() uint32 {
	return cpu.registers[cpu.currentSpsr]
}

func (cpu *CPU) setSpsr(value uint32) {
	cpu.registers[cpu.currentSpsr] = value
}

func (cpu *CPU) tick() {

}

func (cpu *CPU) executeOpcode() {

}
