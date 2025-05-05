package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Constants for email formatting
const (
	MaxLineLength = 76 // RFC 5322 recommended line length
)

// Constants for character sets
const (
	CharsetUTF8      = "utf-8"
	CharsetISO2022JP = "ISO-2022-JP"
)

// EmailCharset represents the character encoding used for emails
type EmailCharset string

// String returns the string representation of the charset
func (char EmailCharset) String() string {
	return string(char)
}

// MaxRunesPerLine returns the maximum number of characters that should be included
// in a single line for the given charset to avoid exceeding MIME line limits
func (char EmailCharset) MaxRunesPerLine() int {
	switch char {
	case CharsetUTF8:
		return 990 / 6 // Unicode characters can be up to 6 bytes
	case CharsetISO2022JP:
		return 36
	default:
		return 990 / 6 // Default to UTF-8 behavior
	}
}

// ToISO2022JP converts UTF-8 text to ISO-2022-JP encoding
func ToISO2022JP(text string) (string, error) {
	r := transform.NewReader(strings.NewReader(text), japanese.ISO2022JP.NewEncoder())
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToISO2022JPWithFallback converts text to ISO-2022-JP with character replacement
// Characters that can't be converted will be replaced with a question mark
func ToISO2022JPWithFallback(text string) (string, error) {
	body := &bytes.Buffer{}
	w := &RuneWriter{W: transform.NewWriter(body, japanese.ISO2022JP.NewEncoder())}
	_, err := io.Copy(w, strings.NewReader(text))
	if err != nil {
		return "", err
	}
	return body.String(), nil
}

// EncodeBase64 encodes text to base64 with optional charset conversion
func EncodeBase64(text string, charset EmailCharset) (string, error) {
	var processed string
	var err error

	if charset == CharsetISO2022JP {
		processed, err = ToISO2022JP(text)
		if err != nil {
			return "", err
		}
	} else {
		processed = text
	}

	return base64.StdEncoding.EncodeToString([]byte(processed)), nil
}

// EncodeMIMEHeader encodes text as a MIME header value with the specified charset
func EncodeMIMEHeader(text string, charset EmailCharset) (string, error) {
	encoded, err := EncodeBase64(text, charset)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("=?%s?B?%s?=", charset.String(), encoded), nil
}

// SplitTextByCharset splits a string into chunks based on the maximum line length
// for the specified charset to avoid encoding issues
func SplitTextByCharset(text string, charset EmailCharset) []string {
	var result []string
	buffer := &bytes.Buffer{}

	switch charset {
	case CharsetUTF8:
		for k, c := range strings.Split(text, "") {
			buffer.WriteString(c)
			if (k+1)%charset.MaxRunesPerLine() == 0 {
				result = append(result, buffer.String())
				buffer.Reset()
			}
		}
	case CharsetISO2022JP:
		for c := range strings.SplitSeq(text, "") {
			if buffer.Len()+len(c) > charset.MaxRunesPerLine() {
				result = append(result, buffer.String())
				buffer.Reset()
			}
			buffer.WriteString(c)
		}
	}

	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}
	return result
}

// EncodeSubject encodes an email subject line per RFC standards
func EncodeSubject(subject string, charset EmailCharset) (string, error) {
	buffer := &bytes.Buffer{}
	splits := SplitTextByCharset(subject, charset)

	for i, line := range splits {
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString("=?" + charset.String() + "?B?")
		encodedLine, err := EncodeBase64(line, charset)
		if err != nil {
			return "", err
		}
		buffer.WriteString(encodedLine)
		buffer.WriteString("?=")
		if i != len(splits)-1 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String(), nil
}

// WrapBase64Lines wraps base64-encoded content to ensure MIME compliance
// Adds CRLF line breaks at the appropriate positions
func WrapBase64Lines(content string) string {
	buffer := bytes.Buffer{}
	for k, c := range strings.Split(content, "") {
		buffer.WriteString(c)
		if (k+1)%MaxLineLength == 0 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}
