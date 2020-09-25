# Downloader
.PHONY: docker_build_downloader
docker_build_downloader:
	docker build -f Dockerfile.downloader -t myarik/demo-downloader:node .

.PHONY: run_downloader
run_downloader:
	docker run --rm myarik/demo-downloader:node

# Slow Consumer
.PHONY: docker_build_slow_consumer
docker_build_slow_consumer:
	docker build -f Dockerfile.slowconsumer -t myarik/demo-slowconsumer:node .

.PHONY: run_slow_consumer
run_slow_consumer:
	docker run --rm myarik/demo-slowconsumer:node

# First response
.PHONY: docker_build_first_response
docker_build_first_response:
	docker build -f Dockerfile.firstresponse -t myarik/demo-firstresponse:node .

.PHONY: run_first_response
run_first_response:
	docker run --rm myarik/demo-firstresponse:node