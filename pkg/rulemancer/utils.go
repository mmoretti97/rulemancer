/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"io"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "0123456789abcdef"

type writer struct {
	io.Writer
	timeFormat string
}

func (w writer) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(w.timeFormat)), b...))
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func sanitizeFacts(facts string) string {
	// Remove any \n characters
	facts = strings.ReplaceAll(facts, "\n", "")
	return strings.TrimSpace(facts)
}

func yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func purple(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func cyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}

func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func isInSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
