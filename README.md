> [!NOTE]
> This is a high-level simulation and does not directly manage OS-level memory

# memory

This project is an experimental simulation of memory management in Go, inspired by the behavior of `alloc`, `malloc`, and `free` functions in C. It provides a simple framework for allocating, initializing, and freeing memory in a controlled environment.

The project is designed to explore how memory management could work in a language like Go, which has its own garbage collection and memory management system.

## Installation

```shell
go get github.com/ezrantn/memory
```

## Usage

### Basic Memory Operation

```go
// Initialize memory manager
mem := memory.NewMemory()

// Allocate memory
addr1, err := mem.Malloc(100)
if err != nil {
    log.Fatal(err)
}

// Write data
data := []byte("Hello, World!")
err = mem.Write(addr1, data)
if err != nil {
    log.Fatal(err)
}

// Read data
readData, err := mem.Read(addr1, len(data))
if err != nil {
    log.Fatal(err)
}

// Free memory
err = mem.Free(addr1)
if err != nil {
    log.Fatal(err)
}
```

### Memory Fragmentation Example

```go
mem := memory.NewMemory()

// Allocate multiple blocks
addr1, _ := mem.Malloc(100)
addr2, _ := mem.Malloc(200)
addr3, _ := mem.Malloc(300)

// Free non-contiguous blocks
mem.Free(addr1)
mem.Free(addr3)
// Memory automatically coalesces free blocks when possible
```

See `examples` directory for details.

## Limitations

1. Fixed pool size
2. No support for resizing allocations
3. Simple implementation of memory coalescing
4. No support for alignment requirements
5. No garbage collection interaction

## Performance Considerations

- Best-fit allocation: O(n) where n is the number of free blocks
- Coalescing: O(n log n) due to sorting
- Memory fragmentation can occur with repeated allocations/deallocations

## Thread Safety

All operations are thread-safe through the use of a RWMutex. Read operations use shared locks while write operations use exclusive locks.

## Testing
 
Run the tests suite:

```shell
make test
```

Run tests with code coverage:

```shell
make cov
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request. For major changes, please open an issue first to discuss what you would like to change.