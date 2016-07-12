package pkg

import (
	"fmt"
	"os"
)

// Must states that the program stops if the error occures
func Must(err error) {
	if err == nil {
		return
	}

	fmt.Println(err)
	os.Exit(1)
}
