package main

import (
	"context"
	"fmt"
	"order/order"
	"order/worker"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	stats := order.Stats{}
	orderSlice := make([]order.Order, 0)
	orderContext, orderCancel := context.WithCancel(context.Background())
	workerContext, workerCancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3 * time.Second)
		orderCancel()

	}()

	go func() {

		time.Sleep(6 * time.Second)
		workerCancel()

	}()

	orderCh := order.OrderPool(orderContext, 1)
	processedCh := worker.WorkerPool(workerContext, 2, orderCh, &stats)

	startTime := time.Now()

	wg.Add(1)
	go func() {

		for v := range processedCh {
			orderSlice = append(orderSlice, v)
		}

		wg.Done()
	}()

	wg.Wait()

	for _, v := range orderSlice {
		fmt.Println("=== Заказ номер ", v.ID, " ===")
		fmt.Println("Айди пользователя: ", v.UserID)
		fmt.Println("Цена: ", v.Amount)
		fmt.Println("======")
	}

	fmt.Println("=== Статистика === ")
	fmt.Println("Обработанных заказов: ", stats.GetStats())
	fmt.Println("Пройденное время: ", time.Since(startTime))
	fmt.Println("======")

}
