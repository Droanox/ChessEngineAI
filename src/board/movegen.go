package board

import "fmt"

func (cb ChessBoard) IsSquareAttackedBySide(square int, side int) bool {
	if side == White {
		switch {
		case (pawnAttacks[Black][square] & cb.WhitePawns) != EmptyBoard:
			return true
		case (knightAttacks[square] & cb.WhiteKnights) != EmptyBoard:
			return true
		case (GetBishopAttacks(square, cb.WhitePieces|cb.BlackPieces) & cb.WhiteBishops) != EmptyBoard:
			return true
		case (GetRookAttacks(square, cb.WhitePieces|cb.BlackPieces) & cb.WhiteRooks) != EmptyBoard:
			return true
		case (GetQueenAttacks(square, cb.WhitePieces|cb.BlackPieces) & cb.WhiteQueen) != EmptyBoard:
			return true
		case (kingAttacks[square] & cb.WhiteKing) != EmptyBoard:
			return true
		}
	}
	if side == Black {
		switch {
		case (pawnAttacks[White][square] & cb.BlackPawns) != EmptyBoard:
			return true
		case (knightAttacks[square] & cb.BlackKnights) != EmptyBoard:
			return true
		case (GetBishopAttacks(square, cb.WhitePieces|cb.BlackPieces) & cb.BlackBishops) != EmptyBoard:
			return true
		case (GetRookAttacks(square, cb.WhitePieces|cb.BlackPieces) & cb.BlackRooks) != EmptyBoard:
			return true
		case (GetQueenAttacks(square, cb.WhitePieces|cb.BlackPieces) & cb.BlackQueen) != EmptyBoard:
			return true
		case (kingAttacks[square] & cb.BlackKing) != EmptyBoard:
			return true
		}
	}
	return false
}

func (cb *ChessBoard) GenerateMoves() {
	var start, end, numMovesMade int
	var attacks uint64
	var allPieces uint64 = cb.WhitePieces | cb.BlackPieces
	var knights, bishops, rooks, queen, otherSide, target uint64

	if SideToMove == White {
		knights = cb.WhiteKnights
		bishops = cb.WhiteBishops
		rooks = cb.WhiteRooks
		queen = cb.WhiteQueen
		otherSide = cb.BlackPieces
		target = ^cb.WhitePieces
	} else {
		knights = cb.BlackKnights
		bishops = cb.BlackBishops
		rooks = cb.BlackRooks
		queen = cb.BlackQueen
		otherSide = cb.WhitePieces
		target = ^cb.BlackPieces
	}
	// White pawns and White castling moves
	if SideToMove == White {
		for bitboard := cb.WhitePawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			start = BitScanForward(bitboard)
			if !isBitOn(allPieces, start+8) {
				if start >= 48 && start <= 55 {
					numMovesMade += promotePiece(start, start+8)
				} else {
					fmt.Printf("%s%s  pawn push\n", IntToSquare[start], IntToSquare[start+8])
					numMovesMade++
					if (start >= 8 && start <= 15) && !isBitOn(allPieces, start+16) {
						fmt.Printf("%s%s  pawn double push\n", IntToSquare[start], IntToSquare[start+16])
						numMovesMade++
					}
				}
			}
			for attacks := pawnAttacks[White][start] & cb.BlackPieces; attacks != EmptyBoard; {
				end = BitScanForward(attacks)
				if start >= 48 && start <= 55 {
					numMovesMade += promotePiece(start, end)
				} else {
					fmt.Printf("%s%s  pawn capture\n", IntToSquare[start], IntToSquare[end])
					numMovesMade++
				}
				popBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[White][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					fmt.Printf("%s%s  Enpassant capture\n", IntToSquare[start], IntToSquare[BitScanForward(attacks)])
					numMovesMade++
				}
			}
		}
		// Castling moves
		if CastleRights&WhiteKingSide != 0 &&
			(allPieces&0x60) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e1"], Black) &&
			!cb.IsSquareAttackedBySide(SquareToInt["f1"], Black) {
			fmt.Printf("O-O\n")
			numMovesMade++
		}
		if CastleRights&WhiteQueenSide != 0 &&
			(allPieces&0xE) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e1"], Black) &&
			!cb.IsSquareAttackedBySide(SquareToInt["d1"], Black) {
			fmt.Printf("O-O-O\n")
			numMovesMade++
		}
		// Black pawn and Black castling moves
	} else {
		for bitboard := cb.BlackPawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			start = BitScanForward(bitboard)
			if !isBitOn(allPieces, start-8) {
				if start >= 8 && start <= 15 {
					numMovesMade += promotePiece(start, start-8)
				} else {
					fmt.Printf("%s%s  pawn push\n", IntToSquare[start], IntToSquare[start-8])
					numMovesMade++
					if (start >= 48 && start <= 55) && !isBitOn(allPieces, start-16) {
						fmt.Printf("%s%s  pawn double push\n", IntToSquare[start], IntToSquare[start-16])
						numMovesMade++
					}
				}
			}
			for attacks := pawnAttacks[Black][start] & cb.WhitePieces; attacks != EmptyBoard; {
				end = BitScanForward(attacks)
				if start >= 8 && start <= 15 {
					numMovesMade += promotePiece(start, end)
				} else {
					fmt.Printf("%s%s  pawn capture\n", IntToSquare[start], IntToSquare[end])
					numMovesMade++
				}
				popBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[Black][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					fmt.Printf("%s%s  Enpassant capture\n", IntToSquare[start], IntToSquare[BitScanForward(attacks)])
					numMovesMade++
				}
			}
		}
		// Castling moves
		if CastleRights&BlackKingSide != 0 &&
			(allPieces&0x6000000000000000) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e8"], White) &&
			!cb.IsSquareAttackedBySide(SquareToInt["f8"], White) {
			fmt.Printf("O-O\n")
			numMovesMade++
		}
		if CastleRights&BlackQueenSide != 0 &&
			(allPieces&0xE00000000000000) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e8"], White) &&
			!cb.IsSquareAttackedBySide(SquareToInt["d8"], White) {
			fmt.Printf("O-O-O\n")
			numMovesMade++
		}
	}
	// Generate knight moves
	for bitboard := knights; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := knightAttacks[start] & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				fmt.Printf("N%s%s  Piece capture\n", IntToSquare[start], IntToSquare[end])
			} else {
				fmt.Printf("N%s%s  Piece move\n", IntToSquare[start], IntToSquare[end])
			}
			numMovesMade++
		}
	}
	// Generate Bishop moves
	for bitboard := bishops; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetBishopAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				fmt.Printf("B%s%s  Piece capture\n", IntToSquare[start], IntToSquare[end])
			} else {
				fmt.Printf("B%s%s  Piece move\n", IntToSquare[start], IntToSquare[end])
			}
			numMovesMade++
		}
	}
	// Generate Rook moves
	for bitboard := rooks; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetRookAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				fmt.Printf("R%s%s  Piece capture\n", IntToSquare[start], IntToSquare[end])
			} else {
				fmt.Printf("R%s%s  Piece move\n", IntToSquare[start], IntToSquare[end])
			}
			numMovesMade++
		}
	}
	// Generate Queen moves
	for bitboard := queen; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetQueenAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				fmt.Printf("Q%s%s  Piece capture\n", IntToSquare[start], IntToSquare[end])
			} else {
				fmt.Printf("Q%s%s  Piece move\n", IntToSquare[start], IntToSquare[end])
			}
			numMovesMade++
		}
	}
}

func promotePiece(start int, end int) int {
	fmt.Printf("pawn promotion: %s%sQ\n", IntToSquare[start], IntToSquare[end])
	fmt.Printf("pawn promotion: %s%sR\n", IntToSquare[start], IntToSquare[end])
	fmt.Printf("pawn promotion: %s%sB\n", IntToSquare[start], IntToSquare[end])
	fmt.Printf("pawn promotion: %s%sN\n", IntToSquare[start], IntToSquare[end])
	return 4
}
