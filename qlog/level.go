// Package qlog defines a custom error level type and a set of constants representing different levels of errors.
package qlog

// LevelError is a custom type used to represent different levels of errors.
type LevelError string

// The following constants represent different levels of errors.
// Each constant is of type LevelError and is assigned a unique string value.
const (
	DebugLevel  LevelError = "debug"
	InfoLevel   LevelError = "info"
	WarnLevel   LevelError = "warn"
	ErrorLevel  LevelError = "error"
	DPanicLevel LevelError = "dpanic"
	PanicLevel  LevelError = "panic"
	FatalLevel  LevelError = "fatal"
)
