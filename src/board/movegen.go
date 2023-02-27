package board

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

// Move represents a move in the game
type Move struct {
	Move  int
	Score int
}

// IsSquareAttackedBySide returns true if the square is attacked by the side
func (cb ChessBoard) IsSquareAttackedBySide(square int, side int) bool {
	both := cb.WhitePieces | cb.BlackPieces
	if side == White {
		switch {
		case (pawnAttacks[Black][square] & cb.WhitePawns) != EmptyBoard:
			return true
		case (knightAttacks[square] & cb.WhiteKnights) != EmptyBoard:
			return true
		case (GetBishopAttacks(square, both) & cb.WhiteBishops) != EmptyBoard:
			return true
		case (GetRookAttacks(square, both) & cb.WhiteRooks) != EmptyBoard:
			return true
		case (GetQueenAttacks(square, both) & cb.WhiteQueen) != EmptyBoard:
			return true
		case (kingAttacks[square] & cb.WhiteKing) != EmptyBoard:
			return true
		}
	} else {
		switch {
		case (pawnAttacks[White][square] & cb.BlackPawns) != EmptyBoard:
			return true
		case (knightAttacks[square] & cb.BlackKnights) != EmptyBoard:
			return true
		case (GetBishopAttacks(square, both) & cb.BlackBishops) != EmptyBoard:
			return true
		case (GetRookAttacks(square, both) & cb.BlackRooks) != EmptyBoard:
			return true
		case (GetQueenAttacks(square, both) & cb.BlackQueen) != EmptyBoard:
			return true
		case (kingAttacks[square] & cb.BlackKing) != EmptyBoard:
			return true
		}
	}
	return false
}

// CopyBoard copies the current board to the next ply
func (oldBoard ChessBoard) CopyBoard() {
	chessBoardCopies[Ply+1] = oldBoard
	aspectsCopies[Ply+1] = [5]int{SideToMove, CastleRights, Enpassant, HalfMoveClock, FullMoveCounter}
	Ply++
}

// MakeBoard copies the current board to the current ply
func (newBoard *ChessBoard) MakeBoard() {
	*newBoard = chessBoardCopies[Ply]

	SideToMove = aspectsCopies[Ply][0]
	CastleRights = aspectsCopies[Ply][1]
	Enpassant = aspectsCopies[Ply][2]
	HalfMoveClock = aspectsCopies[Ply][3]
	FullMoveCounter = aspectsCopies[Ply][4]
	Ply--
}

///////////////////////////////////////////////////////////////////
// Encoding/Decoding
///////////////////////////////////////////////////////////////////

// EncodeMove encodes a move into an int
func EncodeMove(start int, end int, startPiece int, capturedPiece int, flags int) int {
	return start | end<<6 | startPiece<<12 | capturedPiece<<16 | flags<<20
}

// GetMoveStart decodes the start square of a move
func (move Move) GetMoveStart() int {
	return (move.Move & 0x3f)
}

// GetMoveEnd decodes the end square of a move
func (move Move) GetMoveEnd() int {
	return (move.Move & 0xfc0) >> 6
}

// GetMoveStartPiece decodes the start piece of a move
func (move Move) GetMoveStartPiece() int {
	return (move.Move & 0xf000) >> 12
}

// GetMoveCapturedPiece decodes the captured piece of a move
func (move Move) GetMoveCapturedPiece() int {
	return (move.Move & 0xf0000) >> 16
}

// GetMoveFlags decodes the flags of a move
// (see constsvars.go for the encodings)
func (move Move) GetMoveFlags() int {
	return move.Move >> 20
}

// AddMove adds a move to the move list
func AddMove(moveList *[]Move, move int) {
	*moveList = append(*moveList, Move{Move: move})
}

///////////////////////////////////////////////////////////////////
// Making moves
///////////////////////////////////////////////////////////////////

// MakeMove makes a move on the board
func (cb *ChessBoard) MakeMove(move Move) bool {
	cb.CopyBoard()

	start := move.GetMoveStart()
	end := move.GetMoveEnd()
	startPiece := move.GetMoveStartPiece() + (6 * SideToMove)
	CapturedPiece := move.GetMoveCapturedPiece() + (6 * (1 - SideToMove))
	flags := move.GetMoveFlags()

	pieceArr := []*uint64{
		1: &cb.WhitePawns, 2: &cb.WhiteKnights, 3: &cb.WhiteBishops, 4: &cb.WhiteRooks, 5: &cb.WhiteQueen, 6: &cb.WhiteKing,
		7: &cb.BlackPawns, 8: &cb.BlackKnights, 9: &cb.BlackBishops, 10: &cb.BlackRooks, 11: &cb.BlackQueen, 12: &cb.BlackKing,
	}

	// Reset enpassant
	Enpassant = 64

	*pieceArr[startPiece] ^= indexMasks[start] | indexMasks[end]
	// Capture
	if (flags & MoveCaptures) != 0 {
		*pieceArr[CapturedPiece] ^= indexMasks[end]
	}
	// Promotion
	if flags >= MoveKnightPromotion {
		*pieceArr[startPiece] ^= indexMasks[end]
		setBit(pieceArr[PromotionToPiece[flags]+(6*SideToMove)], end)
	}
	// Enpassant capture
	if flags == MoveEnpassantCapture {
		*pieceArr[CapturedPiece] ^= indexMasks[end] | indexMasks[end+offsetBySide[SideToMove]]
	}
	// Double pawn push
	if flags == MoveDoublePawn {
		Enpassant = end + offsetBySide[SideToMove]
	}
	// Castling
	switch flags {
	case MoveKingCastle:
		*pieceArr[startPiece-2] ^= (indexMasks[end+1] | indexMasks[end-1])
	case MoveQueenCastle:
		*pieceArr[startPiece-2] ^= (indexMasks[end-2] | indexMasks[end+1])
	}

	// Update bitboards
	cb.WhitePieces = cb.WhitePawns | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteRooks | cb.WhiteQueen | cb.WhiteKing
	cb.BlackPieces = cb.BlackPawns | cb.BlackKnights | cb.BlackBishops | cb.BlackRooks | cb.BlackQueen | cb.BlackKing

	// Check for check
	if cb.IsSquareAttackedBySide(BitScanForward(*pieceArr[6+(6*SideToMove)]), 1-SideToMove) {
		cb.MakeBoard()
		return false
	}

	// Update castling rights
	CastleRights &= castleRightsUpdate[start]
	CastleRights &= castleRightsUpdate[end]
	// Switch side to move
	SideToMove = 1 - SideToMove

	return true
}

// MakeCapture makes a capture on the board
// (used for quiescence search)
func (cb *ChessBoard) MakeCapture(move Move) bool {
	cb.CopyBoard()

	start := move.GetMoveStart()
	end := move.GetMoveEnd()
	startPiece := move.GetMoveStartPiece() + (6 * SideToMove)
	CapturedPiece := move.GetMoveCapturedPiece() + (6 * (1 - SideToMove))
	flags := move.GetMoveFlags()

	pieceArr := []*uint64{
		1: &cb.WhitePawns, 2: &cb.WhiteKnights, 3: &cb.WhiteBishops, 4: &cb.WhiteRooks, 5: &cb.WhiteQueen, 6: &cb.WhiteKing,
		7: &cb.BlackPawns, 8: &cb.BlackKnights, 9: &cb.BlackBishops, 10: &cb.BlackRooks, 11: &cb.BlackQueen, 12: &cb.BlackKing,
	}

	// Reset enpassant
	Enpassant = 64

	*pieceArr[startPiece] ^= indexMasks[start] | indexMasks[end]
	// Capture
	*pieceArr[CapturedPiece] ^= indexMasks[end]
	// Promotion
	if flags >= MoveKnightPromotionCapture {
		*pieceArr[startPiece] ^= indexMasks[end]
		setBit(pieceArr[PromotionToPiece[flags]+(6*SideToMove)], end)
	}
	// Enpassant capture
	if flags == MoveEnpassantCapture {
		*pieceArr[CapturedPiece] ^= indexMasks[end] | indexMasks[end+offsetBySide[SideToMove]]
	}

	// Update bitboards
	cb.WhitePieces = cb.WhitePawns | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteRooks | cb.WhiteQueen | cb.WhiteKing
	cb.BlackPieces = cb.BlackPawns | cb.BlackKnights | cb.BlackBishops | cb.BlackRooks | cb.BlackQueen | cb.BlackKing

	// Check for check
	if cb.IsSquareAttackedBySide(BitScanForward(*pieceArr[6+(6*SideToMove)]), 1-SideToMove) {
		cb.MakeBoard()
		return false
	}

	// Switch side to move
	SideToMove = 1 - SideToMove

	return true
}

///////////////////////////////////////////////////////////////////
// Move Generation
///////////////////////////////////////////////////////////////////

// GenerateMoves generates all pseudo-legal moves for the current position
// and stores them in the moveList slice (which is passed by reference)
func (cb *ChessBoard) GenerateMoves(moveList *[]Move) {
	var start, end int
	var attacks, knights, bishops, rooks, queen, king, otherSide, target uint64
	var allPieces uint64 = cb.WhitePieces | cb.BlackPieces

	// Update vars based on side to move
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

	// Generate pawn moves
	if SideToMove == White {
		// generate white pawn moves
		for bitboard := cb.WhitePawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			// get start square (next bit)
			start = BitScanForward(bitboard)
			if !isBitOn(allPieces, start+8) {
				if (indexMasks[start] & Rank7On) != EmptyBoard {
					// promotion
					promotePiece(moveList, start, start+8, EmptyPiece, MoveQuiet)
				} else {
					// single pawn push
					AddMove(moveList, EncodeMove(start, start+8, Pawn, EmptyPiece, MoveQuiet))
					if ((indexMasks[start] & Rank2On) != EmptyBoard) && !isBitOn(allPieces, start+16) {
						// double pawn push
						AddMove(moveList, EncodeMove(start, start+16, Pawn, EmptyPiece, MoveDoublePawn))
					}
				}
			}
			for attacks := pawnAttacks[White][start] & cb.BlackPieces; attacks != EmptyBoard; {
				// get end square (next bit)
				end = BitScanForward(attacks)
				if (indexMasks[start] & Rank7On) != EmptyBoard {
					// promotion capture
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					// pawn capture
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				// remove bit to allow next iteration
				PopBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[White][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					// enpassant capture
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
		// Generate white castling moves for king side
		if CastleRights&whiteKingSide != 0 &&
			// Check if there are pieces between king and rook
			(allPieces&0x60) == 0 &&
			// Check if king and rook are not attacked
			!cb.IsSquareAttackedBySide(4, Black) &&
			!cb.IsSquareAttackedBySide(5, Black) {
			// King side castle
			AddMove(moveList, EncodeMove(4, 6, King, EmptyPiece, MoveKingCastle))
		}
		// Generate white castling moves for queen side
		if CastleRights&whiteQueenSide != 0 &&
			// Check if there are pieces between king and rook
			(allPieces&0xE) == 0 &&
			// Check if king and rook are not attacked
			!cb.IsSquareAttackedBySide(4, Black) &&
			!cb.IsSquareAttackedBySide(3, Black) {
			// Queen side castle
			AddMove(moveList, EncodeMove(4, 2, King, EmptyPiece, MoveQueenCastle))
		}
	} else {
		// Generate black pawn moves
		for bitboard := cb.BlackPawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			// get start square (next bit)
			start = BitScanForward(bitboard)
			if !isBitOn(allPieces, start-8) {
				if (indexMasks[start] & Rank2On) != EmptyBoard {
					// promotion
					promotePiece(moveList, start, start-8, EmptyPiece, MoveQuiet)
				} else {
					// single pawn push
					AddMove(moveList, EncodeMove(start, start-8, Pawn, EmptyPiece, MoveQuiet))
					if ((indexMasks[start] & Rank7On) != EmptyBoard) && !isBitOn(allPieces, start-16) {
						// double pawn push
						AddMove(moveList, EncodeMove(start, start-16, Pawn, EmptyPiece, MoveDoublePawn))
					}
				}
			}
			for attacks := pawnAttacks[Black][start] & cb.WhitePieces; attacks != EmptyBoard; {
				// get end square (next bit)
				end = BitScanForward(attacks)
				if (indexMasks[start] & Rank2On) != EmptyBoard {
					// promotion capture
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					// pawn capture
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				// remove bit to allow next iteration
				PopBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[Black][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					// enpassant capture
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
		// Generate black castling moves for king side
		if CastleRights&blackKingSide != 0 &&
			// Check if there are pieces between king and rook
			(allPieces&0x6000000000000000) == 0 &&
			// Check if king and rook are not attacked
			!cb.IsSquareAttackedBySide(60, White) &&
			!cb.IsSquareAttackedBySide(61, White) {
			// King side castle
			AddMove(moveList, EncodeMove(60, 62, King, EmptyPiece, MoveKingCastle))
		}
		// Generate black castling moves for queen side
		if CastleRights&blackQueenSide != 0 &&
			// Check if there are pieces between king and rook
			(allPieces&0xE00000000000000) == 0 &&
			// Check if king and rook are not attacked
			!cb.IsSquareAttackedBySide(60, White) &&
			!cb.IsSquareAttackedBySide(59, White) {
			// Queen side castle
			AddMove(moveList, EncodeMove(60, 58, King, EmptyPiece, MoveQueenCastle))
		}
	}
	// Generate knight moves
	for bitboard := knights; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := knightAttacks[start] & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				// capture
				AddMove(moveList, EncodeMove(start, end, Knight, cb.GetPieceType(end), MoveCaptures))
			} else {
				// quiet move
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
				// capture
				AddMove(moveList, EncodeMove(start, end, Bishop, cb.GetPieceType(end), MoveCaptures))
			} else {
				// quiet move
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
				// capture
				AddMove(moveList, EncodeMove(start, end, Rook, cb.GetPieceType(end), MoveCaptures))
			} else {
				// quiet move
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
				// capture
				AddMove(moveList, EncodeMove(start, end, Queen, cb.GetPieceType(end), MoveCaptures))
			} else {
				// quiet move
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
				// capture
				AddMove(moveList, EncodeMove(start, end, King, cb.GetPieceType(end), MoveCaptures))
			} else {
				// quiet move
				AddMove(moveList, EncodeMove(start, end, King, EmptyPiece, MoveQuiet))
			}
		}
	}
}

// GenerateCaptures generates all captures for the current side to move
// to be used in quiescence search
//
// Separate function to avoid generating all moves and then filtering out the captures
func (cb *ChessBoard) GenerateCaptures(moveList *[]Move) {
	var start, end int
	var attacks, knights, bishops, rooks, queen, king, otherSide, target uint64
	var allPieces uint64 = cb.WhitePieces | cb.BlackPieces

	// Update vars based on side to move
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

	// White pawn moves
	if SideToMove == White {
		for bitboard := cb.WhitePawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			// Get the start square (next bit)
			start = BitScanForward(bitboard)
			for attacks := pawnAttacks[White][start] & cb.BlackPieces; attacks != EmptyBoard; {
				// Get the end square (next bit)
				end = BitScanForward(attacks)
				if (indexMasks[start] & Rank7On) != EmptyBoard {
					// promotion capture
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					// capture
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				// remove bit to allow next iteration
				PopBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[White][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					// enpassant capture
					AddMove(moveList, EncodeMove(start, BitScanForward(attacks), Pawn, Pawn, MoveEnpassantCapture))
				}
			}
		}
		// Black pawn moves
	} else {
		for bitboard := cb.BlackPawns; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			// Get the start square (next bit)
			start = BitScanForward(bitboard)
			for attacks := pawnAttacks[Black][start] & cb.WhitePieces; attacks != EmptyBoard; {
				// Get the end square (next bit)
				end = BitScanForward(attacks)
				if (indexMasks[start] & Rank2On) != EmptyBoard {
					// promotion capture
					promotePiece(moveList, start, end, cb.GetPieceType(end), MoveCaptures)
				} else {
					// capture
					AddMove(moveList, EncodeMove(start, end, Pawn, cb.GetPieceType(end), MoveCaptures))
				}
				// remove bit to allow next iteration
				PopBit(&attacks, end)
			}
			if Enpassant != -1 {
				attacks = pawnAttacks[Black][start] & (1 << Enpassant)
				if attacks != EmptyBoard {
					// enpassant capture
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
				// capture
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
				// capture
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
				// capture
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
				// capture
				AddMove(moveList, EncodeMove(start, end, Queen, cb.GetPieceType(end), MoveCaptures))
			}
		}
	}
	// Generate King moves
	for bitboard := king; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
		start = BitScanForward(bitboard)
		for attacks := kingAttacks[start] & target; attacks != 0; attacks &= attacks - 1 {
			end = BitScanForward(attacks)
			if isBitOn(otherSide, end) {
				// capture
				AddMove(moveList, EncodeMove(start, end, King, cb.GetPieceType(end), MoveCaptures))
			}
		}
	}
}

// Generate all promotions for a pawn move
// Idea from CounterGo:
// https://github.com/ChizhovVadim/CounterGo/blob/master/pkg/common/movegen.go
func promotePiece(moveList *[]Move, start int, end int, captured int, promoteFlag int) {
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveKnightPromotion))
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveBishopPromotion))
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveRookPromotion))
	AddMove(moveList, EncodeMove(start, end, Pawn, captured, promoteFlag|MoveQueenPromotion))
}
