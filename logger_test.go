package logger

import (
    "strings"
    "testing"
)

func TestNew_default(t *testing.T) {
    l, err := New()
    if err != nil {
        t.Errorf("New() error, want: *zap.Logger, got: %v, err: %v", l, err)
    }
    l.Info("test")
}

func TestNew_file(t *testing.T) {
    l, err := New(&Option{
        Output:       []*Output{{Writer: "file", Option: map[string]interface{}{"path": "./testdata", "filename": "file.log"}}, {Writer: "console"}},
        Level:        "error",
        LevelEncoder: "capitalColor",
    })
    if err != nil {
        t.Errorf("New() error, want: *zap.Logger, got: %v, err: %v", l, err)
    }
    l.Error("test")

    l, err = New(&Option{
        Output: []*Output{{Writer: "file", Option: map[string]interface{}{"path": "./testdata", "filename": "file1.log"}, LevelEncoder: "lower"}, {Writer: "console"}},
        Level:  "error",
        Name:   "test",
    })
    if err != nil {
        t.Errorf("New() error, want: *zap.Logger, got: %v, err: %v", l, err)
    }
    l.Error("test")
}

func TestNew_fileRotate(t *testing.T) {
    l, err := New(&Option{
        Output:       []*Output{{Writer: "file", Option: map[string]interface{}{"path": "./testdata", "filename": "file_rotate.log", "maxSize":1}, Level: "error", Format: "json", LevelEncoder: "capital"}, {Writer: "console"}},
        Level:        "debug",
        LevelEncoder: "capitalColor",
    })
    if err != nil {
        t.Errorf("New() error, want: *zap.Logger, got: %v, err: %v", l, err)
    }
    l.Error("test")

    l, err = New(&Option{
        Output: []*Output{{Writer: "file", Option: map[string]interface{}{"path": "./testdata", "filename": "file_rotate2.log","maxSize":1,"maxAge":10}, LevelEncoder: "capital"}},
        Level:  "error",
        Name:   "test",
    })
    if err != nil {
        t.Errorf("New() error, want: *zap.Logger, got: %v, err: %v", l, err)
    }
    msg := strings.Repeat("demo", 204800)
    l.Error(msg)
}
