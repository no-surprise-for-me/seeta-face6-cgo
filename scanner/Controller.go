package scanner

import (
	"fmt"
	"gitee.com/no_surprise_for_me/seeta-face-cgo"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strconv"
)

func init() {
	GetEngine().POST("/search", SearchController)
	GetEngine().GET("/init", InitController)
	GetEngine().GET("/status", StatusQueryController)
	searchFaceEngine = NewFaceEngine()
}

var searchFaceEngine *FaceEngine

func SearchController(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		panic(err)
	}
	open, err := file.Open()
	if err != nil {
		panic(err)
	}
	cimg, err := seeta.NewSeetaFaceImageFromReader(open)
	if err != nil {
		panic(err)
	}
	faceInfo, pointFs := searchFaceEngine.PreHandle(cimg)
	if faceInfo == nil {
		ctx.JSON(200, gin.H{
			"msg": "上传图片中未检测到人脸",
		})
		return
	}
	check := searchFaceEngine.Check(cimg, &faceInfo.Pos, pointFs)
	if check != "" {
		ctx.JSON(200, gin.H{
			"msg": check,
		})
		return
	}
	features := searchFaceEngine.Extract(cimg, pointFs)
	path, num, similarity := Search(features)
	ctx.JSON(200, gin.H{
		"msg": fmt.Sprintf("从%v张图片中找到最相似人脸：%v,相似度为：%f", num, path, similarity),
	})
}
func InitController(ctx *gin.Context) {
	path := ctx.Query("path")
	thread := ctx.Query("thread")
	if path == "" {
		ctx.JSON(200, gin.H{
			"msg": "path参数不能为空",
		})
		return
	}
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(200, gin.H{
				"msg": "路径不存在：" + path,
			})
			return
		} else {
			panic(err)
		}
	}
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			panic(err)
		}
	}
	go func() {
		threadNum, err := strconv.Atoi(thread)
		if err != nil {
			threadNum = 1
		}
		DoScan(path, threadNum)
	}()
	ctx.JSON(200, gin.H{
		"msg": "开始初始化：" + path,
	})
}
func StatusQueryController(ctx *gin.Context) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	StatusInfo.FaceNum = len(FaceInfo)
	ctx.JSON(200, StatusInfo)
}
