# Portbridge

Portbridge is a user-space port-forwarding tool with cross-platform support.

# Features

- Cross-Platform Support (Linux / Windows / Darwin)
- TCP and UDP Forward Support
- Bandwidth Limit Support for TCP
- Batch Port Forwarding Rules Support

# Usage

Portbridge Options:

```
  -s, --source=             Source address and port to bind locally
  -d, --destination=        Destination address and port to connect remotely
  -p, --protocol=           Specify the source protocol type
  -b, --bandwidth-limit=    TCP Bandwidth limit in KiB (default: 0)
      --udp-buffer-size=    UDP data forwarding buffer size in bytes (default: 1024)
      --udp-timeout-second= UDP data forwarding time out in second (default: 5)
  -f, --rule-file=          Batch port forwarding file path
  -h, --help                Show help message
  -v, --version             Print the version number
```

Example:

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

- Execute the above examples in `rules_example.json`(or `rules_example.yaml`)

```shell
portbridge -f rules_example.json
```