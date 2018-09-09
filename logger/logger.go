package logger

import (
	"os"

	// This is the logger we use as it provider level logging, while the default golang logger doesn't
	"github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// Log exports the logger for the applciation to use
var Log = logging.MustGetLogger("example")

// Init initialises the logger
func Init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.AddModuleLevel(backendFormatter)
	logging.SetBackend(backendFormatter)
}
