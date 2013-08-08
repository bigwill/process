package core

func MonitorQuit() MonitorMessage {
	return MonitorMessage{Code: Quit}
}

func MonitorError(err error) MonitorMessage {
	return MonitorMessage{Code: Error, Error: err}
}