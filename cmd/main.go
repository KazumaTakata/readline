package main

import (
	"fmt"
	"github.com/KazumaTakata/readline"
)

func echo(input []byte) {
	fmt.Printf("%s", input)
}

func main() {

	readline.Readline(">>", echo)

}
