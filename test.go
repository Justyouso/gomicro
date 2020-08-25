package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/time/rate"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//阻塞等待
func waitn() {
	r := rate.NewLimiter(1, 5)
	ctx := context.Background()
	for {
		err := r.WaitN(ctx, 3)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second)
	}
}

func allow() {
	r := rate.NewLimiter(1, 5)
	for {
		if r.AllowN(time.Now(), 2) {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("to many request")
		}
		time.Sleep(time.Second)
	}
}

var r = rate.NewLimiter(1, 5)

func MyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !r.Allow() {
			http.Error(writer, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 {
		time.Sleep(time.Second * 5)
	}
	return Product{
		ID:    101,
		Title: "Golang",
		Price: 10,
	}, nil
}

func resProduct() (Product, error) {
	return Product{
		ID:    100,
		Title: "推荐商品",
		Price: 100,
	}, nil
}

func main() {
	//阻塞等待
	//waitn()
	//allow()
	//限流模式
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	writer.Write([]byte("OK!!!"))
	//})
	//http.ListenAndServe(":8080", MyLimit(mux))
	rand.Seed(time.Now().UnixNano())

	configA := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,
		RequestVolumeThreshold: 4,
		ErrorPercentThreshold:  50,
		SleepWindow:            10,
	}
	hystrix.ConfigureCommand("get_prod", configA)

	//hystrix.ConfigureCommand("get_prod2",hystrix.CommandConfig{
	//	Timeout: 2000,
	//})

	//resultChan := make(chan Product, 1)
	//wg := sync.WaitGroup{}

	////使用协程模拟并发
	//for i := 0; i < 10; i++ {
	//	go (func() {
	//		wg.Add(1)
	//		defer wg.Done()
	//		for {
	//			errs := hystrix.Go("get_prod", func() error {
	//				p, _ := getProduct()
	//				resultChan <- p
	//				return nil
	//			}, func(e error) error {
	//				fmt.Println(e)
	//				rec, err := resProduct()
	//				resultChan <- rec
	//				//fmt.Println(rec)
	//				//fmt.Println(resProduct())
	//				return err
	//			})
	//			select {
	//			case getProd := <-resultChan:
	//				fmt.Println(getProd)
	//			case err := <-errs:
	//				fmt.Println(err)
	//			}
	//			time.Sleep(time.Second * 1)
	//		}
	//	})()
	//}
	//wg.Wait()

	//取出熔断器的状态
	c, _, _ := hystrix.GetCircuit("get_prod")
	for i := 0; i < 100; i++ {
		errs := hystrix.Do("get_prod", func() error {
			p, _ := getProduct()
			fmt.Println(p)
			return nil
		}, func(e error) error {
			//fmt.Println(e)
			rcp, err := resProduct()
			fmt.Println(rcp)
			return err
		})
		if errs != nil {
			//fmt.Println(errs)
		}
		fmt.Println(c.IsOpen())
		time.Sleep(time.Second * 1)
	}

}
