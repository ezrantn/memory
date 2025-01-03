package memory

import "testing"

func TestNewMemory(t *testing.T) {
	mem := NewMemory()
	if len(mem.memoryPool) != PoolSize {
		t.Errorf("Expected memory pool size %d, got %d", PoolSize, len(mem.memoryPool))
	}
	if len(mem.freeBlocks) != 1 {
		t.Errorf("Expected 1 free block, got %d", len(mem.freeBlocks))
	}
	if mem.freeBlocks[0].size != PoolSize {
		t.Errorf("Expected free block size %d, got %d", PoolSize, mem.freeBlocks[0].size)
	}
}

func TestAlloc(t *testing.T) {
	mem := NewMemory()

	// Test successful allocation
	addr1, err := mem.Alloc(100)
	if err != nil || addr1 != 0 {
		t.Errorf("Failed to allocate memory: addr=%d, err=%v", addr1, err)
	}

	// Test allocation larger than remaining space
	_, err = mem.Alloc(PoolSize + 1)
	if err == nil {
		t.Error("Expected out of memory error")
	}

	// Test fragmentation handling
	addr2, _ := mem.Alloc(50)
	mem.Free(addr1)
	addr3, _ := mem.Alloc(75)

	if addr3 == addr2 {
		t.Error("New allocation should not overlap with existing allocation")
	}
}

func TestMalloc(t *testing.T) {
	mem := NewMemory()

	// Test successful malloc with zero initialization
	addr, err := mem.Malloc(100)
	if err != nil {
		t.Errorf("Malloc failed: %v", err)
	}

	data, err := mem.Read(addr, 100)
	if err != nil {
		t.Errorf("Failed to read malloc'd memory: %v", err)
	}

	for i, b := range data {
		if b != 0 {
			t.Errorf("Expected zero at index %d, got %d", i, b)
		}
	}
}

func TestFree(t *testing.T) {
	mem := NewMemory()

	// Test basic free
	addr, _ := mem.Alloc(100)
	err := mem.Free(addr)
	if err != nil {
		t.Errorf("Free failed: %v", err)
	}

	// Test double free
	err = mem.Free(addr)
	if err == nil {
		t.Error("Expected error on double free")
	}

	// Test invalid address
	err = mem.Free(PoolSize + 1)
	if err == nil {
		t.Error("Expected error on invalid free address")
	}
}

func TestCoalesceFreeBlocks(t *testing.T) {
	mem := NewMemory()

	// Create three adjacent blocks
	addr1, _ := mem.Alloc(100)
	addr2, _ := mem.Alloc(100)
	addr3, _ := mem.Alloc(100)

	// Free them in random order
	mem.Free(addr2)
	mem.Free(addr1)
	mem.Free(addr3)

	// Should be coalesced into one block
	if len(mem.freeBlocks) != 1 {
		t.Errorf("Expected 1 coalesced block, got %d", len(mem.freeBlocks))
	}

	if mem.freeBlocks[0].size != PoolSize {
		t.Errorf("Expected coalesced block size %d, got %d", PoolSize, mem.freeBlocks[0].size)
	}
}

func TestReadWrite(t *testing.T) {
	mem := NewMemory()

	// Test write and read
	addr, _ := mem.Alloc(100)
	testData := []byte("Hello, World!")

	err := mem.Write(addr, testData)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}

	readData, err := mem.Read(addr, len(testData))
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	for i := range testData {
		if readData[i] != testData[i] {
			t.Errorf("Data mismatch at index %d: expected %d, got %d", i, testData[i], readData[i])
		}
	}

	// Test bounds checking
	err = mem.Write(addr, make([]byte, 101))
	if err == nil {
		t.Error("Expected error on write beyond allocated size")
	}

	_, err = mem.Read(addr, 101)
	if err == nil {
		t.Error("Expected error on read beyond allocated size")
	}
}

func TestConcurrency(t *testing.T) {
	mem := NewMemory()
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			addr, err := mem.Alloc(100)
			if err == nil {
				mem.Write(addr, []byte("test"))
				mem.Read(addr, 4)
				mem.Free(addr)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
