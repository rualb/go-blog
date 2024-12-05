package utilmd

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdown = goldmark.New(

	goldmark.WithExtensions(
		extension.GFM,
	),
)

func MD2HTML(md string) (string, error) {
	var buf bytes.Buffer
	buf.Grow(len(md))

	if err := markdown.Convert([]byte(md), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
