package main

import (
	"context"
	"fmt"
	"go-test/csv-reader"
	"log"

	"google.golang.org/grpc"

	pb "go-test/proto"
)

func main() {
	r, err := csv_reader.NewCSVReader("../../data/data.csv")
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
	for ; index < 10; index++ {
		m := make(map[string]string)

		if err := r.Parse(&m); err != nil {
			log.Fatal(index, err)
		}

		if _, err := c.SendMessage(
			context.Background(),
			&pb.Message{Index: index, Data: m},
		); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v \n", m)
	}
}
