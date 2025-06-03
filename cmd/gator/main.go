package main

import (
	"fmt"
	"os"

	"github.com/luism2302/gator/internal/config"
)

func main() {
	test, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(test)
	err = test.SetUser("luis")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	test, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(test)
}
