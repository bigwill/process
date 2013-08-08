package core

type monMsg struct {
	code int16
	paramerName string
}

func (m *monMsg) Code() int16 {
	return m.code
}

func (m *monMsg) ParamerName() string {
	return m.paramerName
}

type errorMonMsg struct {
	monMsg
	err error
}

func (m *errorMonMsg) Err() error {
	return m.err
}

func MonitorError(paramerName string, err error) ErrorMonitorMessage {
	return &errorMonMsg{monMsg{Error, paramerName}, err}
}