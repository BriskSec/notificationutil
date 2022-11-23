package notificationutil

import "log"

func NotifyAbnormalCondition(message string, err error) {
	if err != nil {
		log.Printf("Error occurred: %s %v \n", message, err)
	} else {
		log.Printf("Error occurred: %s \n", message)
	}
}
