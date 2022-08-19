package seeta

// #include <Seetaface6CGO.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type FaceTracker struct {
	ptr C.int64_t
}

func (ft FaceTracker) getPtr() C.int64_t {
	if ft.ptr == 0 {
		panic("FaceTracker对象已关闭")
	}
	return ft.ptr
}

//int64_t NewFaceTracker(const char** models, int device, int id, int video_width, int video_height);
func NewFaceTracker(modelSetting SeetaModelSetting, videoWidth, videoHeight int) *FaceTracker {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewFaceTracker(setting, C.int(videoWidth), C.int(videoHeight))
	return &FaceTracker{
		ptr: ptr,
	}
}

//void SetInterval(int64_t ptr, int interval);
func (ft FaceTracker) SetInterval(interval int) {
	C.SetInterval(ft.getPtr(), C.int(interval))
}

/**
耗时14~17ms,但从耗时上讲，比不上人脸检测将最小人脸设置成40+时的性能，
但人脸追踪应该是为了解决重复计算是否活体和抽取特征值等操作的，这些计算一次就够了，不需要重复计算，人脸最终可以解决这个问题。
在视频处理中，人脸追踪可以省略很多不必要的计算。
*/
//SeetaTrackingFaceInfoArray Track(int64_t ptr, const SeetaImageData image);
func (ft FaceTracker) Track(image *CSeetaImageData) []SeetaTrackingFaceInfo {
	return cSeetaTrackingFaceInfoArrayToSeetaTrackingFaceInfoSlice(C.Track(ft.getPtr(), C.SeetaImageData(*image)))
}

//SeetaTrackingFaceInfoArray TrackWithFrameNo(int64_t ptr, const SeetaImageData image, int frame_no);
func (ft FaceTracker) TrackWithFrameNo(image *CSeetaImageData, frameNo int) []SeetaTrackingFaceInfo {
	return cSeetaTrackingFaceInfoArrayToSeetaTrackingFaceInfoSlice(C.TrackWithFrameNo(ft.getPtr(), C.SeetaImageData(*image), C.int(frameNo)))
}

//void SetMinFaceSize(int64_t ptr, int32_t size);
func (ft FaceTracker) SetMinFaceSize(size int) {
	C.SetMinFaceSize(ft.getPtr(), C.int32_t(size))
}

//int32_t GetMinFaceSize(int64_t ptr);
func (ft FaceTracker) GetMinFaceSize() int {
	return int(C.GetMinFaceSize(ft.getPtr()))
}

//void SetFaceTrackeThreshold(int64_t ptr, float thresh);
func (ft FaceTracker) SetThreshold(thresh float32) {
	C.SetFaceTrackeThreshold(ft.getPtr(), C.float(thresh))
}

//float GetFaceTrackeThreshold(int64_t ptr);
func (ft FaceTracker) GetThreshold() float32 {
	return float32(C.GetFaceTrackeThreshold(ft.getPtr()))
}

//void SetVideoStable(int64_t ptr, bool stable);
func (ft FaceTracker) SetVideoStable(stable bool) {
	if stable {
		C.SetVideoStable(ft.getPtr(), C.int(1))
	} else {
		C.SetVideoStable(ft.getPtr(), C.int(0))
	}
}

//bool GetVideoStable(int64_t ptr);
func (ft FaceTracker) GetVideoStable() bool {
	return int(C.GetVideoStable(ft.getPtr())) != 0
}

//void SetVideoSize(int64_t ptr, int vidwidth, int vidheight);
func (ft FaceTracker) SetVideoSize(videoWidth, videoHeight int) {
	C.SetVideoSize(ft.getPtr(), C.int(videoWidth), C.int(videoHeight))
}

//void Reset(int64_t ptr);
func (ft FaceTracker) Reset() {
	C.Reset(ft.getPtr())
}

//void DeleteFaceTracker(int64_t ptr);
func (ft *FaceTracker) Close() {
	if ft.ptr != 0 {
		C.DeleteFaceTracker(ft.ptr)
		ft.ptr = 0
	}
}

// 将C中的SeetaTrackingFaceInfoArray结构体转换成Go中的 []SeetaTrackingFaceInfo
func cSeetaTrackingFaceInfoArrayToSeetaTrackingFaceInfoSlice(infoArray C.SeetaTrackingFaceInfoArray) []SeetaTrackingFaceInfo {
	size := int(infoArray.size)
	var cTrackingFaceInfos []C.SeetaTrackingFaceInfo
	header := (*reflect.SliceHeader)(unsafe.Pointer(&cTrackingFaceInfos))
	header.Data = uintptr(unsafe.Pointer(infoArray.data)) // cSeetaTrackingFaceInfoArray.data 这个对应C中的SeetaFaceInfo* 据官方文档所说，这个是"借用"对象，不需要外部释放内存
	header.Cap = size
	header.Len = size
	trackingFaceInfos := make([]SeetaTrackingFaceInfo, size)
	for i, info := range cTrackingFaceInfos {
		trackingFaceInfos[i] = SeetaTrackingFaceInfo{
			Pos: SeetaRect{
				X:      int(info.pos.x),
				Y:      int(info.pos.y),
				Width:  int(info.pos.width),
				Height: int(info.pos.height),
			},
			Score:   float32(info.score),
			FrameNo: int(info.frame_no),
			PID:     int(info.PID),
			Step:    int(info.step),
		}
	}
	return trackingFaceInfos
}
