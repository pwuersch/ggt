package lib

import (
	"fmt"
	"os"
)

func ExitWithError(err error) {
	fmt.Printf("Error: %s", err)
	os.Exit(1)
}
