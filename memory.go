package memory

import (
	"errors"
	"sync"
)

const PoolSize = 1024 * 1024 // 1 MB memory pool

type Memory struct {
	memoryPool []byte
	freeBlocks []block
	allocated  map[int]int // Maps start address to block size
	mu         sync.RWMutex
}

type block struct {
	start int
	size  int
}

// NewMemory: Initialize the memory
func NewMemory() *Memory {
	return &Memory{
		memoryPool: make([]byte, PoolSize),
		freeBlocks: []block{{start: 0, size: PoolSize}},
		allocated:  make(map[int]int),
	}
}

// Simulate alloc: Allocate memory without initializing
func (mm *Memory) Alloc(size int) (int, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	for i, freeBlock := range mm.freeBlocks {
		if freeBlock.size >= size {
			start := freeBlock.start
			mm.allocated[start] = size

			// Update the free block
			if freeBlock.size > size {
				mm.freeBlocks[i] = block{
					start: start + size,
					size:  freeBlock.size - size,
				}
			} else {
				mm.freeBlocks = append(mm.freeBlocks[:i], mm.freeBlocks[i+1:]...)
			}

			return start, nil
		}
	}

	return -1, errors.New("out of memory")
}

// Simulate malloc: Allocate and zero-initialize memory
func (mm *Memory) Malloc(size int) (int, error) {
	addr, err := mm.Alloc(size)
	if err != nil {
		return -1, err
	}

	// Zero initialize the memory
	for i := 0; i < size; i++ {
		mm.memoryPool[addr+i] = 0
	}

	return addr, nil
}

// Simulate free: Deallocate memory
func (mm *Memory) Free(addr int) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	size, exists := mm.allocated[addr]
	if !exists {
		return errors.New("invalid free")
	}

	delete(mm.allocated, addr)
	mm.freeBlocks = append(mm.freeBlocks, block{start: addr, size: size})
	mm.coalesceFreeBlock()
	return nil
}

// Merge adjacent free blocks to reduce fragmentation
func (mm *Memory) coalesceFreeBlock() {

}
