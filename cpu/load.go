package cpu

import (
	"fmt"

	"github.com/damilolarandolph/casper/bits"
)

type addressingIndex int

const (
	noIndexing addressingIndex = iota
	preIndexing
	postIndexing
)

type indexFunc func(uBit bool, condsPassed bool, rnVal uint32, rhs uint32) (uint32, uint32)

func preIndexFunc(uBit bool, condsPassed bool, rnVal uint32, rhs uint32) (uint32, uint32) {
	var address uint32
	if uBit {
		address = rnVal + rhs
	} else {
		address = rnVal - rhs
	}

	if condsPassed {
		return address, address
	}
	return address, rnVal
}

func postIndexFunc(uBit bool, condsPassed bool, rnVal uint32, rhs uint32) (uint32, uint32) {
	var address uint32
	if uBit {
		address = rnVal + rhs
	} else {
		address = rnVal - rhs
	}

	if condsPassed {
		return rnVal, address
	}
	return address, address
}

func immOffLdStr(cpu ArmCPU, indexing indexFunc) uint32 {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	rnVal := cpu.readReg(rn)
	offset12 := bits.GetBits(instruction, 11, 0)
	uBit := bits.GetBit(instruction, 23)
	address, newRnVal := indexing(uBit, testCondition(cpu), rnVal, offset12)
	if bits.GetBit(instruction, 21) {
		cpu.setReg(rn, newRnVal)
	}
	return address
}

func regOffLdStr(cpu ArmCPU, indexing indexFunc) uint32 {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	rnVal := cpu.readReg(rn)

	uBit := bits.GetBit(instruction, 23)
	rm := reg(bits.GetBits(instruction, 3, 0))
	rmVal := cpu.readReg(rm)
	address, newRnVal := indexing(uBit, testCondition(cpu), rnVal, rmVal)
	if bits.GetBit(instruction, 21) {
		cpu.setReg(rn, newRnVal)
	}
	return address
}

func scaledRegOff(cpu ArmCPU, indexing indexFunc) uint32 {
	instruction := cpu.currentInstruction()
	shift := shiftType(bits.GetBits(instruction, 6, 5))
	rn := reg(bits.GetBits(instruction, 19, 16))
	rnVal := cpu.readReg(rn)
	var index uint32
	if shift == lsl {
		index, _ = lli(cpu)
	} else if shift == lsr {
		index, _ = lri(cpu)
	} else if shift == asr {
		index, _ = ari(cpu)
	} else if shift == ror {
		index, _ = rri(cpu)
	}

	uBit := bits.GetBit(instruction, 23)
	address, newRnVal := indexing(uBit, testCondition(cpu), rnVal, index)
	if bits.GetBit(instruction, 21) {
		cpu.setReg(rn, newRnVal)
	}
	return address
}

func miscImmOffLdStr(cpu ArmCPU, indexing indexFunc) uint32 {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	rnVal := cpu.readReg(rn)
	offset8 := (bits.GetBits(instruction, 11, 8) << 4) | bits.GetBits(instruction, 3, 0)
	uBit := bits.GetBit(instruction, 23)
	address, newRnVal := indexing(uBit, testCondition(cpu), rnVal, offset8)
	if bits.GetBit(instruction, 21) {
		cpu.setReg(rn, newRnVal)
	}
	return address
}

func ldStrMullIncAfter(cpu ArmCPU) (uint32, uint32) {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	regList := bits.GetBits(instruction, 15, 0)
	wBit := bits.GetBit(instruction, 21)
	rnVal := cpu.readReg(rn)

	endAddress := rnVal + (uint32(bits.NumSetBits(regList)) * 4) - 4

	if testCondition(cpu) && wBit {
		cpu.setReg(rn, rnVal+(uint32(bits.NumSetBits(regList))*4))
	}
	return rnVal, endAddress
}
func ldStrMullIncBefore(cpu ArmCPU) (uint32, uint32) {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	regList := bits.GetBits(instruction, 15, 0)
	wBit := bits.GetBit(instruction, 21)
	rnVal := cpu.readReg(rn)
	startAddress := rnVal + 4
	endAddress := rnVal + (uint32(bits.NumSetBits(regList)) * 4)

	if testCondition(cpu) && wBit {
		cpu.setReg(rn, rnVal+(uint32(bits.NumSetBits(regList))*4))
	}
	return startAddress, endAddress
}

func ldStrMullDecAfter(cpu ArmCPU) (uint32, uint32) {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	regList := bits.GetBits(instruction, 15, 0)
	wBit := bits.GetBit(instruction, 21)
	rnVal := cpu.readReg(rn)
	startAddress := rnVal - (uint32(bits.NumSetBits(regList)) * 4) + 4
	endAddress := rnVal

	if testCondition(cpu) && wBit {
		cpu.setReg(rn, rnVal-(uint32(bits.NumSetBits(regList))*4))
	}
	return startAddress, endAddress
}
func ldStrMullDecBefore(cpu ArmCPU) (uint32, uint32) {
	instruction := cpu.currentInstruction()
	rn := reg(bits.GetBits(instruction, 19, 16))
	regList := bits.GetBits(instruction, 15, 0)
	wBit := bits.GetBit(instruction, 21)
	rnVal := cpu.readReg(rn)
	startAddress := rnVal - (uint32(bits.NumSetBits(regList)) * 4)
	endAddress := rnVal - 4

	if testCondition(cpu) && wBit {
		cpu.setReg(rn, rnVal-(uint32(bits.NumSetBits(regList))*4))
	}
	return startAddress, endAddress
}

func ldr(cpu ArmCPU, address uint32) {
	endBits := address & 0x3
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	var value uint32
	cpu.Bus().SetSequencial(false)
	if endBits == 0 {
		value = cpu.Bus().ReadData32(address)
	} else {
		value = cpu.Bus().ReadData32(address)
		rotateRight(&value, int(endBits*8), false)
	}

	if rd == rPc {
		if cpu.Architecture() >= V5 {
			cpu.setPc(value & 0xFFFFFFFE)
			cpu.setFlag(thumbMode, (value&0x1) != 0)
			return
		}
		cpu.setPc(value & 0xFFFFFFFC)
		return
	}
	cpu.setReg(rd, value)

}
func ldrt(cpu ArmCPU, address uint32) {
	endBits := address & 0x3
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	var value uint32
	cpu.Bus().SetSequencial(false)
	currentMode := cpu.mode()
	cpu.setMode(user)
	if endBits == 0 {
		value = cpu.Bus().ReadData32(address)
	} else {
		value = cpu.Bus().ReadData32(address)
		rotateRight(&value, int(endBits*8), false)
	}
	cpu.setMode(currentMode)
	cpu.setReg(rd, value)

}

func ldrb(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	cpu.setReg(rd, cpu.Bus().ReadData8(address))
}

func ldrbt(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	currentMode := cpu.mode()
	cpu.setMode(user)
	cpu.Bus().SetSequencial(false)
	cpu.setReg(rd, cpu.Bus().ReadData8(address))
	cpu.setMode(currentMode)
}

func ldrh(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	cpu.setReg(rd, cpu.Bus().ReadData16(address))
}

func ldrsb(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	data := cpu.Bus().ReadData8(address)
	data <<= 32 - 8
	bits.ShiftRightSigned(&data, 32-8)
	cpu.setReg(rd, data)
}
func ldrsh(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	data := cpu.Bus().ReadData8(address)
	data <<= 32 - 16
	bits.ShiftRightSigned(&data, 32-16)
	cpu.setReg(rd, data)
}
func str(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	cpu.Bus().WriteData32(address, cpu.readReg(rd))
}

func strt(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	currentMode := cpu.mode()
	cpu.setMode(currentMode)
	cpu.Bus().SetSequencial(false)
	cpu.Bus().WriteData32(address, cpu.readReg(rd))
	cpu.setMode(user)
}

func strb(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	cpu.Bus().WriteData8(address, cpu.readReg(rd))
}

func strbt(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	currentMode := cpu.mode()
	cpu.setMode(user)
	cpu.Bus().SetSequencial(false)
	cpu.Bus().WriteData8(address, cpu.readReg(rd))
	cpu.setMode(currentMode)
}

func strh(cpu ArmCPU, address uint32) {
	instruction := cpu.currentInstruction()
	rd := reg(bits.GetBits(instruction, 15, 12))
	cpu.Bus().SetSequencial(false)
	cpu.Bus().WriteData16(address, cpu.readReg(rd))
}

func ldm(cpu ArmCPU, startAddress uint32, endAddress uint32) {
	instruction := cpu.currentInstruction()
	regList := bits.GetBits(instruction, 15, 0)

	cpu.Bus().SetSequencial(false)
	for a := 0; a < 15; a++ {
		bitSet := (regList & 0x1) != 0
		if bitSet {
			cpu.setReg(reg(a), cpu.Bus().ReadData32(startAddress))
			startAddress += 4
			regList >>= 1
		}

	}
	if (regList >> 14) != 0 {
		value := cpu.Bus().ReadData32(startAddress)
		if cpu.Architecture() >= V5 {
			cpu.setPc(value & 0xFFFFFFFE)
			cpu.setFlag(thumbMode, (value&0x1) != 0)
		} else {
			cpu.setPc(value & 0xFFFFFFFC)
		}
		startAddress += 4
	}

	if endAddress != (startAddress - 4) {
		fmt.Println("LDM Assert failed")
	}

}

func ldmUser(cpu ArmCPU, startAddress uint32, endAddress uint32) {
	instruction := cpu.currentInstruction()
	regList := bits.GetBits(instruction, 15, 0)
	if bits.GetBit(regList, 15) {
		ldmSPSR(cpu, startAddress, endAddress)
		return
	}
	currentMode := cpu.mode()
	cpu.setMode(user)
	cpu.Bus().SetSequencial(false)
	for a := 0; a < 15; a++ {
		bitSet := (regList & 0x1) != 0
		if bitSet {
			cpu.setReg(reg(a), cpu.Bus().ReadData32(startAddress))
			startAddress += 4
			regList >>= 1
		}

	}
	cpu.setMode(currentMode)
}

func ldmSPSR(cpu ArmCPU, startAddress uint32, endAddress uint32) {
	instruction := cpu.currentInstruction()
	regList := bits.GetBits(instruction, 15, 0)

	cpu.Bus().SetSequencial(false)
	for a := 0; a < 15; a++ {
		bitSet := (regList & 0x1) != 0
		if bitSet {
			cpu.setReg(reg(a), cpu.Bus().ReadData32(startAddress))
			startAddress += 4
			regList >>= 1
		}
	}
	cpu.setCpsr(cpu.readSpsr())

	value := cpu.Bus().ReadData32(startAddress)

	if cpu.Architecture() > V4 && cpu.isFlag(thumbMode) {
		cpu.setPc(value & 0xFFFFFFFE)
	} else {
		cpu.setPc(value & 0xFFFFFFFC)
	}
	startAddress = startAddress + 4
}

func stm(cpu ArmCPU, startAddress uint32, endAddress uint32) {
	instruction := cpu.currentInstruction()
	regList := bits.GetBits(instruction, 15, 0)

	cpu.Bus().SetSequencial(false)
	for a := 0; a < 15; a++ {
		bitSet := (regList & 0x1) != 0
		if bitSet {
			cpu.Bus().WriteData32(startAddress, cpu.readReg(reg(a)))
			startAddress += 4
			regList >>= 1
		}
	}
}

func stmUser(cpu ArmCPU, startAddress uint32, endAddress uint32) {
	currentMode := cpu.mode()
	cpu.setMode(user)
	instruction := cpu.currentInstruction()
	regList := bits.GetBits(instruction, 15, 0)

	cpu.Bus().SetSequencial(false)
	for a := 0; a < 15; a++ {
		bitSet := (regList & 0x1) != 0
		if bitSet {
			cpu.Bus().WriteData32(startAddress, cpu.readReg(reg(a)))
			startAddress += 4
			regList >>= 1
		}
	}

	cpu.setMode(currentMode)
}
