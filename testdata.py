from random import Random

def R(seed, iterations):
    r = Random(seed)
    for i in xrange(iterations):
        r.random()
    return r

def td(v):
    return ("%.11f" % v)[0:-1]

def randomData():
    print td(R(123, 1000).random())
    print td(R(0, 987654).random())
    print td(R(0xfffe, 5927).random())
    print td(R(0xffff, 5927).random())
    print td(R(0x10000, 5927).random())
    print td(R(654321, 0).random())
    
    print
    
    print td(R(0x1234567890deadbeefcafe1337600df00d, 0).random())
    print td(R(0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210, 7777).random())
    
    print
    
    print td(R(0xa37b3f09a188e, 12345).random())
    print td(R(0xffffffffffffffff, 999).random())
    print td(R(432153415134, 986).random())
    
def randbelowData():
    print R(117624834567, 5678)._randbelow(2000)
    
    print R(6513265496841, 4567)._randbelow(0xfffffffd)
    
    print R(21684, 1111)._randbelow(0xfffffffe)
    print R(65132495874231, 12288)._randbelow(0xffffffff)
    print R(987651354, 16587)._randbelow(0x100000000)
    print R(1684651512, 3486)._randbelow(0x100000001)
    
    print R(17209, 68133)._randbelow(0xfffffffffffffffe)
    print R(555555, 17009)._randbelow(0xffffffffffffffff)

def randbitsData():    
    print R(0,0).getrandbits(8)
    print R(0,0).getrandbits(32)
    x = R(0,0).getrandbits(33)
    print [x & 0xffffffff, x >> 32]
    x = R(0,0).getrandbits(63)
    print [x & 0xffffffff, x >> 32]
    x = R(0,0).getrandbits(64)
    print [x & 0xffffffff, x >> 32]
    
def randintdata():
    print R(519876, 8956).randint(13, 97)
    print R(432153415134, 986).randint(-12307, -803)

    
randintdata()