package infrastructure

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	traceing "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Z *zap.SugaredLogger
)

var (
	Tracer traceing.Tracer
)

func InitTracer(name string) traceing.Tracer {

	Tracer = otel.Tracer(name)
	tp := trace.NewTracerProvider(trace.WithSampler(trace.NeverSample()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New())
	return Tracer

}

func InitLogger(logfile, tracerName string, level int8) {
	writeSyncer := Writer(logfile)
	encoder := Encoder()

	lvl := zapcore.Level(level)
	core := zapcore.NewCore(encoder, writeSyncer, lvl)
	zp := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l := zp.Sugar()
	defer func() {
		if err := l.Sync(); err != nil {
			fmt.Printf("Error syncing log: %v\n", err)
		}
	}()

	Z = l

	if tracerName != "" {
		InitTracer(tracerName)
	}
}

func NewSpan(ctx context.Context, name string) context.Context {
	c, _ := Tracer.Start(ctx, name)
	return c
}

func EndSpan(ctx context.Context) {
	span := traceing.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.End()
	}
}

func TraceID(ctx context.Context) string {
	span := traceing.SpanContextFromContext(ctx)
	if span.IsValid() {
		return span.TraceID().String()
	}
	return ""
}

func SpanID(ctx context.Context) string {
	span := traceing.SpanContextFromContext(ctx)

	if span.IsValid() {
		return span.SpanID().String()
	}
	return ""
}

func Encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.CallerKey = "caller" // 显示调用者（如果需要）
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Writer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    10, // 10M
		MaxBackups: 5,  // 5个
		MaxAge:     30, // 最多30天
		Compress:   false,
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberJackLogger))
}

func Debug(args ...interface{}) {
	Z.Debug(args...)
}

func Info(args ...interface{}) {
	Z.Info(args...)
}

func Warn(args ...interface{}) {
	Z.Warn(args)
}

func Error(args ...interface{}) {
	Z.Error(args...)
}

func Fatal(args ...interface{}) {
	Z.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	Z.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Z.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Z.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Z.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Z.Fatalf(template, args...)
}

func IsDebugEnable() bool {
	return Z.Level() == zapcore.DebugLevel
}
