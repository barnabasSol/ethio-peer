gen:
	@protoc \
		--proto_path=internal/protobuf "internal/protobuf/peer.proto" \
		--go_out=internal/genproto/peer --go_opt=paths=source_relative \
		--go-grpc_out=internal/genproto/peer --go-grpc_opt=paths=source_relative

