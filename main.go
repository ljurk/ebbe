package main

import (
	"github.com/ljurk/ebbe/cmd"
	_ "github.com/ljurk/ebbe/cmd/redis"
)

func main() {
	cmd.Execute()
}
