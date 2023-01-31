package board

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

type Move struct {
	Move  int
	Index int
}

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

func (cb ChessBoard) CopyBoard() {
	ChessBoardCopy = cb
	AspectsCopy = [5]int{SideToMove, CastleRights, Enpassant, HalfMoveClock, FullMoveCounter}
}

func (cb *ChessBoard) MakeBoard() {
	*cb = ChessBoardCopy

	SideToMove = AspectsCopy[0]
	CastleRights = AspectsCopy[1]
	Enpassant = AspectsCopy[2]
	HalfMoveClock = AspectsCopy[3]
	FullMoveCounter = AspectsCopy[4]
}

///////////////////////////////////////////////////////////////////
// Encoding/Decoding
///////////////////////////////////////////////////////////////////

func EncodeMove(start int, end int, startPiece int, endPiece int, flags int) int {
	return start | end<<6 | startPiece<<12 | endPiece<<16 | flags<<20
}
func (move Move) GetMoveStart() int {
	return (move.Move & 0x3f)
}
func (move Move) GetMoveEnd() int {
	return (move.Move & 0xfc0) >> 6
}
func (move Move) GetMoveStartPiece() int {
	return (move.Move & 0xf000) >> 12
}
func (move Move) GetMoveEndPiece() int {
	return (move.Move & 0xf0000) >> 16
}
func (move Move) GetMoveFlags() int {
	return move.Move >> 20
}

func AddMove(moveList *[]Move, move int) {
	*moveList = append(*moveList, Move{Move: move, Index: len(*moveList)})
}

///////////////////////////////////////////////////////////////////
// Making moves
///////////////////////////////////////////////////////////////////

func (cb *ChessBoard) MakeMove(move Move, moveListType int) int {
	cb.CopyBoard()

	start := move.GetMoveStart()
	end := move.GetMoveEnd()
	startPiece := move.GetMoveStartPiece() + (6 * SideToMove)
	CapturedPiece := move.GetMoveEndPiece() + (6 * (1 - SideToMove))
	flags := move.GetMoveFlags()

	pieceMap := map[int]*uint64{
		1: &cb.WhitePawns, 2: &cb.WhiteKnights, 3: &cb.WhiteBishops, 4: &cb.WhiteRooks, 5: &cb.WhiteQueen, 6: &cb.WhiteKing,
		7: &cb.BlackPawns, 8: &cb.BlackKnights, 9: &cb.BlackBishops, 10: &cb.BlackRooks, 11: &cb.BlackQueen, 12: &cb.BlackKing,
	}

	if moveListType == 0 {
		cb.CopyBoard()

		Enpassant = 64

		*pieceMap[startPiece] ^= indexMasks[start]
		setBit(pieceMap[startPiece], end)
		if (flags & MoveCaptures) != 0 {
			*pieceMap[CapturedPiece] ^= indexMasks[end]
		}
		if flags >= MoveKnightPromotion {
			*pieceMap[startPiece] ^= indexMasks[end]
			setBit(pieceMap[PromotionToPiece[flags]+(6*SideToMove)], end)
		}
		if flags == MoveEnpassantCapture {
			*pieceMap[CapturedPiece] ^= indexMasks[end]
			*pieceMap[CapturedPiece] ^= indexMasks[end+SideToOffset[SideToMove]]
		}
		if flags == MoveDoublePawn {
			Enpassant = end + SideToOffset[SideToMove]
		}
		switch flags {
		case MoveKingCastle:
			*pieceMap[startPiece-2] ^= (indexMasks[end+1] | indexMasks[end-1])
		case MoveQueenCastle:
			*pieceMap[startPiece-2] ^= (indexMasks[end-2] | indexMasks[end+1])
		}
		CastleRights &= CastleRightsUpdate[start]
		CastleRights &= CastleRightsUpdate[end]

		cb.WhitePieces = cb.WhiteRooks | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteQueen | cb.WhiteKing | cb.WhitePawns
		cb.BlackPieces = cb.BlackRooks | cb.BlackKnights | cb.BlackBishops | cb.BlackQueen | cb.BlackKing | cb.BlackPawns

		SideToMove ^= 1

		if cb.IsSquareAttackedBySide(BitScanForward(*pieceMap[6+(6*(1-SideToMove))]), SideToMove) {
			cb.MakeBoard()
			return 0
		}
		return 1

	} else {

	}
	return 1
}

///////////////////////////////////////////////////////////////////
// Move Generation
///////////////////////////////////////////////////////////////////

func (cb *ChessBoard) GenerateMoves(moveList *[]Move) {
	var start, end int
	var attacks uint64
	var allPieces uint64 = cb.WhitePieces | cb.BlackPieces
	var knights, bishops, rooks, queen, king, otherSide, target uint64

	if SideToMove == White {
		knights = cb.WhiteKnights
		bishops = cb.WhiteBishops
		rooks = cb.WhiteRooks
		queen = cb.WhiteQueen
		king = cb.WhiteKing
		otherSide = cb.BlackPieces
		target = ^cb.WhitePieces
	} else {
		knights = cb.BlackKnights
		bishops = cb.BlackBishops
		rooks = cb.BlackRooks
		queen = cb.BlackQueen
		king = cb.BlackKing
		otherSide = cb.WhitePieces
		target = ^cb.BlackPieces
	}
	// White pawns and White castling moves
	if SideToMove == White {
		for bitboard := cb.WhitePawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			start = BitScanForward(bitboard)
			if !isBitOn(allPieces, start+8) {
				if start >= 48 && start <= 55 {
					promotePiece(moveList, start, start+8, EmptyPiece, MoveQuiet)
				} else {
					AddMove(moveList, EncodeMove(start, start+8, Pawn, EmptyPiece, MoveQuiet))
					if (start >= 8 && start <= 15) && !isBitOn(allPieces, start+16) {
						AddMove(moveList, EncodeMove(start, start+16, Pawn, EmptyPiece, MoveDoublePawn))
					}
				}
			}
			for attacks := pawnAttacks[White][start] & cb.BlackPieces; attacks != EmptyBoard; {
				end = BitScanForward(attacks)
				if start >= 48 && start <= 55 {
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				popBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[White][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
		// Castling moves
		if CastleRights&WhiteKingSide != 0 &&
			(allPieces&0x60) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e1"], Black) &&
			!cb.IsSquareAttackedBySide(SquareToInt["f1"], Black) {
			AddMove(moveList, EncodeMove(SquareToInt["e1"], SquareToInt["g1"], King, EmptyPiece, MoveKingCastle))
		}
		if CastleRights&WhiteQueenSide != 0 &&
			(allPieces&0xE) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e1"], Black) &&
			!cb.IsSquareAttackedBySide(SquareToInt["d1"], Black) {
			AddMove(moveList, EncodeMove(SquareToInt["e1"], SquareToInt["c1"], King, EmptyPiece, MoveQueenCastle))
		}
		// Black pawn and Black castling moves
	} else {
		for bitboard := cb.BlackPawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			start = BitScanForward(bitboard)
			if !isBitOn(allPieces, start-8) {
				if start >= 8 && start <= 15 {
					promotePiece(moveList, start, start-8, EmptyPiece, MoveQuiet)
				} else {
					AddMove(moveList, EncodeMove(start, start-8, Pawn, EmptyPiece, MoveQuiet))
					if (start >= 48 && start <= 55) && !isBitOn(allPieces, start-16) {
						AddMove(moveList, EncodeMove(start, start-16, Pawn, EmptyPiece, MoveDoublePawn))
					}
				}
			}
			for attacks := pawnAttacks[Black][start] & cb.WhitePieces; attacks != EmptyBoard; {
				end = BitScanForward(attacks)
				if start >= 8 && start <= 15 {
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				popBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[Black][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
		// Castling moves
		if CastleRights&BlackKingSide != 0 &&
			(allPieces&0x6000000000000000) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e8"], White) &&
			!cb.IsSquareAttackedBySide(SquareToInt["f8"], White) {
			AddMove(moveList, EncodeMove(SquareToInt["e8"], SquareToInt["g8"], King, EmptyPiece, MoveKingCastle))
		}
		if CastleRights&BlackQueenSide != 0 &&
			(allPieces&0xE00000000000000) == 0 &&
			!cb.IsSquareAttackedBySide(SquareToInt["e8"], White) &&
			!cb.IsSquareAttackedBySide(SquareToInt["d8"], White) {
			AddMove(moveList, EncodeMove(SquareToInt["e8"], SquareToInt["c8"], King, EmptyPiece, MoveQueenCastle))
		}
	}
	// Generate knight moves
	for bitboard := knights; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := knightAttacks[start] & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Knight, cb.GetPieceType(end), MoveCaptures))
			} else {
				AddMove(moveList, EncodeMove(start, end, Knight, EmptyPiece, MoveQuiet))
			}
		}
	}
	// Generate Bishop moves
	for bitboard := bishops; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetBishopAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Bishop, cb.GetPieceType(end), MoveCaptures))
			} else {
				AddMove(moveList, EncodeMove(start, end, Bishop, EmptyPiece, MoveQuiet))
			}
		}
	}
	// Generate Rook moves
	for bitboard := rooks; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetRookAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Rook, cb.GetPieceType(end), MoveCaptures))
			} else {
				AddMove(moveList, EncodeMove(start, end, Rook, EmptyPiece, MoveQuiet))
			}
		}
	}
	// Generate Queen moves
	for bitboard := queen; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetQueenAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Queen, cb.GetPieceType(end), MoveCaptures))
			} else {
				AddMove(moveList, EncodeMove(start, end, Queen, EmptyPiece, MoveQuiet))
			}
		}
	}
	// Generate King moves
	for bitboard := king; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := kingAttacks[start] & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, King, cb.GetPieceType(end), MoveCaptures))
			} else {
				AddMove(moveList, EncodeMove(start, end, King, EmptyPiece, MoveQuiet))
			}
		}
	}
}

func (cb *ChessBoard) GenerateCaptures(moveList *[]Move) {
	var start, end int
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
			for attacks := pawnAttacks[White][start] & cb.BlackPieces; attacks != EmptyBoard; {
				end = BitScanForward(attacks)
				if start >= 48 && start <= 55 {
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				popBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[White][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
		// Black pawn and Black castling moves
	} else {
		for bitboard := cb.BlackPawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			start = BitScanForward(bitboard)
			for attacks := pawnAttacks[Black][start] & cb.WhitePieces; attacks != EmptyBoard; {
				end = BitScanForward(attacks)
				if start >= 8 && start <= 15 {
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				popBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[Black][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
	}
	// Generate knight moves
	for bitboard := knights; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := knightAttacks[start] & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Knight, cb.GetPieceType(end), MoveCaptures))
			}
		}
	}
	// Generate Bishop moves
	for bitboard := bishops; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetBishopAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Bishop, cb.GetPieceType(end), MoveCaptures))
			}
		}
	}
	// Generate Rook moves
	for bitboard := rooks; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetRookAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Rook, cb.GetPieceType(end), MoveCaptures))
			}
		}
	}
	// Generate Queen moves
	for bitboard := queen; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := GetQueenAttacks(start, allPieces) & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				AddMove(moveList, EncodeMove(start, end, Queen, cb.GetPieceType(end), MoveCaptures))
			}
		}
	}
}

func promotePiece(moveList *[]Move, start int, end int, captured int, promoteFlag int) {
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveKnightPromotion))
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveBishopPromotion))
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveRookPromotion))
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveQueenPromotion))
}