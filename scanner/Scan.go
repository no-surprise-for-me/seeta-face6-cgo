package scanner

import (
	"fmt"
	"gitee.com/no_surprise_for_me/seeta-face-cgo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 绝对路径 --》 人脸特征
var FaceInfo = make(map[string][]byte)
var cacheLock = &sync.RWMutex{}

func Search(source []float32) (string, int, float32) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	var max float32
	var path string
	for s, f := range FaceInfo {
		floats := seeta.CompareBytesAndFloats(f, source)
		if floats > max {
			max = floats
			path = s
		}
	}
	return path, len(FaceInfo), max
}

var StatusInfo = struct {
	ScanEndTime   string
	ScanStartTime string

	FaceNum   int
	MissFace  int
	PoolQu    int
	HandleErr int
}{}

func DoScan(absPath string, threadNum int) {
	StatusInfo.ScanStartTime = time.Now().Format("2006-01-02_15:04:05")
	defer func() {
		StatusInfo.ScanEndTime = time.Now().Format("2006-01-02_15:04:05")
	}()
	for i := 0; i < threadNum; i++ {
		select {
		case goCountLimitChan <- NewFaceEngine():
		default:
			threadNum = i //防止threadNum参数过大，超过goCountLimitChan的长度就停止
		}
	}
	defer func() {
		group.Wait()
		for i := 0; i < threadNum; i++ {
			(<-goCountLimitChan).Close()
		}
	}()
	Scan(absPath, func(dir, filename string) {
		ext := strings.ToLower(filepath.Ext(filename))
		if ext == ".jpg" ||
			ext == ".jpeg" ||
			ext == ".png" ||
			ext == ".gif" {
			abs := filepath.Join(dir, filename)
			scanFaceEngine := <-goCountLimitChan
			group.Add(1)
			go Cache(abs, scanFaceEngine)
		}
	})
}
func Scan(absPath string, handle func(dir, filename string)) {
	dir, err := ioutil.ReadDir(absPath)
	if err != nil {
		fmt.Println("读取目录出错：" + err.Error())
		return
	}
	fmt.Println("开始处理目录：", absPath, "，共", len(dir), "个文件")
	for _, d := range dir {
		if d.IsDir() {
			Scan(filepath.Join(absPath, d.Name()), handle)
		} else {
			name := d.Name()
			handle(absPath, name)
		}
	}
}

var group = &sync.WaitGroup{}
var goCountLimitChan = make(chan *FaceEngine, 12) // 并发不会超过12个
func Cache(abs string, scanFaceEngine *FaceEngine) {
	defer group.Done()
	defer func() {
		goCountLimitChan <- scanFaceEngine
	}()
	faceFile := abs[:strings.LastIndexByte(abs, '.')] + ".FaceFeature"
	file, err := ioutil.ReadFile(faceFile)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		cacheLock.Lock()
		FaceInfo[abs] = file
		cacheLock.Unlock()
		return
	}
	cimg, err := seeta.NewSeetaFaceImageFromFile(abs)
	if err != nil {
		fmt.Println("处理：", abs, "文件出错")
		StatusInfo.HandleErr++
		return
	}
	faceInfo, pointFs := searchFaceEngine.PreHandle(cimg)
	if faceInfo == nil {
		fmt.Println(abs, "中未检测到人脸")
		StatusInfo.MissFace++
		return
	}
	check := searchFaceEngine.Check(cimg, &faceInfo.Pos, pointFs)
	if check != "" {
		fmt.Println(abs, "  ", check)
		StatusInfo.PoolQu++
		return
	}
	features := searchFaceEngine.Extract(cimg, pointFs)
	file = seeta.Float64ToByte(features)
	go func() {
		err = ioutil.WriteFile(faceFile, file, 0644)
		if err != nil {
			fmt.Println("人脸特征缓存失败")
		}
	}()
	cacheLock.Lock()
	FaceInfo[abs] = file
	cacheLock.Unlock()
}
