gen-peer:
	@protoc \
		--proto_path=internal/protobuf "internal/protobuf/peer.proto" \
		--go_out=internal/genproto/peer --go_opt=paths=source_relative \
		--go-grpc_out=internal/genproto/peer --go-grpc_opt=paths=source_relative


gen-auth:
	@protoc \
		--proto_path=internal/protobuf "internal/protobuf/user.proto" \
		--go_out=internal/genproto/user --go_opt=paths=source_relative \
		--go-grpc_out=internal/genproto/user --go-grpc_opt=paths=source_relative
	


gen-resource:
	@protoc \
		--proto_path=internal/protobuf "internal/protobuf/resource.proto" \
		--go_out=internal/genproto/resource --go_opt=paths=source_relative \
		--go-grpc_out=internal/genproto/resource --go-grpc_opt=paths=source_relative