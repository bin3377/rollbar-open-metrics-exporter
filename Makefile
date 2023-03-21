all: clean build-docker-image

clean:
	go clean

build-binary:
	go get -d -v && go build

build-docker-image:
	docker build -t bin3377/rollbar-open-metrics-exporter:latest -f Dockerfile .
