package main

import "C"
import (
	"flag"
	"fmt"
	seeta "gitee.com/no_surprise_for_me/seeta-face-cgo"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetSize(file string) string {
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("文件不存在：" + file + "," + err.Error())
			os.Exit(1)
		} else {
			panic(err)
		}
	}
	size := float64(stat.Size())
	unit := " B"
	if size > 1024 {
		size = size / 1024
		unit = " KB"
		if size > 1024 {
			size = size / 1024
			unit = " MB"
		}
	}
	return strconv.FormatFloat(size, 'f', 3, 64) + unit
}
func main() {
	path1 := flag.String("path1", "", "第一张图片路径")
	path2 := flag.String("path2", "", "第二张图片路径")
	minFaceSize := flag.Int("s", 20, "人脸检测最小人脸大小，默认20，推荐60")
	num := flag.Int("n", 1, "计算相似度次数")
	flag.Parse()
	if *num <= 0 {
		*num = 1
	}
	fmt.Println("第一张人脸图片：" + *path1 + "\t大小：" + GetSize(*path1))
	fmt.Println("第二张人脸图片：" + *path2 + "\t大小：" + GetSize(*path1))
	width, height, pix := ReadImgBgr(*path1)
	cimg1 := seeta.SeetaImageData{
		Width:    width,
		Height:   height,
		Channels: 3,
		Data:     pix,
	}.CSeetaImageData()
	width, height, pix = ReadImgBgr(*path2)
	cimg2 := seeta.SeetaImageData{
		Width:    width,
		Height:   height,
		Channels: 3,
		Data:     pix,
	}.CSeetaImageData()
	detector := seeta.GetDetector(FindModel("models/face_detector.csta"))
	defer detector.Close()
	detector.Set(seeta.PROPERTY_MIN_FACE_SIZE, float64(*minFaceSize))
	fmt.Println("人脸检测最小人脸大小:", detector.Get(seeta.PROPERTY_MIN_FACE_SIZE))
	landmarker := seeta.GetLandMarker(FindModel("models/face_landmarker_pts5.csta"))
	defer landmarker.Close()
	recognizer := seeta.GetFaceRecognizer(FindModel("models/face_recognizer.csta"))
	defer recognizer.Close()
	faceAntiSpoofing := seeta.GetFaceAntiSpoofing([]string{FindModel("models/fas_first.csta")})
	defer faceAntiSpoofing.Close()
	agePredictor := seeta.NewAgePredictor(seeta.SeetaModelSetting{Model: []string{FindModel("models/age_predictor.csta")}})
	defer agePredictor.Close()
	genderPredictor := seeta.NewGenderPredictor(seeta.SeetaModelSetting{Model: []string{FindModel("models/gender_predictor.csta")}})
	defer genderPredictor.Close()
	maskDetector := seeta.NewMaskDetector(seeta.SeetaModelSetting{Model: []string{FindModel("models/mask_detector.csta")}})
	defer maskDetector.Close()
	stateDetector := seeta.NewEyeStateDetector(seeta.SeetaModelSetting{Model: []string{FindModel("models/eye_state.csta")}})
	defer stateDetector.Close()
	brightness := seeta.NewQualityOfBrightness()
	defer brightness.Close()
	ofClarity := seeta.NewQualityOfClarity()
	defer ofClarity.Close()
	clarityEx := seeta.NewQualityOfClarityEx(FindModel("models/quality_lbn.csta"), FindModel("models/face_landmarker_pts68.csta"))
	defer clarityEx.Close()
	integrity := seeta.NewQualityOfIntegrity()
	defer integrity.Close()
	pose := seeta.NewQualityOfPose()
	defer pose.Close()
	qualityOfPoseEx := seeta.NewQualityOfPoseEx(seeta.SeetaModelSetting{Model: []string{FindModel("models/pose_estimation.csta")}})
	defer qualityOfPoseEx.Close()
	resolution := seeta.NewQualityOfResolution()
	defer resolution.Close()
	fmt.Println("下面以第一张人脸作为底库，第二张为要对比的图片进行测试：")
	start := time.Now().UnixNano()
	faceInfos1 := detector.Detect(cimg1)
	fmt.Println("人脸检测耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
	if len(faceInfos1) == 0 {
		fmt.Printf("第一张图片无人脸")
		return
	}
	info1 := faceInfos1[0]
	start = time.Now().UnixNano()
	pointFs1 := landmarker.Mark(cimg1, info1.Pos)
	fmt.Println("人脸特征点标记耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
	checkAll(cimg1, info1, pointFs1, brightness, ofClarity, clarityEx, integrity, pose, qualityOfPoseEx, resolution)
	start = time.Now().UnixNano()
	features1 := recognizer.Extract(cimg1, pointFs1)
	fmt.Println("人脸特征抽取耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
	fmt.Println("================")
	start = time.Now().UnixNano()
	total := start
	faceInfos2 := detector.Detect(cimg2)
	fmt.Println("人脸检测耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
	if len(faceInfos2) == 0 {
		fmt.Printf("第二张图片无人脸")
		return
	}
	info2 := faceInfos2[0]
	start = time.Now().UnixNano()
	pointFs2 := landmarker.Mark(cimg2, info2.Pos)
	fmt.Println("人脸特征点标记耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
	checkAll(cimg2, info2, pointFs2, brightness, ofClarity, clarityEx, integrity, pose, qualityOfPoseEx, resolution)
	start = time.Now().UnixNano()
	features2 := recognizer.Extract(cimg2, pointFs2)
	fmt.Println("人脸特征抽取耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
	var similarity float32
	for i := 0; i < *num; i++ {
		similarity = seeta.Compare(features1, features2)
	}
	fmt.Printf("人脸对比总耗时%v毫秒,相似度：%v\n", (time.Now().UnixNano()-total)/1000000, similarity*100)
	fmt.Println("====人脸属性检测====")
	predict := faceAntiSpoofing.Predict(cimg2, info2.Pos, pointFs2)
	var result string
	switch predict {
	case seeta.REAL:
		result = "REAL"
	case seeta.SPOOF:
		result = "SPOOF"
	case seeta.FUZZY:
		result = "FUZZY"
	case seeta.DETECTING:
		result = "DETECTING"
	default:
		panic("活体检测结果范围外")
	}
	fmt.Println("活体检测结果：" + result)
	age := agePredictor.PredictAgeWithCrop(cimg2, pointFs2)
	fmt.Printf("年龄预测结果：%v\n", age)
	gender := genderPredictor.PredictGenderWithCrop(cimg2, pointFs2)
	switch gender {
	case seeta.Male:
		result = "Male"
	case seeta.Female:
		result = "Female"
	case seeta.Unknown:
		result = "Unknown"
	default:
		panic("性别检测结果范围外")
	}
	fmt.Printf("性别预测结果：%v\n", result)
	mask := maskDetector.DetectMask(cimg2, info2.Pos)
	fmt.Printf("口罩检测结果：%v，可信度：%v\n", mask.Mask, mask.Score)
	state := stateDetector.DetectEyeState(cimg2, pointFs2)
	fmt.Printf("眼睛状态检测，左眼：%v，右眼：%v\n", EyeStateString(state.Left), EyeStateString(state.Right))
}
func checkAll(cimg *seeta.CSeetaImageData,
	info seeta.SeetaFaceInfo,
	points seeta.SeetaPointFs,
	brightness *seeta.QualityOfBrightness,
	clarity *seeta.QualityOfClarity,
	clarityEx *seeta.QualityOfClarityEx,
	integrity *seeta.QualityOfIntegrity,
	pose *seeta.QualityOfPose,
	poseEx *seeta.QualityOfPoseEx,
	resolution *seeta.QualityOfResolution) {
	fmt.Println("========")
	CheckPic(brightness, cimg, info, points, "亮度检测：")
	CheckPic(clarity, cimg, info, points, "清晰度检测：")
	CheckPic(clarityEx, cimg, info, points, "清晰度检测(深度)：")
	CheckPic(integrity, cimg, info, points, "完整度检测：")
	CheckPic(pose, cimg, info, points, "姿态检测：")
	CheckPic(poseEx, cimg, info, points, "姿态检测(深度)：")
	CheckPic(resolution, cimg, info, points, "分辨率检测：")
	fmt.Println("========")
}
func CheckPic(qr seeta.QualityRule, cimg *seeta.CSeetaImageData, info seeta.SeetaFaceInfo, pfs seeta.SeetaPointFs, text string) {
	start := time.Now().UnixNano()
	qualityResult := qr.Check(cimg, info.Pos, pfs, 5)
	fmt.Println(text, LevelString(qualityResult.Level), "耗时：", (time.Now().UnixNano()-start)/1000000, "毫秒")
}
func EyeStateString(state seeta.EyeState) string {
	switch state {
	case seeta.EYE_OPEN:
		return "EYE_OPEN"
	case seeta.EYE_CLOSE:
		return "EYE_CLOSE"
	case seeta.EYE_RANDOM:
		return "EYE_RANDOM"
	case seeta.EYE_UNKNOWN:
		return "EYE_UNKNOWN"
	default:
		panic("人眼状态范围外")
	}
}
func LevelString(level seeta.QualityLevel) string {
	switch level {
	case seeta.HIGH:
		return "HIGH"
	case seeta.MEDIUM:
		return "MEDIUM"
	case seeta.LOW:
		return "LOW"
	default:
		panic("质量等级范围外")
	}
}
func ReadImgBgr(filename string) (width, height int, pix []uint8) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	rect := img.Bounds()
	x1 := rect.Min.X
	x2 := rect.Max.X
	y1 := rect.Min.Y
	y2 := rect.Max.Y
	height = y2 - y1
	width = x2 - x1
	pix = make([]uint8, height*width*3)
	curr := 0
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ { // Color 转换为 BGR
			nrgba := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			pix[curr] = nrgba.B
			curr++
			pix[curr] = nrgba.G
			curr++
			pix[curr] = nrgba.R
			curr++
		}
	}
	return
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
