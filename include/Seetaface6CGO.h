#pragma once
typedef struct SeetaFaceInfoArray SeetaFaceInfoArray;
typedef struct SeetaFaceInfo SeetaFaceInfo;
typedef struct SeetaPointF SeetaPointF;
typedef struct SeetaRect SeetaRect;
typedef struct SeetaModelSetting SeetaModelSetting;
typedef struct SeetaPointFs SeetaPointFs;
typedef struct SeetaImageData SeetaImageData;
typedef struct SeetaTrackingFaceInfoArray SeetaTrackingFaceInfoArray;
typedef struct SeetaTrackingFaceInfo SeetaTrackingFaceInfo;
typedef struct MaskFace MaskFace;
typedef struct QualityResult QualityResult;
typedef struct SeetaBuffer SeetaBuffer;
typedef struct EyeStates EyeStates;

#include <seeta/CFaceInfo.h>
#include <seeta/CTrackingFaceInfo.h>
struct SeetaPointFs 
{
	SeetaPointF* data;
	int size;
};

// 人脸检测模块
int64_t NewFaceDetector(SeetaModelSetting model);
void Set(int64_t ptr, int property, double value);
double Get(int64_t ptr, int property);
SeetaFaceInfoArray Detect(int64_t ptr, SeetaImageData image);
void DeleteFaceDetector(int64_t ptr);

// 人脸特征点标记
int64_t NewFaceLandmarker(SeetaModelSetting model);
int GetMarkPointNumber(int64_t ptr);
void mark(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points);

// 标记人脸特征点同时判断人脸是否被遮挡
void markWithMask(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points, int32_t* mask);
void DeleteFaceLandmarker(int64_t ptr);

// 人脸特征抽取
int64_t NewFaceRecognizer(SeetaModelSetting model);
int GetExtractFeatureSize(int64_t ptr);
void Extract(int64_t ptr, SeetaImageData image, SeetaPointF* pf, float* features);
void DeleteFaceRecognizer(int64_t ptr);

// 活体检测
int64_t NewFaceAntiSpoofing(SeetaModelSetting model);
void DeleteFaceAntiSpoofing(int64_t ptr);
int Predict(int64_t ptr,const SeetaImageData image, const SeetaRect face, const SeetaPointF* points);
int PredictVideo(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points);
void ResetVideo(int64_t ptr);
void GetPreFrameScore(int64_t ptr, float* clarity, float* reality);
void SetVideoFrameCount(int64_t ptr, int32_t number);
int32_t GetVideoFrameCount(int64_t ptr);
void SetThreshold(int64_t ptr, float clarity, float reality);
void GetThreshold(int64_t ptr, float* clarity, float* reality);
void SetBoxThresh(int64_t ptr, float box_thresh);
float GetBoxThresh(int64_t ptr);

// 人脸跟踪

int64_t NewFaceTracker(SeetaModelSetting model, int video_width, int video_height);
void SetInterval(int64_t ptr, int interval);
SeetaTrackingFaceInfoArray Track(int64_t ptr, const SeetaImageData image);
SeetaTrackingFaceInfoArray TrackWithFrameNo(int64_t ptr, const SeetaImageData image, int frame_no);
void SetMinFaceSize(int64_t ptr, int32_t size);
int32_t GetMinFaceSize(int64_t ptr);
void SetFaceTrackeThreshold(int64_t ptr, float thresh);
float GetFaceTrackeThreshold(int64_t ptr);
void SetVideoStable(int64_t ptr, int stable);
int GetVideoStable(int64_t ptr);
void SetVideoSize(int64_t ptr, int vidwidth, int vidheight);
void Reset(int64_t ptr);
void DeleteFaceTracker(int64_t ptr);

// 图片质量评估
struct QualityResult
{
	int   level;  // 0,1,2对应 LOW ,MEDIUM ,HIGH 
	float score; 
};

// 亮度评估
int64_t NewQualityOfBrightness();
int64_t NewQualityOfBrightnessWithParam(float v0, float v1, float v2, float v3);
// 完整度评估
int64_t NewQualityOfIntegrity();
int64_t NewQualityOfIntegrityWithParam(float low, float height);
// 清晰度评估
int64_t NewQualityOfClarity();
int64_t NewQualityOfClarityWithParam(float low, float height);
// 清晰度评估深度，代码 清晰度评估深度评估代码后两个参数是没有使用的，可以传NULL和68，固定值
int64_t NewQualityOfClarityEx(const char* quality_lbn_model_path, const char* landmark_pts68_model_path);
int64_t NewQualityOfClarityExWithParam(const char* quality_lbn_model_path, const char* landmark_pts68_model_path, float blur_thresh);
// 姿态评估
int64_t NewQualityOfPose();
// 姿态评估（深度）
int64_t NewQualityOfPoseEx(SeetaModelSetting model);
float getQualityOfPoseExProperty(int64_t ptr, int property);
void setQualityOfPoseExProperty(int64_t ptr, int property, float value);
// 分辨率评估
int64_t NewQualityOfResolution();
int64_t NewQualityOfResolutionWithParam(float low, float height);
QualityResult check(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points, int32_t N);
void DeleteQualityRule(int64_t ptr);

// 年龄预测
int64_t NewAgePredictor(SeetaModelSetting model);
int PredictAgeWithCrop(int64_t ptr, const SeetaImageData image, const SeetaPointF* points);
void DeleteAgePredictor(int64_t ptr);
// 性别预测
int64_t NewGenderPredictor(SeetaModelSetting model);
// 1 男 2 女 0 未知
int PredictGenderWithCrop(int64_t ptr, const SeetaImageData image, const SeetaPointF* points);
void DeleteGenderPredictor(int64_t ptr);

// 口罩检测
struct MaskFace {
	float score; //  戴了口罩的置信度，score超过0.5，则认为是检测带上了口罩
	int   mask;  // 0 没戴口罩  1 戴口罩了
};

int64_t NewMaskDetector(SeetaModelSetting model);
MaskFace DetectMask(int64_t ptr, const SeetaImageData image, const SeetaRect face);
void DeleteMaskDetector(int64_t ptr);


// 眼睛状态检测
// enum EYE_STATE { EYE_CLOSE, EYE_OPEN, EYE_RANDOM, EYE_UNKNOWN };
struct EyeStates {
	int left;
	int right;
};
int64_t NewEyeStateDetector(SeetaModelSetting model);
EyeStates DetectEyeState(int64_t ptr,SeetaImageData image,SeetaPointF* points);
void DeleteEyeStateDetector(int64_t ptr);