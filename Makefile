all:
	goimports -w *.go
	go test
	go build life.go
	go build evolution.go
