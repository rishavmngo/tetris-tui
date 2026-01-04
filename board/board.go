// Package board
package board

import (
	"fmt"
	"math/rand"
	"strings"

	"rishavmngo/tetris-tui-v2/shape"
	cells "rishavmngo/tetris-tui-v2/shared"
)

type Board struct {
	Grid        [20][10]cells.Cell
	ActiveShape *shape.Shape
	AShapeX     int
	AShapeY     int
	Height      int
	Width       int
	GameOver    bool
	Score       int
}

func NewBoard() *Board {
	return &Board{
		Grid:        [20][10]cells.Cell{},
		ActiveShape: nil,
		Height:      20,
		Width:       10,
		AShapeX:     0,
		AShapeY:     0,
		GameOver:    false,
		Score:       0,
	}
}

func (b *Board) isValid(shape *shape.Shape, x, y int) bool {
	for i, row := range shape.Grid {
		for j, col := range row {

			if !col {
				continue
			}

			r := x + i
			c := y + j

			if r >= len(b.Grid) || c >= len(b.Grid[0]) {
				return false
			}

			if r < 0 || c < 0 {
				return false
			}

			if b.Grid[r][c] != cells.Empty {
				return false
			}

		}
	}
	return true
}

func (b *Board) LockPosition() {
	for i, row := range b.ActiveShape.Grid {
		for j, col := range row {

			if !col {
				continue
			}
			r := b.AShapeX + i
			c := b.AShapeY + j
			b.Grid[r][c] = b.ActiveShape.ShapeType
		}
	}

	b.ClearLines()
	b.SpwanPiece()
}

func (b *Board) ClearLines() {
	linesCleared := 0
	for i := len(b.Grid) - 1; i >= 0; i-- {

		full := true
		for _, cell := range b.Grid[i] {
			if cell == cells.Empty {
				full = false
				break
			}
		}

		if full {
			for k := i; k > 0; k-- {
				b.Grid[k] = b.Grid[k-1]
			}
			b.Grid[0] = [10]cells.Cell{}
			i++
			linesCleared++
		}
	}
	b.Score = b.Score + (linesCleared * 100)
}

func (b *Board) SpwanPiece() {
	pieces := []shape.ShapeType{
		shape.TypeI, shape.TypeO, shape.TypeJ, shape.TypeL, shape.TypeS, shape.TypeT, shape.TypeZ,
	}
	newPieceType := pieces[rand.Intn(len(pieces))]
	sh := shape.NewShape(newPieceType)
	if !b.isValid(sh, 0, 3) {
		b.GameOver = true
		return
	}

	b.ActiveShape = sh
	b.AShapeX = 0
	b.AShapeY = 3
}

func (b *Board) MoveLeft() {
	if b.isValid(b.ActiveShape, b.AShapeX, b.AShapeY-1) {
		b.AShapeY--
	}
}

func (b *Board) MoveRight() {
	if b.isValid(b.ActiveShape, b.AShapeX, b.AShapeY+1) {
		b.AShapeY++
	}
}

func (b *Board) MoveDown() bool {
	if b.isValid(b.ActiveShape, b.AShapeX+1, b.AShapeY) {
		b.AShapeX++
		return true
	}
	return false
}

func (b *Board) Rotate() {
	sh := b.ActiveShape.Clone()

	sh.RotateClockWise()

	if b.isValid(sh, b.AShapeX, b.AShapeY) {
		b.ActiveShape = sh
		return
	}

	kicks := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}

	for _, kick := range kicks {
		if b.isValid(sh, b.AShapeX+kick[0], b.AShapeY+kick[1]) {
			b.ActiveShape = sh
			b.AShapeX = b.AShapeX + kick[0]
			b.AShapeY = b.AShapeY + kick[1]
			return
		}
	}
}

func (b *Board) GetCompositeLayer() [20][10]cells.Cell {
	composite := b.Grid

	if b.ActiveShape != nil {
		for i, row := range b.ActiveShape.Grid {
			for j, col := range row {

				if !col {
					continue
				}
				r := b.AShapeX + i
				c := b.AShapeY + j
				composite[r][c] = b.ActiveShape.ShapeType
			}
		}
	}

	return composite
}

func (b *Board) Render() {
	composite := b.GetCompositeLayer()

	var sb strings.Builder
	sb.WriteString("\033[H") // Move cursor to top

	// Top border
	sb.WriteString("+")
	sb.WriteString(strings.Repeat("-", b.Width*2))
	sb.WriteString("+\r\n")

	for _, row := range composite {
		sb.WriteString("|")
		for _, cell := range row {
			if cell == cells.Empty {
				sb.WriteString(". ")
			} else {
				// Color based on piece type
				sb.WriteString(colorize(cell))
			}
		}
		sb.WriteString("|\r\n")
	}

	// Bottom border + stats
	sb.WriteString("+")
	sb.WriteString(strings.Repeat("-", b.Width*2))
	sb.WriteString("+\r\n")
	sb.WriteString(fmt.Sprintf("Score: %d \r\n", b.Score))

	fmt.Print(sb.String())
}

func colorize(cell cells.Cell) string {
	colors := map[cells.Cell]string{
		cells.TypeI: "\033[36m[]\033[0m",       // Cyan
		cells.TypeO: "\033[33m[]\033[0m",       // Yellow
		cells.TypeT: "\033[35m[]\033[0m",       // Magenta
		cells.TypeS: "\033[32m[]\033[0m",       // Green
		cells.TypeZ: "\033[31m[]\033[0m",       // Red
		cells.TypeJ: "\033[34m[]\033[0m",       // Blue
		cells.TypeL: "\033[38;5;208m[]\033[0m", // Orange
	}

	if color, exists := colors[cell]; exists {
		return color
	}
	return "[]"
}
