package core

type ctrlMsg struct {
	code int16
}

func (m *ctrlMsg) Code() int16 {
	return m.code
}

type paramCtrlMsg struct {
	ctrlMsg
	index Integer
	pos   Quantity
}

func (m *paramCtrlMsg) Index() Integer {
	return m.index
}

func (m *paramCtrlMsg) Pos() Quantity {
	return m.pos
}

func NewQuitControlMessage() ControlMessage {
	return &ctrlMsg{Quit}
}

func NewParamControlMessage(index Integer, pos Quantity) ParamControlMessage {
	return &paramCtrlMsg{ctrlMsg{ParamChange}, index, pos}
}
