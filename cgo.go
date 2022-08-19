package seeta

/*
linux下编译要先设置环境变量LD_LIBRARY_PATH和LIBRARY_PATH
指定编译时库路径和运行时库路径
windows下编译后运行需要和依赖库放在同一目录下

在linux设置相对路径要使用 -Wl,-rpath,\$ORIGIN/lib:\$ORIGIN/libs -Wl,--disable-new-dtags
原本直接用-Wl,-rpath=\$ORIGIN/lib:\$ORIGIN/libs 就可以了，但是编译器升级后rpath的含义好像变了，要用--disable-new-dtags表示使用旧的行为

linux端没有编译动态链接库，直接编译C++语法进go包内，所以要开启  CXXFLAGS: -std=c++11 且使用CPPFLAGS: -I./include来引入头文件
windows端因为库文件是使用MSVC编译的，跨编译器只能使用纯C头文件，且调用必须通过动态链接库，所以使用CFLAGS: -I./include引入头文件
windows下没法设置运行时库的路径，只能搜索系统路径及当前目录等几个特定的路径
*/

/*
#cgo linux CPPFLAGS: -I./include
#cgo linux LDFLAGS: -lSeetaAgePredictor600 -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600  -lSeetaMaskDetector200  -lSeetaPoseEstimation600  -lSeetaQualityAssessor300
#cgo linux CXXFLAGS: -std=c++11
#cgo linux LDFLAGS: -Wl,-rpath,\$ORIGIN/lib:\$ORIGIN/libs -Wl,--disable-new-dtags
#cgo linux,arm64 LDFLAGS: -Wl,--no-as-needed -ldl
#cgo windows CFLAGS: -I./include
#cgo windows LDFLAGS: -L${SRCDIR}/lib -lSeetaface6CGO
*/
import "C"
