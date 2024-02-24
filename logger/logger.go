package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	Logger *log.Logger
}

var colorCodes = map[string]string{
	"yellow":    "33",
	"green":     "32",
	"red":       "31",
	"blue":      "34",
	"white":     "37",
	"blueBg":    "44",
	"yellowBg":  "43",
	"redBg":     "41",
	"greenBg":   "42",
	"gray":      "90", // Gray text
	"lightBlue": "94", // Light blue text
}

func (l *Logger) WebsocketConnect(msg string) {
	colorCode := color.New(color.FgBlue).SprintFunc()
	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s][%s]", colorCode("+"), colorCode("websocket"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) WebsocketDisconnect(msg string) {
	colorCode := color.New(color.FgRed).SprintFunc()
	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s][%s]", colorCode("-"), colorCode("websocket"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) WebsocketError(msg string) {
	colorCode := color.New(color.FgHiRed).SprintFunc()
	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s][%s]", colorCode("error"), colorCode("websocket"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) WebsocketInfo(msg string) {
	colorCode := color.New(color.FgBlue).SprintFunc()

	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s][%s]", colorCode("info"), colorCode("websocket"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) Info(msg string) {
	colorCode := color.New(color.FgBlue).SprintFunc()

	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s]", colorCode("info"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) Error(msg string) {
	colorCode := color.New(color.FgRed).SprintFunc()

	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s]", colorCode("error"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) Warn(msg string) {
	colorCode := color.New(color.FgYellow).SprintFunc()
	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s]", colorCode("warn"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}

func (l *Logger) Success(msg string) {
	colorCode := color.New(color.FgGreen).SprintFunc()
	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	prefix := fmt.Sprintf("[%s]", colorCode("success"))
	l.Logger.SetPrefix(prefix)
	l.Logger.Printf("%v - %v\n", " "+formattedTime, msg)
}
