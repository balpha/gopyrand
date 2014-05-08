package pyrand

import "testing"
import "fmt"

func tendigits(v float64) string {
	s := fmt.Sprintf("%.11f", v)
	return s[:len(s)-1]
}

// Yeah, this isn't very go-ish, and doesn't give you good line numbers
// on failure.

func assertEqual(a string, b string, t *testing.T) {
	if a != b {
		t.Error(a + " != " + b)
	}
}

func assertEqual64(a uint64, b uint64, t *testing.T) {
	if a != b {
		t.Error(fmt.Sprintf("%v != %v", a, b))
	}
}

func assertEqual32(a uint32, b uint32, t *testing.T) {
	if a != b {
		t.Error(fmt.Sprintf("%v != %v", a, b))
	}
}

func assertEqualInt(a int, b int, t *testing.T) {
	if a != b {
		t.Error(fmt.Sprintf("%v != %v", a, b))
	}
}

func assertEqual32s(a []uint32, b []uint32, t *testing.T) {
	if len(a) != len(b) {
		t.Error(fmt.Sprintf("slice length differs, %v != %v", len(a), len(b)))
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Error(fmt.Sprintf("slice[%v] differs, %v != %v", i, a[i], b[i]))
		}
	}
}

func makeR(seed interface{}, iterations int) *Random {
	r := NewRandom()
	switch seed := seed.(type) {
	default:
		panic("cannot seed from this")
	case uint32:
		r.SeedFromUInt32(seed)
	case uint64:
		r.SeedFromUInt64(seed)
	case int:
		r.SeedFromUInt64(uint64(seed))
	case string:
		e := r.SeedFromHexString(seed)
		if e != nil {
			panic("failed to seed from hex string " + seed + ": " + e.Error())
		}
	}
	for i := 0; i < iterations; i++ {
		r.Random()
	}
	return r
}

func TestRandom(t *testing.T) {
	assertEqual(tendigits(makeR(uint32(123), 1000).Random()), "0.0638474579", t)
	assertEqual(tendigits(makeR(uint32(0), 987654).Random()), "0.7203062140", t)
	assertEqual(tendigits(makeR(uint32(0xfffe), 5927).Random()), "0.5279272538", t)
	assertEqual(tendigits(makeR(uint32(0xffff), 5927).Random()), "0.7647091554", t)
	assertEqual(tendigits(makeR(uint32(0x10000), 5927).Random()), "0.8890962216", t)
	assertEqual(tendigits(makeR(uint32(654321), 0).Random()), "0.0657799204", t)

	assertEqual(tendigits(makeR("1234567890deadbeefcafe1337600df00d", 0).Random()), "0.9094618764", t)
	assertEqual(tendigits(makeR("fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210", 7777).Random()), "0.9053125602", t)

	assertEqual(tendigits(makeR(uint64(0xa37b3f09a188e), 12345).Random()), "0.6162433684", t)
	assertEqual(tendigits(makeR(uint64(0xffffffffffffffff), 999).Random()), "0.9009945166", t)
	assertEqual(tendigits(makeR(432153415134, 986).Random()), "0.7026873940", t)
}

func TestRandBelow(t *testing.T) {
	assertEqual64(makeR(uint64(117624834567), 5678).randBelow(2000), 1453, t)
	assertEqual64(makeR(uint64(6513265496841), 4567).randBelow(0xfffffffd), 2688309836, t)
	assertEqual64(makeR(uint64(65132495874231), 12288).randBelow(0xffffffff), 848139872, t)
	assertEqual64(makeR(uint64(987651354), 16587).randBelow(0x100000000), 617983553, t)
	assertEqual64(makeR(uint64(1684651512), 3486).randBelow(0x100000001), 3726269297, t)
	assertEqual64(makeR(uint64(17209), 68133).randBelow(0xfffffffffffffffe), 17889265393449113490, t)
	assertEqual64(makeR(uint64(555555), 17009).randBelow(0xffffffffffffffff), 14674416218734170714, t)
}

func TestRandBits(t *testing.T) {
	assertEqual32s(makeR(0, 0).RandBits(8), []uint32{216}, t)
	assertEqual32s(makeR(0, 0).RandBits(32), []uint32{3626764237}, t)
	assertEqual32s(makeR(0, 0).RandBits(33), []uint32{3626764237, 0}, t)
	assertEqual32s(makeR(0, 0).RandBits(63), []uint32{3626764237, 827307999}, t)
	assertEqual32s(makeR(0, 0).RandBits(64), []uint32{3626764237, 1654615998}, t)
	assertEqual32s(makeR(21684, 1111).RandBits(33), []uint32{1651504065, 1}, t)
}

func TestRandInt(t *testing.T) {
	assertEqualInt(makeR(519876, 8956).RandInt(13, 97), 84, t)
	assertEqualInt(makeR(432153415134, 986).RandInt(-12307, -803), -4223, t)
}
