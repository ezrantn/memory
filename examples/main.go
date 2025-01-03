package main

import (
	"fmt"
	"github.com/ezrantn/memory"
	"log"
)

func main() {
	mem := memory.NewMemory()

	// Allocate memory
	addr1, err := mem.Malloc(100)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Allocated address:", addr1)

	// Write data
	data := []byte("Hello World!")
	err = mem.Write(addr1, data)
	if err != nil {
		log.Fatal(err)
	}

	// Read data
	readData, err := mem.Read(addr1, len(data))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Read data:", string(readData))

	// Free memory
	if err = mem.Free(addr1); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Memory free:", addr1)
}
