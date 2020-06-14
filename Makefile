all: mayhem

characters.go: characters.yaml utils/character_gen.go
	go run utils/character_gen.go > characters.go

mayhem: characters.go *.go */*.go go.mod go.sum
	go build
