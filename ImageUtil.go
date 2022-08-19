package seeta

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func ReadImgBgr(img image.Image) (width, height int, pix []uint8) {
	rect := img.Bounds()
	x1 := rect.Min.X
	x2 := rect.Max.X
	y1 := rect.Min.Y
	y2 := rect.Max.Y
	height = y2 - y1
	width = x2 - x1
	pix = make([]uint8, height*width*3)
	curr := 0
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ { // Color 转换为 BGR
			nrgba := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			pix[curr] = nrgba.B
			curr++
			pix[curr] = nrgba.G
			curr++
			pix[curr] = nrgba.R
			curr++
		}
	}
	return
}
