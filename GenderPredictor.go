package seeta

// #include <Seetaface6CGO.h>
import "C"
import "unsafe"

// 性别预测
type GenderPredictor struct {
	ptr C.int64_t
}

func (gp GenderPredictor) getPtr() C.int64_t {
	if gp.ptr == 0 {
		panic("GenderPredictor对象已关闭")
	}
	return gp.ptr
}

//void DeleteGenderPredictor(int64_t ptr);
func (gp *GenderPredictor) Close() {
	if gp.ptr != 0 {
		C.DeleteGenderPredictor(gp.ptr)
		gp.ptr = 0
	}
}

//int64_t NewGenderPredictor(SeetaModelSetting model);
func NewGenderPredictor(model SeetaModelSetting) *GenderPredictor {
	setting, free := model.CSeetaModelSetting()
	defer free()
	return &GenderPredictor{
		ptr: C.NewGenderPredictor(setting),
	}
}

//// 1 男 2 女 0 未知
//int PredictGenderWithCrop(int64_t ptr, const SeetaImageData& image, const SeetaPointF* points);
func (gp GenderPredictor) PredictGenderWithCrop(img *CSeetaImageData, pfs SeetaPointFs) Gender {
	return Gender(C.PredictGenderWithCrop(gp.getPtr(), C.SeetaImageData(*img), (*C.SeetaPointF)(unsafe.Pointer(&pfs.data[0]))))
}

type Gender int

const (
	Unknown = iota
	Male
	Female
)
