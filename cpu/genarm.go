package cpu

import (
	"fmt"
	"os"

	"github.com/damilolarandolph/casper/bits"
)

type Generator struct {
	rows        int
	cols        int
	maskTargets []MaskTarget
}

type MaskHandler func(row int, col int) bool

func NewGenerator(rows int, cols int) *Generator {
	return &Generator{
		rows: rows,
		cols: cols,
	}

}

type MaskTarget struct {
	rowMask string
	colMask string
	name    string
	handler MaskHandler
}

func (gen *Generator) start() {
	foundLog, _ := os.Create("foundlog.txt")
	notFoundLog, _ := os.Create("notfoundlog.txt")
	defer foundLog.Close()
	defer notFoundLog.Close()
	foundHandlers := 0
	for row := 0; row <= gen.rows; row++ {
		var found bool
		for col := 0; col <= gen.cols; col++ {

			for _, item := range gen.maskTargets {
				found = item.tryRun(row, col)
				if found {
					foundHandlers++
					fmt.Fprintln(foundLog, "Found handler", item.name, "for row", row, "coloumn", col)
					break
				}
			}

			if !found {

				fmt.Fprintln(notFoundLog, "No handler for row", row, "column", col)
			}
		}
	}

	fmt.Println("Found: ", foundHandlers, "Didn't find: ", (gen.rows*gen.cols)-foundHandlers)
}

func (target *MaskTarget) tryRun(row int, col int) bool {
	orginalRow := row
	originalCol := col
	for index := len(target.rowMask) - 1; index >= 0; index-- {

		if target.rowMask[index] == '0' && (row&0x1) != 0 {
			return false
		}

		if target.rowMask[index] == '1' && (row&0x1) != 1 {
			return false
		}

		row >>= 1
	}

	for index := len(target.colMask) - 1; index >= 0; index-- {

		if target.colMask[index] == '0' && (col&0x1) != 0 {
			return false
		}

		if target.colMask[index] == '1' && (col&0x1) != 1 {
			return false
		}

		col >>= 1
	}

	return target.handler(orginalRow, originalCol)

}

func (gen *Generator) AddMaskTarget(target MaskTarget) {
	gen.maskTargets = append(gen.maskTargets, target)
}

func emitDataOpcode(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	if bits.GetBits(opcodeFrag, 27, 26) != 0 {
		return false
	}
	isImmediate := bits.GetBit(opcodeFrag, 25)

	aluCode := bits.GetBits(opcodeFrag, 24, 21)
	var aluFunc func(addressingMode arthAddrMode, cpu ArmCPU)

	switch aluCode {
	case 0:
		aluFunc = and
	case 1:
		aluFunc = eor
	case 2:
		aluFunc = sub
	case 3:
		aluFunc = rsb
	case 4:
		aluFunc = add
	case 5:
		aluFunc = addC
	case 6:
		aluFunc = sbc
	case 7:
		aluFunc = rsc
	case 8:
		aluFunc = tst
	case 9:
		aluFunc = teq
	case 0xa:
		aluFunc = cmp
	case 0xb:
		aluFunc = cmn
	case 0xc:
		aluFunc = orr
	case 0xd:
		aluFunc = mov
	case 0xe:
		aluFunc = bic
	case 0xf:
		aluFunc = mvn
	}

	setConds := bits.GetBit(opcodeFrag, 20)

	if !setConds && (aluCode >= 8 && aluCode <= 0xb) {
		return false
	}

	if isImmediate {
		armInstructions[row][col] = func(cpu ArmCPU) {
			aluFunc(imm, cpu)
		}

		return true
	}

	shiftByRegister := bits.GetBit(opcodeFrag, 4)
	if shiftByRegister {
		if bits.GetBit(opcodeFrag, 7) {
			return false
		}
	}
	shift := shiftType(bits.GetBits(opcodeFrag, 6, 5))
	var addressMode arthAddrMode
	if shift == lsl {
		if !shiftByRegister {
			addressMode = lli
		} else {
			addressMode = llr
		}
	} else if shift == lsr {
		if !shiftByRegister {
			addressMode = lri
		} else {
			addressMode = lrr
		}
	} else if shift == asr {
		if !shiftByRegister {
			addressMode = ari
		} else {
			addressMode = arr
		}
	} else if shift == ror {
		if !shiftByRegister {
			addressMode = rri
		} else {
			addressMode = rrr
		}
	} else {
		return false
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		aluFunc(addressMode, cpu)
	}

	return true

}

func getAluFuc(code uint32) func(addressingMode arthAddrMode, cpu ArmCPU) {
	switch code {
	case 0:
		return and
	case 1:
		return eor
	case 2:
		return sub
	case 3:
		return rsb
	case 4:
		return add
	case 5:
		return addC
	case 6:
		return sbc
	case 7:
		return rsc
	case 8:
		return tst
	case 9:
		return teq
	case 0xa:
		return cmp
	case 0xb:
		return cmn
	case 0xc:
		return orr
	case 0xd:
		return mov
	case 0xe:
		return bic
	case 0xf:
		return mvn
	default:
		return nil
	}

}

func emitDataProcessImmShift(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	aluCode := bits.GetBits(opcodeFrag, 24, 21)
	var aluFunc func(addressingMode arthAddrMode, cpu ArmCPU)

	aluFunc = getAluFuc(aluCode)
	setConds := bits.GetBit(opcodeFrag, 20)

	if !setConds && (aluCode >= 8 && aluCode <= 0xb) {
		return false
	}
	shift := shiftType(bits.GetBits(opcodeFrag, 6, 5))
	var addressMode arthAddrMode

	if shift == lsl {
		addressMode = lli
	} else if shift == lsr {
		addressMode = lri
	} else if shift == asr {
		addressMode = ari
	} else if shift == ror {
		addressMode = rri
	} else {
		return false
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		aluFunc(addressMode, cpu)
	}

	return true
}

func emitDataImm(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	aluCode := bits.GetBits(opcodeFrag, 24, 21)
	var aluFunc func(addressingMode arthAddrMode, cpu ArmCPU)

	aluFunc = getAluFuc(aluCode)
	setConds := bits.GetBit(opcodeFrag, 20)

	if !setConds && (aluCode >= 8 && aluCode <= 0xb) {
		return false
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		aluFunc(imm, cpu)
	}

	return true
}

func emitDataRegShift(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	aluCode := bits.GetBits(opcodeFrag, 24, 21)
	var aluFunc func(addressingMode arthAddrMode, cpu ArmCPU)

	aluFunc = getAluFuc(aluCode)
	setConds := bits.GetBit(opcodeFrag, 20)

	if !setConds && (aluCode >= 8 && aluCode <= 0xb) {
		return false
	}
	shift := shiftType(bits.GetBits(opcodeFrag, 6, 5))
	var addressMode arthAddrMode

	if shift == lsl {
		addressMode = llr
	} else if shift == lsr {
		addressMode = lrr
	} else if shift == asr {
		addressMode = arr
	} else if shift == ror {
		addressMode = rrr
	} else {
		return false
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		aluFunc(addressMode, cpu)
	}

	return true
}

func emitMRSOpcode(row int, col int) bool {
	armInstructions[row][col] = mrs
	return true
}
func emitMSROpcode(row int, col int) bool {

	armInstructions[row][col] = func(cpu ArmCPU) {
		msr(cpu, lli)
	}
	return true
}

func emitMSRImmOpcode(row int, col int) bool {

	armInstructions[row][col] = func(cpu ArmCPU) {
		msr(cpu, imm)
	}
	return true
}

func emitBranchLinkOpcode(row int, col int) bool {

	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	opcode := bits.GetBit(opcodeFrag, 24)

	if opcode {
		armInstructions[row][col] = func(cpu ArmCPU) {
			condition := bits.GetBits(cpu.currentInstruction(), 31, 28)
			if condition == 0xf {
				blxImm(cpu)
				return
			}
			branch(false, cpu)
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			condition := bits.GetBits(cpu.currentInstruction(), 31, 28)
			if condition == 0xf {
				blxImm(cpu)
				return
			}
			branch(true, cpu)
		}
	}
	return true
}

func emitBranchExchangeOpcode(row int, col int) bool {

	armInstructions[row][col] = bx

	return true
}

func emitCLZOpcode(row int, col int) bool {

	armInstructions[row][col] = clz

	return true
}

func emitBLXRegOpcode(row int, col int) bool {
	armInstructions[row][col] = blxReg
	return true
}

func emitDSPAddSub(row int, col int) bool {

	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	op := bits.GetBits(opcodeFrag, 23, 20)
	if op == 0 {
		armInstructions[row][col] = qadd
	} else if op == 2 {
		armInstructions[row][col] = qsub
	} else if op == 4 {
		armInstructions[row][col] = qdadd
	} else if op == 6 {
		armInstructions[row][col] = qdsub
	} else {
		return false
	}

	return true
}
func emitDSPMultiplyOpcodes(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	op := bits.GetBits(opcodeFrag, 24, 21)

	switch op {
	case 8:
		armInstructions[row][col] = smlaXY
	case 9:
		if bits.GetBit(opcodeFrag, 5) {
			armInstructions[row][col] = smulwY
		} else {
			armInstructions[row][col] = smlawY
		}

	case 10:
		armInstructions[row][col] = smlalXY
	case 11:
		armInstructions[row][col] = smulXY
	default:
		return false
	}

	return true
}

func emitMultiplyAccumulate(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	op := bits.GetBits(opcodeFrag, 24, 21)

	if op == 0 {
		armInstructions[row][col] = mul
	} else if op == 1 {
		armInstructions[row][col] = mla
	} else {
		return false
	}

	return true

}

func emitMultiplyAccumulateLong(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	op := bits.GetBits(opcodeFrag, 24, 21)

	if op == 4 {
		armInstructions[row][col] = umull
	} else if op == 5 {
		armInstructions[row][col] = umlal
	} else if op == 6 {
		armInstructions[row][col] = smull
	} else if op == 7 {
		armInstructions[row][col] = smlal
	} else {
		return false
	}

	return true

}

func emitSwapOpcodes(row int, col int) bool {

	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	byteBit := bits.GetBit(opcodeFrag, 22)

	if byteBit {
		armInstructions[row][col] = swpb
	} else {
		armInstructions[row][col] = swp
	}

	return true

}

func emitLdrStrHReg(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	L := bits.GetBit(opcodeFrag, 20)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if L {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrh(cpu, scaledRegOff(cpu, indexingFunc))
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			strh(cpu, scaledRegOff(cpu, indexingFunc))
		}
	}

	return true
}

func emitLdrStrHImm(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	L := bits.GetBit(opcodeFrag, 20)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if L {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrh(cpu, miscImmOffLdStr(cpu, indexingFunc))
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			strh(cpu, miscImmOffLdStr(cpu, indexingFunc))
		}
	}

	return true
}
func emitLdrStrDReg(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	S := bits.GetBit(opcodeFrag, 5)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if !S {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrd(cpu, scaledRegOff(cpu, indexingFunc))
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			strd(cpu, scaledRegOff(cpu, indexingFunc))
		}
	}

	return true

}

func emitLdrStrDImm(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	S := bits.GetBit(opcodeFrag, 5)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if !S {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrd(cpu, miscImmOffLdStr(cpu, indexingFunc))
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			strd(cpu, miscImmOffLdStr(cpu, indexingFunc))
		}
	}

	return true

}

func emitLdrSBReg(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	H := bits.GetBit(opcodeFrag, 5)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}
	if H {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrsh(cpu, scaledRegOff(cpu, indexingFunc))
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrsb(cpu, scaledRegOff(cpu, indexingFunc))
		}
	}

	return true
}

func emitLdrSBImm(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	H := bits.GetBit(opcodeFrag, 5)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}
	if H {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrsh(cpu, miscImmOffLdStr(cpu, indexingFunc))
		}
	} else {
		armInstructions[row][col] = func(cpu ArmCPU) {
			ldrsb(cpu, miscImmOffLdStr(cpu, indexingFunc))
		}
	}

	return true
}

func emitLdrStrImm(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	B := bits.GetBit(opcodeFrag, 22)
	TW := bits.GetBit(opcodeFrag, 21)
	L := bits.GetBit(opcodeFrag, 20)

	var loadStoreFunc func(cpu ArmCPU, address uint32)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if L {
		if B && TW && !P {
			loadStoreFunc = ldrbt
		} else if B {
			loadStoreFunc = ldrb
		} else if TW && !P {
			loadStoreFunc = ldrt
		} else {
			loadStoreFunc = ldr
		}
	} else {
		if B && TW && !P {
			loadStoreFunc = strbt
		} else if B {
			loadStoreFunc = strb
		} else if TW && !P {
			loadStoreFunc = strt
		} else {
			loadStoreFunc = str
		}
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		loadStoreFunc(cpu, immOffLdStr(cpu, indexingFunc))
	}

	return true
}

func emitLdrStrReg(row int, col int) bool {
	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	P := bits.GetBit(opcodeFrag, 24)
	B := bits.GetBit(opcodeFrag, 22)
	TW := bits.GetBit(opcodeFrag, 21)
	L := bits.GetBit(opcodeFrag, 20)

	var loadStoreFunc func(cpu ArmCPU, address uint32)
	var indexingFunc indexFunc
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if L {
		if B && TW && !P {
			loadStoreFunc = ldrbt
		} else if B {
			loadStoreFunc = ldrb
		} else if TW && !P {
			loadStoreFunc = ldrt
		} else {
			loadStoreFunc = ldr
		}
	} else {
		if B && TW && !P {
			loadStoreFunc = strbt
		} else if B {
			loadStoreFunc = strb
		} else if TW && !P {
			loadStoreFunc = strt
		} else {
			loadStoreFunc = str
		}
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		loadStoreFunc(cpu, scaledRegOff(cpu, indexingFunc))
	}

	return true
}

func emitLoadStoreSingleOpcode(row int, col int) bool {

	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	I := bits.GetBit(opcodeFrag, 25)
	P := bits.GetBit(opcodeFrag, 24)
	B := bits.GetBit(opcodeFrag, 22)
	TW := bits.GetBit(opcodeFrag, 21)
	L := bits.GetBit(opcodeFrag, 20)

	/*	if !P && !TW {
		return false
	} */
	var loadStoreFunc func(cpu ArmCPU, address uint32)
	var indexingFunc indexFunc
	var addressingMode func(ArmCPU, indexFunc) uint32
	if !I {
		addressingMode = immOffLdStr
	} else {
		addressingMode = scaledRegOff
		if bits.GetBit(opcodeFrag, 4) {
			return false
		}
	}
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}

	if L {
		if B && TW && !P {
			loadStoreFunc = ldrbt
		} else if B {
			loadStoreFunc = ldrb
		} else if TW && !P {
			loadStoreFunc = ldrt
		} else {
			loadStoreFunc = ldr
		}
	} else {
		if B && TW && !P {
			loadStoreFunc = strbt
		} else if B {
			loadStoreFunc = strb
		} else if TW && !P {
			loadStoreFunc = strt
		} else {
			loadStoreFunc = str
		}
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		loadStoreFunc(cpu, addressingMode(cpu, indexingFunc))
	}

	return true
}

func emitLoadStoreHDOpcode(row int, col int) bool {

	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)
	I := bits.GetBit(opcodeFrag, 25)
	P := bits.GetBit(opcodeFrag, 24)
	//TW := bits.GetBit(opcodeFrag, 21)
	L := bits.GetBit(opcodeFrag, 20)

	/*if !P && TW {
		return false
	} */

	var loadStoreFunc func(cpu ArmCPU, address uint32)
	var indexingFunc indexFunc
	var addressingMode func(ArmCPU, indexFunc) uint32
	if I {
		addressingMode = miscImmOffLdStr
	} else {
		addressingMode = scaledRegOff
	}
	if P {
		indexingFunc = preIndexFunc
	} else {
		indexingFunc = postIndexFunc
	}
	opcode := bits.GetBits(opcodeFrag, 6, 5)
	if L {
		if opcode == 0 {
			// Reserved
			return false
		} else if opcode == 1 {
			loadStoreFunc = ldrh
		} else if opcode == 2 {
			loadStoreFunc = ldrsb
		} else if opcode == 3 {
			loadStoreFunc = ldrsh
		}
	} else {
		if opcode == 0 {
			// reserved for swap
			return false
		} else if opcode == 1 {
			loadStoreFunc = strh
		} else if opcode == 2 {
			// TODO: implement LDRD
			loadStoreFunc = strh
		} else if opcode == 3 {
			// TODO: implement STRD
			loadStoreFunc = strh
		}
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		loadStoreFunc(cpu, addressingMode(cpu, indexingFunc))
	}

	return true
}

func emitLdrStrMOpcode(row int, col int) bool {

	opcodeFrag := (uint32(row) << 20) | (uint32(col) << 4)

	P := bits.GetBit(opcodeFrag, 24)
	U := bits.GetBit(opcodeFrag, 23)
	S := bits.GetBit(opcodeFrag, 22)
	L := bits.GetBit(opcodeFrag, 20)

	var addressing func(ArmCPU) (uint32, uint32)
	var loadFunc func(ArmCPU, uint32, uint32)
	if P && U {
		addressing = ldStrMullIncBefore
	} else if !P && U {
		addressing = ldStrMullIncAfter
	} else if P && !U {
		addressing = ldStrMullDecBefore
	} else {
		addressing = ldStrMullDecAfter
	}

	if L {
		if S {
			loadFunc = ldmUser
		} else {
			loadFunc = ldm
		}
	} else {
		if S {
			loadFunc = stmUser
		} else {
			loadFunc = stm
		}
	}

	armInstructions[row][col] = func(cpu ArmCPU) {
		condition := bits.GetBits(cpu.currentInstruction(), 31, 28)
		if condition == 0xf {
			fmt.Println("Undefined Instruction NV ldm/stm")
			return
		}
		startAddress, endAddress := addressing(cpu)
		loadFunc(cpu, startAddress, endAddress)
	}

	return true

}

func init() {
	gen := NewGenerator(0xff, 0xf)
	gen.AddMaskTarget(MaskTarget{
		rowMask: "000xxxx0",
		name:    "Data processing immediate shift",
		colMask: "xxx0",
		handler: emitDataProcessImmShift,
	})

	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010x00",
		name:    "Move status register to register",
		colMask: "0000",
		handler: emitMRSOpcode,
	})
	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010x10",
		name:    "Move register to status to register",
		colMask: "0000",
		handler: emitMSROpcode,
	})
	gen.AddMaskTarget(MaskTarget{
		name:    "Data processing register shift",
		rowMask: "000xxxxx",
		colMask: "0xx1",
		handler: emitDataRegShift,
	})
	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010010",
		colMask: "0001",
		name:    "Branch/Exchange Instruction",
		handler: emitBranchExchangeOpcode,
	})
	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010110",
		name:    "Count leading zeros",
		colMask: "0001",
		handler: emitCLZOpcode,
	})
	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010010",
		name:    "Branch Link Exchange Instruction",
		colMask: "0011",
		handler: emitBLXRegOpcode,
	})

	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010xx0",
		name:    "Enhanced DSP add/subtracts",
		colMask: "0101",
		handler: emitDSPAddSub,
	})

	gen.AddMaskTarget(MaskTarget{
		rowMask: "00010010",
		name:    "Software Interrupt",
		colMask: "0111",
		handler: func(row int, col int) bool {
			//TODO implement software interrupts
			return true
		},
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Enchanced DSP Multiplies",
		rowMask: "00010xx0",
		colMask: "0111",
		handler: emitDSPMultiplyOpcodes,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Multiply (accumulate)",
		rowMask: "000000xx",
		colMask: "1001",
		handler: emitMultiplyAccumulate,
	})
	gen.AddMaskTarget(MaskTarget{
		name:    "Multiply (accumulate) long",
		rowMask: "00001xxx",
		colMask: "1001",
		handler: emitMultiplyAccumulateLong,
	})
	gen.AddMaskTarget(MaskTarget{
		name:    "Swap/swap byte",
		rowMask: "00010x00",
		colMask: "1001",
		handler: emitSwapOpcodes,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store halfword register offset",
		rowMask: "000xx0xx",
		colMask: "1011",
		handler: emitLdrStrHReg,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store halfword immediate offset",
		rowMask: "000xx1xx",
		colMask: "1011",
		handler: emitLdrStrHImm,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store two words register offset",
		rowMask: "000xx0x0",
		colMask: "11x1",
		handler: emitLdrStrDReg,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load signed halfword/byte register offset",
		rowMask: "000xx0x1",
		colMask: "11x1",
		handler: emitLdrSBReg,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store two words immediate offset",
		rowMask: "000xx1x0",
		colMask: "11x1",
		handler: emitLdrStrDImm,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load signed halfword/byte immediate offset",
		rowMask: "000xx1x1",
		colMask: "11x1",
		handler: emitLdrStrDImm,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Data processing immediate",
		rowMask: "001xxxxx",
		colMask: "xxxx",
		handler: emitDataImm,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Undefined Instruction",
		rowMask: "00110x00",
		colMask: "xxxx",
		handler: func(row int, col int) bool {
			return true
		},
	})
	gen.AddMaskTarget(MaskTarget{
		name:    "Move immediate to status register",
		rowMask: "00110x10",
		colMask: "xxxx",
		handler: emitMSRImmOpcode,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store immediate offset",
		rowMask: "010xxxxx",
		colMask: "xxxx",
		handler: emitLdrStrImm,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store register offset",
		rowMask: "011xxxxx",
		colMask: "xxx0",
		handler: emitLdrStrReg,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Undefined Instruction",
		rowMask: "011xxxxx",
		colMask: "xxx1",
		handler: func(row int, col int) bool {
			armInstructions[row][col] = func(cpu ArmCPU) {
				if bits.GetBits(cpu.currentInstruction(), 31, 28) == 0xf {
					fmt.Println("Undefined Instruction with NV")
				}
			}
			return true
		},
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Undefined Instruction",
		rowMask: "0xxxxxxx",
		colMask: "xxxx",
		handler: func(row int, col int) bool {
			return true
		},
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Load/store multiple",
		rowMask: "100xxxxx",
		colMask: "xxxx",
		handler: emitLdrStrMOpcode,
	})

	gen.AddMaskTarget(MaskTarget{
		name:    "Branch and Branch with Link / Change to thumb",
		rowMask: "101xxxxx",
		colMask: "xxxx",
		handler: emitBranchLinkOpcode,
	})

	gen.start()
}
