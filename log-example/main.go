package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"go.uber.org/zap"
)

func main() {

	//=================================================
	// SLOG
	//=================================================
	log.Print("info log")
	slog.Info("info log")

	// 输出
	// 2024/08/19 17:01:39 info log
	// 2024/08/19 17:01:39 INFO info log

	// 通过JsonHandler，将日志格式化成JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	//=================================================
	// ZAP
	//=================================================
	// 三种快速创建 logger 的方法: zap.NewProduction()，zap.NewDevelopment()，zap.NewExample()
	zap_logger, _ := zap.NewProduction()
	defer zap_logger.Sync() // flushes buffer, if any
	sugar := zap_logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "http://localhost",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", "http://localhost")

	zap_logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "aaa"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
