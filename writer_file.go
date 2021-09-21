package logger

import (
    "encoding/json"
    "os"
    "path"

    "github.com/sqmt/helper"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

type writerFileOption struct {
    Path       string `json:"path" yaml:"path"`
    Filename   string `json:"filename" yaml:"filename"`
    MaxSize    int    `json:"maxSize" yaml:"maxSize"` // 如果大小为0，则关闭rotate
    MaxAge     int    `json:"maxAge" yaml:"maxAge"`
    MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
    LocalTime  bool   `json:"localtime" yaml:"localtime"`
    Compress   bool   `json:"compress" yaml:"compress"`
}

func writerFile(option map[string]interface{}) (zapcore.WriteSyncer, error) {
    o := &writerFileOption{
        Path:       os.TempDir(),
        Filename:   "logger.log",
        MaxSize:    0,
        MaxAge:     0,
        MaxBackups: 0,
        LocalTime:  true,
        Compress:   true,
    }
    if option != nil && len(option) > 0 {
        b, _ := json.Marshal(option)
        _ = json.Unmarshal(b, o)
    }
    if o.Path != "" && !helper.FileExists(o.Path) {
        if err := os.MkdirAll(o.Path, os.ModePerm); err != nil {
            return nil, err
        }
    }
    if o.MaxSize > 0 {
        return writerFileRotate(o)
    }
    filename := path.Join(o.Path, o.Filename)
    if f, err := os.Create(filename); err != nil {
        return nil, err
    } else {
        return zapcore.AddSync(f), nil
    }
}

func writerFileRotate(o *writerFileOption) (zapcore.WriteSyncer, error) {
    writer := &lumberjack.Logger{}
    if o.Filename != "" {
        writer.Filename = path.Join(o.Path, o.Filename)
    }
    writer.MaxSize = o.MaxSize
    writer.MaxAge = o.MaxAge
    writer.MaxBackups = o.MaxBackups
    writer.LocalTime = o.LocalTime
    writer.Compress = o.Compress
    return zapcore.AddSync(writer), nil
}
