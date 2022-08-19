package seeta

// #include <Seetaface6CGO.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type Property int32

const (
	PROPERTY_MIN_FACE_SIZE    Property = 0
	PROPERTY_THRESHOLD        Property = 1
	PROPERTY_MAX_IMAGE_WIDTH  Property = 2
	PROPERTY_MAX_IMAGE_HEIGHT Property = 3
	PROPERTY_NUMBER_THREADS   Property = 4

	PROPERTY_ARM_CPU_MODE Property = 0x101
)

type FaceDetector struct {
	ptr C.int64_t
}

func (fd FaceDetector) getPtr() C.int64_t {
	if fd.ptr == 0 {
		panic("FaceDetector对象已关闭")
	}
	return fd.ptr
}

/**
最小人脸20耗时31~38ms，
最小人脸40耗时8~9ms，
最小人脸60耗时4~6ms， 推荐设置成60
最小人脸80耗时3~4ms   开始出现检测不到人脸的现象
*/
func (fd FaceDetector) Detect(img *CSeetaImageData) []SeetaFaceInfo {
	infoArray := C.Detect(fd.getPtr(), C.SeetaImageData(*img))
	size := int(infoArray.size)
	var result = make([]SeetaFaceInfo, size)
	var data []C.SeetaFaceInfo //C中的结构体指针没法直接读取，先封装成切片，然后遍历
	header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	header.Data = uintptr(unsafe.Pointer(infoArray.data)) // infoArray.data 这个对应C中的SeetaFaceInfo* 据官方文档所说，这个是"借用"对象，不需要外部释放内存
	header.Len = size
	header.Cap = size
	for i, info := range data {
		result[i] = SeetaFaceInfo{
			Pos: SeetaRect{
				X:      int(info.pos.x),
				Y:      int(info.pos.y),
				Width:  int(info.pos.width),
				Height: int(info.pos.height),
			},
			Score: float32(info.score),
		}
	}
	return result
}
func (fd FaceDetector) Get(pro Property) float64 {
	return float64(C.Get(fd.getPtr(), C.int(pro)))
}
func (fd FaceDetector) Set(pro Property, value float64) {
	C.Set(fd.getPtr(), C.int(pro), C.double(value))
}
func (fd *FaceDetector) Close() {
	if fd.ptr != 0 {
		C.DeleteFaceDetector(fd.ptr)
		fd.ptr = 0
	}
}

func NewFaceDetector(modelSetting SeetaModelSetting) *FaceDetector {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewFaceDetector(setting)
	return &FaceDetector{
		ptr: ptr,
	}
}
