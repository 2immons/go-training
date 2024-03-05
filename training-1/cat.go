// cat.go
package main

import (
	"fmt"
)

type Cat struct {
	speed int
}

func (c *Cat) Move() {
	if c.speed == 0 {
		fmt.Println("Error: speed <= 0")
		return
	}
	fmt.Print("Cat is moving with speed = ", c.speed, "\n")
}

func (c *Cat) Jump() {
	fmt.Println("Cat is jumping")
}

func (c *Cat) SetSpeed(speed int) {
	if speed >= 0 {
		c.speed = speed
	} else {
		fmt.Print("Wrong input: speed <= 0 or undefined (", speed, ")\n")
	}
}

func (c *Cat) GetSpeed() int {
	return c.speed
}
