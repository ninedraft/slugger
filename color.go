package main

import (
	"bytes"
	"fmt"
)

func ColoredText(text string, codes ...ASCIICode) string {
	var buf = bytes.NewBufferString(text)
	buf.Reset()
	buf.WriteString("\x1b[")
	for i, code := range codes {
		fmt.Fprintf(buf, "%d", code)
		if i != len(codes)-1 {
			buf.WriteString(";")
		}
	}
	buf.WriteString("m")
	buf.WriteString(text)
	buf.WriteString("\x1b[0m")
	return buf.String()
}

type ASCIICode int

const (
	// All attributes off
	ResetAll     ASCIICode = iota
	Bold                   // Bold text
	Underscore             // Underscore (on monochrome display adapter only)
	Blink                  // Blinking text
	ReverseVideo           // Reverse video on
	Concealed              // Concealed text
)
const (
	// Foreground colors
	Black   ASCIICode = 30 + iota // Black text foreground
	Red                           // Red text foreground
	Green                         // Green text foreground
	Yellow                        // Yellow text foreground
	Blue                          // Blue text foreground
	Magenta                       // Magenta text foreground
	Cyan                          // Cyan text foreground
	White                         // White text foreground
)

const (
	// Background colors
	BackInBlack   = 40 + iota // Black background colors
	BackInRed                 // Red background colors
	BackInGreen               // Green background colors
	BackInYellow              // Yellow background colors
	BackInBlue                // Blue background colors
	BackInMagenta             // Magenta background colors
	BackInCyan                // Cyan background colors
	BackInWhite               // White background colors
)
