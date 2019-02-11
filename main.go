package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-audio/wav"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	serverAddr := "127.0.0.1:6565"

	var wg sync.WaitGroup
	err := filepath.Walk("./wav",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}
			wg.Add(1)
			fmt.Println(path, info.Size())
			go send(serverAddr, &wg, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()
}

func send(serverAddr string, wg *sync.WaitGroup, filePath string) {
	fmt.Println(filePath)
	defer wg.Done()

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := NewDetectorServiceClient(conn)
	c := context.Background()
	w, err := client.DetectAnswerMachine(c)

	if err != nil {
		log.Fatal(err)
	}

	step := 0

	var iwg sync.WaitGroup
	iwg.Add(1)
	go func() {
		var res DetectionResult
		w.RecvMsg(&res)
		fmt.Printf("result at step %d: %v", step, res)
		iwg.Done()
	}()

	start := time.Now()

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := wav.NewDecoder(f)
	fmt.Println("test")
	i := 0
	for {
		if decoder.EOF() {
			break
		}

		chunk, err := decoder.NextChunk()
		if err != nil {
			break
		}

		data := make([]byte, 128)
		_, err = chunk.Read(data)
		if err != nil {
			break
			// log.Fatalf("here: %v", err)
		}

		s := Sample{Data: data, CallDate: "2019-02-01T12:12:12+0010"}
		if err := w.Send(&s); err != nil {
			log.Printf("cannot send: %v\n", err)
			break
		}

		fmt.Printf("Sent: %d\n", i)
		i++
	}

	w.CloseSend()
	iwg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Send + response: %v\n", elapsed)
}
