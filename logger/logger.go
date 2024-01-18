package logger

import (
	"log"
	"time"
)

type Logger struct {
	Logger *log.Logger
}

var colorCodes = map[string]string{
	"yellow":   "33",
	"green":    "32",
	"red":      "31",
	"blue":     "34",
	"white":    "37",
	"blueBg":   "4",
	"yellowBg": "3",
	"redBg":    "1",
	"greenBg":  "2",
}

func (l *Logger) Warn(msg string) {
	colorCode := colorCodes["yellow"]
	bgColorCode := colorCodes["yellowBg"]

	now := time.Now()

	formattedTime := now.Format("2006/01/02 15:04:05")

	l.Logger.SetPrefix("\x1b[1m\x1b[48;5;" + bgColorCode + "m\x1b[37m[WARNING]\x1b[0m ")
	l.Logger.Printf("\x1b[%sm%v - %v\x1b[0m\n", colorCode, formattedTime, msg)
}

func (l *Logger) Success(msg string) {
	colorCode := colorCodes["green"]
	bgColorCode := colorCodes["greenBg"]

	now := time.Now()

	formattedTime := now.Format("2006/01/02 15:04:05")

	l.Logger.SetPrefix("\x1b[1m\x1b[48;5;" + bgColorCode + "m\x1b[37m[SUCCESS]\x1b[0m ")
	l.Logger.Printf("\x1b[%sm%v - %v\x1b[0m\n", colorCode, formattedTime, msg)
}

func (l *Logger) Info(msg string) {
	colorCode := colorCodes["blue"]
	bgColorCode := colorCodes["blueBg"]

	now := time.Now()

	formattedTime := now.Format("2006/01/02 15:04:05")

	l.Logger.SetPrefix("\x1b[1m\x1b[48;5;" + bgColorCode + "m\x1b[37m[INFO]\x1b[0m ")
	l.Logger.Printf("\x1b[%sm%v - %v\x1b[0m\n", colorCode, formattedTime, msg)
}

func (l *Logger) Error(msg string) {
	colorCode := colorCodes["red"]
	bgColorCode := colorCodes["redBg"]

	now := time.Now()

	formattedTime := now.Format("2006/01/02 15:04:05")

	l.Logger.SetPrefix("\x1b[1m\x1b[48;5;" + bgColorCode + "m\x1b[37m[ERROR]\x1b[0m ")
	l.Logger.Printf("\x1b[%sm%v - %v\x1b[0m\n", colorCode, formattedTime, msg)
}
