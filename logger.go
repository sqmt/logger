package logger

import (
    "fmt"
    "strings"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    defaultOption = &Option{
        Output:        []*Output{{Writer: "console"}},
        Name:          "",
        Level:         "info",
        Format:        "",
        ShowLine:      false,
        MessageKey:    "message",
        LevelKey:      "level",
        TimeKey:       "time",
        NameKey:       "name",
        CallerKey:     "caller",
        FunctionKey:   "function",
        StacktraceKey: "stacktrace",
        LineEnding:    zapcore.DefaultLineEnding,
        LevelEncoder:  "lower",
        TimeFormat:    "2006/01/02 15:04:05.000",
    }
    writers = map[string]Writer{
        "console": writerConsole,
        "file":    writerFile,
    }
)

type Output struct {
    Writer        string                 `json:"writer" yaml:"writer"`
    Level         string                 `json:"level" yaml:"level"`
    Format        string                 `json:"format" yaml:"format"`
    MessageKey    string                 `json:"messageKey" yaml:"messageKey"`
    LevelKey      string                 `json:"levelKey" yaml:"levelKey"`
    TimeKey       string                 `json:"timeKey" yaml:"timeKey"`
    NameKey       string                 `json:"nameKey" yaml:"nameKey"`
    CallerKey     string                 `json:"callerKey" yaml:"callerKey"`
    FunctionKey   string                 `json:"functionKey" yaml:"functionKey"`
    StacktraceKey string                 `json:"stacktraceKey" yaml:"stacktraceKey"`
    LineEnding    string                 `json:"lineEnding" yaml:"lineEnding"`
    LevelEncoder  string                 `json:"levelEncoder" yaml:"levelEncoder"` // capital, capitalColor, color, lower(default)
    TimeFormat    string                 `json:"timeFormat" yaml:"timeFormat"`
    Option        map[string]interface{} `json:"option" yaml:"option"`
}

type Option struct {
    Output         []*Output `mapstructure:"output" json:"output" yaml:"output"`
    Name           string    `json:"name" yaml:"name"`
    ShowLine       bool      `json:"showLine" yaml:"showLine"`
    ShowStacktrace bool      `json:"ShowStacktrace" yaml:"ShowStacktrace"`
    // 公共配置
    Level string `json:"level" yaml:"level"`
    // Color         bool   `json:"color" yaml:"color"`
    Format        string `json:"format" yaml:"format"`
    MessageKey    string `json:"messageKey" yaml:"messageKey"`
    LevelKey      string `json:"levelKey" yaml:"levelKey"`
    TimeKey       string `json:"timeKey" yaml:"timeKey"`
    NameKey       string `json:"nameKey" yaml:"nameKey"`
    CallerKey     string `json:"callerKey" yaml:"callerKey"`
    FunctionKey   string `json:"functionKey" yaml:"functionKey"`
    StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
    LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
    LevelEncoder  string `json:"levelEncoder" yaml:"levelEncoder"` // capital,capitalColor, lower(default), lowerColor
    TimeFormat    string `json:"timeFormat" yaml:"timeFormat"`
}

type Writer func(option map[string]interface{}) (zapcore.WriteSyncer, error)

// New 返回zap.Logger
func New(opts ...*Option) (*zap.Logger, error) {
    o := defaultOption
    if len(opts) > 0 {
        o = opts[0]
    }
    optionSet(o)
    level := convertLevelStr(o.Level)
    cores := make([]zapcore.Core, 0)
    for _, output := range o.Output {
        if core, l, err := newOutput(o, output); core != nil {
            cores = append(cores, core)
            if l == zapcore.DebugLevel || l == zapcore.ErrorLevel {
                level = l
            }
        } else {
            return nil, err
        }
    }
    zapLogger := zap.New(zapcore.NewTee(cores...))
    if level == zapcore.DebugLevel || level == zapcore.ErrorLevel {
        zapLogger = zapLogger.WithOptions(zap.AddStacktrace(level))
    }
    if o.ShowLine {
        zapLogger = zapLogger.WithOptions(zap.AddCaller())
    }
    return zapLogger.Named(o.Name), nil
}

// newOutput 创建一个zapcore.Core
func newOutput(o *Option, output *Output) (zapcore.Core, zapcore.Level, error) {
    f, ok := writers[output.Writer]
    if !ok {
        return nil, zapcore.InfoLevel, fmt.Errorf("writer %s not support", output.Writer)
    }
    outputOptionSet(o, output)
    writer, err := f(output.Option)
    if err != nil {
        return nil, zapcore.InfoLevel, err
    }
    level := convertLevelStr(output.Level)
    enc := getZapCoreEncoder(output)

    return zapcore.NewCore(enc, writer, level), level, nil
}

// optionSet 设置option默认值
func optionSet(o *Option) {
    if len(o.Output) == 0 {
        o.Output = defaultOption.Output
    }
    if o.Level == "" {
        o.Level = defaultOption.Level
    }
    if o.MessageKey == "" {
        o.MessageKey = defaultOption.MessageKey
    }
    if o.LevelKey == "" {
        o.LevelKey = defaultOption.LevelKey
    }
    if o.TimeKey == "" {
        o.TimeKey = defaultOption.TimeKey
    }
    if o.NameKey == "" {
        o.NameKey = defaultOption.NameKey
    }
    if o.CallerKey == "" {
        o.CallerKey = defaultOption.CallerKey
    }
    if o.FunctionKey == "" {
        o.FunctionKey = defaultOption.FunctionKey
    }
    if o.StacktraceKey == "" {
        o.StacktraceKey = defaultOption.StacktraceKey
    }
    if o.LineEnding == "" {
        o.LineEnding = defaultOption.LineEnding
    }
    if o.LevelEncoder == "" {
        o.LevelEncoder = defaultOption.LevelEncoder
    }
    if o.TimeFormat == "" {
        o.TimeFormat = defaultOption.TimeFormat
    }
}

// outputOptionSet 输出配置
func outputOptionSet(o *Option, opt *Output) {
    if opt.Level == "" {
        opt.Level = o.Level
    }
    if opt.Format == "" {
        opt.Format = o.Format
    }
    if opt.MessageKey == "" {
        opt.MessageKey = o.MessageKey
    }
    if opt.LevelKey == "" {
        opt.LevelKey = o.LevelKey
    }
    if opt.TimeKey == "" {
        opt.TimeKey = o.TimeKey
    }
    if opt.NameKey == "" {
        opt.NameKey = o.NameKey
    }
    if opt.CallerKey == "" {
        opt.CallerKey = o.CallerKey
    }
    if opt.FunctionKey == "" {
        opt.FunctionKey = o.FunctionKey
    }
    if opt.StacktraceKey == "" {
        opt.StacktraceKey = o.StacktraceKey
    }
    if opt.LineEnding == "" {
        opt.LineEnding = o.LineEnding
    }
    if opt.LevelEncoder == "" {
        opt.LevelEncoder = o.LevelEncoder
    }
    if opt.TimeFormat == "" {
        opt.TimeFormat = o.TimeFormat
    }
}

// convertLevelStr 转换level值
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
    case "dpanic":
        return zapcore.DPanicLevel
    case "panic":
        return zapcore.PanicLevel
    default:
        return zapcore.InfoLevel
    }
}

// getZapCoreEncoderConfig 获取zapcore.EncoderConfig配置
func getZapCoreEncoderConfig(o *Output) zapcore.EncoderConfig {
    config := zapcore.EncoderConfig{
        MessageKey:     o.MessageKey,
        LevelKey:       o.LevelKey,
        TimeKey:        o.TimeKey,
        NameKey:        o.NameKey,
        CallerKey:      o.CallerKey,
        LineEnding:     o.LineEnding,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,
    }
    switch o.LevelEncoder {
    case "capital":
        config.EncodeLevel = zapcore.CapitalLevelEncoder
    case "capitalColor":
        config.EncodeLevel = zapcore.CapitalColorLevelEncoder
    case "lowerColor":
        config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
    default:
        config.EncodeLevel = zapcore.LowercaseLevelEncoder
    }
    if o.TimeFormat != "" {
        config.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
            encoder.AppendString(time.Format(o.TimeFormat))
        }
    }

    return config
}

// getZapCoreEncoder 获取 zapcore.Encoder
func getZapCoreEncoder(o *Output) zapcore.Encoder {
    if strings.ToLower(o.Format) == "json" {
        return zapcore.NewJSONEncoder(getZapCoreEncoderConfig(o))
    }
    return zapcore.NewConsoleEncoder(getZapCoreEncoderConfig(o))
}
