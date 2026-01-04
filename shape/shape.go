// Package shape
package shape

import (
	"slices"

	cells "rishavmngo/tetris-tui-v2/shared"
)

const (
	TypeI ShapeType = "I"
	TypeJ ShapeType = "J"
	TypeL ShapeType = "L"
	TypeO ShapeType = "O"
	TypeS ShapeType = "S"
	TypeT ShapeType = "T"
	TypeZ ShapeType = "Z"
)

type (
	Grid      [4][4]bool
	ShapeType string
)

type ShapeDef struct {
	grid      Grid
	size      int
	shapeType cells.Cell
}

type Shape struct {
	Grid      Grid
	Size      int
	ShapeType cells.Cell
}

func NewShape(shapeType ShapeType) *Shape {
	def := registry[shapeType]

	return &Shape{
		Grid:      def.grid,
		Size:      def.size,
		ShapeType: def.shapeType,
	}
}

func (s *Shape) Clone() *Shape {
	return &Shape{
		Grid:      s.Grid,
		Size:      s.Size,
		ShapeType: s.ShapeType,
	}
}

func (s *Shape) transpose() {
	for i := range s.Size {
		for j := i + 1; j < s.Size; j++ {
			temp := s.Grid[i][j]
			s.Grid[i][j] = s.Grid[j][i]
			s.Grid[j][i] = temp
		}
	}
}

func (s *Shape) reverse() {
	slices.Reverse(s.Grid[:s.Size])
}

func (s *Shape) RotateClockWise() {
	if s.Size <= 2 {
		return
	}
	s.reverse()
	s.transpose()
}

var registry = map[ShapeType]ShapeDef{
	TypeI: {
		size:      4,
		shapeType: cells.TypeI,
		grid: Grid{
			{true, true, true, true},
			{false, false, false, false},
			{false, false, false, false},
			{false, false, false, false},
		},
	},
	TypeJ: {
		size:      3,
		shapeType: cells.TypeJ,
		grid: Grid{
			{false, true, false, false},
			{false, true, false, false},
			{true, true, false, false},
			{false, false, false, false},
		},
	},
	TypeL: {
		size:      3,
		shapeType: cells.TypeL,
		grid: Grid{
			{true, false, false, false},
			{true, false, false, false},
			{true, true, false, false},
			{false, false, false, false},
		},
	},
	TypeO: {
		size:      2,
		shapeType: cells.TypeO,
		grid: Grid{
			{true, true, false, false},
			{true, true, false, false},
			{false, false, false, false},
			{false, false, false, false},
		},
	},
	TypeS: {
		size:      3,
		shapeType: cells.TypeS,
		grid: Grid{
			{false, true, true, false},
			{true, true, false, false},
			{false, false, false, false},
			{false, false, false, false},
		},
	},
	TypeT: {
		size:      3,
		shapeType: cells.TypeT,
		grid: Grid{
			{false, true, false, false},
			{true, true, true, false},
			{false, false, false, false},
			{false, false, false, false},
		},
	},
	TypeZ: {
		size:      3,
		shapeType: cells.TypeZ,
		grid: Grid{
			{true, true, false, false},
			{false, true, true, false},
			{false, false, false, false},
			{false, false, false, false},
		},
	},
}
