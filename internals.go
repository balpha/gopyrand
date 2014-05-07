package pyrand

import "math"

// randBelow returns a uint64 strictly smaller than the given
// upper boundary.
func (r *Random) randBelow(n uint64) uint64 {

	// This is what Python 2 does. Python 3 counts the
	// actual number of bits of n and thus has
	// different behavior.
	bits := uint(1.00001 + math.Log2(float64(n)-1))

	two := bits > 32

	// we will never need more than 64 bits, but the bits formular above
	// may cause more to be requested
	three := bits > 64

	v := n
	for v >= n {
		s := r.RandBits(bits)

		// since we're constrained to uint64, and s[2] > 0 means
		// there's a 1 beyond the 64th bit of the number represented
		// by s, this is the v >= n case, and thus we retry
		if three && s[2] > 0 {
			v = n
			continue
		}

		v = uint64(s[0])
		if two {
			v |= uint64(s[1]) << 32
		}
	}

	return v
}
