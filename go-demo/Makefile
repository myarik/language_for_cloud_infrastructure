BINPATH = bin

.PHONY: build_firstresponse
build_firstresponse:
	go build -ldflags="-s -w" -o ${BINPATH}/firstresponse firstresponse/main.go

.PHONY: build_slowconsumer
build_slowconsumer:
	go build -ldflags="-s -w" -o ${BINPATH}/slowconsumer slowconsumer/main.go

.PHONY: build_downloader
build_downloader:
	go build -ldflags="-s -w" -o ${BINPATH}/downloader downloader/main.go

# Downloader
.PHONY: docker_build_downloader
docker_build_downloader:
	docker build -f Dockerfile.downloader -t myarik/demo-downloader:golang .

.PHONY: run_downloader
run_downloader:
	docker run --rm myarik/demo-downloader:golang

.PHONY: run_downloader_debug
run_downloader_debug:
	docker run --rm -e DEBUG=True myarik/demo-downloader:golang

# Slow Consumer
.PHONY: docker_build_slow_consumer
docker_build_slow_consumer:
	docker build -f Dockerfile.slowconsumer -t myarik/demo-slowconsumer:golang .

.PHONY: run_slow_consumer
run_slow_consumer:
	docker run --rm myarik/demo-slowconsumer:golang

.PHONY: run_slow_consumer_debug
run_slow_consumer_debug:
	docker run --rm -e DEBUG=True myarik/demo-slowconsumer:golang

# First response
.PHONY: docker_build_first_response
docker_build_first_response:
	docker build -f Dockerfile.firstresponse -t myarik/demo-firstresponse:golang .

.PHONY: run_first_response
run_first_response:
	docker run --rm myarik/demo-firstresponse:golang

.PHONY: run_first_response_debug
run_first_response_debug:
	docker run --rm -e DEBUG=True myarik/demo-firstresponse:golang

.PHONY: build
build: docker_build_downloader docker_build_slow_consumer docker_build_first_response

.PHONY: clean
clean:
	rm -rf ./bin

rebuild: clean build
