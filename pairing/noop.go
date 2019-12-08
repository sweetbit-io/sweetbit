package pairing

// check compliance to interface during compile time
var _ Controller = (*NoopController)(nil)

type NoopController struct {
}

func NewNoopController() *NoopController {
	return &NoopController{}
}

func (c *NoopController) Advertise() error {
	return nil
}

func (c *NoopController) Start() error {
	return nil
}

func (c *NoopController) Stop() error {
	return nil
}
