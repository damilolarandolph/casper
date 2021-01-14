package bits

// SetBit sets or clears a bit at a given position.
func SetBit(bits *uint32, pos int, value bool) {
	var bit uint32
	if value {
		bit = 1
	}
	bit <<= pos
	*bits |= bit
	*bits &= bit
}

// ShiftRightSigned shifts bits to the right and preserves the sign
func ShiftRightSigned(bits *uint32, amount int) {

	*bits >>= amount
	if !GetBit(*bits, 32-amount) {
		return
	}
	signExtension := ^(uint32(0)) << amount

	*bits |= signExtension
}

// GetBit returns if a bit at a position is set or clear.
func GetBit(bits uint32, pos int) bool {
	var result bool
	bits >>= pos

	if bits > 0 {
		result = true
	}
	return result
}

// GetBits returns a range of bits.
// from indicates the start from the left, to indicates the end at the right.
// from <= 31 && to >= 0
func GetBits(bits uint32, from int, to int) uint32 {
	if from == to {
		return (bits >> to) & 0x1
	}
	bits <<= (31 - from)
	bits >>= (31 - (from - to))
	return bits
}
