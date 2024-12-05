package web

import (
	"embed"
	"io/fs"
)

//go:embed go-blog/blog/assets/*
var fsBlogAssetsFS embed.FS

func MustBlogAssetsFS() fs.FS {
	res, err := fs.Sub(fsBlogAssetsFS, "go-blog/blog/assets")
	if err != nil {
		panic(err)
	}
	return res
}

//go:embed  go-blog/index.html
var fsBlogIndexHTML embed.FS

func MustBlogIndexHTML() string {
	data, err := fsBlogIndexHTML.ReadFile("go-blog/index.html")
	if err != nil {
		panic(err)
	}
	return string(data)
}
