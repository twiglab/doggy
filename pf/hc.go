package pf

type (
	hc struct {
		countHandler   CountHandler
		densityHandler DensityHandler
		deviceRegister DeviceRegister
	}

	Option func(*hc)
)

func SetCountHandler(h CountHandler) Option {
	return func(c *hc) {
		c.countHandler = h
	}
}

func SetDensityHandler(h DensityHandler) Option {
	return func(c *hc) {
		c.densityHandler = h
	}
}

func SetDeviceRegister(h DeviceRegister) Option {
	return func(c *hc) {
		c.deviceRegister = h
	}
}
