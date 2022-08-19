package seeta

// #include <Seetaface6CGO.h>
import "C"

type SeetaFaceInfo struct {
	Pos   SeetaRect
	Score float32
}

type SeetaRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (sr SeetaRect) CSeetaRect() C.SeetaRect {
	var cRect C.SeetaRect
	cRect.x = C.int(sr.X)
	cRect.y = C.int(sr.Y)
	cRect.width = C.int(sr.Width)
	cRect.height = C.int(sr.Height)
	return cRect
}
