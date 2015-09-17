all:
	goimports -w *.go
	go build weblife.go
	go build evolution.go
	go test
