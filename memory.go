package memory

import (
	"errors"
	"sort"
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

	// Find best-fit block
	bestFitIdx := -1
	bestFitSize := PoolSize + 1

	for i, block := range mm.freeBlocks {
		if block.size >= size && block.size < bestFitSize {
			bestFitIdx = i
			bestFitSize = block.size
		}
	}

	if bestFitIdx == -1 {
		return -1, errors.New("out of memory")
	}

	freeBlock := mm.freeBlocks[bestFitIdx]
	start := freeBlock.start
	mm.allocated[start] = size

	// Update or remove the free block
	if freeBlock.size > size {
		mm.freeBlocks[bestFitIdx] = block{
			start: start + size,
			size:  freeBlock.size - size,
		}
	} else {
		mm.freeBlocks = append(mm.freeBlocks[:bestFitIdx], mm.freeBlocks[bestFitIdx+1:]...)
	}

	return start, nil
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
	if len(mm.freeBlocks) <= 1 {
		return
	}

	// Sort blocks by start address
	sort.Slice(mm.freeBlocks, func(i, j int) bool {
		return mm.freeBlocks[i].start < mm.freeBlocks[j].start
	})

	// Merge adjacent blocks
	newBlocks := []block{mm.freeBlocks[0]}
	for i := 1; i < len(mm.freeBlocks); i++ {
		curr := mm.freeBlocks[i]
		prev := &newBlocks[len(newBlocks)-1]

		if prev.start+prev.size == curr.start {
			// Merge blocks
			prev.size += curr.size
		} else {
			newBlocks = append(newBlocks, curr)
		}
	}

	mm.freeBlocks = newBlocks
}

func (mm *Memory) Read(addr, size int) ([]byte, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	allocSize, exists := mm.allocated[addr]
	if !exists {
		return nil, errors.New("reading from unallocated memory")
	}

	if size > allocSize {
		return nil, errors.New("reading beyond allocated memory")
	}

	result := make([]byte, size)
	copy(result, mm.memoryPool[addr:addr+size])
	return result, nil
}

func (mm *Memory) Write(addr int, data []byte) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	allocSize, exists := mm.allocated[addr]
	if !exists {
		return errors.New("writing to unallocated memory")
	}

	if len(data) > allocSize {
		return errors.New("writing beyond allocated memory")
	}

	copy(mm.memoryPool[addr:addr+len(data)], data)
	return nil
}
