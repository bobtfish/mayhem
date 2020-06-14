package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

const HEREDOC = string('`')

func main() {
	fmt.Printf("package main\nconst sprite_sheet_base64 = " + HEREDOC)
	dat, err := ioutil.ReadFile("sprite_sheet.png")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", base64.StdEncoding.EncodeToString([]byte(dat)))
	fmt.Printf(HEREDOC + "")
}
