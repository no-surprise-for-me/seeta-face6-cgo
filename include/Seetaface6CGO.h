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

// �������ģ��
int64_t NewFaceDetector(SeetaModelSetting model);
void Set(int64_t ptr, int property, double value);
double Get(int64_t ptr, int property);
SeetaFaceInfoArray Detect(int64_t ptr, SeetaImageData image);
void DeleteFaceDetector(int64_t ptr);

// ������������
int64_t NewFaceLandmarker(SeetaModelSetting model);
int GetMarkPointNumber(int64_t ptr);
void mark(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points);

// �������������ͬʱ�ж������Ƿ��ڵ�
void markWithMask(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points, int32_t* mask);
void DeleteFaceLandmarker(int64_t ptr);

// ����������ȡ
int64_t NewFaceRecognizer(SeetaModelSetting model);
int GetExtractFeatureSize(int64_t ptr);
void Extract(int64_t ptr, SeetaImageData image, SeetaPointF* pf, float* features);
void DeleteFaceRecognizer(int64_t ptr);

// ������
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

// ��������

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

// ͼƬ��������
struct QualityResult
{
	int   level;  // 0,1,2��Ӧ LOW ,MEDIUM ,HIGH 
	float score; 
};

// ��������
int64_t NewQualityOfBrightness();
int64_t NewQualityOfBrightnessWithParam(float v0, float v1, float v2, float v3);
// ����������
int64_t NewQualityOfIntegrity();
int64_t NewQualityOfIntegrityWithParam(float low, float height);
// ����������
int64_t NewQualityOfClarity();
int64_t NewQualityOfClarityWithParam(float low, float height);
// ������������ȣ����� ��������������������������������û��ʹ�õģ����Դ�NULL��68���̶�ֵ
int64_t NewQualityOfClarityEx(const char* quality_lbn_model_path, const char* landmark_pts68_model_path);
int64_t NewQualityOfClarityExWithParam(const char* quality_lbn_model_path, const char* landmark_pts68_model_path, float blur_thresh);
// ��̬����
int64_t NewQualityOfPose();
// ��̬��������ȣ�
int64_t NewQualityOfPoseEx(SeetaModelSetting model);
float getQualityOfPoseExProperty(int64_t ptr, int property);
void setQualityOfPoseExProperty(int64_t ptr, int property, float value);
// �ֱ�������
int64_t NewQualityOfResolution();
int64_t NewQualityOfResolutionWithParam(float low, float height);
QualityResult check(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points, int32_t N);
void DeleteQualityRule(int64_t ptr);

// ����Ԥ��
int64_t NewAgePredictor(SeetaModelSetting model);
int PredictAgeWithCrop(int64_t ptr, const SeetaImageData image, const SeetaPointF* points);
void DeleteAgePredictor(int64_t ptr);
// �Ա�Ԥ��
int64_t NewGenderPredictor(SeetaModelSetting model);
// 1 �� 2 Ů 0 δ֪
int PredictGenderWithCrop(int64_t ptr, const SeetaImageData image, const SeetaPointF* points);
void DeleteGenderPredictor(int64_t ptr);

// ���ּ��
struct MaskFace {
	float score; //  ���˿��ֵ����Ŷȣ�score����0.5������Ϊ�Ǽ������˿���
	int   mask;  // 0 û������  1 ��������
};

int64_t NewMaskDetector(SeetaModelSetting model);
MaskFace DetectMask(int64_t ptr, const SeetaImageData image, const SeetaRect face);
void DeleteMaskDetector(int64_t ptr);


// �۾�״̬���
// enum EYE_STATE { EYE_CLOSE, EYE_OPEN, EYE_RANDOM, EYE_UNKNOWN };
struct EyeStates {
	int left;
	int right;
};
int64_t NewEyeStateDetector(SeetaModelSetting model);
EyeStates DetectEyeState(int64_t ptr,SeetaImageData image,SeetaPointF* points);
void DeleteEyeStateDetector(int64_t ptr);