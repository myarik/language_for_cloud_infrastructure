.PHONY: build clean

BINPATH = bin

build_firstresponse:
	go build -ldflags="-s -w" -o ${BINPATH}/firstresponse firstresponse/main.go

build_slowconsumer:
	go build -ldflags="-s -w" -o ${BINPATH}/slowconsumer slowconsumer/main.go

build_downloader:
	go build -ldflags="-s -w" -o ${BINPATH}/downloader downloader/main.go

build: build_downloader build_slowconsumer build_firstresponse

clean:
	rm -rf ./bin

rebuild: clean build
