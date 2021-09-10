package logger

import (
    "encoding/json"
    "os"
    "path"

    "github.com/sqmt/helper"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

type fileWriterOption struct {
    Path       string `json:"path" yaml:"path"`
    Filename   string `json:"filename" yaml:"filename"`
    MaxSize    int    `json:"maxsize" yaml:"maxsize"`
    MaxAge     int    `json:"maxage" yaml:"maxage"`
    MaxBackups int    `json:"maxbackups" yaml:"maxbackups"`
    LocalTime  bool   `json:"localtime" yaml:"localtime"`
    Compress   bool   `json:"compress" yaml:"compress"`
}

func fileWriteSyncer(option map[string]interface{}) (zapcore.WriteSyncer, error) {
    if option == nil {
        option = map[string]interface{}{}
    }

    writer := &lumberjack.Logger{}
    if option != nil {
        b, _ := json.Marshal(option)
        var o fileWriterOption
        _ = json.Unmarshal(b, &o)
        if o.Path != "" && !helper.FileExists(o.Path) {
            os.MkdirAll(o.Path, os.ModePerm)
        }
        if o.Filename != "" {
            writer.Filename = path.Join(o.Path, o.Filename)
        }
        writer.MaxSize = o.MaxSize
        writer.MaxAge = o.MaxAge
        writer.MaxBackups = o.MaxBackups
        writer.LocalTime = o.LocalTime
        writer.Compress = o.Compress
    }

    return zapcore.AddSync(writer), nil
}
