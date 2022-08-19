package seeta

// #include <Seetaface6CGO.h>
import "C"

type SeetaPointFs struct {
	size int
	data []C.SeetaPointF // 因为特征抽取需要这个切片底层的数组，还是保留的好
}

func (pfs SeetaPointFs) slice() []SeetaPointF {
	size := pfs.size
	var result = make([]SeetaPointF, size)
	for i, info := range pfs.data {
		result[i] = SeetaPointF{
			X: int(info.x),
			Y: int(info.y),
		}
	}
	return result
}

type SeetaPointF struct {
	X int
	Y int
}
