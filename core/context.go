package core

type context struct {
	sampleRate  Quantity
	numChannels Index
	framePool   SampleFramePool
}

func NewContext(sampleRate Quantity, numChannels Index, framePoolSize Index) Context {
	return &context{sampleRate: sampleRate,
		numChannels: numChannels,
		framePool:   NewFramePool(framePoolSize, numChannels)}
}

func (c *context) SampleRate() Quantity {
	return c.sampleRate
}

func (c *context) NumChannels() Index {
	return c.numChannels
}

func (c *context) FramePool() SampleFramePool {
	return c.framePool
}
