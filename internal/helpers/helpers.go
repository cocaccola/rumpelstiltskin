package helpers

import (
	"fmt"
	"os"
)

func PanicWithConfigError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing a config value: %v\n", err.Error())
		panic(err)
	}
}
