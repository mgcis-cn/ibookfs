// Package worker provides asynchronous image processing.
package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// ImageWorker processes images asynchronously using a worker pool pattern.
type ImageWorker struct {
	queue    chan uint
	service  Processor // Use interface to avoid circular dependency
	workers  int
	wg       sync.WaitGroup
	cancel   context.CancelFunc
	stopChan chan struct{}
	mu       sync.Mutex
	running  bool
}

// Processor defines the interface for processing images.
// This avoids circular dependency with service package.
type Processor interface {
	ProcessImage(ctx context.Context, imageID uint) error
}

// Config holds worker configuration.
type Config struct {
	// Concurrent is the number of concurrent workers.
	Concurrent int

	// QueueSize is the maximum queue size.
	QueueSize int

	// RetryTimes is the number of retry attempts on failure.
	RetryTimes int

	// RetryInterval is the delay between retries.
	RetryInterval time.Duration
}

// DefaultConfig returns the default worker configuration.
func DefaultConfig() Config {
	return Config{
		Concurrent:    3,
		QueueSize:     1000,
		RetryTimes:    3,
		RetryInterval: 60 * time.Second,
	}
}

// NewImageWorker creates a new image processing worker.
func NewImageWorker(svc Processor, cfg Config) *ImageWorker {
	return &ImageWorker{
		queue:    make(chan uint, cfg.QueueSize),
		service:  svc,
		workers:  cfg.Concurrent,
		stopChan: make(chan struct{}),
	}
}

// SetService sets the processor service. This allows lazy initialization.
func (w *ImageWorker) SetService(svc any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if proc, ok := svc.(Processor); ok {
		w.service = proc
	}
}

// Enqueue adds an image ID to the processing queue.
func (w *ImageWorker) Enqueue(imageID uint) (err error) {
	// Recover from panic if channel is closed
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("worker is not running")
		}
	}()

	w.mu.Lock()
	running := w.running
	w.mu.Unlock()

	if !running {
		return fmt.Errorf("worker is not running")
	}

	select {
	case w.queue <- imageID:
		return nil
	default:
		return fmt.Errorf("worker queue is full")
	}
}

// Start begins processing images from the queue.
func (w *ImageWorker) Start() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return
	}

	w.running = true
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	// Start worker goroutines
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.worker(ctx, i)
	}

	log.Printf("Image worker started with %d workers", w.workers)
}

// Stop gracefully shuts down the worker.
func (w *ImageWorker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return
	}

	log.Println("Stopping image worker...")

	// Cancel context
	if w.cancel != nil {
		w.cancel()
	}

	// Close stop channel
	close(w.stopChan)

	// Wait for workers to finish
	w.wg.Wait()

	// Drain queue
	close(w.queue)
	for range w.queue {
	}

	w.running = false
	log.Println("Image worker stopped")
}

// IsRunning returns whether the worker is currently running.
func (w *ImageWorker) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// worker processes images from the queue.
func (w *ImageWorker) worker(ctx context.Context, id int) {
	defer w.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d received stop signal", id)
			return

		case <-w.stopChan:
			log.Printf("Worker %d stopping via stop channel", id)
			return

		case imageID, ok := <-w.queue:
			if !ok {
				log.Printf("Worker %d: queue closed, exiting", id)
				return
			}

			w.processImage(ctx, imageID, id)
		}
	}
}

// processImage processes a single image with retry logic.
func (w *ImageWorker) processImage(ctx context.Context, imageID uint, workerID int) {
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			log.Printf("Worker %d: retrying image %d (attempt %d/%d)", workerID, imageID, attempt+1, maxRetries)
			select {
			case <-time.After(60 * time.Second):
			case <-ctx.Done():
				return
			}
		}

		log.Printf("Worker %d: processing image %d", workerID, imageID)

		if err := w.service.ProcessImage(ctx, imageID); err != nil {
			lastErr = err
			log.Printf("Worker %d: failed to process image %d: %v", workerID, imageID, err)
			continue
		}

		log.Printf("Worker %d: successfully processed image %d", workerID, imageID)
		return
	}

	// All retries exhausted
	log.Printf("Worker %d: gave up on image %d after %d attempts: %v", workerID, imageID, maxRetries, lastErr)
}

// QueueSize returns the current queue size.
func (w *ImageWorker) QueueSize() int {
	return len(w.queue)
}

// QueueCapacity returns the maximum queue capacity.
func (w *ImageWorker) QueueCapacity() int {
	return cap(w.queue)
}
