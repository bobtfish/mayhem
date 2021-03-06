.PHONY: all
all: mayhem

sprite_sheet.go: sprite_sheet.png utils/sprite_sheet_gen/main.go
	go run utils/sprite_sheet_gen/main.go > sprite_sheet.go

character/yaml.go: character/characters.yaml utils/character_gen/main.go
	go run utils/character_gen/main.go > character/yaml.go

mayhem: sprite_sheet.go character/yaml.go *.go */*.go */*/*.go go.mod go.sum
	go build

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm sprite_sheet.go character/yaml.go
