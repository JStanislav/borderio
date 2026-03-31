package utils

type GridPosition struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}

type WallPosition struct {
	CellA GridPosition `json:"cellA"`
	CellB GridPosition `json:"cellB"`
}

func (w WallPosition) Orientation() string {
	if w.CellA.Row == w.CellB.Row {
		return VerticalLine
	} else {
		return HorizontalLine
	}
}
