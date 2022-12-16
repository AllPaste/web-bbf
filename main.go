package main

import (
	"fmt"
	"sync"
	"time"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Person")
			return new(Person)
		},
	}
}

func main() {
	initPool()

	go demo1(pool)

	go demo3(pool)

	go demo2(pool)

	for {
	}
}

func demo1(pool *sync.Pool) {
	p := pool.Get().(*Person)

	p.Name = "first"

	time.Sleep(3 * time.Second)

	fmt.Printf("demo1 p.Name = %s\n", p.Name)

	pool.Put(p)
}

func demo2(pool *sync.Pool) {
	p := pool.Get().(*Person)

	p.Name = "second"
	fmt.Printf("demo2 p.Name = %s\n", p.Name)

	pool.Put(p)
}

func demo3(pool *sync.Pool) {
	p := pool.Get().(*Person)

	p.Name = "three"
	fmt.Printf("demo3 p.Name = %s\n", p.Name)

	pool.Put(p)
}
