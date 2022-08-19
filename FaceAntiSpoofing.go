package seeta

// #include <Seetaface6CGO.h>
import "C"
import (
	"unsafe"
)

type Status int

const (
	REAL      Status = 0 ///< 真实人脸
	SPOOF     Status = 1 ///< 攻击人脸（假人脸）
	FUZZY     Status = 2 ///< 无法判断（人脸成像质量不好）
	DETECTING Status = 3 ///< 正在检测
)

/**
 * \brief 加载模型文件
 * \param setting 模型文件, 0-局部活体检测文件fas_first.csta（必选），1-全局活体检测文件fas_second.csta（可选）
 *  按官方文档的例子，最好不要改变顺序
 */
func NewFaceAntiSpoofing(modelSetting SeetaModelSetting) *FaceAntiSpoofing {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewFaceAntiSpoofing(setting)
	return &FaceAntiSpoofing{
		ptr: ptr,
	}
}

type FaceAntiSpoofing struct {
	ptr C.int64_t
}

func (fas FaceAntiSpoofing) getPtr() C.int64_t {
	if fas.ptr == 0 {
		panic("FaceAntiSpoofing对象已关闭")
	}
	return fas.ptr
}

// void DeleteFaceAntiSpoofing(int64_t ptr)
func (fas *FaceAntiSpoofing) Close() {
	if fas.ptr != 0 {
		C.DeleteFaceAntiSpoofing(fas.ptr)
		fas.ptr = 0
	}
}

/**
 * \brief 检测活体
 * \param [in] image 输入图像，需要 RGB 彩色通道
 * \param [in] face 要识别的人脸位置
 * \param [in] points 要识别的人脸特征点
 * \return 人脸状态 @see Status
 * \note 此函数不支持多线程调用，在多线程环境下需要建立对应的 FaceAntiSpoofing 的对象分别调用检测函数
 * \note 当前版本可能返回 REAL, SPOOF, FUZZY
 * \see SeetaImageData, SeetaRect, PointF, Status
// 耗时97~110ms  图片不清晰的情况下测不出是否活体，或测出为非活体
*/
//Status Predict(int64_t ptr,const SeetaImageData& image, const SeetaRect& face, const SeetaPointF* points);
func (fas FaceAntiSpoofing) Predict(img *CSeetaImageData, rect SeetaRect, points SeetaPointFs) Status {
	return Status(C.Predict(fas.getPtr(), C.SeetaImageData(*img), rect.CSeetaRect(), (*C.SeetaPointF)(unsafe.Pointer(&points.data[0]))))
}

/**
 * \brief 检测活体（Video模式）
 * \param [in] image 输入图像，需要 RGB 彩色通道
 * \param [in] face 要识别的人脸位置
 * \param [in] points 要识别的人脸特征点
 * \return 人脸状态 @see Status
 * \note 此函数不支持多线程调用，在多线程环境下需要建立对应的 FaceAntiSpoofing 的对象分别调用检测函数
 * \note 需要输入连续帧序列，当需要输入下一段视频是，需要调用 ResetVideo 重置检测状态
 * \note 当前版本可能返回 REAL, SPOOF, DETECTION
 * \see SeetaImageData, SeetaRect, PointF, Status
// 耗时93~115ms   效果比单帧检测要好，但是需要的检测样本图片数量较多
*/
//Status PredictVideo(int64_t ptr, const SeetaImageData& image, const SeetaRect& face, const SeetaPointF* points);
func (fas FaceAntiSpoofing) PredictVideo(img *CSeetaImageData, rect SeetaRect, points SeetaPointFs) Status {
	return Status(C.PredictVideo(fas.getPtr(), C.SeetaImageData(*img), rect.CSeetaRect(), (*C.SeetaPointF)(unsafe.Pointer(&points.data[0]))))
}

/**
 * \brief 重置 Video，开始下一次 PredictVideo 识别
 */
//void ResetVideo(int64_t ptr);
func (fas FaceAntiSpoofing) ResetVideo() {
	C.ResetVideo(fas.getPtr())
}

/**
 * \brief 获取活体检测内部分数
 * \param [out] clarity 输出人脸质量分数
 * \param [out] reality 真实度
 * \note 获取的是上一次调用 Predict 或 PredictVideo 接口后内部的阈值
 */
//void GetPreFrameScore(int64_t ptr, float* clarity, float* reality);
func (fas FaceAntiSpoofing) GetPreFrameScore() (clarity, reality float32) {
	C.GetPreFrameScore(fas.getPtr(), (*C.float)(&clarity), (*C.float)(&reality))
	return
}

/**
 * 设置 Video 模式中，识别视频帧数，当输入帧数为该值以后才会有返回值
 * \param [in] number 视频帧数
 */
//void SetVideoFrameCount(int64_t ptr, int32_t number);
func (fas FaceAntiSpoofing) SetVideoFrameCount(number int) {
	C.SetVideoFrameCount(fas.getPtr(), C.int32_t(number))
}

//int32_t GetVideoFrameCount(int64_t ptr);
func (fas FaceAntiSpoofing) GetVideoFrameCount() int {
	return int(C.GetVideoFrameCount(fas.getPtr()))
}

/**
 * 设置阈值
 * \param [in] clarity 清晰度阈值
 * \param [in] reality 活体阈值
 * \note clarity 越高要求输入的图像质量越高，reality 越高对识别要求越严格
 * \note 默认阈值为 0.3, 0.8
 */
//void SetThreshold(int64_t ptr, float clarity, float reality);
func (fas FaceAntiSpoofing) SetThreshold(clarity, reality float32) {
	C.SetThreshold(fas.getPtr(), C.float(clarity), C.float(reality))
}

//void GetThreshold(int64_t ptr, float* clarity, float* reality);
func (fas FaceAntiSpoofing) GetThreshold() (clarity, reality float32) {
	C.GetThreshold(fas.getPtr(), (*C.float)(&clarity), (*C.float)(&reality))
	return
}

/**
 * 设置全局阈值
 * \param [in] box_thresh 全局检测阈值
 * \note 默认阈值为 0.8
默认为0.8，这个是攻击介质存在的分数阈值，该阈值越高，表示对攻击介质的要求越严格，一般的疑似就不会认为是攻击介质。这个一般不进行调整
*/
//void SetBoxThresh(int64_t ptr, float box_thresh);
func (fas FaceAntiSpoofing) SetBoxThresh(boxThresh float32) {
	C.SetBoxThresh(fas.getPtr(), C.float(boxThresh))
}

//float GetBoxThresh(int64_t ptr);
func (fas FaceAntiSpoofing) GetBoxThresh() float32 {
	return float32(C.GetBoxThresh(fas.getPtr()))
}
