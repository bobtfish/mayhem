all: mayhem

sprite_sheet.go: sprite_sheet.png utils/sprite_sheet_gen.go
	go run utils/sprite_sheet_gen.go > sprite_sheet.go

characters.go: characters.yaml utils/character_gen.go
	go run utils/character_gen.go > characters.go

mayhem: sprite_sheet.go characters.go *.go */*.go go.mod go.sum
	go build
