package cpu

import "github.com/damilolarandolph/casper/bits"

// Branch or Branch with Link Instruction.
// Linking the return address to R14 is determined by
// the withLink parameter.
func branch(withLink bool, cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}
	if bits.GetBits(cpu.currentInstruction(), 31, 28) == 0xf {
		blxImm(cpu)
		return
	}
	branchAddress := bits.GetBits(cpu.currentInstruction(), 23, 0)
	bits.ShiftRightSigned(&branchAddress, 12)
	branchAddress <<= 2
	if withLink {
		cpu.setReg(r14, cpu.nextInstruction())
	}
	var newPc uint32 = uint32(int32(cpu.readPc()) + int32(branchAddress))
	cpu.setPc(newPc)
}

// Branch Link Exchange instruction with immediate operand.
func blxImm(cpu ArmCPU) {
	branchAddress := bits.GetBits(cpu.currentInstruction(), 23, 0)
	bits.ShiftRightSigned(&branchAddress, 12)
	branchAddress <<= 2
	bits.SetBit(&branchAddress, 1, bits.GetBit(cpu.currentInstruction(), 24))
	cpu.setFlag(thumbMode, true)
	var newPc uint32 = uint32(int32(cpu.readPc()) + int32(branchAddress))
	cpu.setPc(newPc)
}

// Branch Link Exchange instruction with register operand.
func blxReg(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}
	branchReg := reg(cpu.currentInstruction() & 0xf)
	branchAddress := cpu.readReg(branchReg)
	cpu.setReg(r14, cpu.nextInstruction())
	cpu.setFlag(thumbMode, (branchAddress&0x1) != 0)
	branchAddress &= 0xfffffffe
	cpu.setPc(branchAddress)
}

// Determines which branch link exchange instruction to run.
// This is done by inspecting a region of the instruction and
// comparing it to a known constant.
func blx(cpu ArmCPU) {
	var blxRegConstant uint32 = 0x12fff
	if bits.GetBits(cpu.currentInstruction(), 27, 8) == blxRegConstant {
		blxReg(cpu)
		return
	}
	blxImm(cpu)
}

func bx(cpu ArmCPU) {
	if !testCondition(cpu) {
		return
	}
	// The branch address register is stored at bits 3 - 0
	// of the instruction
	branchReg := reg(cpu.currentInstruction() & 0xf)
	branchAddress := cpu.readReg(branchReg)
	cpu.setFlag(thumbMode, (branchAddress&0x1) != 0)
	branchAddress &= 0xfffffffe
	cpu.setPc(branchAddress)

}

var br = func(cpu ArmCPU) {

	branch(false, cpu)
}

var bl = func(cpu ArmCPU) {
	branch(true, cpu)
}
