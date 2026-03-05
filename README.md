# go-blog

`go-blog` is a microservice responsible for serving the public-facing blog content. It is built with Go and the [Echo](https://echo.labstack.com/) framework.

## Features

- **Content Delivery:** Serves published blog posts to users.
- **Markdown Parsing:** Renders markdown into HTML via `github.com/yuin/goldmark`.
- **Safe HTML:** Uses `github.com/microcosm-cc/bluemonday` to ensure rendered content is safe.
- **Database:** Uses GORM with PostgreSQL to retrieve posts.
- **Metrics:** Exposes Prometheus metrics for traffic monitoring.

## Prerequisites

- Go 1.26+
- Python 3.x

## Build and Run

```sh
# Run tests
python Makefile.py test

# Build binary for Linux
python Makefile.py linux
```
