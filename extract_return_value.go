package main

import (
	"fmt"
)

func extractReturnValue(e error) (int, error) {
	s := e.Error()

	var i int
	_, err := fmt.Sscanf(s, "exit status %d", &i)

	if err != nil {
		return 0, err
	}
	return i, nil
}
