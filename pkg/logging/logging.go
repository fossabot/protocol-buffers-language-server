// Copyright 2019 The Protocol Buffers Language Server Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/micnncim/protocol-buffers-language-server/pkg/config"
)

// NewLogger returns *zap.Logger initialized by provided log level and []zap.Option.
func NewLogger(cfg config.Log, opts ...zap.Option) (*zap.Logger, error) {
	c, err := newConfig(cfg)
	if err != nil {
		return nil, err
	}
	return c.Build(opts...)
}

func newConfig(cfg config.Log) (zap.Config, error) {
	c := zap.Config{
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     logTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := parseLogLevel(cfg.Level)
	if err != nil {
		return zap.Config{}, err
	}
	c.Level = zap.NewAtomicLevelAt(l)

	if f := cfg.File; f != "" {
		c.OutputPaths = []string{f}
		c.ErrorOutputPaths = []string{f}
	}

	return c, nil
}

func logTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05"))
}

func parseLogLevel(levelStr string) (zapcore.Level, error) {
	switch strings.ToUpper(levelStr) {
	case zapcore.DebugLevel.CapitalString():
		return zapcore.DebugLevel, nil
	case zapcore.InfoLevel.CapitalString():
		return zapcore.InfoLevel, nil
	case zapcore.WarnLevel.CapitalString():
		return zapcore.WarnLevel, nil
	case zapcore.ErrorLevel.CapitalString():
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("undefined log level: %v", levelStr)
	}
}
