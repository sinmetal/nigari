// +build ignore

// The emoji list from https://github.com/googlefonts/noto-emoji.
// ls noto-emoji/svg/*.svg | xargs basename | sed 's/emoji_u//g' | sed 's/.svg//g' | grep -v "_" > emoji.txt

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

var tmpl = template.Must(template.New("").Parse(`// Code generated by emojilistgen.go; DO NOT EDIT.
package nigari

//go:generate go run emojilistgen.go

var emojilist = map[rune]bool {
	{{range .}}0x{{.}}: true,
{{end}}}
`))

func main() {
	f, err := os.Open("emoji.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer f.Close()

	var emojis []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		emojis = append(emojis, s.Text())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, emojis); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	dst, err := os.Create("emojilist.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}()

	if _, err := io.Copy(dst, bytes.NewReader(src)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
