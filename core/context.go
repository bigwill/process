package core

type context struct {
	sampleRate  Quantity
	numChannels Integer
	framePool   SampleFramePool
}

func NewContext(sampleRate Quantity, numChannels Integer, framePoolSize Integer) Context {
	return &context{sampleRate: sampleRate,
		numChannels: numChannels,
		framePool:   NewFramePool(framePoolSize, numChannels)}
}

func (c *context) SampleRate() Quantity {
	return c.sampleRate
}

func (c *context) NumChannels() Integer {
	return c.numChannels
}

func (c *context) FramePool() SampleFramePool {
	return c.framePool
}
