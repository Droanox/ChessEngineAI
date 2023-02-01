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
	rand.Seed(50000)
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
	1477180952748916768,
	144133338611851524,
	108095805627302272,
	324267969532200064,
	7241792600066261002,
	108096329619604480,
	288234843001728136,
	11673334958619001090,
	140877074808840,
	36380709480308736,
	9801662513954230272,
	1266946642280576,
	288793343554093088,
	7138768395915231745,
	1125934334738690,
	36310280591933696,
	31666484707067936,
	5332298243764158720,
	142937052680706,
	9223945983031775264,
	9242513439078746112,
	1153204079129001992,
	108090789372363010,
	108088590147272833,
	4611756406498951200,
	45053591680991264,
	585476748725977154,
	49540156394832432,
	379437066846732928,
	14357424983244928,
	432354377508995344,
	4611968051749912676,
	585538320914712713,
	9227875774193745920,
	54184070464147457,
	3459755208165298176,
	154252696029759488,
	5800777074746593792,
	4901043429498163202,
	295549017637519489,
	288371114713972776,
	18014952565178400,
	4039764050391892096,
	70403674472465,
	85587084751994885,
	2306124570092306456,
	288239515867545602,
	288239464311029809,
	35734199272192,
	81100528494904576,
	1152994072912232960,
	1407409511997696,
	74309947969505408,
	9011597334839424,
	4925820682436736,
	9223442682641383936,
	9300214988973023250,
	9529654196001839234,
	9809121678813356033,
	2449975789610082561,
	9223935004525465602,
	1688888917889062,
	9232467199221301508,
	11530341084496660746,
}
var bishopMagicNumber = [64]uint64{
	2252916639417600,
	72631548335227904,
	4612838343137430528,
	4789474873606248,
	565217701396736,
	4789482316890224,
	290305496846848,
	4611758620831518728,
	1765762910602397336,
	1129207165944064,
	1455806174048835584,
	9223521589780816257,
	18016606391107584,
	9224499053491020832,
	1190675439801074729,
	229685789154608128,
	585468227040649761,
	4625759904736872981,
	2612162585327575073,
	2252091906596864,
	4684869521249403136,
	1134704623420424,
	10134269553029188,
	2379167276181000320,
	144749056664536064,
	761257870883881220,
	11552876570892861520,
	5770245818670514240,
	4612816316452110401,
	5919805148268544,
	2738787257612240960,
	4755873791459070216,
	2287001551438848,
	650779569381704705,
	1152939440399188098,
	145452763303117312,
	18032541005119520,
	38852385963053056,
	13629547480154384,
	2308112964124704832,
	1163420808253572,
	4909244655591383056,
	18085042676830208,
	71606634286336,
	2320550145617314337,
	81073658711244864,
	12106846263400464640,
	9297646569160770,
	9512732719584383016,
	282127885287424,
	1166434504727070881,
	577586656572932224,
	10268207425419086080,
	153267591627014145,
	9261124879449776192,
	4901121554571724032,
	9571271285973008,
	18018951213089794,
	595601086220600352,
	4899916669465658388,
	11316173943611904,
	18032025089368576,
	590043089840701504,
	578785190218956928,
}
