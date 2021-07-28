package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// Checks length of the argument
func CheckLen(expected int, args []string) error {
	if len(args) < expected {
		mes := fmt.Sprintf(
			"not enough number of arguments: %d given, %d expected",
			len(args),
			expected,
		)
		return errors.New(mes)
	}
	return nil
}

func Log(keyword string, prefix int, message string) {
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
	switch prefix {
	case 1:
		logger.SetPrefix(keyword + " [INFO]: ")
	case 2:
		logger.SetPrefix(keyword + " [ERROR]: ")
	default:
		logger.SetPrefix(keyword + " [WARNING]: ")
	}
}
