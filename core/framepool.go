package core

type frame struct {
	channels []Quantity
}

func newFrame(numChannels Integer) SampleFrame {
	return &frame{channels: make([]Quantity, numChannels)}
}

func (f *frame) NumChannels() Integer {
	return Integer(len(f.channels))
}

func (f *frame) ChannelVal(i Integer) Quantity {
	return f.channels[i]
}

func (f *frame) SetChannelVal(i Integer, v Quantity) {
	f.channels[i] = v
}

type framePool struct {
	frameChannel chan SampleFrame
}

func NewFramePool(poolSize Integer, numChannels Integer) SampleFramePool {
	frameChannel := make(chan SampleFrame, poolSize)
	for i := Integer(0); i < poolSize; i++ {
		frameChannel <- newFrame(numChannels)
	}

	return &framePool{frameChannel: frameChannel}
}

func (fp *framePool) Size() Integer {
	return Integer(cap(fp.frameChannel))
}

func (fp *framePool) NumAvailable() Integer {
	return Integer(len(fp.frameChannel))
}

func (fp *framePool) DequeueFrame() SampleFrame {
	return <-fp.frameChannel
}

func (fp *framePool) EnqueueFrame(frame SampleFrame) {
	fp.frameChannel <- frame
}
