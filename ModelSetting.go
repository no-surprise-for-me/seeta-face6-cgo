package seeta

/*
#include <Seetaface6CGO.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type SeetaDevice uint32

const (
	SEETA_DEVICE_AUTO SeetaDevice = 0
	SEETA_DEVICE_CPU  SeetaDevice = 1
	SEETA_DEVICE_GPU  SeetaDevice = 2
)

type SeetaModelSetting struct {
	Device SeetaDevice
	Id     int      // when Device is GPU, Id means GPU Id
	Model  []string // Model string terminate with nullptr
}

func (sms *SeetaModelSetting) CSeetaModelSetting() (C.SeetaModelSetting, func()) {
	var cSeetaModelSetting C.SeetaModelSetting
	cSeetaModelSetting.id = C.int(sms.Id)
	cSeetaModelSetting.device = uint32(sms.Device) // Go中的uint32自动转换成C中的枚举

	var models []*C.char
	for _, s := range sms.Model {
		var cs *C.char = C.CString(s)
		models = append(models, cs)
	}
	models = append(models, nil)
	cSeetaModelSetting.model = (**C.char)(unsafe.Pointer(&models[0]))

	return cSeetaModelSetting, func() {
		for _, m := range models {
			if m != nil {
				C.free(unsafe.Pointer(m))
			}
		}
	}
}
