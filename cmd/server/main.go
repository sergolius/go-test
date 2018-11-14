package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"go-test/mysql"
	pb "go-test/proto"
	"go-test/sql-builder"
)

var columns = []string{
	"id", "name", "email", "mobile_number",
}

type Server struct {
	table string
	db    *sql.DB
}

func (s *Server) SendMessage(ctx context.Context, message *pb.Message) (*empty.Empty, error) {
	var exists bool
	querySelect := sql_builder.GetExistsSQL(s.table, message.Data["id"])
	err := s.db.QueryRow(querySelect).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return &empty.Empty{}, err
	}

	var query string
	if exists {
		query = sql_builder.GetUpdateSQL(s.table, columns, message.Data)
	} else {
		query = sql_builder.GetInsertSQL(s.table, columns, message.Data)
	}

	if _, err := s.db.Exec(query); err != nil {
		return &empty.Empty{}, err
	}

	if exists {
		fmt.Printf("Record with id %q updated \n", message.Data["id"])
	} else {
		fmt.Printf("Record with id %q created \n", message.Data["id"])
	}

	return &empty.Empty{}, nil
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := mysql.Connect(
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.dbname"),
	)
	if err != nil {
		log.Fatal(err)
	}

	initialSQL := sql_builder.GetCreateTableSQL(viper.GetString("mysql.table"))
	if _, err := db.Exec(initialSQL); err != nil {
		log.Fatal(err)
	}

	server := Server{
		table: viper.GetString("mysql.table"),
		db:    db,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterServiceServer(grpcServer, &server)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
