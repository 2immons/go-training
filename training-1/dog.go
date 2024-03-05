// dog.go
package main

import (
	"fmt"
)

type Dog struct {
	speed int
}

func (d *Dog) Move() {
	if d.speed == 0 {
		fmt.Println("Error: speed <= 0")
		return
	}
	fmt.Print("Dog is moving with speed = ", d.speed, "\n")
}

func (d *Dog) Jump() {
	fmt.Println("Dog is jumping")
}

func (d *Dog) SetSpeed(speed int) {
	if speed >= 0 {
		d.speed = speed
	} else {
		fmt.Print("Wrong input: speed <= 0 or undefined (", speed, ")\n")
	}
}

func (d *Dog) GetSpeed() int {
	return d.speed
}
