# socketcli

An aggregate event calculation cli

## Deploy

### Docker-Compose

1. `git clone https://github.com/antoniomika/socketcli`
2. `docker-compose -f deploy/docker-compose.yml up`

### Docker

1. `docker run -it --rm antoniomika/socketcli`

## CLI Flags

```text
The socketcli command

Usage:
  socketcli [flags]

Flags:
  -c, --config string                 Config file (default "config.yml")
      --debug                         Enable debugging information
  -h, --help                          help for socketcli
      --log-to-file                   Enable writing log output to file, specified by log-to-file-path
      --log-to-file-compress          Enable compressing log output files
      --log-to-file-max-age int       The maxium number of days to store log output in a file (default 28)
      --log-to-file-max-backups int   The maxium number of rotated logs files to keep (default 3)
      --log-to-file-max-size int      The maximum size of outputed log files in megabytes (default 500)
      --log-to-file-path string       The file to write log output to (default "/tmp/socketcli.log")
      --log-to-stdout                 Enable writing log output to stdout (default true)
      --stop-after                    Enable stopping the program after stop-after-time
      --stop-after-time duration      The duration to stop the program after if stop-after is set (default 1m0s)
      --time-format string            The time format to use for general log messages (default "2006/01/02 - 15:04:05")
  -v, --version                       version for socketcli
  -w, --websocket-address string      The websocket address to connect to (default "wss://chaotic-stream.herokuapp.com")
```
