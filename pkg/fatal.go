package pkg

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// Fatal states that the program stops if the error occures
func Fatal(err error) {
	if err == nil {
		return
	}

	logrus.Fatal(err)
	os.Exit(1)
}
