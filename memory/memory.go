package memory

const (
	mainRAMSize   int = 0x400000 // 4MB
	wramSize      int = 0x400    // 96KB
	vramSize      int = 0xa4000  // 656KB
	oamPalSize    int = 0x1000   // 4KB
	_3DMemorySize int = 0x3e000  // 248KB
	// Matrix stack size is unknown according to GBATEK.
	// So I just set it to 4MB like main ram :-p
	matrixStackSize   int = mainRAMSize
	wifiRAMSize       int = oamPalSize * 2 // 8KB
	firmwareFlashSize int = 0x40000        // 256KB
	biosSize          int = 0x9000         // 36KB
)

type MemorySlot int

const (
	MainRAM MemorySlot = iota
	WRAM
	VRAM
	OAMPAL
	GFX
	MatrixStack
	WIFIRAM
	Firmware
	BIOS
	ROM
)

type InternalMemory struct {
	MainRAM     []uint8
	Wram        []uint8
	Vram        []uint8
	OamPal      []uint8
	Gfx         []uint8
	MatrixStack []uint8
	WifiRAM     []uint8
	Firmware    []uint8
	BIOS        []uint8
	ROM         []uint8
}

func NewInternalMemory() *InternalMemory {

	return &InternalMemory{
		MainRAM:     make([]uint8, mainRAMSize),
		Wram:        make([]uint8, wramSize),
		Vram:        make([]uint8, vramSize),
		OamPal:      make([]uint8, oamPalSize),
		Gfx:         make([]uint8, _3DMemorySize),
		MatrixStack: make([]uint8, matrixStackSize),
		WifiRAM:     make([]uint8, wifiRAMSize),
		Firmware:    make([]uint8, firmwareFlashSize),
		BIOS:        make([]uint8, biosSize),
	}
}
