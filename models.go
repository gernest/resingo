package resingo

//DeviceType is the identity of the the device that is supported by resin.
type DeviceType int

// supported devices
const (
	Artik10 DeviceType = iota
	Artik5
	BeagleboneBlack
	HumingBoard
	IntelAdison
	IntelNuc
	Nitrogen6x
	OdroidC1
	OdroidXu4
	Parallella
	RaspberryPi
	RaspberryPi2
	RaspberryPi3
	Ts4900
	Ts700
	ViaVabx820Quad
	ZyncXz702
)

func (d DeviceType) String() string {
	switch d {
	case Artik10:
		return "artik10"
	case Artik5:
		return "artik5"
	case BeagleboneBlack:
		return "beaglebone-black"
	case HumingBoard:
		return "hummingboard"
	case IntelAdison:
		return "intel-edison"
	case IntelNuc:
		return "intel-nuc"
	case Nitrogen6x:
		return "nitrogen6x"
	case OdroidC1:
		return "odroid-c1"
	case OdroidXu4:
		return "odroid-xu4"
	case Parallella:
		return "parallella"
	case RaspberryPi:
		return "raspberry-pi"
	case RaspberryPi2:
		return "raspberry-pi2"
	case RaspberryPi3:
		return "raspberrypi3"
	case Ts4900:
		return "ts4900"
	case Ts700:
		return "ts7700"
	case ViaVabx820Quad:
		return "via-vab820-quad"
	case ZyncXz702:
		return "zynq-xz702"
	}
	return "Unknown device"
}

//Repository is a resin remote repository
type Repository struct {
	URL    string
	Commit string
}

//User a resin user
type User struct {
	ID       int64 `json:"__id"`
	Metadata struct {
		URI string `json:"uri"`
	} `json:"__deferred"`
}
