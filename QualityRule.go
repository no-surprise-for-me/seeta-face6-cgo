package seeta

// #include <Seetaface6CGO.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

/*
照片质量评估从
亮度、完整度（人脸是否完全进入摄像头）、是否遮挡、分辨率
姿态评估、姿态评估（深度）、清晰度评估、清晰度评估（深度）
等方面进行评估
在应用中，往往需要建立多个QualityRule，根据实际规则，往往需要多个QualityRule全部返回HIGH才认为图像合格。
*/

type QualityLevel int

const (
	LOW QualityLevel = iota
	MEDIUM
	HIGH
)

type QualityResult struct {
	Level QualityLevel // 0,1,2对应 LOW ,MEDIUM ,HIGH
	Score float32

	Desc string
}
type QualityRule interface {
	// QualityResult check(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points, int32_t N);
	Check(img *CSeetaImageData, rect SeetaRect, points SeetaPointFs, N int) QualityResult
	Close()
}
type qualityRuleCommonImpl struct {
	ptr  C.int64_t
	name string
	Desc string
	id   int
}

func (qr qualityRuleCommonImpl) getPtr() C.int64_t {
	if qr.ptr == 0 {
		panic(qr.name + "对象已关闭")
	}
	return qr.ptr
}
func (qr *qualityRuleCommonImpl) Close() {
	if qr.ptr != 0 {
		C.DeleteQualityRule(qr.ptr)
		qr.ptr = 0
	}
}

// QualityResult check(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points, int32_t N);
func (qr qualityRuleCommonImpl) Check(img *CSeetaImageData, rect SeetaRect, pfs SeetaPointFs, N int) QualityResult {
	var result C.QualityResult
	if pfs.data == nil {
		result = C.check(qr.getPtr(), C.SeetaImageData(*img), rect.CSeetaRect(), nil, C.int32_t(N))
	} else {
		result = C.check(qr.getPtr(), C.SeetaImageData(*img), rect.CSeetaRect(), (*C.SeetaPointF)(unsafe.Pointer(&pfs.data[0])), C.int32_t(N))
	}
	return QualityResult{
		Level: QualityLevel(result.level),
		Score: float32(result.score),
		Desc:  qr.Desc,
	}
}

// 亮度评估
type QualityOfBrightness struct {
	qualityRuleCommonImpl
}

//int64_t NewQualityOfBrightness();
func NewQualityOfBrightness() *QualityOfBrightness {
	var ptr C.int64_t = C.NewQualityOfBrightness()
	return &QualityOfBrightness{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfBrightness",
			Desc: "亮度评估",
			id:   0,
		},
	}
}

//int64_t NewQualityOfBrightnessWithParam(float v0, float v1, float v2, float v3);
func NewQualityOfBrightnessWithParam(v0, v1, v2, v3 float32) *QualityOfBrightness {
	var ptr C.int64_t = C.NewQualityOfBrightnessWithParam(C.float(v0), C.float(v1), C.float(v2), C.float(v3))
	return &QualityOfBrightness{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfBrightness",
			Desc: "亮度评估",
			id:   0,
		},
	}
}

// 完整度评估
type QualityOfIntegrity struct {
	qualityRuleCommonImpl
}

//int64_t NewQualityOfIntegrity();
func NewQualityOfIntegrity() *QualityOfIntegrity {
	var ptr C.int64_t = C.NewQualityOfIntegrity()
	return &QualityOfIntegrity{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfIntegrity",
			Desc: "完整度评估",
			id:   1,
		},
	}
}

//int64_t NewQualityOfIntegrityWithParam(float low, float height);
func NewQualityOfIntegrityWithParam(low, height float32) *QualityOfIntegrity {
	var ptr C.int64_t = C.NewQualityOfIntegrityWithParam(C.float(low), C.float(height))
	return &QualityOfIntegrity{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfIntegrity",
			Desc: "完整度评估",
			id:   1,
		},
	}
}

// 清晰度评估
type QualityOfClarity struct {
	qualityRuleCommonImpl
}

//int64_t NewQualityOfClarity();
func NewQualityOfClarity() *QualityOfClarity {
	var ptr C.int64_t = C.NewQualityOfClarity()
	return &QualityOfClarity{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfClarity",
			Desc: "清晰度评估",
			id:   2,
		},
	}
}

//int64_t NewQualityOfClarityWithParam(float low, float height);
func NewQualityOfClarityWithParam(low, height float32) *QualityOfClarity {
	var ptr C.int64_t = C.NewQualityOfClarityWithParam(C.float(low), C.float(height))
	return &QualityOfClarity{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfClarity",
			Desc: "清晰度评估",
			id:   2,
		},
	}
}

// 清晰度评估深度，代码 清晰度评估深度评估代码后两个参数是没有使用的，可以传NULL和68，固定值
type QualityOfClarityEx struct {
	qualityRuleCommonImpl
}

func (qcex QualityOfClarityEx) CheckSimple(img *CSeetaImageData, rect SeetaRect) QualityResult {
	return qcex.Check(img, rect, SeetaPointFs{}, 68)
}

// 检测耗时14~17ms
//int64_t NewQualityOfClarityEx(const char* quality_lbn_model_path, const char* landmark_pts68_model_path);
func NewQualityOfClarityEx(qualityLbnModelPath, landmarkPts68ModelPath string) *QualityOfClarityEx {
	s1 := C.CString(qualityLbnModelPath)
	s2 := C.CString(landmarkPts68ModelPath)
	var ptr C.int64_t = C.NewQualityOfClarityEx(s1, s2)
	C.free(unsafe.Pointer(s1))
	C.free(unsafe.Pointer(s2))

	return &QualityOfClarityEx{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfClarityEx",
			Desc: "清晰度评估（深度）",
			id:   3,
		},
	}
}

//int64_t NewQualityOfClarityExWithParam(const char* quality_lbn_model_path, const char* landmark_pts68_model_path, float blur_thresh);
func NewQualityOfClarityExWithParam(qualityLbnModelPath, landmarkPts68ModelPath string, blurThresh float32) *QualityOfClarityEx {
	s1 := C.CString(qualityLbnModelPath)
	s2 := C.CString(landmarkPts68ModelPath)
	var ptr C.int64_t = C.NewQualityOfClarityExWithParam(s1, s2, C.float(blurThresh))
	C.free(unsafe.Pointer(s1))
	C.free(unsafe.Pointer(s2))
	return &QualityOfClarityEx{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfClarityEx",
			Desc: "清晰度评估（深度）",
			id:   3,
		},
	}
}

// 姿态评估
type QualityOfPose struct {
	qualityRuleCommonImpl
}

//int64_t NewQualityOfPose();
func NewQualityOfPose() *QualityOfPose {
	var ptr C.int64_t = C.NewQualityOfPose()
	return &QualityOfPose{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfPose",
			Desc: "姿态评估",
			id:   4,
		},
	}
}

// 姿态评估（深度）
type QualityOfPoseEx struct {
	qualityRuleCommonImpl
}
type QualityOfPoseExProperty int

const (
	YAW_LOW_THRESHOLD QualityOfPoseExProperty = iota
	YAW_HIGH_THRESHOLD
	PITCH_LOW_THRESHOLD
	PITCH_HIGH_THRESHOLD
	ROLL_LOW_THRESHOLD
	ROLL_HIGH_THRESHOLD
)

//float getQualityOfPoseExProperty(int64_t ptr, int property);
func (qpex QualityOfPoseEx) GetProperty(property QualityOfPoseExProperty) float32 {
	return float32(C.getQualityOfPoseExProperty(qpex.getPtr(), C.int(property)))
}

//void setQualityOfPoseExProperty(int64_t ptr, int property, float value);
func (qpex QualityOfPoseEx) SetProperty(property QualityOfPoseExProperty, value float32) {
	C.setQualityOfPoseExProperty(qpex.getPtr(), C.int(property), C.float(value))
}

//int64_t NewQualityOfPoseEx(const SeetaModelSetting setting);
func NewQualityOfPoseEx(modelSetting SeetaModelSetting) *QualityOfPoseEx {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewQualityOfPoseEx(setting)
	return &QualityOfPoseEx{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfPoseEx",
			Desc: "姿态评估（深度）",
			id:   5,
		},
	}
}

// 分辨率评估
type QualityOfResolution struct {
	qualityRuleCommonImpl
}

//int64_t NewQualityOfResolution();
func NewQualityOfResolution() *QualityOfResolution {
	var ptr C.int64_t = C.NewQualityOfResolution()
	return &QualityOfResolution{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfResolution",
			Desc: "分辨率评估",
			id:   6,
		},
	}
}

//int64_t NewQualityOfResolutionWithParam(float low, float height);
func NewQualityOfResolutionWithParam(low, height float32) *QualityOfResolution {
	var ptr C.int64_t = C.NewQualityOfResolutionWithParam(C.float(low), C.float(height))
	return &QualityOfResolution{
		qualityRuleCommonImpl{
			ptr:  ptr,
			name: "QualityOfResolution",
			Desc: "分辨率评估",
			id:   6,
		},
	}
}
