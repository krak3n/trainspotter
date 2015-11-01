image:
	docker build --force-rm -t krak3n/trainspotter:latest .

run:
	go run ./trainspotter/main.go

build:
	go build -o $(GOPATH)/bin/trainspotter ./trainspotter

install:
	go install ./trainspotter
