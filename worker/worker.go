package worker

import (
	"context"
	"fmt"
	"math/rand"
	"order/order"
	"sync"
	"time"
)

func worker(ctx context.Context, orderCh <-chan order.Order, processedCh chan<- order.Order, wg *sync.WaitGroup, stats *order.Stats, workerID int) {

	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Конец работы воркеров с заказами!!")
			return
		case v, ok := <-orderCh:
			if !ok {
				return
			}
			fmt.Println("Обрабатываю заказ номер", v.ID)
			time.Sleep(time.Duration(100+rand.Intn(300)) * time.Millisecond)
		
			select {
			case <-ctx.Done():
				fmt.Println("Конец работы воркеров с заказами!!")
				return
			case processedCh <- v:
				fmt.Println("Добавляю заказ в группу обработанных.")
				stats.Inc()
			}
		}
	}

}

func WorkerPool(ctx context.Context, workerCount int, orderCh <-chan order.Order, stats *order.Stats) <-chan order.Order {
	wg := sync.WaitGroup{}
	processedCh := make(chan order.Order)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(ctx, orderCh, processedCh, &wg, stats, i)
	}

	go func() {
		wg.Wait()
		close(processedCh)
	}()

	return processedCh

}
