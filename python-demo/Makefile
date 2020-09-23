# Downloader
.PHONY: docker_build_downloader
docker_build_downloader:
	docker build -f Dockerfile.downloader -t myarik/demo-downloader:python .

.PHONY: run_downloader
run_downloader:
	docker run --rm myarik/demo-downloader:python

.PHONY: run_downloader_debug
run_downloader_debug:
	docker run --rm -e DEBUG=True myarik/demo-downloader:python

# Slow Consumer
.PHONY: docker_build_slow_consumer
docker_build_slow_consumer:
	docker build -f Dockerfile.slowconsumer -t myarik/demo-slowconsumer:python .

.PHONY: run_slow_consumer
run_slow_consumer:
	docker run --rm myarik/demo-slowconsumer:python

.PHONY: run_slow_consumer_debug
run_slow_consumer_debug:
	docker run --rm -e DEBUG=True myarik/demo-slowconsumer:python

# First response
.PHONY: docker_build_first_response
docker_build_first_response:
	docker build -f Dockerfile.firstresponse -t myarik/demo-firstresponse:python .

.PHONY: run_first_response
run_first_response:
	docker run --rm myarik/demo-firstresponse:python

.PHONY: run_first_response_debug
run_first_response_debug:
	docker run --rm -e DEBUG=True myarik/demo-firstresponse:python