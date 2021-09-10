package logger

import (
    "fmt"
    "os"
    "strings"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var writers map[string]Writer

type Writer func(option map[string]interface{}) (zapcore.WriteSyncer, error)

type Option struct {
    Writer        string                 `json:"writer" yaml:"writer"`
    Console       bool                   `json:"console" yaml:"console"`
    Level         string                 `json:"level" yaml:"level"`
    Format        string                 `json:"format" yaml:"format"` // json text
    ShowLine      bool                   `json:"showLine" yaml:"showLine"`
    WriterOptions map[string]interface{} `json:"writerOptions" yaml:"writerOptions"`

    // Encoder
    MessageKey       string `json:"messageKey" yaml:"messageKey"`
    LevelKey         string `json:"levelKey" yaml:"levelKey"`
    TimeKey          string `json:"timeKey" yaml:"timeKey"`
    NameKey          string `json:"nameKey" yaml:"nameKey"`
    CallerKey        string `json:"callerKey" yaml:"callerKey"`
    FunctionKey      string `json:"functionKey" yaml:"functionKey"`
    StacktraceKey    string `json:"stacktraceKey" yaml:"stacktraceKey"`
    LineEnding       string `json:"lineEnding" yaml:"lineEnding"`
    LevelEncoder     string `json:"levelEncoder" yaml:"levelEncoder"` // capital, capitalColor, color, lower(default)
    EncodeTimeFormat string `json:"encodeTimeFormat" yaml:"encodeTimeFormat"`
    ColorOnlyConsole bool   `json:"colorOnlyConsole" yaml:"colorOnlyConsole"`
}

func init() {
    SetWriter("file", fileWriteSyncer)
}

func SetWriter(key string, f Writer) {
    if writers == nil {
        writers = map[string]Writer{}
    }
    writers[key] = f
}

func defaultOption() *Option {
    return &Option{
        Console:          true,
        Level:            "info",
        MessageKey:       "message",
        LevelKey:         "level",
        TimeKey:          "time",
        NameKey:          "logger",
        CallerKey:        "caller",
        FunctionKey:      "function",
        StacktraceKey:    "stack",
        LineEnding:       zapcore.DefaultLineEnding,
        EncodeTimeFormat: "2006/01/02 15:04:05.000",
    }
}

func getOption(o *Option) *Option {
    if o.MessageKey == "" {
        o.MessageKey = "message"
    }
    if o.LevelKey == "" {
        o.LevelKey = "level"
    }
    if o.TimeKey == "" {
        o.TimeKey = "time"
    }
    if o.NameKey == "" {
        o.NameKey = "logger"
    }
    if o.CallerKey == "" {
        o.CallerKey = "caller"
    }
    if o.FunctionKey == "" {
        o.FunctionKey = "function"
    }
    if o.StacktraceKey == "" {
        o.StacktraceKey = "stack"
    }
    if o.LineEnding == "" {
        o.LineEnding = zapcore.DefaultLineEnding
    }
    if o.EncodeTimeFormat == "" {
        o.EncodeTimeFormat = "2006/01/02 15:04:05.000"
    }
    return o
}

func convertLevelStr(level string) zapcore.Level {
    switch strings.ToLower(level) {
    case "debug":
        return zapcore.DebugLevel
    case "info":
        return zapcore.InfoLevel
    case "warn":
        return zapcore.WarnLevel
    case "error":
        return zapcore.ErrorLevel
    case "dpanic", "d-panic", "d_panic":
        return zapcore.DPanicLevel
    case "panic":
        return zapcore.PanicLevel
    default:
        return zapcore.InfoLevel
    }
}

// New returns a new logger management object.
func New(option ...*Option) *zap.Logger {
    o := defaultOption()
    if len(option) > 0 && option[0] != nil {
        o = getOption(option[0])
    }

    return NewWithOption(o)
}

// NewWithOption returns a new logger management object.
func NewWithOption(option *Option) *zap.Logger {
    o := getOption(option)
    level := convertLevelStr(o.Level)

    return newZapLogger(option, level)
}

func getZapCoreEncoderConfig(option *Option) zapcore.EncoderConfig {
    config := zapcore.EncoderConfig{
        MessageKey:     option.MessageKey,
        LevelKey:       option.LevelKey,
        TimeKey:        option.TimeKey,
        NameKey:        option.NameKey,
        CallerKey:      option.CallerKey,
        LineEnding:     option.LineEnding,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,
    }
    switch option.LevelEncoder {
    case "capital":
        config.EncodeLevel = zapcore.CapitalLevelEncoder
    case "capitalColor", "capitalcolor", "capital-color", "capital_color":
        config.EncodeLevel = zapcore.CapitalColorLevelEncoder
    case "color":
        config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
    default:
        config.EncodeLevel = zapcore.LowercaseLevelEncoder
    }
    if option.EncodeTimeFormat != "" {
        config.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
            encoder.AppendString(time.Format(option.EncodeTimeFormat))
        }
    }

    return config
}

func getZapCoreEncoder(option *Option) zapcore.Encoder {
    if strings.ToLower(option.Format) == "json" {
        return zapcore.NewJSONEncoder(getZapCoreEncoderConfig(option))
    }
    return zapcore.NewConsoleEncoder(getZapCoreEncoderConfig(option))
}

func newZapLogger(option *Option, level zapcore.Level) *zap.Logger {
    var (
        writer zapcore.WriteSyncer
        err    error
        core   zapcore.Core
    )
    if f, ok := writers[strings.ToLower(option.Writer)]; ok {
        writer, err = f(option.WriterOptions)
    } else {
        err = fmt.Errorf("not support writer: %s", option.Writer)
    }
    if err != nil {
        fmt.Println("get WriterSyncer failed", err)
    }
    core = newTee(option, level, writer)
    zapLogger := zap.New(core)
    if level == zapcore.DebugLevel || level == zapcore.ErrorLevel {
        zapLogger = zapLogger.WithOptions(zap.AddStacktrace(level))
    }
    if option.ShowLine {
        zapLogger = zapLogger.WithOptions(zap.AddCaller())
    }

    return zapLogger
}

func newTee(option *Option, level zapcore.Level, writer zapcore.WriteSyncer) zapcore.Core {
    cores := make([]zapcore.Core, 0)
    encoder := getZapCoreEncoder(option)
    if writer == nil || option.Console {
        cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
    }
    if writer != nil {
        if option.ColorOnlyConsole && strings.Contains(strings.ToLower(option.LevelEncoder), "color") {
            option.LevelEncoder = strings.ReplaceAll(strings.ToLower(option.LevelEncoder), "color", "")
            encoder = getZapCoreEncoder(option)
        }
        cores = append(cores, zapcore.NewCore(encoder, writer, level))
    }
    return zapcore.NewTee(cores...)
}
