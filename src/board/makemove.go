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

	if Enpassant != 64 {
		HashKey ^= enpassantKeys[Enpassant]
	}

	// Reset enpassant
	Enpassant = 64

	// Move piece to new square and update hash
	*pieceArr[startPiece] ^= indexMasks[start] | indexMasks[end]
	HashKey ^= pieceKeys[startPiece-1][start] ^ pieceKeys[startPiece-1][end]

	// Capture
	if (flags & MoveCaptures) != 0 {
		*pieceArr[CapturedPiece] ^= indexMasks[end]
		HashKey ^= pieceKeys[CapturedPiece-1][end]
	}
	// Promotion
	if flags >= MoveKnightPromotion {
		*pieceArr[startPiece] ^= indexMasks[end]
		setBit(pieceArr[PromotionToPiece[flags]+(6*SideToMove)], end)
		HashKey ^= pieceKeys[startPiece-1][end] ^ pieceKeys[PromotionToPiece[flags]+(6*SideToMove)-1][end]
	}
	// Enpassant capture
	if flags == MoveEnpassantCapture {
		*pieceArr[CapturedPiece] ^= indexMasks[end] | indexMasks[end+offsetBySide[SideToMove]]
		HashKey ^= pieceKeys[CapturedPiece-1][end] ^ pieceKeys[CapturedPiece-1][end+offsetBySide[SideToMove]]
	}
	// Double pawn push
	if flags == MoveDoublePawn {
		Enpassant = end + offsetBySide[SideToMove]
		HashKey ^= enpassantKeys[Enpassant]
	}
	// Castling
	switch flags {
	case MoveKingCastle:
		*pieceArr[startPiece-2] ^= (indexMasks[end+1] | indexMasks[end-1])
		HashKey ^= pieceKeys[startPiece-3][end+1] ^ pieceKeys[startPiece-3][end-1]
	case MoveQueenCastle:
		*pieceArr[startPiece-2] ^= (indexMasks[end-2] | indexMasks[end+1])
		HashKey ^= pieceKeys[startPiece-3][end-2] ^ pieceKeys[startPiece-3][end+1]
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
	HashKey ^= castleKeys[CastleRights]
	CastleRights &= castleRightsUpdate[start]
	CastleRights &= castleRightsUpdate[end]
	HashKey ^= castleKeys[CastleRights]

	// Switch side to move and update hash
	SideToMove = 1 - SideToMove
	HashKey ^= sideKey

	// Uncomment this to check for hash mismatch
	/*
		var checkHash uint64 = GenHash(*cb)
		if checkHash != HashKey {
			cb.PrintChessBoard()
			fmt.Print("MakeMove: Hash mismatch\n")
			fmt.Printf("%0x    Received\n", checkHash)
			fmt.Printf("%0x    Expected\n", HashKey)
			cb.MakeBoard()
			cb.PrintChessBoard()
			os.Exit(1)
		}
	*/
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

func (cb *ChessBoard) MakeMoveNull() {
	cb.CopyBoard()

	// Reset enpassant
	Enpassant = 64

	// Switch side to move
	SideToMove = 1 - SideToMove
}
