package board

import (
	"fmt"
	"math/rand"
)

/*
Functions for generating magic bitboards for sliding pieces.
These functions are based on the following resources:
https://www.chessprogramming.org/Magic_Bitboards
https://www.youtube.com/watch?v=KqWeOVyOoyU
https://www.youtube.com/watch?v=UnEu5GOiSEs
https://www.youtube.com/watch?v=1lAM8ffBg0A
*/

// maskMagicBishopAttacks returns a bitboard of the squares attacked by a bishop
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

// maskMagicRookAttacks returns a bitboard of the squares attacked by a rook
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

// maskBishopAttacks returns a bitboard of the squares attacked by a bishop
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

// maskRookAttacks returns a bitboard of the squares attacked by a rook
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
	/*
		equivalent to :
		occupancy &= bishopMasks[square]
		occupancy *= bishopMagicNumber[square]
		occupancy >>= 64 - bishopBits[square]
	*/
	return bishopAttacks[square][((occupancy&bishopMasks[square])*bishopMagicNumber[square])>>(64-bishopBits[square])]
}

func GetRookAttacks(square int, occupancy uint64) uint64 {
	/*
		equivalent to:
		occupancy &= rookMasks[square]
		occupancy *= rookMagicNumber[square]
		occupancy >>= 64 - rookBits[square]
	*/
	return rookAttacks[square][((occupancy&rookMasks[square])*rookMagicNumber[square])>>(64-rookBits[square])]
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
		PopBit(&mask, square)
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

// MagicInit is a function that prints out the magic numbers for rooks and bishops
func MagicInit(seed int64) {
	rand.Seed(seed)
	fmt.Printf("var rookMagicNumber = [64]uint64{\n")
	for squareRank := 0; squareRank < 8; squareRank++ {
		for squareFile := 0; squareFile < 8; squareFile++ {
			fmt.Printf("%d, ", findMagic((8*squareRank)+squareFile, rookBits[(8*squareRank)+squareFile], false))
		}
		fmt.Println()
	}
	fmt.Printf("}\n")
	fmt.Printf("\n")
	fmt.Printf("var bishopMagicNumber = [64]uint64{\n")
	for squareRank := 0; squareRank < 8; squareRank++ {
		for squareFile := 0; squareFile < 8; squareFile++ {
			fmt.Printf("%d, ", findMagic((8*squareRank)+squareFile, bishopBits[(8*squareRank)+squareFile], true))
		}
		fmt.Println()
	}
	fmt.Printf("}")
}

// MagicInitWithReturn is the same as MagicInit but returns the magic numbers instead of printing them
// used to generate possibly better magics by comparing operation times
func MagicInitWithReturn(seed int64) {
	rand.Seed(seed)
	for square := 0; square < 64; square++ {
		rookMagicNumber[square] = findMagic(square, rookBits[square], false)
	}
	for square := 0; square < 8; square++ {
		bishopMagicNumber[square] = findMagic(square, bishopBits[square], true)
	}
}
