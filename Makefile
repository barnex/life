all:
	goimports -w *.go
	go test
