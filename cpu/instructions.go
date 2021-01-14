package cpu

import "github.com/damilolarandolph/casper/bits"

var armInstructions = [0xff][0xf]instructionHandler{}

type condition func(cpu *CPU) bool

var conditions = []condition{
	//Equal to flag - EQ
	func(cpu *CPU) bool { return cpu.isFlag(zero) },

	//not equal to flag - NE
	func(cpu *CPU) bool { return !cpu.isFlag(zero) },

	//unsigned carry or same - CS/HS
	func(cpu *CPU) bool { return cpu.isFlag(carry) },

	//unsigned lower - CC/LO
	func(cpu *CPU) bool { return !cpu.isFlag(carry) },

	//signed negative - MI
	func(cpu *CPU) bool { return cpu.isFlag(negative) },

	//signed positive - PL
	func(cpu *CPU) bool { return !cpu.isFlag(negative) },

	//signed overflow - VS
	func(cpu *CPU) bool { return cpu.isFlag(overflow) },

	//signed no overflow - VC
	func(cpu *CPU) bool { return !cpu.isFlag(overflow) },

	//unsiged higher - HI
	func(cpu *CPU) bool { return cpu.isFlag(carry) && !cpu.isFlag(zero) },

	//unsigned lower or same - LS
	func(cpu *CPU) bool { return !cpu.isFlag(carry) && cpu.isFlag(zero) },

	//signed greater or equal - GE
	func(cpu *CPU) bool { return cpu.isFlag(negative) == cpu.isFlag(overflow) },

	//signed less than - LT
	func(cpu *CPU) bool { return cpu.isFlag(negative) != cpu.isFlag(overflow) },

	//signed greater than - GT
	func(cpu *CPU) bool { return !cpu.isFlag(zero) && (cpu.isFlag(negative) == cpu.isFlag(overflow)) },

	//signed less or equal - LE
	func(cpu *CPU) bool { return !cpu.isFlag(zero) || (cpu.isFlag(negative) != cpu.isFlag(overflow)) },

	//always - AL
	func(cpu *CPU) bool { return true },

	//never - NV
	func(cpu *CPU) bool { return false },
}

func populateInstruction(rows ...int) func(...int) func(instructionHandler) {
	return func(cols ...int) func(instructionHandler) {
		return func(ih instructionHandler) {

		}
	}
}

func testCondition(cpu *CPU) bool {
	instruction := cpu.currentInstruction
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
	populateInstruction(1, 2, 3)(1, 2, 3)(func(cpu *CPU) {})
}
