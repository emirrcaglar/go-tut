package main

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Process()
}

type EmailTask struct {
	Email       string
	Subject     string
	MessageBody string
}

// Way to process the EmailTask
func (t *EmailTask) Process() {
	fmt.Printf("Sending email to: %s\n", t.Email)
	time.Sleep(2 * time.Second)
}

// Image processing task
type ImageProcessingTask struct {
	ImageUrl string
}

func (t *ImageProcessingTask) Process() {
	fmt.Printf("Processing image: %s\n", t.ImageUrl)
	time.Sleep(5 * time.Second)
}

// Worker pool definition
type WorkerPool struct {
	Tasks       []Task         // The list of the tasks
	concurrency int            // Number of workers to run simultaneously
	tasksChan   chan Task      // The channel which tasks will be sent from
	wg          sync.WaitGroup // Wait for the tasks to be completed
}

// Functions to execute the worker pool
func (wp *WorkerPool) worker() {
	for task := range wp.tasksChan {
		task.Process()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run() {
	// Initialize the tasks channel
	wp.tasksChan = make(chan Task, len(wp.Tasks))

	// Start workers
	for i := 0; i < wp.concurrency; i++ {
		go wp.worker()
	}

	// Send tasks to the tasks channel
	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}
	close(wp.tasksChan)

	// Wait for all tasks to finish
	wp.wg.Wait()
}
