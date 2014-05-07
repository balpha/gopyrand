/* This file is a pretty straightforward port to Go of the relevant C code in
 * http://hg.python.org/cpython/file/faef1da30c6d/Modules/_randommodule.c.
 * Code and comments have been only slightly adjusted.
 *
 * What follows is the preamble from that file; note that not the full content
 * of the original C file was ported here.
 *
 */

/* ------------------------------------------------------------------
   The code in this module was based on a download from:
      http://www.math.keio.ac.jp/~matumoto/MT2002/emt19937ar.html

   It was modified in 2002 by Raymond Hettinger as follows:

    * the principal computational lines untouched.

    * renamed genrand_res53() to random_random() and wrapped
      in python calling/return code.

    * genrand_int32() and the helper functions, init_genrand()
      and init_by_array(), were declared static, wrapped in
      Python calling/return code.  also, their global data
      references were replaced with structure references.

    * unused functions from the original were deleted.
      new, original C python code was added to implement the
      Random() interface.

   The following are the verbatim comments from the original code:

   A C-program for MT19937, with initialization improved 2002/1/26.
   Coded by Takuji Nishimura and Makoto Matsumoto.

   Before using, initialize the state by using init_genrand(seed)
   or init_by_array(init_key, key_length).

   Copyright (C) 1997 - 2002, Makoto Matsumoto and Takuji Nishimura,
   All rights reserved.

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:

     1. Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.

     2. Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in the
    documentation and/or other materials provided with the distribution.

     3. The names of its contributors may not be used to endorse or promote
    products derived from this software without specific prior written
    permission.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR
   CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
   EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
   PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
   PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
   LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
   NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
   SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


   Any feedback is very welcome.
   http://www.math.keio.ac.jp/matumoto/emt.html
   email: matumoto@math.keio.ac.jp
*/

/* ---------------------------------------------------------------*/

package pyrand

const (
	cN                 = 624
	cM                 = 397
	cMATRIX_A   uint32 = 0x9908b0df // constant vector a
	cUPPER_MASK uint32 = 0x80000000 // most significant w-r bits
	cLOWER_MASK uint32 = 0x7fffffff // least significant r bits
)

var mag01 = [2]uint32{0, cMATRIX_A}

type Random struct {
	state [cN]uint32
	index uint32
}

/* genrandRes53 generates a random number on [0,1) with 53-bit resolution; note that
 * 9007199254740992 == 2**53; I assume they're spelling "/2**53" as
 * multiply-by-reciprocal in the (likely vain) hope that the compiler will
 * optimize the division away at compile-time.  67108864 is 2**26.  In
 * effect, a contains 27 random bits shifted left 26, and b fills in the
 * lower 26 bits of the 53-bit numerator.
 * The orginal code credited Isaku Wada for this algorithm, 2002/01/09.
 */
func (r *Random) genrandRes53() float64 {
	a := float64(r.genRandInt32() >> 5)
	b := float64(r.genRandInt32() >> 6)
	return (a*float64(67108864.0) + b) * (1.0 / float64(9007199254740992.0))
}

func (r *Random) genRandInt32() uint32 {
	var y uint32
	mt := &r.state
	if r.index >= cN { // generate N words at one time
		var kk uint32
		for kk = 0; kk < cN-cM; kk++ {
			y = (mt[kk] & cUPPER_MASK) | (mt[kk+1] & cLOWER_MASK)
			mt[kk] = mt[kk+cM] ^ (y >> 1) ^ mag01[y&1]
		}
		for ; kk < cN-1; kk++ {
			y = (mt[kk] & cUPPER_MASK) | (mt[kk+1] & cLOWER_MASK)
			mt[kk] = mt[kk+cM-cN] ^ (y >> 1) ^ mag01[y&1]
		}
		y = (mt[cN-1] & cUPPER_MASK) | (mt[0] & cLOWER_MASK)
		mt[cN-1] = mt[cM-1] ^ (y >> 1) ^ mag01[y&1]

		r.index = 0
	}

	y = mt[r.index]
	r.index++
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)
	return y
}

// initGenrand initializes mt[N] with a seed.
func (r *Random) initGenrand(s uint32) {
	var mti uint32
	mt := &r.state
	mt[0] = s
	for mti = 1; mti < cN; mti++ {
		/* See Knuth TAOCP Vol2. 3rd Ed. P.106 for multiplier. */
		/* In the previous versions, MSBs of the seed affect   */
		/* only MSBs of the array mt[].                        */
		/* 2002/01/09 modified by Makoto Matsumoto             */
		mt[mti] = uint32(1812433253)*(mt[mti-1]^(mt[mti-1]>>30)) + mti
	}
	r.index = mti
}

// initBySlice initializes by a slice. initKey is the slice for initializing keys
func (r *Random) initBySlice(initKey []uint32) {
	keyLength := uint32(len(initKey))
	mt := &r.state
	r.initGenrand(19650218)
	var i, j uint32 = 1, 0
	var k uint32
	if cN > keyLength {
		k = cN
	} else {
		k = keyLength
	}
	for ; k > 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 30)) * 1664525)) + initKey[j] + j // non linear
		i++
		j++
		if i >= cN {
			mt[0] = mt[cN-1]
			i = 1
		}
		if j >= keyLength {
			j = 0
		}
	}
	for k = cN - 1; k > 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 30)) * 1566083941)) - i // non linear
		i++
		if i >= cN {
			mt[0] = mt[cN-1]
			i = 1
		}
	}
	mt[0] = 0x80000000 // MSB is 1; assuring non-zero initial array
}
