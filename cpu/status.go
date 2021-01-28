package cpu

type flag int
type cpuMode uint8

const (
	negative flag = 31 - iota
	zero
	carry
	overflow
	overflowSat
)

const (
	irqDisable flag = 7 - iota
	fiqDisable
	thumbMode
)

const (
	user       cpuMode = 0b10000
	fiq        cpuMode = 0b10001
	irq        cpuMode = 0b10010
	supervisor cpuMode = 0b10011
	abort      cpuMode = 0b10111
	undefined  cpuMode = 0b11011
	system     cpuMode = 0b11111
)

func (cpu *Arm7) isFlag(fl flag) bool {
	result := (cpu.readCpsr() >> fl) & 0x1
	return result == 1
}

func (cpu *Arm7) setFlag(fl flag, value bool) {

	var mask uint32 = 1 << fl

	if value {
		cpu.setCpsr(cpu.readCpsr() | mask)
		return
	}

	mask = ^mask
	cpu.setCpsr(cpu.readCpsr() & mask)
}

func (cpu *Arm7) mode() cpuMode {
	return cpuMode(cpu.readCpsr() & 0x1f)
}

func (cpu *Arm7) setMode(mode cpuMode) {
	var mask uint32
	mask = (^mask) | uint32(mode)
	cpu.setCpsr(cpu.readCpsr() & mask)
	cpu.bankedRegisters = bankedRegMap[mode]
	val, ok := bankedSpsrMap[mode]
	if ok {
		cpu.currentSpsr = val
	}
}
