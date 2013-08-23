package core

type frame struct {
	channels []Quantity
}

func newFrame(numChannels Index) SampleFrame {
	return &frame{channels: make([]Quantity, numChannels)}
}

func (f *frame) NumChannels() Index {
	return Index(len(f.channels))
}

func (f *frame) ChannelVal(i Index) Quantity {
	return f.channels[i]
}

func (f *frame) SetChannelVal(i Index, v Quantity) {
	f.channels[i] = v
}

type framePool struct {
	frameChannel chan SampleFrame
}

func NewFramePool(poolSize Index, numChannels Index) SampleFramePool {
	frameChannel := make(chan SampleFrame, poolSize)
	for i := Index(0); i < poolSize; i++ {
		frameChannel <- newFrame(numChannels)
	}

	return &framePool{frameChannel: frameChannel}
}

func (fp *framePool) DequeueFrame() SampleFrame {
	return <-fp.frameChannel
}

func (fp *framePool) EnqueueFrame(frame SampleFrame) {
	fp.frameChannel <- frame
}
