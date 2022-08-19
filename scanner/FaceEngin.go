package scanner

import (
	seeta "gitee.com/no_surprise_for_me/seeta-face-cgo"
	"os"
	"path/filepath"
	"sync"
)

type FaceEngine struct {
	detector   *seeta.FaceDetector
	landmarker *seeta.FaceLandmarker
	recognizer *seeta.FaceRecognizer

	//brightness *seeta.QualityOfBrightness
	//clarityEx  *seeta.QualityOfClarityEx
	//poseEx     *seeta.QualityOfPoseEx
	//
	//qrLock            sync.Locker
	recognizerLock    sync.Locker
	facePreHandleLock sync.Locker
}

func NewFaceEngine() *FaceEngine {
	fe := new(FaceEngine)

	fe.detector = seeta.GetDetector(FindModel("face_detector.csta"))
	fe.landmarker = seeta.GetLandMarker(FindModel("face_landmarker_pts5.csta"))
	fe.recognizer = seeta.GetFaceRecognizer(FindModel("face_recognizer.csta"))
	fe.detector.Set(seeta.PROPERTY_MIN_FACE_SIZE, 40)

	//fe.brightness = seeta.NewQualityOfBrightness()
	//fe.clarityEx = seeta.NewQualityOfClarityEx(FindModel("quality_lbn.csta"), FindModel("face_landmarker_pts68.csta"))
	//fe.poseEx = seeta.NewQualityOfPoseEx(seeta.SeetaModelSetting{Model: []string{FindModel("pose_estimation.csta")}})

	//fe.qrLock = &sync.Mutex{}
	fe.recognizerLock = &sync.Mutex{}
	fe.facePreHandleLock = &sync.Mutex{}

	return fe
}
func (fe *FaceEngine) Close() {
	fe.facePreHandleLock.Lock()
	fe.detector.Close()
	fe.landmarker.Close()

	//fe.qrLock.Lock()
	//fe.poseEx.Close()
	//fe.clarityEx.Close()
	//fe.brightness.Close()

	fe.recognizerLock.Lock()
	fe.recognizer.Close()
}

// 多张人脸时取最大人脸
func (fe FaceEngine) PreHandle(cimg *seeta.CSeetaImageData) (*seeta.SeetaFaceInfo, *seeta.SeetaPointFs) {
	fe.facePreHandleLock.Lock()
	defer fe.facePreHandleLock.Unlock()
	detect := fe.detector.Detect(cimg)
	var faceInfo *seeta.SeetaFaceInfo
	if len(detect) > 1 {
		var max, index int
		for i, face := range detect {
			if face.Pos.Width > max {
				index = i
				max = face.Pos.Width
			}
		}
		faceInfo = &detect[index]
	} else if len(detect) == 0 {
		return nil, nil
	} else {
		faceInfo = &detect[0]
	}
	pointFs := fe.landmarker.Mark(cimg, faceInfo.Pos)
	return faceInfo, &pointFs
}
func (fe FaceEngine) Extract(cimg *seeta.CSeetaImageData, points *seeta.SeetaPointFs) []float32 {
	fe.recognizerLock.Lock()
	defer fe.recognizerLock.Unlock()
	return fe.recognizer.Extract(cimg, *points)
}
func (fe FaceEngine) Check(cimg *seeta.CSeetaImageData, pos *seeta.SeetaRect, pfs *seeta.SeetaPointFs) string {
	//fe.qrLock.Lock()
	//defer fe.qrLock.Unlock()
	//check := fe.brightness.Check(cimg, *pos, *pfs, 5)
	//if check.Level != seeta.HIGH {
	//	return "亮度评估不合格"
	//}
	//check = fe.poseEx.Check(cimg, *pos, *pfs, 5)
	//if check.Level != seeta.HIGH {
	//	return "姿态评估不合格"
	//}
	//check = fe.clarityEx.Check(cimg, *pos, *pfs, 5)
	//if check.Level != seeta.HIGH {
	//	return "清晰度评估不合格"
	//}
	return ""
}

func FindModel(path string) string {
	var pathList = make([]string, 9)
	pathList[0] = path
	root := filepath.Dir(os.Args[0])
	if !filepath.IsAbs(path) {
		pathList[1] = filepath.Join(root, path)
		pathList[2] = filepath.Join(root, "..", path)
	}
	_, file := filepath.Split(path)
	pathList[3] = filepath.Join(root, "models", file)
	pathList[4] = filepath.Join(root, "model", file)
	pathList[5] = filepath.Join(root, "..", "models", file)
	pathList[6] = filepath.Join(root, "..", "model", file)
	pathList[7] = filepath.Join(root, file)
	pathList[8] = filepath.Join(root, "..", file)
	for _, p := range pathList {
		if p == "" {
			continue
		}
		_, err := os.Stat(p)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		} else {
			return p
		}
	}
	panic("模型文件不存在：" + path)
}
