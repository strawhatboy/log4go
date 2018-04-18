// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
	"time"
	"github.com/fatih/color"
)

var stdout io.Writer = os.Stdout

var (
	debugColor func(a ...interface{}) string
	traceColor func(a ...interface{}) string
	infoColor func(a ...interface{}) string
	warnColor func(a ...interface{}) string
	errorColor func(a ...interface{}) string
	criticalColor func(a ...interface{}) string
)

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	format string
	w      chan *LogRecord
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() *ConsoleLogWriter {
	if (debugColor == nil) {
		debugColor = color.New(color.FgCyan).SprintFunc()
		traceColor = color.New(color.FgBlue).SprintFunc()
		infoColor = color.New(color.FgMagenta).SprintFunc()
		warnColor = color.New(color.FgYellow).SprintFunc()
		errorColor = color.New(color.FgRed).SprintFunc()
		criticalColor = color.New(color.FgHiRed).SprintFunc()
	}
	consoleWriter := &ConsoleLogWriter{
		format: "[%T %D] [%C] [%L] (%S) %M",
		w:      make(chan *LogRecord, LogBufferLength),
	}
	go consoleWriter.run(stdout)
	return consoleWriter
}
func (c *ConsoleLogWriter) SetFormat(format string) {
	c.format = format
}
func (c *ConsoleLogWriter) run(out io.Writer) {
	for rec := range c.w {
		outString := FormatLogRecord(c.format, rec)
		switch levelStrings[rec.Level] {
		case "TRAC":
			outString = traceColor(outString)
			break;
		case "DEBG":
			outString = debugColor(outString)
			break;
		case "INFO":
			//outString = infoColor(outString)
			break;
		case "WARN":
			outString = warnColor(outString)
			break;
		case "EROR":
			outString = errorColor(outString)
			break;
		case "CRIT":
			outString = criticalColor(outString)
			break;
		default:
			break;
		}
		fmt.Fprint(out, outString)
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (c *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	c.w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (c *ConsoleLogWriter) Close() {
	close(c.w)
	time.Sleep(50 * time.Millisecond) // Try to give console I/O time to complete
}
