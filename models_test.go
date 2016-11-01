package resingo

import "testing"

func TestDeviceType(t *testing.T) {
	sample := []struct {
		typ    DeviceType
		expect string
	}{
		{Artik10, "artik10"},
		{Artik5, "artik5"},
		{BeagleboneBlack, "beaglebone-black"},
		{HumingBoard, "hummingboard"},
		{IntelAdison, "intel-edison"},
		{IntelNuc, "intel-nuc"},
		{Nitrogen6x, "nitrogen6x"},
		{OdroidC1, "odroid-c1"},
		{OdroidXu4, "odroid-xu4"},
		{Parallella, "parallella"},
		{RaspberryPi, "raspberry-pi"},
		{RaspberryPi2, "raspberry-pi2"},
		{RaspberryPi3, "raspberrypi3"},
		{Ts4900, "ts4900"},
		{Ts700, "ts7700"},
		{ViaVabx820Quad, "via-vab820-quad"},
		{ZyncXz702, "zynq-xz702"},
	}

	for _, v := range sample {
		if v.typ.String() != v.expect {
			t.Errorf("expetcted %s got %v", v.expect, v.typ)
		}
	}
}
