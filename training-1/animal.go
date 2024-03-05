// animal.go
package main

type Animal interface {
	Move()
	Jump()
	SetSpeed(speed int)
	GetSpeed() int
}

func Move(animal Animal) {
	animal.Move()
}
