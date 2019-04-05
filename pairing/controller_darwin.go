package pairing

type Controller struct {
	log Logger
}

func NewController(config *Config) (*Controller, error) {
	controller := &Controller{}

	if config.Logger != nil {
		controller.log = config.Logger
	} else {
		controller.log = noopLogger{}
	}

	return controller, nil
}

func (c *Controller) Start() error {
	c.log.Infof("Won't start pairing controller: It's currently only supported on Linux.")
	return nil
}

func (c *Controller) Stop() error {
	c.log.Infof("Won't stop pairing controller: It's currently only supported on Linux.")
	return nil
}
