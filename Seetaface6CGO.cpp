#if defined(_WIN32)
	#define IGNORE_Seetaface6CGO_CPP  // windows下编译go时不需要这些代码，忽略
	#if defined(_MSC_VER)
	// Windows平台  Visual Statuio 编译器编译动态库时需要这些代码，不忽略
		#undef IGNORE_Seetaface6CGO_CPP
	#endif
#endif

#ifndef IGNORE_Seetaface6CGO_CPP
	extern "C" {
	#include "Seetaface6CGO.h"
	}

	#include "seeta/FaceDetector.h"
	#include "seeta/FaceLandmarker.h"
	#include "seeta/FaceRecognizer.h"
	#include "seeta/FaceAntiSpoofing.h"
	#include "seeta/FaceTracker.h"
	#include "seeta/QualityOfBrightness.h"  // 亮度评估
	#include "seeta/QualityOfIntegrity.h"   // 完整度评估
	#include "seeta/QualityOfClarity.h"   // 清晰度评估

	#include "seeta/QualityStructure.h"  // 清晰度评估（深度）
	#include "seeta/QualityOfLBN.h"      // 清晰度评估（深度）

	#include "seeta/QualityOfPose.h"     // 姿态评估
	#include "seeta/QualityOfPoseEx.h"   // 姿态评估（深度）

	#include "seeta/QualityOfResolution.h"   // 分辨率评估

	#include "seeta/AgePredictor.h"         //  年龄预测
	#include "seeta/GenderPredictor.h"      // 性别预测
	#include "seeta/MaskDetector.h"         // 口罩检测

	#include <seeta/EyeStateDetector.h>    // 眼睛状态检测

	#include <iostream>


	// 清晰度评估深度，代码
	namespace seeta {
		class QualityOfClarityEx : public QualityRule {
		public:
			QualityOfClarityEx(const char* quality_lbn_model_path, const char* landmark_pts68_model_path) {
				m_lbn = std::make_shared<QualityOfLBN>(ModelSetting(quality_lbn_model_path));
				m_marker = std::make_shared<FaceLandmarker>(ModelSetting(landmark_pts68_model_path));
			}
			QualityOfClarityEx(const char* quality_lbn_model_path, const char* landmark_pts68_model_path, float blur_thresh) {
				m_lbn = std::make_shared<QualityOfLBN>(ModelSetting(quality_lbn_model_path));
				m_marker = std::make_shared<FaceLandmarker>(ModelSetting(landmark_pts68_model_path));
				m_lbn->set(QualityOfLBN::PROPERTY_BLUR_THRESH, blur_thresh);
			}
			QualityResult check(const SeetaImageData& image, const SeetaRect& face, const SeetaPointF* points, int32_t N) override {
				// assert(N == 68);
				auto points68 = m_marker->mark(image, face);
				int light, blur, noise;
				m_lbn->Detect(image, points68.data(), &light, &blur, &noise);
				if (blur == QualityOfLBN::BLUR) {
					return { QualityLevel::LOW, 0 };
				}
				else {
					return { QualityLevel::HIGH, 1 };
				}
			}
		private:
			std::shared_ptr<QualityOfLBN> m_lbn;
			std::shared_ptr<FaceLandmarker> m_marker;
		};
	}
	// 清晰度评估深度，代码,这块代码是文档里的，开源版是没有的
	// 遮挡评估 代码
	//namespace seeta {
	//	class QualityOfNoMask : public QualityRule {
	//	public:
	//		QualityOfNoMask(int64_t ptr) {
	//			m_marker = (seeta::FaceLandmarker*)ptr;
	//		}
	//		QualityResult check(const SeetaImageData& image, const SeetaRect& face, const SeetaPointF* points, int32_t N) override {
	//			auto mask_points = m_marker->mark_v2(image, face);
	//			int mask_count = 0;
	//			for (auto point : mask_points) {
	//				if (point.mask) mask_count++;
	//			}
	//			QualityResult result;
	//			if (mask_count > 0) {
	//				return { QualityLevel::LOW, 1 - float(mask_count) / mask_points.size() };
	//			}
	//			else {
	//				return { QualityLevel::HIGH, 1 };
	//			}
	//		}
	//	private:
	//		seeta::FaceLandmarker* m_marker;
	//	};
	//}
	// 遮挡评估 代码


	// 人脸检测模块
	int64_t NewFaceDetector(SeetaModelSetting model) {
		return int64_t(new seeta::FaceDetector(model));
	}
	void Set(int64_t ptr, int property, double value) {
		((seeta::FaceDetector*)ptr)->set(seeta::FaceDetector::Property(property), value);
	}
	double Get(int64_t ptr, int property) {
		return ((seeta::FaceDetector*)ptr)->get(seeta::FaceDetector::Property(property));
	}
	SeetaFaceInfoArray Detect(int64_t ptr, SeetaImageData image) {
		return ((seeta::FaceDetector*)ptr)->detect(image);
	}
	void DeleteFaceDetector(int64_t ptr){
		delete (seeta::FaceDetector*)ptr;
	}

	// 人脸特征点标记
	int64_t NewFaceLandmarker(SeetaModelSetting model) {
		return int64_t(new seeta::FaceLandmarker(model));
	}
	int GetMarkPointNumber(int64_t ptr) {
		return ((seeta::FaceLandmarker*)ptr)->number();
	}
	void mark(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points) {
		((seeta::FaceLandmarker*)ptr)->mark(image, sr, points);
	}
	//  遮挡评估   遮挡评估功能由FaceLandmarker对象提供，要使用5点标记模型，调用重载方法中带mask参数的那个，长度和标记点长度一致（5），返回值直接赋值给bool，0 false 无遮挡 其他是true，有遮挡
	// 标记人脸特征点同时判断人脸是否被遮挡
	void markWithMask(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points, int32_t* mask) {
		((seeta::FaceLandmarker*)ptr)->mark(image, sr, points,mask);
	}
	void DeleteFaceLandmarker(int64_t ptr) {
		delete (seeta::FaceLandmarker*)ptr;
	}
	// 人脸特征抽取
	int64_t NewFaceRecognizer(SeetaModelSetting model) {
		return int64_t(new seeta::FaceRecognizer(model));
	}
	int GetExtractFeatureSize(int64_t ptr) {
		return ((seeta::FaceRecognizer*)ptr)->GetExtractFeatureSize();
	}
	void Extract(int64_t ptr, SeetaImageData image, SeetaPointF* pf,float* features) {
		((seeta::FaceRecognizer*)ptr)->Extract(image, pf, features);
	}
	void DeleteFaceRecognizer(int64_t ptr) {
		delete (seeta::FaceRecognizer*)ptr;
	}
	// 活体检测
	int64_t NewFaceAntiSpoofing(SeetaModelSetting model) {
		return int64_t(new seeta::FaceAntiSpoofing(model));
	}
	void DeleteFaceAntiSpoofing(int64_t ptr) {
		delete (seeta::FaceAntiSpoofing*)ptr;
	}
	int Predict(int64_t ptr,const SeetaImageData image, const SeetaRect face, const SeetaPointF* points) {
		return (int)(((seeta::FaceAntiSpoofing*)ptr)->Predict(image, face, points));
	}
	int PredictVideo(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points) {
		return (int)(((seeta::FaceAntiSpoofing*)ptr)->PredictVideo(image, face, points));
	}
	void ResetVideo(int64_t ptr) {
		((seeta::FaceAntiSpoofing*)ptr)->ResetVideo();
	}
	void GetPreFrameScore(int64_t ptr,float* clarity, float* reality) {
		((seeta::FaceAntiSpoofing*)ptr)->GetPreFrameScore(clarity, reality);
	}
	void SetVideoFrameCount(int64_t ptr, int32_t number) {
		((seeta::FaceAntiSpoofing*)ptr)->SetVideoFrameCount(number);
	}
	int32_t GetVideoFrameCount(int64_t ptr) {
		return ((seeta::FaceAntiSpoofing*)ptr)->GetVideoFrameCount();
	}
	void SetThreshold(int64_t ptr, float clarity, float reality) {
		((seeta::FaceAntiSpoofing*)ptr)->SetThreshold(clarity, reality);
	}
	void GetThreshold(int64_t ptr, float* clarity, float* reality) {
		((seeta::FaceAntiSpoofing*)ptr)->GetThreshold(clarity, reality);
	}
	void SetBoxThresh(int64_t ptr, float box_thresh) {
		((seeta::FaceAntiSpoofing*)ptr)->SetBoxThresh(box_thresh);
	}
	float GetBoxThresh(int64_t ptr) {
		return ((seeta::FaceAntiSpoofing*)ptr)->GetBoxThresh();
	}
	// 人脸跟踪
	int64_t NewFaceTracker(SeetaModelSetting model, int video_width, int video_height) {
		return int64_t(new seeta::FaceTracker(model,video_width,video_height));
	}
	void SetInterval(int64_t ptr, int interval) {
		((seeta::FaceTracker*)ptr)->SetInterval(interval);
	}
	SeetaTrackingFaceInfoArray Track(int64_t ptr, const SeetaImageData image) {
		return ((seeta::FaceTracker*)ptr)->Track(image);
	}
	SeetaTrackingFaceInfoArray TrackWithFrameNo(int64_t ptr, const SeetaImageData image, int frame_no) {
		return ((seeta::FaceTracker*)ptr)->Track(image, frame_no);
	}
	void SetMinFaceSize(int64_t ptr, int32_t size) {
		((seeta::FaceTracker*)ptr)->SetMinFaceSize(size);
	}
	int32_t GetMinFaceSize(int64_t ptr) {
		return ((seeta::FaceTracker*)ptr)->GetMinFaceSize();
	}
	void SetFaceTrackeThreshold(int64_t ptr,float thresh) {
		((seeta::FaceTracker*)ptr)->SetThreshold(thresh);
	}
	float GetFaceTrackeThreshold(int64_t ptr) {
		return ((seeta::FaceTracker*)ptr)->GetThreshold();
	}
	void SetVideoStable(int64_t ptr,int stable) {
		((seeta::FaceTracker*)ptr)->SetVideoStable(stable != 0);
	}
	int GetVideoStable(int64_t ptr) {
		return ((seeta::FaceTracker*)ptr)->GetVideoStable() ? 1 : 0;
	}
	void SetVideoSize(int64_t ptr,int vidwidth, int vidheight) {
		((seeta::FaceTracker*)ptr)->SetVideoSize(vidwidth, vidheight);
	}
	void Reset(int64_t ptr) {
		((seeta::FaceTracker*)ptr)->Reset();
	}
	void DeleteFaceTracker(int64_t ptr) {
		delete (seeta::FaceTracker*)ptr;
	}
	// 亮度评估
	int64_t NewQualityOfBrightness() {
		return (int64_t)new seeta::QualityOfBrightness();
	}
	int64_t NewQualityOfBrightnessWithParam(float v0, float v1, float v2, float v3) {
		return (int64_t)new seeta::QualityOfBrightness( v0,  v1,  v2,  v3);
	}
	// 完整度评估
	int64_t NewQualityOfIntegrity() {
		return (int64_t)new seeta::QualityOfIntegrity();
	}
	int64_t NewQualityOfIntegrityWithParam(float low, float height) {
		return (int64_t)new seeta::QualityOfIntegrity(low, height);
	}
	// 清晰度评估
	int64_t NewQualityOfClarity() {
		return (int64_t)new seeta::QualityOfClarity();
	}
	int64_t NewQualityOfClarityWithParam(float low, float height) {
		return (int64_t)new seeta::QualityOfClarity(low, height);
	}
	// 清晰度评估深度，代码 清晰度评估深度评估代码后两个参数是没有使用的，可以传NULL和68，固定值
	int64_t NewQualityOfClarityEx(const char* quality_lbn_model_path, const char* landmark_pts68_model_path) {
		return (int64_t)new seeta::QualityOfClarityEx(quality_lbn_model_path, landmark_pts68_model_path);
	}
	int64_t NewQualityOfClarityExWithParam(const char* quality_lbn_model_path, const char* landmark_pts68_model_path, float blur_thresh) {
		return (int64_t)new seeta::QualityOfClarityEx(quality_lbn_model_path ,landmark_pts68_model_path, blur_thresh);
	}
	// 姿态评估
	int64_t NewQualityOfPose() {
		return (int64_t)new seeta::QualityOfPose();
	}
	// 姿态评估（深度）
	int64_t NewQualityOfPoseEx(SeetaModelSetting model) {
		auto qa = new seeta::QualityOfPoseEx(model);
		qa->set(seeta::QualityOfPoseEx::YAW_LOW_THRESHOLD, 25);
		qa->set(seeta::QualityOfPoseEx::YAW_HIGH_THRESHOLD, 10);
		qa->set(seeta::QualityOfPoseEx::PITCH_LOW_THRESHOLD, 20);
		qa->set(seeta::QualityOfPoseEx::PITCH_HIGH_THRESHOLD, 10);
		qa->set(seeta::QualityOfPoseEx::ROLL_LOW_THRESHOLD, 33.33f);
		qa->set(seeta::QualityOfPoseEx::ROLL_HIGH_THRESHOLD, 16.67f);
		return (int64_t)qa;
	}
	float getQualityOfPoseExProperty(int64_t ptr , int property) {
		return ((seeta::QualityOfPoseEx*)ptr)->get(seeta::QualityOfPoseEx::PROPERTY(property));
	}
	void setQualityOfPoseExProperty(int64_t ptr, int property, float value) {
		((seeta::QualityOfPoseEx*)ptr)->set(seeta::QualityOfPoseEx::PROPERTY(property), value);
	}
	// 分辨率评估
	int64_t NewQualityOfResolution() {
		return (int64_t)new seeta::QualityOfResolution();
	}
	int64_t NewQualityOfResolutionWithParam(float low, float height) {
		return (int64_t)new seeta::QualityOfResolution(low, height);
	}
	QualityResult check(int64_t ptr, const SeetaImageData image, const SeetaRect face, const SeetaPointF* points, int32_t N) {
		auto result = ((seeta::QualityRule*)ptr)->check(image, face, points, N);
		return  QualityResult{
			(int)result.level,
			result.score 
		};
	}
	void DeleteQualityRule(int64_t ptr) {
		delete (seeta::QualityRule*)ptr;
	 }
	// 人脸属性检测
	// 年龄预测
	int64_t NewAgePredictor(SeetaModelSetting model) {
		return int64_t(new seeta::AgePredictor(model));
	}
	int PredictAgeWithCrop(int64_t ptr,const SeetaImageData image, const SeetaPointF* points) {
		int age = -1;
		((seeta::AgePredictor*)ptr)->PredictAgeWithCrop(image, points, age);
		return age;
	}
	void DeleteAgePredictor(int64_t ptr) {
		delete (seeta::AgePredictor*)ptr;
	}
	// 性别预测
	int64_t NewGenderPredictor(SeetaModelSetting model) {
		return int64_t(new seeta::GenderPredictor(model));
	}
	// 1 男 2 女 0 未知
	int PredictGenderWithCrop(int64_t ptr, const SeetaImageData image, const SeetaPointF* points) {
		seeta::GenderPredictor::GENDER gender;
		if (((seeta::GenderPredictor*)ptr)->PredictGenderWithCrop(image, points, gender)) {
			return gender == seeta::GenderPredictor::GENDER::MALE ? 1 : gender == seeta::GenderPredictor::GENDER::MALE ? 2 : 0;
		}else {
			return 0;
		}
	}
	void DeleteGenderPredictor(int64_t ptr) {
		delete (seeta::GenderPredictor*)ptr;
	}
	// 口罩检测
	int64_t NewMaskDetector(SeetaModelSetting model) {
		return int64_t(new seeta::MaskDetector(model));
	}
	MaskFace DetectMask(int64_t ptr, const SeetaImageData image, const SeetaRect face) {
		float score;
		auto b = ((seeta::MaskDetector*)ptr)->detect(image,face,&score);
		return MaskFace{ score,b ? 1 : 0 };
	}
	void DeleteMaskDetector(int64_t ptr) {
		delete (seeta::MaskDetector*)ptr;
	}
	// 眼睛状态检测
	int64_t NewEyeStateDetector(SeetaModelSetting model) {
		return int64_t(new seeta::EyeStateDetector(model));
	}
	EyeStates DetectEyeState(int64_t ptr, const SeetaImageData image, SeetaPointF* points) {
		seeta::EyeStateDetector::EYE_STATE left_eye, right_eye;
		((seeta::EyeStateDetector*)ptr)->Detect(image, points, left_eye, right_eye);
		return EyeStates{ left_eye,right_eye };
	}
	void DeleteEyeStateDetector(int64_t ptr) {
		delete (seeta::EyeStateDetector*)ptr;
	}
#endif
