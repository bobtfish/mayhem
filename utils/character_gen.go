package main

import (
	"fmt"
	"io/ioutil"
)

const HEREDOC = string('`')

func main() {
	fmt.Printf("package main\nconst character_yaml = " + HEREDOC)
	dat, err := ioutil.ReadFile("characters.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", dat)
	fmt.Printf(HEREDOC + "")
}
