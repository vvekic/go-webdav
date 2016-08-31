# go-webdav
A thin and fast WebDAV server written in Go.

## Features

To use these features, set the following environment variables in either Docker run command, or in docker-compose file.
TLS and authentication are not included - I recommend putting go-webdav behind a reverse-proxy for security. See [taskinen/https-portal](https://hub.docker.com/r/taskinen/https-portal).

### WebDAV prefix

Prefix is the URL path prefix to strip from WebDAV resource paths. Empty by default.

`WEBDAV_PREFIX=/prefix`

### In-memory filesystem

If you want to host WebDAV from an in-memory storage (no files are saved to disk):

`WEBDAV_INMEMORY=true`

### Log level

By default, go-webdav will log only server and request errors. You can also log all HTTP requests, or disable logging altogether. 

`WEBDAV_LOGLEVEL=ALL`

`WEBDAV_LOGLEVEL=NONE`

## Volumes

go-webdav uses two volumes, one for the root filesystem and one for writing the log file:

`/webdav/root`

`/webdav/log`

You can mount any of these volumes to a local directory to persist your WebDAV data and/or log file.

## Examples

These examples will setup WebDAV with verbose logging using the local ~/data/webdav/root directory as the root filesystem, and exposing the interface at localhost port 8080. Also, the log volume is mounted at ~/data/webdav/log.

### docker run

`docker run -d -e "WEBDAV_LOGLEVEL=ALL" -p 127.0.0.1:8080:80 -v ~/data/webdav/log:/webdav/log -v ~/data/webdav/root:/webdav/root bishof/go-webdav`

### docker-compose

```
webdav:
  container_name: webdav
  image: bishof/go-webdav
  ports:
    - "8080:80"
  volumes:
    - ~/data/webdav/log:/webdav/log
    - ~/data/webdav/root:/webdav/root
  environment:
    - WEBDAV_LOGLEVEL=ALL
```