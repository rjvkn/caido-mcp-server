package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunDecode_Base64RawURL(t *testing.T) {
	// A typical JWT segment uses RawURLEncoding (no padding, -_ charset)
	// e.g., "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" decodes to {"alg":"HS256","typ":"JWT"}
	val := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	expected := `{"alg":"HS256","typ":"JWT"}`

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runDecode(&cobra.Command{}, []string{"base64", val})

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strings.TrimSpace(buf.String()) != expected {
		t.Fatalf("expected %q, got %q", expected, strings.TrimSpace(buf.String()))
	}
}
