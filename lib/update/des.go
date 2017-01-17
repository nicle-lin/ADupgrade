package update

const EN0 = 0 /*MODE == encrypt */
const DN1 = 1 /*MODE == decrypt */

var KnL = [32]uint32{0}

var bytebit = [8]uint32{01, 02, 04, 10, 20, 040, 100, 200}

var bigbyte = [24]uint32{
	0x800000, 0x400000, 0x200000, 0x100000,
	0x80000, 0x40000, 0x20000, 0x10000,
	0x8000, 0x4000, 0x2000, 0x1000,
	0x800, 0x400, 0x200, 0x100,
	0x80, 0x40, 0x20, 0x10,
	0x8, 0x4, 0x2, 0x1,
}

var pc1 = [56]uint8{
	56, 48, 40, 32, 24, 16, 8, 0, 57, 49, 41, 33, 25, 17,
	9, 1, 58, 50, 42, 34, 26, 18, 10, 2, 59, 51, 43, 35,
	62, 54, 46, 38, 30, 22, 14, 6, 61, 53, 45, 37, 29, 21,
	13, 5, 60, 52, 44, 36, 28, 20, 12, 4, 27, 19, 11, 3,
}

var totrot = [16]uint8{1, 2, 4, 6, 8, 10, 12, 14, 15, 17, 19, 21, 23, 25, 27, 28}

var pc2 = [48]uint8{
	13, 16, 10, 23, 0, 4, 2, 27, 14, 5, 20, 9,
	22, 18, 11, 3, 25, 7, 15, 6, 26, 19, 12, 1,
	40, 51, 30, 36, 46, 54, 29, 39, 50, 44, 32, 47,
	43, 48, 38, 55, 33, 52, 45, 41, 49, 35, 28, 31,
}

func deskey(key *int, edf int) {
	var i, j, l, m, n int
	var pc1m, pcr [56]uint8
	var kn [32]uint32

	for j = 0; j < 56; j++ {
		l = pc1[j]
		m = l & 07
		if key[l>>3]&bytebit[m] > 0 {
			pc1m[j] = 1
		} else {
			pc1m[j] = 0
		}
	}

	for i = 0; i < 16; i++ {
		if edf == DN1 {
			m = (15 - i) << 1
		} else {
			m = i << 1
		}
		n = m + 1
		kn[m] = 0
		kn[n] = 0

		for j = 0; j < 28; j++ {
			l = j + totrot[i]
			if l < 28 {
				pcr[j] = pc1m[l]
			} else {
				pcr[j] = pc1m[l-28]
			}
		}

		for j = 28; j < 56; j++ {
			l = j + totrot[i]
			if l < 56 {
				pcr[j] = pc1m[l]
			} else {
				pcr[j] = pc1m[l-28]
			}
		}

		for j = 0; j < 24; j++ {
			if pcr[pc2[j]] {
				kn[m] |= bigbyte[j]
			}
			if pcr[pc2[j+24]] {
				kn[n] |= bigbyte[j]
			}
		}
	}
	cookey(&kn[0])
	return
}

func cookey(raw1 *uint32) {

}
