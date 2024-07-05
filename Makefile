include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	mkdir -p pkg/swagger
	make generate-note-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include="*.css,*.html,*.js,*.json,*.png"
	make generate-access-api
	make generate-auth-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 --proto_path vendor.protogen \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/note_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--grpc-gateway_out=pkg/note_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/note_v1/note.proto

generate-other-note-api:
	mkdir -p pkg/other_note_v1
	protoc --proto_path api/other_note_v1 --proto_path vendor.protogen \
	--go_out=pkg/other_note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/other_note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/other_note_v1/other_note.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 \
		--go_out=pkg/access_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		api/access_v1/access.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
    	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
    	--plugin=protoc-gen-go=bin/protoc-gen-go \
    	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
    	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
    	api/auth_v1/auth.proto


local-migrations-status:
	${LOCAL_BIN}/goose -dir $(LOCAL_MIGRATION_DIR) postgres ${PG_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} down -v


test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=di_container/internal/service/...,di_container/internal/api... -count 5


test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=di_container/internal/service/...,di_container/internal/api... -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

gen-cert:
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
    		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext

grpc-load-test:
	ghz \
		--proto api/note_v1/note.proto \
		--call note_v1.NoteV1.Get \
		--data '{"id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/note_v1/note.proto \
		--call note_v1.NoteV1.Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051