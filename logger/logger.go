package logger

import "log"

const (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Purple = "\033[35m"
	Reset  = "\033[0m"
)

const (
	INFO    string = Green + "INFO" + Reset
	WARNING        = Yellow + "WARNING" + Reset
	ERROR          = Red + "ERROR" + Reset
	FATAL          = Purple + "FATAL" + Reset
)

func PrintInfo(msg string) {
	log.Printf("[%s]: %s\n", INFO, msg)
}

func PrintWarning(msg string) {
	log.Printf("[%s]: %s\n", WARNING, msg)
}

func PrintError(msg string) {
	log.Printf("[%s]: %s\n", ERROR, msg)
}

func PrintFatal(msg string) {
	log.Fatalf("[%s]: %s\n", FATAL, msg)
}
