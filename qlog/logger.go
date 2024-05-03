package qlog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	stg "github.com/correctinho/correct-util-sdk-go/stg"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Keys used for logging context JSON.
const (
	KeyAPIRequestID = "api-request-id"
	KeyDomainName   = "domain-name"
	KeySourceIP     = "source-ip"
	KeyProtocol     = "protocol"
	KeyAPIKey       = "api-key"
	KeyDomainPrefix = "domain-prefix"
	KeyXRequestID   = "x-request-id"
	KeyAccount      = "account"
	KeyService      = "service"
	KeyRequestURI   = "request_uri"
)

// Logger - struct para controle de log
type Logger struct {
	Context interface{}
	Zap     *zap.Logger
}

// NewProduction builds a sensible production Logger that writes InfoLevel and
// above logs to standard error as JSON.
func NewProduction(context interface{}) *Logger {
	cf := zap.NewProductionConfig()
	cf.EncoderConfig.MessageKey = "message"
	log, _ := cf.Build()
	return &Logger{
		Zap:     log,
		Context: context,
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		println(fmt.Sprintf(msg, keysAndValues...))
		return
	}
	nrfs := l.logFromContext(l.Context)
	if len(keysAndValues) > 0 {
		msg = fmt.Sprintf(msg, keysAndValues...)
	}
	l.Zap.Fatal(msg, nrfs...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		println(fmt.Sprintf(msg, keysAndValues...))
		return
	}
	nrfs := l.logFromContext(l.Context)
	if len(keysAndValues) > 0 {
		msg = fmt.Sprintf(msg, keysAndValues...)
	}
	l.Zap.Error(msg, nrfs...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		println(fmt.Sprintf(msg, keysAndValues...))
		return
	}
	nrfs := l.logFromContext(l.Context)
	if len(keysAndValues) > 0 {
		msg = fmt.Sprintf(msg, keysAndValues...)
	}
	l.Zap.Warn(msg, nrfs...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		println(fmt.Sprintf(msg, keysAndValues...))
		return
	}
	nrfs := l.logFromContext(l.Context)
	if len(keysAndValues) > 0 {
		msg = fmt.Sprintf(msg, keysAndValues...)
	}
	l.Zap.Info(msg, nrfs...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		println(fmt.Sprintf(msg, keysAndValues...))
		return
	}
	nrfs := l.logFromContext(l.Context)
	if len(keysAndValues) > 0 {
		msg = fmt.Sprintf(msg, keysAndValues...)
	}
	l.Zap.Debug(msg, nrfs...)
}

// DebugEnabled - Valida modo debug
func (l *Logger) DebugEnabled() bool {
	ce := l.Zap.Check(zap.DebugLevel, "debugging")
	return ce != nil
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() error {
	return l.Zap.Sync()
}

func (l *Logger) logFromContext(ctx interface{}) (fields []zap.Field) {
	if value, ok := ctx.(*gin.Context); ok {
		if uuid := value.GetString("request_id"); !stg.IsEmpty(&uuid) {
			fields = append(fields, zap.String(KeyXRequestID, uuid))
		}
		if service, ok := os.LookupEnv("SERVICE_NAME"); ok {
			fields = append(fields, zap.String(KeyService, service))
		}

	}

	if _, ok := ctx.(*http.Request); ok {
		if service, ok := os.LookupEnv("SERVICE_NAME"); ok {
			fields = append(fields, zap.String(KeyService, service))
		}

	}

	if value, ok := ctx.(*fasthttp.RequestCtx); ok {
		if uuid, ok := value.UserValue("request_id").(string); ok {
			fields = append(fields, zap.String(KeyXRequestID, uuid))
		}
		if service, ok := os.LookupEnv("SERVICE_NAME"); ok {
			fields = append(fields, zap.String(KeyService, service))
		}

		return
	}
	return fields
}

// LoggerExtras - extras keys
type LoggerExtras struct {
	Key    string
	Value  map[string]interface{}
	Filter []string
}

// InfoJSON - print map
func (l *Logger) InfoJSON(msg, jbs string, keys LoggerExtras) {
	nrfs := l.logFromContext(l.Context)
	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		println(msg)
		return
	}
	if !json.Valid([]byte(jbs)) {
		l.Zap.Info(msg, nrfs...)
		return
	}
	if !stg.IsEmpty(&keys.Key) && len(keys.Value) > 0 {
		nrfs = append(nrfs, zap.Any(keys.Key, keys.Value))
	}
	l.Zap.Info(fmt.Sprintf("%s %s", msg, jbs), nrfs...)
}
