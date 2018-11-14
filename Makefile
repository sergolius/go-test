api: proto/api.proto
	@protoc -I proto/ \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:proto \
		proto/api.proto

dep:
	@go get -u github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
