.PHONY: build

build:
	go build -o bin/login.exe cmd/login/main.go
	go build -o bin/cluster.exe cmd/cluster/main.go
	go build -o bin/world.exe cmd/world/main.go
