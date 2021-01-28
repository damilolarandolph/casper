package cpu

import (
	"github.com/damilolarandolph/casper/bits"
)

const (
	_MinInt int32 = -2147483648
	_MaxInt int32 = 2147483647
)

func signedSat(lhs int32, rhs int32, upperLimit int32, lowerLimit int32) (int32, bool) {
	result := lhs + rhs

	if (lhs < 0 && rhs > 0) && result >= 0 {
		return lowerLimit, true
	} else if (lhs > 0 && rhs > 0) && result <= 0 {
		return upperLimit, true
	}

	return result, false
}

/*
The QADD instruction performs integer addition,
saturating the result to the 32-bit signed integer range –2^31 ≤ x ≤ 2^31 – 1.
If saturation actually occurs, the instruction sets the Q flag in the CPSR.
*/
func qadd(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rn, rd, _ := parseDataInstr(instruction)
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := int32(cpu.readReg(rn))
	rmVal := int32(cpu.readReg(rm))
	result, didSat := signedSat(rmVal, rnVal, _MaxInt, _MinInt)
	cpu.setFlag(overflowSat, didSat)
	cpu.setReg(rd, uint32(result))
}

func signedSatM(lhs int32, rhs int32, upperLimit int32, lowerLimit int32) (int32, bool) {
	result := lhs - rhs

	if rhs > 0 && result < lhs {
		return upperLimit, true
	} else if rhs < 0 && result > lhs {
		return lowerLimit, true
	}

	return result, false
}

/*
The QDADD instruction doubles its second operand, then adds
the result to its first operand. Both the doubling and the
addition have their results saturated to the 32-bit signed
integer range –2^31≤ x ≤ 2^31–1. If saturation actually
occurs in either operation, the instruction sets the Q flag in the CPSR.
*/
func qdadd(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rn, rd, _ := parseDataInstr(instruction)
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := int32(cpu.readReg(rn))
	rmVal := int32(cpu.readReg(rm))
	result, didSat := signedSat(rnVal, rnVal, _MaxInt, _MinInt)
	accResult, accDidSat := signedSat(rmVal, result, _MaxInt, _MinInt)
	cpu.setFlag(overflowSat, didSat || accDidSat)
	cpu.setReg(rd, uint32(accResult))
}

/*
The QDSUB instruction doubles its second operand,
then subtracts the result from its first operand.
Both the doubling and the subtraction have their
results saturated to the 32-bit signed integer
range –2^31 ≤ x ≤ 2^31–1. If saturation actually
occurs in either operation, the instruction sets
the Q flag in the CPSR.
*/
func qdsub(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rn, rd, _ := parseDataInstr(instruction)
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := int32(cpu.readReg(rn))
	rmVal := int32(cpu.readReg(rm))
	result, didSat := signedSat(rnVal, rnVal, _MaxInt, _MinInt)
	accResult, accDidSat := signedSatM(rmVal, result, _MaxInt, _MinInt)
	cpu.setFlag(overflowSat, didSat || accDidSat)
	cpu.setReg(rd, uint32(accResult))
}

func qsub(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rn, rd, _ := parseDataInstr(instruction)
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := int32(cpu.readReg(rn))
	rmVal := int32(cpu.readReg(rm))
	result, didSat := signedSatM(rmVal, rnVal, _MaxInt, _MinInt)
	cpu.setFlag(overflowSat, didSat)
	cpu.setReg(rd, uint32(result))
}

func smlaXY(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 19, 16))
	rn := reg(bits.GetBits(instruction, 15, 12))
	rs := reg(bits.GetBits(instruction, 11, 8))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := cpu.readReg(rn)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	x := bits.GetBit(instruction, 5)
	y := bits.GetBit(instruction, 6)

	if !x {
		rmVal &= 0xffff
		rmVal <<= 16
		bits.ShiftRightSigned(&rmVal, 16)
	} else {
		rmVal &= 0xffff0000
		bits.ShiftRightSigned(&rmVal, 16)
	}

	if !y {
		rsVal &= 0xffff
		rsVal <<= 16
		bits.ShiftRightSigned(&rsVal, 16)
	} else {
		rsVal &= 0xffff0000
		bits.ShiftRightSigned(&rsVal, 16)
	}

	result := (int32(rmVal) * int32(rsVal)) + int32(rnVal)
	cpu.setReg(rd, uint32(result))
	if didSignOverflow(uint32(int32(rmVal)*int32(rsVal)), rnVal) {
		cpu.setFlag(overflowSat, true)
	}
}

func smlalXY(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 19, 16))
	rn := reg(bits.GetBits(instruction, 15, 12))
	rs := reg(bits.GetBits(instruction, 11, 8))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rdVal := cpu.readReg(rd)
	rnVal := cpu.readReg(rn)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	x := bits.GetBit(instruction, 5)
	y := bits.GetBit(instruction, 6)
	var operand1 uint32
	var operand2 uint32
	if !x {
		operand1 = rmVal & 0xffff
		operand1 <<= 16
		bits.ShiftRightSigned(&operand1, 16)
	} else {
		operand1 = rmVal & 0xffff0000
		bits.ShiftRightSigned(&operand1, 16)
	}

	if !y {
		operand2 = rsVal & 0xffff
		operand2 <<= 16
		bits.ShiftRightSigned(&operand2, 16)
	} else {
		operand2 = rsVal & 0xffff0000
		bits.ShiftRightSigned(&operand2, 16)
	}

	rdLowHi := (int64(rdVal) << 32 & int64(rnVal)) + (int64(int32(operand1)) * int64(int32(operand2)))
	cpu.setReg(rd, uint32(rdLowHi>>32))
	cpu.setReg(rn, uint32(rdLowHi))
}

func smlawY(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 19, 16))
	rn := reg(bits.GetBits(instruction, 15, 12))
	rs := reg(bits.GetBits(instruction, 11, 8))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rnVal := cpu.readReg(rn)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	y := bits.GetBit(instruction, 6)

	var operand2 uint32

	if y {
		operand2 = rsVal & 0xffff0000
		bits.ShiftRightSigned(&operand2, 16)
	} else {
		operand2 = rsVal & 0xffff
		operand2 <<= 16
		bits.ShiftRightSigned(&operand2, 16)
	}
	prod48 := ((int64(rmVal) * int64(operand2)) >> 16) & 0xffffffff
	result := int32(prod48) + int32(rnVal)
	cpu.setReg(rd, uint32(result))
	cpu.setFlag(overflowSat, didSignOverflow(uint32(prod48), uint32(rnVal)))
}

func smulXY(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 19, 16))
	rs := reg(bits.GetBits(instruction, 11, 8))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	x := bits.GetBit(instruction, 5)
	y := bits.GetBit(instruction, 6)
	var operand1 uint32
	var operand2 uint32
	if !x {
		operand1 = rmVal & 0xffff
		operand1 <<= 16
		bits.ShiftRightSigned(&operand1, 16)
	} else {
		operand1 = rmVal & 0xffff0000
		bits.ShiftRightSigned(&operand1, 16)
	}

	if !y {
		operand2 = rsVal & 0xffff
		operand2 <<= 16
		bits.ShiftRightSigned(&operand2, 16)
	} else {
		operand2 = rsVal & 0xffff0000
		bits.ShiftRightSigned(&operand2, 16)
	}

	cpu.setReg(rd, uint32(int32(operand1)*int32(operand2)))
}

func smulwY(cpu ArmCPU) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 19, 16))
	rs := reg(bits.GetBits(instruction, 11, 8))
	rm := reg(bits.GetBits(instruction, 3, 0))
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	y := bits.GetBit(instruction, 6)

	var operand2 uint32

	if y {
		operand2 = rsVal & 0xffff0000
		bits.ShiftRightSigned(&operand2, 16)
	} else {
		operand2 = rsVal & 0xffff
		operand2 <<= 16
		bits.ShiftRightSigned(&operand2, 16)
	}
	prod48 := ((int64(rmVal) * int64(operand2)) >> 16)
	result := uint32(prod48)
	cpu.setReg(rd, result)
}

func ldrd(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))

	if (rd%2) == 0 && (address&7) == 0 && rd != r14 {
		cpu.setReg(rd, cpu.Bus().ReadData32(address))
		cpu.setReg(rd+1, cpu.Bus().ReadData32(address))
	}

}

func strd(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))

	if (rd%2) == 0 && (address&7) == 0 && rd != r14 {
		cpu.Bus().WriteData32(address, cpu.readReg(rd))
		cpu.Bus().WriteData32(address+4, cpu.readReg(rd+1))
	}
}
