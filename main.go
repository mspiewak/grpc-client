package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	serverAddr := "localhost:6565"

	var wg sync.WaitGroup
	// wg.Add(10)

	start := time.Now()
	send(serverAddr, &wg)

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Total time (connect + send): %v\n", elapsed)

	// for i := 0; i < 10; i++ {
	// 	go send(serverAddr, &wg)
	// }

	// wg.Wait()
}

func send(serverAddr string, wg *sync.WaitGroup) {
	// defer wg.Done()
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := NewDetectorServiceClient(conn)
	w, err := client.DetectAnswerMachine(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	s := Signature{Data: "test"}

	start := time.Now()
	// for i := 0; i < 100; i++ {
	if err := w.Send(&s); err != nil {
		log.Fatal(err)
	}
	// }

	_, err = w.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Send + response: %v\n", elapsed)

	// fmt.Println(res.GetIsAnswerMachine())
}
