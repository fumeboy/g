package util

type Callbacker struct {
	cb []func()error
}

func (c *Callbacker) AddCallback(fn func() error){
	c.cb = append(c.cb, fn)
}

func (c *Callbacker) TODO() func()error {
	if len(c.cb) == 0{
		return nil
	}
	return func() error{
		for _,fn := range c.cb{
			if err := fn(); err != nil{
				return err
			}
		}
		return nil
	}
}

func (c *Callbacker) DO() error {
	for _,fn := range c.cb{
		if err := fn(); err != nil{
			return err
		}
	}
	return nil
}
