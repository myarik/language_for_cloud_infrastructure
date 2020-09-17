.PHONY: build clean

BINPATH = bin

build_firstresponse:
	go build -ldflags="-s -w" -o ${BINPATH}/firstresponse firstresponse/main.go

build_slowconsumer:
	go build -ldflags="-s -w" -o ${BINPATH}/slowconsumer slowconsumer/main.go

build_downloader:
	go build -ldflags="-s -w" -o ${BINPATH}/downloader downloader/main.go

# Downloader
docker_downloader:
	docker build -f Dockerfile.downloader -t myarik/demo-downloader:golang .

run_downloader:
	docker run --rm myarik/demo-downloader:golang

run_downloader_debug:
	docker run --rm -e DEBUG=True myarik/demo-downloader:golang

# Slow Consumer
docker_slow_consumer:
	docker build -f Dockerfile.slowconsumer -t myarik/demo-slowconsumer:golang .

run_slow_consumer:
	docker run --rm myarik/demo-slowconsumer:golang

run_slow_consumer_debug:
	docker run --rm -e DEBUG=True myarik/demo-slowconsumer:golang

# First response
docker_first_response:
	docker build -f Dockerfile.firstresponse -t myarik/demo-firstresponse:golang .

run_first_response:
	docker run --rm myarik/demo-firstresponse:golang

run_first_response_debug:
	docker run --rm -e DEBUG=True myarik/demo-firstresponse:golang

build: docker_downloader docker_slow_consumer docker_first_response

clean:
	rm -rf ./bin

rebuild: clean build
