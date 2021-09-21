package logger

import (
    "os"

    "go.uber.org/zap/zapcore"
)

func writerConsole(option map[string]interface{}) (zapcore.WriteSyncer, error) {
    if len(option) > 0 {
        if v, ok := option["error"]; ok && v.(bool) {
            return os.Stderr, nil
        }
    }
    return os.Stdout, nil
}
