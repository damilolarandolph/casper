package cpu

import "github.com/damilolarandolph/casper/bits"

func parseMult(cpu ArmCPU) (reg, reg, reg, reg) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 19, 16))
	rn := reg(bits.GetBits(instruction, 15, 12))
	rs := reg(bits.GetBits(instruction, 11, 8))
	rm := reg(instruction & 0xf)

	return rd, rn, rs, rm
}

func mla(cpu ArmCPU) {

	if !testCondition(cpu) {
		return
	}
	rd, rn, rs, rm := parseMult(cpu)
	rnVal := cpu.readReg(rn)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	result := (rmVal * rsVal) + rnVal
	cpu.setReg(rd, result)

	if !shouldSetCondition(cpu.currentInstruction()) {
		return
	}

	cpu.setFlag(negative, isNegative(result))
	cpu.setFlag(zero, result == 0)
}

func mul(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}
	rd, _, rs, rm := parseMult(cpu)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	result := rmVal * rsVal
	cpu.setReg(rd, result)
	if !shouldSetCondition(cpu.currentInstruction()) {
		return
	}

	cpu.setFlag(negative, isNegative(result))
	cpu.setFlag(zero, result == 0)

}

func smlal(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}

	rdHi, rdLo, rs, rm := parseMult(cpu)
	rdHiVal := cpu.readReg(rdHi)
	rdLoVal := cpu.readReg(rdLo)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	result64 := int64(rmVal) * int64(rsVal)
	rdComb := (int64(rdHiVal) << 32) & int64(rdLoVal)
	resultComb := rdComb + result64

	cpu.setReg(rdHi, uint32(resultComb>>32))
	cpu.setReg(rdLo, uint32(resultComb&0xffffffff))

	if !shouldSetCondition(cpu.currentInstruction()) {
		return
	}

	cpu.setFlag(negative, resultComb < 0)
	cpu.setFlag(zero, resultComb == 0)
}

func smull(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}

	rdHi, rdLo, rs, rm := parseMult(cpu)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	result64 := int64(rmVal) * int64(rsVal)

	cpu.setReg(rdHi, uint32(result64>>32))
	cpu.setReg(rdLo, uint32(result64&0xffffffff))

	if !shouldSetCondition(cpu.currentInstruction()) {
		return
	}

	cpu.setFlag(negative, result64 < 0)
	cpu.setFlag(zero, result64 == 0)
}

func umlal(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}

	rdHi, rdLo, rs, rm := parseMult(cpu)
	rdHiVal := cpu.readReg(rdHi)
	rdLoVal := cpu.readReg(rdLo)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	result64 := uint64(rmVal) * uint64(rsVal)
	rdComb := (uint64(rdHiVal) << 32) & uint64(rdLoVal)
	resultComb := rdComb + result64

	cpu.setReg(rdHi, uint32(resultComb>>32))
	cpu.setReg(rdLo, uint32(resultComb&0xffffffff))

	if !shouldSetCondition(cpu.currentInstruction()) {
		return
	}

	cpu.setFlag(negative, (resultComb>>63) != 0)
	cpu.setFlag(zero, resultComb == 0)
}

func umull(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}

	rdHi, rdLo, rs, rm := parseMult(cpu)
	rsVal := cpu.readReg(rs)
	rmVal := cpu.readReg(rm)
	result64 := uint64(rmVal) * uint64(rsVal)

	cpu.setReg(rdHi, uint32(result64>>32))
	cpu.setReg(rdLo, uint32(result64&0xffffffff))

	if !shouldSetCondition(cpu.currentInstruction()) {
		return
	}

	cpu.setFlag(negative, (result64>>63) != 0)
	cpu.setFlag(zero, result64 == 0)
}
