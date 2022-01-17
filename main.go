package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unsafe"
)

const NUM_ROWS = 4
const NUM_COLS = 6

var INITIAL_BOARD = Board{
	Rows: NUM_ROWS,
	Cols: NUM_COLS,
	PieceCoords: map[Piece][]Coord{
		TL: {{0, 0}, {0, 1}, {1, 0}},
		TR: {{0, 2}, {0, 3}, {1, 3}},
		BL: {{2, 0}, {3, 0}, {3, 1}},
		BR: {{2, 3}, {3, 2}, {3, 3}},
		SQ: {{1, 4}, {1, 5}, {2, 4}, {2, 5}},
	},
}

var (
	ErrCollision    error = fmt.Errorf("pieces collided")
	ErrInvalidCoord       = fmt.Errorf("invalid coord")
)

type Piece string

var (
	TL Piece = "TL"
	TR Piece = "TR"
	BL Piece = "BL"
	BR Piece = "BR"
	SQ Piece = "SQ"
)

func allPieces() []Piece {
	return []Piece{TL, TR, BL, BR, SQ}
}

type Dir string

var (
	DIR_UL Dir = "up-left"
	DIR_U  Dir = "up"
	DIR_UR Dir = "up-right"
	DIR_R  Dir = "right"
	DIR_DR Dir = "down-right"
	DIR_D  Dir = "down"
	DIR_DL Dir = "down-left"
	DIR_L  Dir = "left"
)

func (d Dir) Vec() Coord {
	switch d {
	case DIR_UL:
		return Coord{Row: -1, Col: -1}
	case DIR_U:
		return Coord{Row: -1, Col: 0}
	case DIR_UR:
		return Coord{Row: -1, Col: 1}
	case DIR_R:
		return Coord{Row: 0, Col: 1}
	case DIR_DR:
		return Coord{Row: 1, Col: 1}
	case DIR_D:
		return Coord{Row: 1, Col: 0}
	case DIR_DL:
		return Coord{Row: 1, Col: -1}
	case DIR_L:
		return Coord{Row: 0, Col: -1}
	default:
		panic("unexpected direction: " + d)
	}
}

func allDirs() []Dir {
	return []Dir{
		// DIR_UL,
		DIR_U,
		// DIR_UR,
		DIR_R,
		// DIR_DR,
		DIR_D,
		// DIR_DL,
		DIR_L,
	}
}

type Coord struct {
	Row int
	Col int
}

func (c Coord) Eq(c2 Coord) bool {
	return c.Row == c2.Row && c.Col == c2.Col
}

func (c Coord) Map(c2 Coord) Coord {
	return Coord{Row: c.Row + c2.Row, Col: c.Col + c2.Col}
}

type Board struct {
	Rows        int
	Cols        int
	PieceCoords map[Piece][]Coord
}

func (b *Board) EncodeLayout() string {
	ret := ""
	for _, piece := range allPieces() {
		enc := b.EncodeCoords(piece)
		ret += fmt.Sprintf("%x", enc)
	}
	return ret
}

func (b *Board) EncodeCoords(piece Piece) int32 {
	coords := b.PieceCoords[piece]
	bytes := []byte{}
	for _, coord := range coords {
		encoding := coord.Row*NUM_COLS + coord.Col
		bytes = append(bytes, *(*byte)(unsafe.Pointer(&encoding)))
	}
	for len(bytes) < 4 {
		bytes = append(bytes, 0)
	}
	return *(*int32)(unsafe.Pointer(&bytes[0]))
}

func (b *Board) Print() {
	matrix := [][]string{}
	for row := 0; row < b.Rows; row++ {
		r := []string{}
		for col := 0; col < b.Cols; col++ {
			r = append(r, "--")
		}
		matrix = append(matrix, r)
	}
	for piece, coords := range b.PieceCoords {
		for _, coord := range coords {
			matrix[coord.Row][coord.Col] = string(piece)
		}
	}
	for _, row := range matrix {
		fmt.Printf("| %s\n", strings.Join(row, " | "))
	}
	fmt.Printf("-----------------------------------\n")
}

func (b *Board) Validate() error {
	for _, coords := range b.PieceCoords {
		for _, coord := range coords {
			if !b.IsValidCoord(coord) {
				return ErrInvalidCoord
			}
			occupiers := b.AllPiecesOn(coord)
			if len(occupiers) > 1 {
				return ErrCollision
			}
		}
	}
	return nil
}

func (b *Board) MovePiece(piece Piece, dir Dir) error {
	currentCoords := b.PieceCoords[piece]
	newCoords := make([]Coord, len(currentCoords))
	vec := dir.Vec()
	for i, coord := range currentCoords {
		newCoord := coord.Map(vec)
		if !b.IsValidCoord(newCoord) {
			return ErrInvalidCoord
		}
		occupiers := b.AllPiecesOn(newCoord)
		if len(occupiers) > 1 || (len(occupiers) == 1 && occupiers[0] != piece) {
			return ErrCollision
		}
		// by, taken := b.Taken(newCoord)
		// if taken && by != piece {
		// 	return ErrCollision
		// }
		newCoords[i] = newCoord
	}
	b.PieceCoords[piece] = newCoords
	return nil
}

func (b *Board) IsValidCoord(c Coord) bool {
	row, col := c.Row, c.Col
	return row >= 0 && row < b.Rows && col >= 0 && col < b.Cols
}

func (b *Board) AllPiecesOn(c Coord) []Piece {
	ret := []Piece{}
	for piece, coords := range b.PieceCoords {
		for _, coord := range coords {
			if coord.Eq(c) {
				ret = append(ret, piece)
			}
		}
	}
	return ret
}

func (b *Board) Taken(c Coord) (Piece, bool) {
	for piece, coords := range b.PieceCoords {
		for _, coord := range coords {
			if coord.Eq(c) {
				return piece, true
			}
		}
	}
	return "", false
}

func (b *Board) Dup() *Board {
	newCoords := map[Piece][]Coord{}
	for _, piece := range allPieces() {
		src := b.PieceCoords[piece]
		dst := make([]Coord, len(src))
		copy(dst, src)
		newCoords[piece] = dst
	}
	return &Board{Rows: b.Rows, Cols: b.Cols, PieceCoords: newCoords}
}

func (b *Board) Eq(b2 *Board) bool {
	return b.EncodeLayout() == b2.EncodeLayout()
}

// Solved becomes true when the entire perimeter of the square
// is covered by other pieces.
func (b *Board) Solved() bool {
	squareCoords := b.PieceCoords[SQ]
	for _, dir := range allDirs() {
		for _, c := range squareCoords {
			newCoord := c.Map(dir.Vec())
			if !b.IsValidCoord(newCoord) {
				return false
			}
			_, taken := b.Taken(newCoord)
			if !taken {
				return false
			}
		}
	}
	return true
}

func WasVisited(s map[string]struct{}, hash string) bool {
	_, ok := s[hash]
	return ok
}

func MarkVisited(s map[string]struct{}, hash string) {
	s[hash] = struct{}{}
}

type Move struct {
	Piece Piece
	Dir   Dir
}
type Step struct {
	Board          *Board
	PrecedingMoves []Move
}

func Solve(board *Board, moves []Move) *Step {
	if moves == nil {
		moves = []Move{}
	}

	visited := map[string]struct{}{}
	queue := []Step{{Board: board, PrecedingMoves: moves}}

	for len(queue) > 0 {
		step := queue[0]
		queue = queue[1:]

		b := step.Board
		if b.Validate() != nil {
			continue
		}

		if b.Solved() {
			return &step
		}

		hash := b.EncodeLayout()
		if WasVisited(visited, hash) {
			continue
		}

		for _, dir := range allDirs() {
			for _, piece := range allPieces() {
				newBoard := b.Dup()
				err := newBoard.MovePiece(piece, dir)
				if err != nil {
					continue
				}
				newMoves := make([]Move, len(step.PrecedingMoves))
				copy(newMoves, step.PrecedingMoves)
				newMoves = append(newMoves, Move{Piece: piece, Dir: dir})
				queue = append(queue, Step{Board: newBoard, PrecedingMoves: newMoves})
			}
		}
		MarkVisited(visited, hash)
	}
	return nil
}

func main() {
	board := INITIAL_BOARD
	cp := board.Dup()
	winner := Solve(cp, nil)

	if winner == nil {
		fmt.Printf("Couldn't win this one.")
		return
	} else {
		fmt.Printf(
			"Found a solution in %d steps\n",
			len(winner.PrecedingMoves),
		)
	}

	for _, mov := range winner.PrecedingMoves {
		fmt.Printf("Move %s %s\n", mov.Piece, mov.Dir)
		board.MovePiece(mov.Piece, mov.Dir)
		board.Print()
	}
	bs, err := json.Marshal(winner)
	must(err)
	fmt.Printf("%s\n", string(bs))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
