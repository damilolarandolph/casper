package cpu

import "github.com/damilolarandolph/casper/bits"

func mrs(cpu ArmCPU) {

	if !testCondition(cpu) {
		return
	}

	isSPSR := bits.GetBit(cpu.currentInstruction(), 22)
	rd := reg(bits.GetBits(cpu.currentInstruction(), 15, 12))
	if isSPSR {
		cpu.setReg(rd, cpu.readSpsr())
		return
	}

	cpu.setReg(rd, cpu.readCpsr())
}

func msr(cpu ArmCPU, addressingMode arthAddrMode) {

	if !testCondition(cpu) {
		return
	}

	if cpu.mode() == user {
		return
	}

	var result uint32
	result, _ = addressingMode(cpu)
	fieldMask := bits.GetBits(cpu.currentInstruction(), 19, 16)
	var maskStart int
	var maskEnd int

	if bits.GetBit(fieldMask, 0) {
		maskStart = 7
		maskEnd = 0
	} else if bits.GetBit(fieldMask, 1) {
		maskStart = 15
		maskEnd = 8
	} else if bits.GetBit(fieldMask, 2) {
		maskStart = 23
		maskEnd = 16
	} else if bits.GetBit(fieldMask, 3) {
		maskStart = 31
		maskEnd = 24
	} else {
		return
	}

	isSPSR := bits.GetBit(cpu.currentInstruction(), 22)

	if !isSPSR {
		current := cpu.readCpsr()
		newVal := bits.GetBits(result, maskStart, maskEnd)
		newVal <<= maskStart
		current |= newVal
		current &= newVal
		cpu.setCpsr(current)
		return
	}

	if cpu.mode() == system {
		return
	}

	current := cpu.readSpsr()
	newVal := bits.GetBits(result, maskStart, maskEnd)
	newVal <<= maskStart
	current |= newVal
	current &= newVal
	cpu.setSpsr(current)
}
