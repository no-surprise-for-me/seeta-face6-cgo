package seeta

// #include <Seetaface6CGO.h>
import "C"
import (
	"image"
	"io"
	"os"
	"unsafe"
)

type CSeetaImageData C.SeetaImageData

// 用于手动构建SeetaImageData
type SeetaImageData struct {
	Width    int
	Height   int
	Channels int     // 3通道
	Data     []uint8 // BGR格式的图片内存数据
}

// 将手动构建的SeetaImageData转换为可以传入C代码的CSeetaImageData
func (s SeetaImageData) CSeetaImageData() *CSeetaImageData {
	c := new(CSeetaImageData)
	c.width = C.int(s.Width)
	c.height = C.int(s.Height)
	c.channels = C.int(s.Channels)
	c.data = (*C.uint8_t)(unsafe.Pointer(&s.Data[0]))
	return c
}
func NewSeetaFaceImage(img image.Image) *CSeetaImageData {
	width, height, pix := ReadImgBgr(img)
	c := new(CSeetaImageData)
	c.width = C.int(width)
	c.height = C.int(height)
	c.channels = C.int(3)
	c.data = (*C.uint8_t)(unsafe.Pointer(&pix[0]))
	return c
}
func NewSeetaFaceImageFromReader(r io.Reader) (*CSeetaImageData, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return NewSeetaFaceImage(img), nil
}
func NewSeetaFaceImageFromFile(filename string) (*CSeetaImageData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewSeetaFaceImageFromReader(file)
}
