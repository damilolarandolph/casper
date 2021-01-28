package cpu

type comparable interface {
	equals(to comparable) bool
}

func compare(val comparable, to comparable) bool {
	return val.equals(to)
}

type exception struct {
	mode          cpuMode
	normalAddress uint32
	highAddress   uint32
}

func (e exception) equals(to comparable) bool {
	val, _ := to.(exception)
	return (e.mode == val.mode) && (e.normalAddress == val.normalAddress)
}

var (
	// ResetEx is a Reset CPU exception
	ResetEx = exception{
		mode:          supervisor,
		normalAddress: 0x0,
		highAddress:   0xffff0000,
	}

	// UndInstrEx is an Undefined Instruction Exception
	UndInstrEx = exception{
		mode:          undefined,
		normalAddress: 0x00000004,
		highAddress:   0xFFFF0004,
	}

	// SwiEx is a Software Exception
	SwiEx = exception{
		mode:          supervisor,
		normalAddress: 0x00000008,
		highAddress:   0xFFFF0008,
	}

	// PrefetchAbtEx is a Prefetch Abort Exception
	PrefetchAbtEx = exception{
		mode:          abort,
		normalAddress: 0x0000000C,
		highAddress:   0xFFFF000C,
	}

	// DataAbtEx is a Data Abort Exception
	DataAbtEx = exception{
		mode:          abort,
		normalAddress: 0x00000010,
		highAddress:   0xFFFF0010,
	}

	// IrqEx is an Interrupt Exception
	IrqEx = exception{
		mode:          irq,
		normalAddress: 0x00000018,
		highAddress:   0xFFFF0018,
	}

	// FiqEx is a Fast Interrupt Exception
	FiqEx = exception{
		mode:          fiq,
		normalAddress: 0x0000001C,
		highAddress:   0xFFFF001C,
	}
)

func (cpu *Arm7) highIrqVectors() bool {
	return cpu.irqHigh
}

func (cpu *Arm7) requestInterrupt(ex exception) {

	if compare(ex, FiqEx) && cpu.isFlag(fiqDisable) {
		return
	}

	if cpu.isFlag(irqDisable) {
		return
	}
	cpu.setMode(supervisor)
	cpu.setFlag(thumbMode, true)
	if compare(ex, ResetEx) || compare(ex, FiqEx) {
		cpu.setFlag(fiqDisable, true)
	}
	cpu.setFlag(irqDisable, true)

	if cpu.irqHigh {
		cpu.setPc(ex.highAddress)
		return
	}

	cpu.setPc(ex.normalAddress)
}
