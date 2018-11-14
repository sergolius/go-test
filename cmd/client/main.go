package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"go-test/csv-reader"
	pb "go-test/proto"
)

func main() {
	flag.Parse()
	filePath := flag.Arg(0)

	r, err := csv_reader.NewCSVReader(filePath) // "../../data/data.csv"
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewServiceClient(conn)

	var index int64
	for {
		index++
		m := make(map[string]string)

		if err := r.Parse(&m); err != nil && err != io.EOF {
			log.Fatal(index, err)
		} else if err == io.EOF {
			break
		}

		if _, err := c.SendMessage(
			context.Background(),
			&pb.Message{Index: index, Data: m},
		); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Record with id %q sent \n", m["id"])
	}
}
