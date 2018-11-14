# go-test
1. Install deps
`go get ./...`
2. Rise MySQL container
`docker-compose up -d`
3. Create Database and set it name in config
4. Start the Server
```
cd cmd/client
go run main.go
```
5. Start the Client(don't forget to specify path to csv file as argument)
```
cd cmd/client
go run main.go ../../data/data.csv
```
6. Generate protobuf: `make api`