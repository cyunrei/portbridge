# Portbridge ![GitHub Release](https://img.shields.io/github/v/release/cyunrei/portbridge) [![Go Report Card](https://goreportcard.com/badge/github.com/cyunrei/portbridge)](https://goreportcard.com/report/github.com/cyunrei/portbridge) ![GitHub License](https://img.shields.io/github/license/cyunrei/portbridge)

Portbridge is a port forwarding tool with cross-platform support.

## Features

- Cross-Platform Support (Linux / Windows / Darwin)
- TCP and UDP Forward Support
- IPv4 and IPv6 Mutual Forward Support
- TCP and UDP Bandwidth Limit Support
- Batch Port Forwarding Rules Support

## Usage

```
Usage:
  portbridge [OPTIONS]

Options:
  -s, --source=          Source address and port to bind locally
  -d, --destination=     Destination address and port to connect remotely
  -p, --protocol=        Source protocol type (e.g., tcp, udp)
  -b, --bandwidth-limit= Bandwidth limit in KiB (default: 0)
      --udp-buffer-size= UDP data forwarding buffer size in bytes (default: 1024)
      --udp-timeout=     UDP data forwarding time out in second (default: 5)
  -f, --rules-file=      Batch port forwarding rules file path
  -g, --gen-rules-file   Generate an example rules file for reference and modification
  -l, --log-file=        Path to the logfile where logs will be written
  -h, --help             Print help message
  -v, --version          Print the version number
```

### Examples:

- Access the Cloudflare DNS (ipv6) via 127.0.0.2:53 with 100 udp buffer size

```shell
portbridge -s 127.0.0.2:53 -d [2606:4700:4700::1111]:53 -p udp --udp-buffer-size=100
```

- Resolve the issue of Terraria not supporting game join via an ipv6 address

```shell
portbridge -s 127.0.0.1:7777 -d [::1]:7777 -p tcp
```

- Expose local TCP port 8080 to 8081 with a bandwidth limit of 1 MiB

```shell
portbridge -s :8081 -d 127.0.0.1:8080 -p tcp -b 1024
```

- Execute the above examples in `rules_example.json`(or `rules_example.yaml`, which in the release files)

```shell
portbridge -f rules_example.json
```