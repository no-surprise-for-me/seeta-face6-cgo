package scanner

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
)

var engine *gin.Engine

func GetEngine() *gin.Engine {
	if engine == nil {
		gin.SetMode(gin.ReleaseMode)
		logRoot := getLogRoot()
		setHttpLogWriter(logRoot)
		setHttpErrorLogWriter(logRoot)
		engine = gin.Default()
	}
	return engine
}
func getLogRoot() string {
	dir := filepath.Dir(os.Args[0])
	if strings.HasSuffix(dir, "bin") {
		dir = filepath.Join(dir, "..", "log")
	} else {
		dir = filepath.Join(dir, "log")
	}
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0644)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return dir
}
func setHttpLogWriter(logRoot string) {
	logfile, err := os.Create(filepath.Join(logRoot, "gin_http.log"))
	if err != nil {
		panic("创建gin_http.log文件失败" + err.Error())
	} else {
		gin.DefaultWriter = logfile
	}
}
func setHttpErrorLogWriter(logRoot string) {
	logfile, err := os.Create(filepath.Join(logRoot, "gin_http_error.log"))
	if err != nil {
		panic("创建gin_http_error.log文件失败" + err.Error())
	} else {
		gin.DefaultErrorWriter = logfile
	}
}

func SetCorsHeaders(ctx *gin.Context) {
	ctx.Header("Content-Security-Policy", "upgrade-insecure-requests")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Content-Security-Policy", "upgrade-insecure-requests")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE, PUT")
	ctx.Header("Access-Control-Max-Age", "3600")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Expose-Headers", "Set-Cookie")
}
