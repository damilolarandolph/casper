package cpu

import "github.com/damilolarandolph/casper/memory"

type fetchType int

const (
	wordFetch     fetchType = 3
	halfWordFetch           = 1
	byteFetch               = 0
)

// DataBus describes a system bus used for data accesses.
type DataBus interface {
	ReadData8(address uint32) uint32
	ReadData16(address uint32) uint32
	ReadData32(address uint32) uint32
	WriteData8(address uint32, val uint32)
	WriteData16(address uint32, val uint32)
	WriteData32(address uint32, val uint32)
	SetSequencial(val bool)
}

// CodeBus describes a system bus used for code accesses.
type CodeBus interface {
	SetSequencial(val bool)
	//	ReadCode8(address uint32) uint32
	ReadCode16(address uint32) uint32
	ReadCode32(address uint32) uint32
	/*WriteCode8(address uint32, val uint32)
	WriteCode16(address uint32, val uint32)
	WriteCode32(address uint32, val uint32)*/
}

// SystemBus describes a bus that can be used for both
// code and data accesses.
type SystemBus interface {
	DataBus
	CodeBus
}

// Arm7Bus implements an Arm7 system bus
type Arm7Bus struct {
	arm7        ArmCPU
	arm9        ArmCPU
	memory      *memory.InternalMemory
	isOpcode    bool
	accessType  fetchType
	sequencial  bool
	dataTimings [][]int
	codeTimings [][]int
}

/* ReadCode8 performs an 8 bit opcode fetch.
func (bus *Arm7Bus) ReadCode8(address uint32) uint32 {
	bus.isOpcode = true
	bus.accessType = byteFetch
	return bus.read(address)
}*/

// ReadCode16 performs a 16 bit opcode fetch.
func (bus *Arm7Bus) ReadCode16(address uint32) uint32 {
	bus.isOpcode = true
	bus.accessType = halfWordFetch
	return bus.read(address)
}

// ReadCode32 performs a 32 bit opcode fetch.
func (bus *Arm7Bus) ReadCode32(address uint32) uint32 {
	bus.isOpcode = true
	bus.accessType = wordFetch
	return bus.read(address)
}

// ReadData8 performs an 8 bit data fetch.
func (bus *Arm7Bus) ReadData8(address uint32) uint32 {
	bus.isOpcode = false
	bus.accessType = byteFetch
	return bus.read(address)
}

// ReadData16 performs a 16 bit data fetch.
func (bus *Arm7Bus) ReadData16(address uint32) uint32 {
	bus.isOpcode = false
	bus.accessType = halfWordFetch
	return bus.read(address)
}

// ReadData32 performs a 32 bit data fetch.
func (bus *Arm7Bus) ReadData32(address uint32) uint32 {
	bus.isOpcode = false
	bus.accessType = wordFetch
	return bus.read(address)
}

/*
func (bus *Arm7Bus) WriteCode8(address uint32, val uint32) {
	bus.isOpcode = true
	bus.accessType = byteFetch
	bus.write(address, val)
}

func (bus *Arm7Bus) WriteCode16(address uint32, val uint32) {
	bus.isOpcode = true
	bus.accessType = halfWordFetch
	bus.write(address, val)
}

func (bus *Arm7Bus) WriteCode32(address uint32, val uint32) {
	bus.isOpcode = true
	bus.accessType = wordFetch
	bus.write(address, val)
} */

// WriteData8 performs an 8 bit write.
func (bus *Arm7Bus) WriteData8(address uint32, val uint32) {
	bus.isOpcode = false
	bus.accessType = byteFetch
	bus.write(address, val)
}

// WriteData16 performs a 16 bit write.
func (bus *Arm7Bus) WriteData16(address uint32, val uint32) {
	bus.isOpcode = false
	bus.accessType = halfWordFetch
	bus.write(address, val)
}

// WriteData32 performs a 32 bit write.
func (bus *Arm7Bus) WriteData32(address uint32, val uint32) {
	bus.isOpcode = false
	bus.accessType = wordFetch
	bus.write(address, val)
}

// SetSequencial set the next memory access an non-sequencial.
func (bus *Arm7Bus) SetSequencial(val bool) {
	bus.sequencial = val
}

func (bus *Arm7Bus) read(address uint32) uint32 {
	if address >= 04000000 && address < 0x04800000 {
		return bus.ioReadBytes(address)
	}
	return bus.memReadBytes(address)
}

func (bus *Arm7Bus) write(address uint32, val uint32) {
	if address >= 04000000 && address < 0x04800000 {
		bus.ioWriteBytes(address, val)
		return
	}
	bus.memWriteBytes(address, val)
	return
}

func (bus *Arm7Bus) getTiming(row int) int {

	var col int

	if bus.sequencial {
		if bus.accessType == wordFetch {
			col = 1
		} else if bus.accessType == halfWordFetch {
			col = 3
		}
	} else {
		if bus.accessType == wordFetch {
			col = 0
		} else if bus.accessType == halfWordFetch {
			col = 2
		}
	}

	if bus.isOpcode {
		return bus.codeTimings[row][col]
	}

	return bus.dataTimings[row][col]

}

func (bus *Arm7Bus) wait(amount int) {
	for ; amount > 0; amount-- {
		bus.arm7.WaitForTick()
	}
}

func (bus *Arm7Bus) getMemRegionParams(address uint32) (memRegion []uint8, translatedAddr uint32) {

	//ARM7 BIOS
	if address >= 0x0 && address < 0x02000000 {
		memRegion = bus.memory.BIOS
		translatedAddr = address
		bus.wait(bus.getTiming(1))
		return
	}
	// MAIN RAM
	if address < 0x03000000 {
		memRegion = bus.memory.MainRAM
		translatedAddr = address - 0x02000000
		bus.wait(bus.getTiming(0))
		return
	}
	// Shared WRAM
	if address < 0x03800000 {
		memRegion = bus.memory.Wram
		translatedAddr = address - 0x03000000
		bus.wait(bus.getTiming(1))
		return
	}
	// WRAM
	if address < 0x04000000 {
		memRegion = bus.memory.Wram
		translatedAddr = address - 0x03800000
		bus.wait(bus.getTiming(1))
		return
	}
	// I/O
	if address < 0x04800000 {
		memRegion = nil
		translatedAddr = address - 0x04000000
		bus.wait(bus.getTiming(1))
		return
	}
	// WIFI RAM
	if address < 0x04808000 {
		memRegion = bus.memory.WifiRAM
		translatedAddr = address - 0x04800000
		bus.wait(bus.getTiming(1))
		return
	}

	// WIFI I/O
	if address < 0x06000000 {
		memRegion = nil
		translatedAddr = address - 0x04808000
		bus.wait(bus.getTiming(1))
		return
	}

	//VRAM
	if address < 0x08000000 {
		memRegion = bus.memory.Vram
		translatedAddr = address - 0x06000000
		bus.wait(bus.getTiming(2))
		return
	}

	//GBA ROM
	if address < 0x0A000000 {
		memRegion = bus.memory.ROM
		translatedAddr = address - 0x08000000
		bus.wait(bus.getTiming(3))
		return
	}

	//GBA RAM
	memRegion = nil
	translatedAddr = address - 0x0A000000
	bus.wait(bus.getTiming(4))
	return

}

func (bus *Arm7Bus) memReadBytes(address uint32) uint32 {

	var value uint32
	memRegion, tranlatedAddress := bus.getMemRegionParams(address)
	var endAddress = tranlatedAddress + uint32(bus.accessType)

	for ; endAddress >= tranlatedAddress; endAddress-- {
		value <<= 8
		value |= uint32(memRegion[endAddress])
		if endAddress == 0 {
			break
		}
	}

	return value

}

func (bus *Arm7Bus) memWriteBytes(address uint32, val uint32) {
	var value uint32
	memRegion, tranlatedAddress := bus.getMemRegionParams(address)
	endAddress := tranlatedAddress + uint32(bus.accessType)
	for ; endAddress >= tranlatedAddress; endAddress-- {
		memRegion[endAddress] = uint8(value)
		value >>= 8
		if endAddress == 0 {
			break
		}
	}
}

func (bus *Arm7Bus) ioWriteBytes(address uint32, val uint32) {}

func (bus *Arm7Bus) ioReadBytes(address uint32) uint32 {
	return 0
}
