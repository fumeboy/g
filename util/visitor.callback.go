package util

type Callbacker struct {
	CB []func()error // 使用 callback 是为了不遍历AST的时候修改 AST，遍历结束了再修改
}

func (c *Callbacker) AddCallback(fn func() error){
	c.CB = append(c.CB, fn)
}

func (c *Callbacker) TODO() func()error {
	if len(c.CB) == 0{
		return nil
	}
	return func() error{
		for _,fn := range c.CB {
			if err := fn(); err != nil{
				return err
			}
		}
		return nil
	}
}

func (c *Callbacker) DO() error {
	for _,fn := range c.CB {
		if err := fn(); err != nil{
			return err
		}
	}
	return nil
}
