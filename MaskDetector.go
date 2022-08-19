package seeta

// #include <Seetaface6CGO.h>
import "C"

type MaskDetector struct {
	ptr C.int64_t
}

func (md MaskDetector) getPtr() C.int64_t {
	if md.ptr == 0 {
		panic("MaskDetector对象已关闭")
	}
	return md.ptr
}

//void DeleteMaskDetector(int64_t ptr);
func (md *MaskDetector) Close() {
	if md.ptr != 0 {
		C.DeleteMaskDetector(md.ptr)
		md.ptr = 0
	}
}

//int64_t NewMaskDetector(SeetaModelSetting model);
func NewMaskDetector(model SeetaModelSetting) *MaskDetector {
	setting, free := model.CSeetaModelSetting()
	defer free()
	return &MaskDetector{
		ptr: C.NewMaskDetector(setting),
	}
}

//MaskFace DetectMask(int64_t ptr, const SeetaImageData image, const SeetaRect face);
func (md MaskDetector) DetectMask(img *CSeetaImageData, rect SeetaRect) MaskFace {
	cmaskFace := C.DetectMask(md.getPtr(), (C.SeetaImageData)(*img), rect.CSeetaRect())
	return MaskFace{
		Score: float32(cmaskFace.score),
		Mask:  int(cmaskFace.mask) == 1, // 0 没戴口罩  1 戴口罩了
	}
}

type MaskFace struct {
	Score float32 //  戴了口罩的置信度，score超过0.5，则认为是检测带上了口罩
	Mask  bool    //  是否戴口罩
}
