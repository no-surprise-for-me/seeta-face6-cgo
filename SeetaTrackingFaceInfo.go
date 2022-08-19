package seeta

type SeetaTrackingFaceInfo struct {
	Pos   SeetaRect
	Score float32 // 人脸跟踪中，这个属性一直是0，不知道为什么。。。虽然不是很重要。

	FrameNo int // C++代码 中字段名为 frame_no
	PID     int
	Step    int
}
