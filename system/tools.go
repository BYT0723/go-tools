package system

type Event struct{}

func MonitorDirEvent() <-chan *Event {
	return nil
}

func GetOSName() {
}
