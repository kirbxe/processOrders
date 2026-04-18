package order

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	UserID int
	Amount int
}

func NewOrder(id int, userID int, amount int) Order {

	return Order{
		ID:     id,
		UserID: userID,
		Amount: amount,
	}

}

// Генерирует заказ с уникальным юзер айди каждый 100-300ms
func orderGeneration(ctx context.Context, orderCh chan<- Order, workerID int, wg *sync.WaitGroup, genID *IDGen) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Конец генерации заказов!!")
			return
		default:
			delay := time.Duration(100+rand.Intn(200)) * time.Millisecond
			fmt.Println("Генерирую айди для нового заказа")
			id := genID.NextGen()
			fmt.Println("Генерирую новый заказ")
			time.Sleep(delay)
			newOrder := NewOrder(id, workerID, rand.Intn(100))
			select {
			case <-ctx.Done():
				return
			case orderCh <- newOrder:
			}

		}
	}
}

func OrderPool(ctx context.Context, orderWorkCount int) <-chan Order {
	wg := sync.WaitGroup{}
	idGen := IDGen{}
	orderCh := make(chan Order)

	for i := 1; i <= orderWorkCount; i++ {
		wg.Add(1)
		go orderGeneration(ctx, orderCh, i, &wg, &idGen)
	}

	go func() {
		wg.Wait()
		close(orderCh)
	}()

	return orderCh
}
