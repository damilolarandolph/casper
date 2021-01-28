package cpu

import "github.com/damilolarandolph/casper/bits"

type instructionHandler func(cpu ArmCPU)

var armInstructions = [0xff + 1][0xf + 1]instructionHandler{}

type condition func(cpu ArmCPU) bool

var conditions = []condition{
	//Equal to flag - EQ
	func(cpu ArmCPU) bool { return cpu.isFlag(zero) },

	//not equal to flag - NE
	func(cpu ArmCPU) bool { return !cpu.isFlag(zero) },

	//unsigned carry or same - CS/HS
	func(cpu ArmCPU) bool { return cpu.isFlag(carry) },

	//unsigned lower - CC/LO
	func(cpu ArmCPU) bool { return !cpu.isFlag(carry) },

	//signed negative - MI
	func(cpu ArmCPU) bool { return cpu.isFlag(negative) },

	//signed positive - PL
	func(cpu ArmCPU) bool { return !cpu.isFlag(negative) },

	//signed overflow - VS
	func(cpu ArmCPU) bool { return cpu.isFlag(overflow) },

	//signed no overflow - VC
	func(cpu ArmCPU) bool { return !cpu.isFlag(overflow) },

	//unsiged higher - HI
	func(cpu ArmCPU) bool { return cpu.isFlag(carry) && !cpu.isFlag(zero) },

	//unsigned lower or same - LS
	func(cpu ArmCPU) bool { return !cpu.isFlag(carry) && cpu.isFlag(zero) },

	//signed greater or equal - GE
	func(cpu ArmCPU) bool { return cpu.isFlag(negative) == cpu.isFlag(overflow) },

	//signed less than - LT
	func(cpu ArmCPU) bool { return cpu.isFlag(negative) != cpu.isFlag(overflow) },

	//signed greater than - GT
	func(cpu ArmCPU) bool { return !cpu.isFlag(zero) && (cpu.isFlag(negative) == cpu.isFlag(overflow)) },

	//signed less or equal - LE
	func(cpu ArmCPU) bool { return !cpu.isFlag(zero) || (cpu.isFlag(negative) != cpu.isFlag(overflow)) },

	//always - AL
	func(cpu ArmCPU) bool { return true },

	//never - NV
	func(cpu ArmCPU) bool { return false },
}

func populateInstruction(rows ...int) func(...int) func(instructionHandler) {
	return func(cols ...int) func(instructionHandler) {
		return func(ih instructionHandler) {

		}
	}
}

func testCondition(cpu ArmCPU) bool {
	instruction := cpu.currentInstruction()
	return conditions[bits.GetBits(instruction, 31, 28)](cpu)
}

func isImmediate(instruction uint32) bool {
	return bits.GetBit(instruction, 25)
}

func shouldSetCondition(instruction uint32) bool {
	return bits.GetBit(instruction, 20)
}

func getSourceReg(instruction uint32) reg {
	return reg(bits.GetBits(instruction, 19, 16))
}

func getDestinationReg(instruction uint32) reg {
	return reg(bits.GetBits(instruction, 15, 12))
}

func init() {
	populateInstruction(1, 2, 3)(1, 2, 3)(func(cpu ArmCPU) {})
}
