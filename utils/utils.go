package utils

type GridPosition struct {
	Column int
	Row    int
}

type WallPosition struct {
	CellA GridPosition
	CellB GridPosition
}
