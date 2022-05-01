package volume

type Volume interface {
	GetRequestGroup() string
}

type volume struct {
	requestGroup string
}

func (v volume) GetRequestGroup() string {
	return v.requestGroup
}

func MakeVolume(requestGroup string) Volume {
	return volume{requestGroup: requestGroup}
}
