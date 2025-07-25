BUILD_DIR ?= build
PROG_NAME ?= bookish-guacamole

default: dev

build:
	@mkdir -p build
	go build -o ${BUILD_DIR}/${PROG_NAME} cmd/main.go

build_test:
	@mkdir -p build
	go build -o ${BUILD_DIR}/${PROG_NAME}_test cmd/test/main.go

run:
	${BUILD_DIR}/${PROG_NAME}

dev:
	go run cmd/main.go

test:
	go run cmd/test/main.go

auto_test:
	go test ./...

build_image: build
	docker build -t ghcr.io/qcrg/bookish-guacamole .

clean:
	${RM} -rf ${BUILD_DIR}

.PHONY: clean build run dev test build_image
