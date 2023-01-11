package board

import (
	"fmt"
	"math/rand"
)

var bishopBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

var rookBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

func maskMagicBishopAttacks(square int) uint64 {
	var attacks uint64
	targetRank := square / 8
	targetFile := square % 8

	for rank, file := targetRank+1, targetFile+1; rank < 7 && file < 7; rank, file = rank+1, file+1 {
		attacks |= uint64(1) << (rank*8 + file)
	}
	for rank, file := targetRank-1, targetFile+1; rank > 0 && file < 7; rank, file = rank-1, file+1 {
		attacks |= uint64(1) << (rank*8 + file)
	}
	for rank, file := targetRank+1, targetFile-1; rank < 7 && file > 0; rank, file = rank+1, file-1 {
		attacks |= uint64(1) << (rank*8 + file)
	}
	for rank, file := targetRank-1, targetFile-1; rank > 0 && file > 0; rank, file = rank-1, file-1 {
		attacks |= uint64(1) << (rank*8 + file)
	}

	return attacks
}

func maskMagicRookAttacks(square int) uint64 {
	var attacks uint64
	targetRank := square / 8
	targetFile := square % 8

	for rank := targetRank + 1; rank < 7; rank++ {
		attacks |= uint64(1) << (rank*8 + targetFile)
	}
	for rank := targetRank - 1; rank > 0; rank-- {
		attacks |= uint64(1) << (rank*8 + targetFile)
	}
	for file := targetFile + 1; file < 7; file++ {
		attacks |= uint64(1) << (targetRank*8 + file)
	}
	for file := targetFile - 1; file > 0; file-- {
		attacks |= uint64(1) << (targetRank*8 + file)
	}

	return attacks
}

func maskBishopAttacks(square int, blockers uint64) uint64 {
	var attacks uint64
	targetRank := square / 8
	targetFile := square % 8

	for rank, file := targetRank+1, targetFile+1; rank < 8 && file < 8; rank, file = rank+1, file+1 {
		attacks |= uint64(1) << (rank*8 + file)
		if (uint64(1)<<(rank*8+file))&blockers != 0 {
			break
		}
	}
	for rank, file := targetRank-1, targetFile+1; rank >= 0 && file < 8; rank, file = rank-1, file+1 {
		attacks |= uint64(1) << (rank*8 + file)
		if (uint64(1)<<(rank*8+file))&blockers != 0 {
			break
		}
	}
	for rank, file := targetRank+1, targetFile-1; rank < 8 && file >= 0; rank, file = rank+1, file-1 {
		attacks |= uint64(1) << (rank*8 + file)
		if (uint64(1)<<(rank*8+file))&blockers != 0 {
			break
		}
	}
	for rank, file := targetRank-1, targetFile-1; rank >= 0 && file >= 0; rank, file = rank-1, file-1 {
		attacks |= uint64(1) << (rank*8 + file)
		if (uint64(1)<<(rank*8+file))&blockers != 0 {
			break
		}
	}

	return attacks
}

func maskRookAttacks(square int, blockers uint64) uint64 {
	var attacks uint64
	targetRank := square / 8
	targetFile := square % 8

	for rank := targetRank + 1; rank < 8; rank++ {
		attacks |= uint64(1) << (rank*8 + targetFile)
		if (uint64(1)<<(rank*8+targetFile))&blockers != 0 {
			break
		}
	}
	for rank := targetRank - 1; rank >= 0; rank-- {
		attacks |= uint64(1) << (rank*8 + targetFile)
		if (uint64(1)<<(rank*8+targetFile))&blockers != 0 {
			break
		}
	}
	for file := targetFile + 1; file < 8; file++ {
		attacks |= uint64(1) << (targetRank*8 + file)
		if (uint64(1)<<(targetRank*8+file))&blockers != 0 {
			break
		}
	}
	for file := targetFile - 1; file >= 0; file-- {
		attacks |= uint64(1) << (targetRank*8 + file)
		if (uint64(1)<<(targetRank*8+file))&blockers != 0 {
			break
		}
	}

	return attacks
}

func GetBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= bishopMasks[square]
	occupancy *= bishopMagicNumber[square]
	occupancy >>= 64 - bishopBits[square]

	return bishopAttacks[square][occupancy]
}

func GetRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= rookMasks[square]
	occupancy *= rookMagicNumber[square]
	occupancy >>= 64 - rookBits[square]

	return rookAttacks[square][occupancy]
}

func GetQueenAttacks(square int, occupancy uint64) uint64 {
	return GetBishopAttacks(square, occupancy) | GetRookAttacks(square, occupancy)
}

// Functions from this point on is Tord Romstad's proposal to find magics,
// With his implementation written in C and converted to Go for use here
func getRandomUINT64() uint64 {
	var r1, r2, r3, r4 uint64

	r1 = (uint64)(rand.Uint32() & 0xFFFF)
	r2 = (uint64)(rand.Uint32() & 0xFFFF)
	r3 = (uint64)(rand.Uint32() & 0xFFFF)
	r4 = (uint64)(rand.Uint32() & 0xFFFF)

	return r1 | (r2 << 16) | (r3 << 32) | (r4 << 48)
}

func generateMagicNumberCandidate() uint64 {
	r1 := getRandomUINT64()
	r2 := getRandomUINT64()
	r3 := getRandomUINT64()
	return r1 & r2 & r3
}

func setMagicOccupancies(start int, end int, mask uint64) uint64 {
	occupancy := uint64(0)
	for i := 0; i < end; i++ {
		square := BitScanForward(mask)
		popBit(&mask, square)
		if start&(1<<i) != 0 {
			occupancy |= (1 << square)
		}
	}
	return occupancy
}

func findMagic(square int, bits int, isBishop bool) uint64 {
	var occupancies, attacks, usedAttacks [4096]uint64
	var mask, magic uint64

	if isBishop {
		mask = maskMagicBishopAttacks(square)
	} else {
		mask = maskMagicRookAttacks(square)
	}
	n := 1 << bits

	for i := 0; i < n; i++ {
		occupancies[i] = setMagicOccupancies(i, bits, mask)
		if isBishop {
			attacks[i] = maskBishopAttacks(square, occupancies[i])
		} else {
			attacks[i] = maskRookAttacks(square, occupancies[i])
		}
	}
	for i := 0; i < 100000000; i++ {
		magic = generateMagicNumberCandidate()
		if BitCount((mask*magic)&0xFF00000000000000) < 6 {
			continue
		}
		for j := 0; j < 4096; j++ {
			usedAttacks[j] = uint64(0)
		}
		fail := false
		for j := 0; !fail && j < n; j++ {
			magicIndex := (occupancies[j] * magic) >> (64 - bits)
			if usedAttacks[magicIndex] == uint64(0) {
				usedAttacks[magicIndex] = attacks[j]
			} else if usedAttacks[magicIndex] != attacks[j] {
				fail = true
			}
		}
		// if the magic number works then return it, may not be the best one
		if !fail {
			return magic
		}
	}
	// no magic number worked so tell the user that
	fmt.Print("No magic number worked")
	return uint64(0)
}

func MagicInit() {
	fmt.Printf("Rook magic numbers: \n")
	for square := 0; square < 64; square++ {
		fmt.Printf("%d,\n", findMagic(square, rookBits[square], false))
	}
	fmt.Printf("\n")
	fmt.Printf("Bishop magic numbers: \n")
	for square := 0; square < 64; square++ {
		fmt.Printf("%d,\n", findMagic(square, bishopBits[square], true))
	}
}

var rookMagicNumber = [64]uint64{
	36028934734757888,
	9241386504085831744,
	72075323667185672,
	36051096500699520,
	144150441200976904,
	4755803405594136577,
	2341943274522017920,
	2413949339656274176,
	369436045445464068,
	1161999210048917504,
	140806209929344,
	1153062276463333376,
	5066687288445440,
	4614078572909428864,
	20547681890074880,
	865254220147327108,
	9265171620721868802,
	9405920105197740032,
	9713635744030720,
	4612821814207325440,
	9241951584475422736,
	9223513324166054400,
	27585651524502016,
	2596336180379812932,
	2882374149588680736,
	3458835440912695560,
	576478346645348352,
	6485192261655068801,
	36037597408003328,
	5773619122483101824,
	3891189260069503248,
	577023985724752036,
	6949195237407326242,
	70370908451072,
	9302202624675684352,
	562984615157824,
	36169637611705345,
	1153062250693526528,
	9548774779494895880,
	36169673028403456,
	9223512775416905760,
	288265561865994242,
	144150372716413056,
	2310346750575575060,
	2256197994446976,
	4612248985829244936,
	1168264921096,
	4611721757366288388,
	9223662857682437376,
	1189232052025819648,
	873707124877689408,
	4688387983941828736,
	1154328896939229440,
	1163054620981887104,
	2306405963529716224,
	140741783339136,
	144133916280758337,
	2381296445551558914,
	2378465477620564066,
	9042452346571009,
	5764889556625653765,
	563139200942870,
	149561550045724,
	18577556785611778,
}
var bishopMagicNumber = [64]uint64{
	1172237727064784960,
	309631331208798336,
	4508000360597506,
	2256820630520852,
	1154612628685947393,
	13660472184164353,
	3459328030789353504,
	578994168845058112,
	4665799619224275072,
	4508004150018571,
	4683763438038876320,
	2305851917053592064,
	4505807509261440,
	1128103311052800,
	1206965018231644224,
	1268950302263316,
	2310634955832624128,
	20478412925960320,
	2251817127777296,
	167759120682752000,
	4644354373190936,
	1171006718689153024,
	9516177223258998080,
	13837389041425187138,
	4512395806377056,
	9246462337774322688,
	432644635991736384,
	73192290641907714,
	9446942884231127040,
	20838496341065856,
	1560783148206401664,
	2702900847276394624,
	77159465430173698,
	901863451931185544,
	2310840358287328256,
	6917531228820734082,
	4504700222178308,
	6057359092073243648,
	2323285249425536,
	285907987005696,
	2296918983395328,
	576605907702657544,
	576469688251385860,
	72515266323451936,
	316702441341953,
	721147892955414784,
	82208568992088576,
	37155409911645696,
	10376368326585632772,
	9799905365784133760,
	12402767181119744,
	1339859101542777344,
	5045439369925888101,
	36108040373928006,
	11547233846955741184,
	18652407379542528,
	577657160814563330,
	4612813045771804736,
	11565253365038465280,
	436889571117237280,
	72066390407907456,
	38281215342547202,
	4611694823383106050,
	18023228970008864,
}
