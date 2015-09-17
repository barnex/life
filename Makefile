all:
	goimports -w *.go
	go test
	go build weblife.go
	go build evolution.go
