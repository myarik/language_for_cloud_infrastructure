.PHONY: build

# Downloader
docker_downloader:
	docker build -f Dockerfile.downloader -t myarik/demo-downloader:python .

run_downloader:
	docker run --rm myarik/demo-downloader:python

run_downloader_debug:
	docker run --rm -e DEBUG=True myarik/demo-downloader:python

# Slow Consumer
docker_slow_consumer:
	docker build -f Dockerfile.slowconsumer -t myarik/demo-slowconsumer:python .

run_slow_consumer:
	docker run --rm myarik/demo-slowconsumer:python

run_slow_consumer_debug:
	docker run --rm -e DEBUG=True myarik/demo-slowconsumer:python

# First response
docker_first_response:
	docker build -f Dockerfile.firstresponse -t myarik/demo-firstresponse:python .

run_first_response:
	docker run --rm myarik/demo-firstresponse:python

run_first_response_debug:
	docker run --rm -e DEBUG=True myarik/demo-firstresponse:python