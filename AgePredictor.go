package seeta

// #include <Seetaface6CGO.h>
import "C"
import "unsafe"

type AgePredictor struct {
	ptr C.int64_t
}

func (ap AgePredictor) getPtr() C.int64_t {
	if ap.ptr == 0 {
		panic("AgePredictor对象已关闭")
	}
	return ap.ptr
}

//void DeleteAgePredictor(int64_t ptr);
func (ap *AgePredictor) Close() {
	if ap.ptr != 0 {
		C.DeleteAgePredictor(ap.ptr)
		ap.ptr = 0
	}
}

//int64_t NewAgePredictor(SeetaModelSetting model);
func NewAgePredictor(model SeetaModelSetting) *AgePredictor {
	setting, free := model.CSeetaModelSetting()
	defer free()
	return &AgePredictor{
		ptr: C.NewAgePredictor(setting),
	}
}

// 耗时34~41ms 建议隔1秒取一次值，取十次，然后取中间三个结果的平均值
//int PredictAgeWithCrop(int64_t ptr, const SeetaImageData& image, const SeetaPointF* points);
func (ap *AgePredictor) PredictAgeWithCrop(img *CSeetaImageData, pfs SeetaPointFs) int {
	return int(C.PredictAgeWithCrop(ap.getPtr(), C.SeetaImageData(*img), (*C.SeetaPointF)(unsafe.Pointer(&pfs.data[0]))))
}
