package seeta

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
	"unsafe"
)

func ExtractFeature(path1 string, detector *FaceDetector, landmarker *FaceLandmarker, faceRecognizer *FaceRecognizer, spoofing *FaceAntiSpoofing) ([]float32, Status) {
	seetaFaceImage1, err := NewSeetaFaceImageFromFile(path1)
	if err != nil {
		panic(err)
	}
	faceInfoArray := detector.Detect(seetaFaceImage1)
	if len(faceInfoArray) == 1 {
		mark := landmarker.Mark(seetaFaceImage1, faceInfoArray[0].Pos)
		feature1 := faceRecognizer.Extract(seetaFaceImage1, mark)
		predict := spoofing.Predict(seetaFaceImage1, faceInfoArray[0].Pos, mark)
		return feature1, predict
	} else {
		return nil, FUZZY
	}
}
func GetDetector(model string) *FaceDetector {
	LocationModel(&model)
	seetaModel := SeetaModelSetting{
		Device: SEETA_DEVICE_AUTO,
		Id:     0,
		Model: []string{
			model,
		},
	}
	return NewFaceDetector(seetaModel)
}
func GetTracker(model string, videoWidth, videoHeight int) *FaceTracker {
	LocationModel(&model)
	seetaModel := SeetaModelSetting{
		Device: SEETA_DEVICE_AUTO,
		Id:     0,
		Model: []string{
			model,
		},
	}
	return NewFaceTracker(seetaModel, videoWidth, videoHeight)
}
func GetLandMarker(model string) *FaceLandmarker {
	LocationModel(&model)
	seetaModel := SeetaModelSetting{
		Device: SEETA_DEVICE_AUTO,
		Id:     0,
		Model: []string{
			model,
		},
	}
	return NewFaceLandmarker(seetaModel)
}
func GetFaceRecognizer(model string) *FaceRecognizer {
	LocationModel(&model)
	seetaModel := SeetaModelSetting{
		Device: SEETA_DEVICE_AUTO,
		Id:     0,
		Model: []string{
			model,
		},
	}
	return NewFaceRecognizer(seetaModel)
}
func GetFaceAntiSpoofing(model []string) *FaceAntiSpoofing {
	for i := range model {
		LocationModel(&model[i])
	}
	seetaModel := SeetaModelSetting{
		Device: SEETA_DEVICE_AUTO,
		Id:     0,
		Model:  model,
	}
	return NewFaceAntiSpoofing(seetaModel)
}
func LocationModel(model0 *string) {
	if model0 == nil || "" == *model0 {
		panic("模型文件未指定")
	}
	model := *model0
	defer func() {
		*model0 = model
	}()
	if !fileExist(model) {
		_, file := filepath.Split(model)
		model = filepath.Join(filepath.Dir(os.Args[0]), file)
	} else {
		return
	}
	if !fileExist(model) {
		_, file := filepath.Split(model)
		model = filepath.Join(filepath.Dir(os.Args[0]), "Model", file)
	} else {
		return
	}
	if !fileExist(model) {
		_, file := filepath.Split(model)
		model = filepath.Join(filepath.Dir(os.Args[0]), "models", file)
	} else {
		return
	}
	if !fileExist(model) {
		_, file := filepath.Split(model)
		fmt.Println("找不到模型文件：" + file)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	} else {
		return
	}
}
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		panic(err)
	}
}
func Compare(f1, f2 []float32) float32 {
	length := len(f1)
	if length != len(f2) {
		if length < len(f2) {
			return -1 // 第一个参数较短
		} else {
			return -2 // 第二个参数较短
		}
	}
	var sum float32
	for i := 0; i < length; i++ {
		sum += f1[i] * f2[i]
	}
	return sum
}
func CompareBytes(f1, f2 []byte) float32 {
	return Compare(ByteToFloat64(f1), ByteToFloat64(f2))
}
func CompareBytesAndFloats(f1 []byte, f2 []float32) float32 {
	return Compare(ByteToFloat64(f1), f2)
}

//用这个转换可能要注意CPU大小端的问题，相同的内存数据，在大小端CPU中的含义是不一样的
//Float64ToByte Float64转byte
func Float64ToByte(floats []float32) (bytes []byte) {
	floatHeader := (*reflect.SliceHeader)(unsafe.Pointer(&floats))
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	byteHeader.Data = floatHeader.Data
	byteHeader.Len = floatHeader.Len * 4
	byteHeader.Cap = floatHeader.Cap * 4
	return
}

//ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) (floats []float32) {
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	floatHeader := (*reflect.SliceHeader)(unsafe.Pointer(&floats))
	floatHeader.Data = byteHeader.Data
	floatHeader.Len = byteHeader.Len / 4
	floatHeader.Cap = byteHeader.Cap / 4
	return
}
