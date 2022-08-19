cmake_minimum_required (VERSION 3.0)
# 项目名称
project (Seetaface6CGO)
set(CMAKE_BUILD_TYPE "Release")
set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS} -Wall -s -std=c++11 -g -O2")
set(CMAKE_INSTALL_PREFIX "../")
include_directories("./include/") 
# ================================================================================
# 引入seetaface6依赖库文件
IF (WIN32)
	IF(CMAKE_CL_64)
		MESSAGE(STATUS "platform: windows,x64")
		set(CMAKE_CONFIGURATION_TYPES "Release")
		# 这里要重新设置自己的seetaface6依赖路径
		link_directories("D:/user/14461/DeskTop/seetafaceJNI/SeetaFace6JNI/lib/win-x64")
	ELSE()
		MESSAGE(FATAL_ERROR "暂不支持的32位系统下编译，可以手动修改CMakeLists.txt")
	ENDIF()
ELSEIF (APPLE)
	MESSAGE(FATAL_ERROR "not support Apple system.")
ELSEIF (UNIX)
	MESSAGE(FATAL_ERROR "linux下不需要编译动态链接库")
ENDIF ()
# ================================================================================
#aux_source_directory(./src SRC_LIST)

add_library(${PROJECT_NAME} SHARED ./Seetaface6CGO.cpp ./lib/Seetaface6CGO.def)

target_link_libraries(${PROJECT_NAME} 
			SeetaAgePredictor600
			SeetaEyeStateDetector200
			SeetaFaceAntiSpoofingX600
			SeetaFaceDetector600
			SeetaFaceLandmarker600
			SeetaFaceRecognizer610
			SeetaFaceTracking600
			SeetaGenderPredictor600
			SeetaMaskDetector200
			SeetaPoseEstimation600
			SeetaQualityAssessor300
)

install(TARGETS ${PROJECT_NAME}
        EXPORT ${PROJECT_NAME}
        LIBRARY DESTINATION lib  # 动态库安装路径
        ARCHIVE DESTINATION lib  # 静态库安装路径
        RUNTIME DESTINATION lib  # 可执行文件安装路径
)

