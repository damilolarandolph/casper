package cpu

import "github.com/damilolarandolph/casper/bits"

// Control Extensions

func clz(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}

	var leadingZeros uint32 = 32

	sourceReg := reg(cpu.currentInstruction() & 0xf)
	sourceVal := cpu.readReg(sourceReg)
	destReg := reg(bits.GetBits(cpu.currentInstruction(), 15, 12))

	for sourceVal != 0 {
		sourceVal >>= 1
		leadingZeros--
	}
	cpu.setReg(destReg, leadingZeros)
}

func swp(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	rd := reg(bits.GetBits(instruction, 15, 12))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := cpu.readReg(rn)
	rnBits := int(rnVal & 3)
	temp := cpu.Bus().ReadData32(rnVal)
	if rnBits > 0 {
		rotateRight(&temp, rnBits*8, false)
	}
	cpu.Bus().WriteData32(rnVal, cpu.readReg(rm))
	cpu.setReg(rd, temp)
}

func swpb(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	rd := reg(bits.GetBits(instruction, 15, 12))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := cpu.readReg(rn)
	temp := cpu.Bus().ReadData8(rnVal)
	cpu.Bus().WriteData8(rnVal, cpu.readReg(rm))
	cpu.setReg(rd, temp)
}
