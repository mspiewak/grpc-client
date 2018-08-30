package main

import (
	"fmt"
	"log"
	"sync"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	serverAddr := "localhost:6565"

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go send(serverAddr, &wg)
	}

	wg.Wait()
}

func send(serverAddr string, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := NewDetectorServiceClient(conn)
	w, err := client.DetectAnswerMachine(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	s := Signature{Data: "test"}

	if err := w.Send(&s); err != nil {
		log.Fatal(err)
	}

	res, err := w.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	defer conn.Close()
}
