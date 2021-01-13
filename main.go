package main

func main() {}

type counter struct {
	n int
}

func (c *counter) inc() {
	c.n++
}

func (c *counter) dec() {
	// comment out to fail PBT tests
	// if c.n > 3 {
	// 	c.n -= 2
	// 	return
	// }

	c.n--
}

func (c *counter) reset() {
	c.n = 0
}
