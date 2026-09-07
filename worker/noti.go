package worker

import (
	"log"
	"raise-child/constants/shared"
	"raise-child/interfaces/repository"
	"raise-child/model/entities"
	"raise-child/util"
	"sync"
)

// import (
// 	"context"
// 	"database/sql"
// 	"log"
// 	"raise-child/constants/shared"
// 	"raise-child/model/entities"
// 	"raise-child/util"
// 	"sync"
// 	"time"
// )

var (
	notiWorker *notificationWorker
	once       sync.Once
)

type notificationWorker struct {
	volunteerNotiJobs chan entities.VolunteerNoti
	leaderNotiJobs    chan entities.LeaderNoti
	volunteerNotiRepo repository.IVolunteerNotiRepository
	leaderNotiRepo    repository.ILeaderNotiRepository
	wg                sync.WaitGroup
	isClosed          bool
	mu                sync.Mutex
	errLogger         *log.Logger
	infoLogger        *log.Logger
}

type INotificationWorker interface {
	EnqueueVolunteerNoti(v entities.VolunteerNoti) bool
	EnqueueLeaderNoti(l entities.LeaderNoti) bool
	Stop()
}

func InitalizeNotificationWorker(
	volunteerNotiRepo repository.IVolunteerNotiRepository,
	leaderNotiRepo repository.ILeaderNotiRepository,
	errLogger *log.Logger,
) INotificationWorker {
	once.Do(func() {
		notiWorker = &notificationWorker{
			volunteerNotiJobs: make(chan entities.VolunteerNoti, 1000),
			errLogger:         errLogger,
			infoLogger:        util.GetLogConfig(shared.INFO_LEVEL),
		}

		for i := 0; i < 3; i++ {
			//go notiWorker.start(i)
		}

	})

	return notiWorker
}

// EnqueueLeaderNoti implements INotificationWorker.
func (n *notificationWorker) EnqueueLeaderNoti(l entities.LeaderNoti) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.isClosed {
		log.Println("Cảnh báo: Worker đã đóng, không thể nhận thêm job")
		return false
	}

	// Gửi job vào channel (không block nếu channel còn chỗ)
	n.leaderNotiJobs <- l
	return true
}

// EnqueueVolunteerNoti implements INotificationWorker.
func (n *notificationWorker) EnqueueVolunteerNoti(v entities.VolunteerNoti) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.isClosed {
		log.Println("Cảnh báo: Worker đã đóng, không thể nhận thêm job")
		return false
	}

	// Gửi job vào channel (không block nếu channel còn chỗ)
	n.volunteerNotiJobs <- v
	return true
}

// Stop implements INotificationWorker.
func (n *notificationWorker) Stop() {
	panic("unimplemented")
}

// func (w *NotificationWorker) start(id int) {
// 	w.wg.Add(1)
// 	defer w.wg.Done()

// 	var batch []entities.Notification
// }

// // Enqueue là method mà tầng Business sẽ gọi
// func (w *NotificationWorker) Enqueue(n entities.Notification) bool {
// 	w.mu.Lock()
// 	defer w.mu.Unlock()

// 	if w.isClosed {
// 		log.Println("Cảnh báo: Worker đã đóng, không thể nhận thêm job")
// 		return false
// 	}

// 	// Gửi job vào channel (không block nếu channel còn chỗ)
// 	w.jobChan <- n
// 	return true
// }

// // Start bắt đầu chạy vòng lặp worker (chạy trong 1 goroutine riêng)
// func (w *NotificationWorker) Start() {
// 	w.wg.Add(1)
// 	defer w.wg.Done()

// 	// ... (Code logic Batch & Ticker như đã thảo luận ở trên)
// }

// func (w *NotificationWorker) start(id int) {
// 	w.wg.Add(1)
// 	defer w.wg.Done()

// 	var batch []entities.Notification
// 	const batchSize = 50

// 	// Thêm một chút ngẫu nhiên để các worker không flush cùng lúc
// 	// Ví dụ: 5s + (0-500ms)
// 	interval := 5*time.Second + time.Duration(id*100)*time.Millisecond
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop()

// 	log.Printf("[Worker %d] Đã khởi động", id)

// 	for {
// 		select {
// 		case job, ok := <-w.jobChan:
// 			if !ok {
// 				if len(batch) > 0 {
// 					w.flushWithRetry(id, batch)
// 				}
// 				return
// 			}
// 			batch = append(batch, job)
// 			if len(batch) >= batchSize {
// 				w.flushWithRetry(id, batch)
// 				batch = nil
// 			}

// 		case <-ticker.C:
// 			if len(batch) > 0 {
// 				w.flushWithRetry(id, batch)
// 				batch = nil
// 			}
// 		}
// 	}
// }

// // Stop thực hiện Graceful Shutdown
// func (w *NotificationWorker) Stop() {
// 	w.mu.Lock()
// 	if w.isClosed {
// 		w.mu.Unlock()
// 		return
// 	}
// 	w.isClosed = true
// 	close(w.jobChan) // Đóng channel để vòng lặp for range trong Start() kết thúc
// 	w.mu.Unlock()

// 	w.wg.Wait() // Đợi worker xử lý nốt batch cuối cùng trong RAM
// 	log.Println("Worker đã dừng an toàn.")
// }

// func (w *NotificationWorker) flushWithRetry(items []entities.Notification) {
// 	// Sử dụng cơ chế Exponential Backoff: thử lại 2s, 4s, 8s...
// 	expBackoff := backoff.NewExponentialBackOff()
// 	expBackoff.MaxElapsedTime = 1 * time.Minute // Tổng thời gian thử lại tối đa

// 	operation := func() error {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		return w.repo.SaveBatch(ctx, items)
// 	}

// 	err := backoff.Retry(operation, expBackoff)
// 	if err != nil {
// 		log.Printf("[Worker] Thất bại vĩnh viễn sau nhiều lần thử: %v", err)
// 		// Tại đây có thể lưu vào file hoặc hệ thống logging tập trung
// 	}
// }
