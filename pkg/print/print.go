package print

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
	Non    = ""
)

type Tag string

const (
	DB              Tag = "DB   💾"
	Config          Tag = "Conf 🖊️ "
	Pet             Tag = "Pet  🐶"
	Synchronization Tag = "Sync 🛜"
	System          Tag = "Sys  ⚙️ "
	File            Tag = "File 📁"
	App             Tag = "App  📦"
)

var (
	infoLogger   = log.New(os.Stdout, "", 0)
	warnLogger   = log.New(os.Stdout, "", 0)
	assertLogger = log.New(os.Stdout, "", 0)
	fatalLogger  = log.New(os.Stderr, "", 0)
)

func formatMessage(color, level, tag Tag, format string, args ...any) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s[%s] [%s] [%s] %s%s", color, timestamp, level, string(tag), msg, Reset)
}

func Info(tag Tag, format string, args ...any) {
	infoLogger.Println(formatMessage(Non, "📄   INFO", tag, format, args...))
}

func Warn(tag Tag, format string, args ...any) {
	warnLogger.Println(formatMessage(Yellow, "⚠️    WARN", tag, format, args...))
}

func Assert(tag Tag, format string, args ...any) {
	assertLogger.Println(formatMessage(Green, "✅ ASSERT", tag, format, args...))
}

func Fatal(tag Tag, format string, args ...any) {
	fatalLogger.Println(formatMessage(Red, "❌  FATAL", tag, format, args...))
	os.Exit(1)
}
