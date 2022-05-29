# miniflux-luma

Atom feed exporter for Miniflux starred items

## Installation

```
go get -u github.com/erdnaxe/miniflux-luma
```

## Usage

Fetch your API token from Miniflux settings and write it in `api_token` file.
This file path can be changed using `-api-token-file` argument.

Then you may start the web service:

```
miniflux-luma -endpoint https://rss.example.com -listen-addr 127.0.0.1:8080
```

