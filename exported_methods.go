package pyrand

import "math/big"
import "fmt"

// Returns a new Random ready to be used
func NewRandom() *Random {
	return &Random{}
}

// SeedFromUInt32 seeds the PRNG with the given seed value,
// causing the same sequence of numbers to be generated as
// when seeding Python's PRNG with the corresponding number.
func (r *Random) SeedFromUInt32(seed uint32) {
	r.SeedFromUInt32s([]uint32{seed})
}

// SeedFromUInt64 seeds the PRNG with the given seed value,
// causing the same sequence of numbers to be generated as
// when seeding Python's PRNG with the corresponding number.
func (r *Random) SeedFromUInt64(seed uint64) {
	if seed <= 0xffffffff {
		r.SeedFromUInt32(uint32(seed))
		return
	}
	r.SeedFromUInt32s([]uint32{uint32(seed), uint32(seed >> 32)})
}

// SeedFromUInt32s seeds the PRNG with the given seed slice.
// Python doesn't expose such a method directly, but behind
// the scenes this is what Python does e.g. when seeding with
// a long integer: Split it into 32-bit chunks (little-endian),
// then seed based on these.
func (r *Random) SeedFromUInt32s(seed []uint32) {
	r.initBySlice(seed)
}

// SeedFromHexString takes a string of hex digits and seeds
// the PRNG with the corresponding number. E.g. seeding
// Pythons PRNG with 0x1234567890deadbeef1337cafe is
// equivalent to calling
//
//     r.SeedFromHexString("1234567890deadbeef1337cafe")
func (r *Random) SeedFromHexString(seed string) error {
	b := new(big.Int)
	_, err := fmt.Sscanf(seed, "%x", b)
	if err != nil {
		return err
	}
	r.SeedFromBig(b)
	return nil
}

// SeedFromBig seeds the PRNG with a bit integer (from the math/big
// package), causing the same sequence of numbers to be generated as
// when seeding Python's PRNG with the corresponding number.
func (r *Random) SeedFromBig(seed *big.Int) {
	seed.Abs(seed)
	bytes := seed.Bytes()
	r.SeedFromBytes(bytes)
}

// SeedFromBytes takes a big-endian slice of bytes and
// seeds the PRNG with it. Python doesn't have such a method,
// but we need it here to implement SeedFromBig, so
// we might as well export it.
func (r *Random) SeedFromBytes(b []byte) {
	sourcelen := uint32(len(b))
	if sourcelen == 0 {
		return
	}
	length := (sourcelen-1)/4 + 1
	uints := make([]uint32, length)
	for i := uint32(0); i < sourcelen; i += 4 {
		pos := sourcelen - 1 - i
		var v uint32 = 0
		for j := uint32(0); j < 4 && i+j < sourcelen; j++ {
			v = v | (uint32(b[pos-j]) << (j * 8))
		}
		uints[i/4] = v
	}
	r.initBySlice(uints)
}

// Random returns the next random floating point number in the
// range [0.0, 1.0).
func (r *Random) Random() float64 {
	return r.genrandRes53()
}

// RandBits returns a slice of uint32 that's filled with random
// bits. If k is not divisible by 32, then the 32-k most significant
// bits of the last slice element will always be zero.
func (r *Random) RandBits(k uint) []uint32 {
	if k == 0 {
		panic("number of bits must be >0")
	}

	quads := ((k-1)/32 + 1)
	result := make([]uint32, quads)

	for i := uint(0); i < quads; i++ {
		v := r.genRandInt32()
		if k < 32 {
			v >>= 32 - k
		}
		result[i] = v

		k -= 32
	}
	return result
}

// RandInt returns a random integer in range [a, b], including both end points.
func (r *Random) RandInt(a int, b int) int {
	return r.RandRange(a, b+1)
}

const maxwidth = 1 << 53

// RandInt returns a random integer in range [start, stop), including the low value but
// excluding the high value.
func (r *Random) RandRange(start int, stop int) int {
	if start >= stop {
		panic("empty range for randrange")
	}
	width := uint64(stop - start)
	if width >= maxwidth {
		return start + int(r.randBelow(width))
	}
	return int(start + int(r.Random()*float64(width)))
}

// Choice is essentially RandRange with a first argument of 0.
// It's provided here as the equivalent to Python's random.choice(),
// where
//
//     l := []byte{42, 666, 13}
//     c := l[r.Choice(len(l))]
//
// is equivalent to Python's
//
//     l := [42, 666, 13]
//     c := r.choice(l)
//
func (r *Random) Choice(length int) int {
	return int(r.Random() * float64(length))
}