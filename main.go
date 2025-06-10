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
	count   int
}

func newPool(input chan string) *WorkerPool {
	return &WorkerPool{
		workers: make(map[int]context.CancelFunc),
		input:   input,
	}
}

func (wp *WorkerPool) addWork() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	id := wp.count
	wp.count++

	wp.workers[id] = cancel

	go func(workerID int, ctx context.Context) {
		fmt.Printf("Worker %d ЗАПУСТИЛСЯ\n", workerID)
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d ОСТАНОВИЛСЯ\n", workerID)
				return
			case data := <-wp.input:
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d прерван до обработки\n", workerID)
					return
				default:
					fmt.Printf("Worker %d обработал: %s\n", workerID, data)
				}
			}
		}
	}(id, ctx)
}

func (wp *WorkerPool) removeWork() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.workers) == 0 {
		fmt.Println("нет воркеров для удаления")
		return
	}

	var lastID int
	for v := range wp.workers {
		if v > lastID {
			lastID = v
		}
	}

	wp.workers[lastID]()
	delete(wp.workers, lastID)
}

func main() {

	input := make(chan string)
	pool := newPool(input)
	var numWorks int
	fmt.Println("Введите количество необходимых воркеров")
	fmt.Scan(&numWorks)
	var tasks int
	fmt.Println("Введите количество задач: ")
	fmt.Scan(&tasks)
	i := 0
	for i < numWorks {
		pool.addWork()
		i++
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < tasks; i++ {
			input <- fmt.Sprintf("задача-%d", i)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(2 * time.Second)

	pool.removeWork()

	//time.Sleep(2 * time.Second)
	//pool.removeWork()

	//pool.removeWork()

	wg.Wait()

	time.Sleep(2 * time.Second)
}
