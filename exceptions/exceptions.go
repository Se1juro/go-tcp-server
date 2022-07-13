package exceptions

import "log"

func ManageError(err error, message string) {
	if err != nil {
		log.Println(message)
	}
}
