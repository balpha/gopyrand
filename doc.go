/*
 * Package pyrand is an implementation of a pseudo-random number generator
 * (PRNG) with the explicit goal that when seeded with identical values,
 * it generates the same numbers as the default implementation in the CPython 2
 * standard library, including support for seeding with arbitrarily large
 * integers (via the math/big package), however not (yet) for *generating*
 * arbitrarily large random numbers. The stress is on the "2": The PRNGs
 * in Python 2 and 3 don't behave 100% identically.
 *
 * This is a fairly narrow use case; if you don't care about it, then this
 * PRNG probably won't give you anything that other implementations don't,
 * including Go's built-in one. In particular, it is a Mersenne Twister,
 * and thus is *not* cryptographically secure. Don't use it to generate
 * SSH keys.
 *
 * Usage
 *
 * pyrand provides a single type, Random. A newly created Random is unseeded
 * and thus useless; there is no auto-seeding because the only point of this
 * implementation is that it predictably gives the same number sequences as Python
 * when seeded identically. So the first thing you should do is seed the Random:
 *
 *     r := Random{}
 *     r.SeedFromUInt32(12345)
 *
 * The following seeding methods are provided. They take numbers
 * in various representations and will cause the Random to generate the same
 * values as a Python 2 random.Random seeded with the same number. The argument
 * to SeedFromBytes is interpreted as big-endian (most significant byte first).
 *
 *     SeedFromUInt32(uint32)
 *     SeedFromUInt64(uint64)
 *     SeedFromBig(*big.Int)
 *     SeedFromHexString(string)
 *     SeedFromBytes([]byte)
 *
 * The following methods to retrieve random numbers are provided, and they
 * mimic the correspondingly named Python methods. Note that, unlike the Python
 * version, RandRange takes no step parameter.
 *
 *     Random() float64
 *     RandBits(uint) []uint32
 *     RandInt(int, int) int
 *     RandRange(int, int) int
 *
 * See the method documentations for more details.
 *
 * License
 *
 * See the accompanying file LICENSE.txt for copyright and licensing information.
 */
package pyrand