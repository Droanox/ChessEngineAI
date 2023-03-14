package board

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
