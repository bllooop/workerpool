package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	mu      sync.Mutex
	workers map[int]context.CancelFunc
	input   chan string
	counter int
}

func NewWorkerPool(input chan string) *WorkerPool {
	return &WorkerPool{
		workers: make(map[int]context.CancelFunc),
		input:   input,
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	id := wp.counter
	wp.counter++

	wp.workers[id] = cancel

	go func(workerID int, ctx context.Context) {
		fmt.Printf("Worker %d ЗАПУСТИЛСЯ\n", workerID)
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d ОСТАНОВИЛСЯ\n", workerID)
				return
			case data := <-wp.input:
				fmt.Printf("Worker %d обработал: %s\n", workerID, data)
			}
		}
	}(id, ctx)
}

func (wp *WorkerPool) RemoveWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.workers) == 0 {
		fmt.Println("нет воркеров для удаления")
		return
	}

	var lastID int
	for id := range wp.workers {
		if id > lastID {
			lastID = id
		}
	}

	wp.workers[lastID]()
	delete(wp.workers, lastID)
}

func main() {

	inputChan := make(chan string)
	pool := NewWorkerPool(inputChan)
	var numWorks int
	fmt.Println("Введите количество необходимых воркеров")
	fmt.Scan(&numWorks)
	var taskCount int
	fmt.Println("Введите количество задач: ")
	fmt.Scan(&taskCount)
	i := 0
	for i < numWorks {
		pool.AddWorker()
		i++
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < taskCount; i++ {
			inputChan <- fmt.Sprintf("задача-%d", i)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(2 * time.Second)

	pool.RemoveWorker()
	wg.Wait()

	time.Sleep(2 * time.Second)
}
