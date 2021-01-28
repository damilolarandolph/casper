package cpu

import "sync"

// Arm7 represents an Arm7 CPU
type Arm7 struct {
	registers          [37]uint32
	bankedRegisters    []reg
	currentInstruction uint32
	currentSpsr        reg
	currentMode        cpuMode
	nextInstruction    uint32
	instructionMap     map[uint32]instructionHandler
	irqHigh            bool
	cyclesChan         chan int
	clockMultiple      int
	ticksPast          int
	waitGroup          sync.WaitGroup
}

// NewArm7 constructs a new Arm7 CPU
func NewArm7(clockMultiple int) *Arm7 {
	return &Arm7{
		registers:       [37]uint32{},
		bankedRegisters: bankedRegMap[user],
		currentMode:     user,
		cyclesChan:      make(chan int, clockMultiple),
		clockMultiple:   clockMultiple,
	}
}

func (cpu *Arm7) readReg(register reg) uint32 {
	if register == rPc {
		return cpu.readPc()
	}

	if register < r8 {
		return cpu.registers[register]
	}

	return cpu.registers[cpu.bankedRegisters[register-8]]
}

func (cpu *Arm7) setReg(register reg, value uint32) {
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

func (cpu *Arm7) readPc() uint32 {
	return cpu.registers[rPc]
}

func (cpu *Arm7) setPc(value uint32) {
	cpu.registers[rPc] = value
}

func (cpu *Arm7) readCpsr() uint32 {
	return cpu.registers[rCpsr]
}

func (cpu *Arm7) setCpsr(value uint32) {
	cpu.registers[rCpsr] = value
}

func (cpu *Arm7) readSpsr() uint32 {
	return cpu.registers[cpu.currentSpsr]
}

func (cpu *Arm7) setSpsr(value uint32) {
	cpu.registers[cpu.currentSpsr] = value
}

// Tick notifies the CPU of a system clock tick.
func (cpu *Arm7) Tick() {
	for a := 0; a < cpu.clockMultiple; a++ {
		cpu.cyclesChan <- 1
	}
}

// WaitForTick blocks the caller till the CPU recieves a tick.
func (cpu *Arm7) WaitForTick() {
	if cpu.ticksPast == cpu.clockMultiple {
		cpu.waitGroup.Done()
		cpu.ticksPast = 0
	}
	<-cpu.cyclesChan
	cpu.ticksPast++
}

// Run starts the CPU execution loop
func (cpu *Arm7) Run() {
	cpu.WaitForTick()
}
