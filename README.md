# Portbridge

Portbridge is a user-space port-forwarding tool with cross-platform support.

# Features

- Cross-Platform Support (Linux / Windows / Darwin)
- TCP and UDP Forward Support
- Bandwidth Limit Support for TCP
- Batch Port Forwarding Rules Support

# Example

Access the Cloudflare DNS (ipv6) via 127.0.0.2:53 with 100 udp buffer size

```shell
portbridge -s 127.0.0.2:53 -d [2606:4700:4700::1111]:53 -p udp --udp-buffer-size=100
```

Resolve the issue of Terraria not supporting game join via an ipv6 address

```shell
portbridge -s 127.0.0.1:7777 -d [::1]:7777 -p tcp
```

Expose local TCP port 8080 to 8081 with a bandwidth limit of 1 MiB

```shell
portbridge -s :8081 -d 127.0.0.1:8080 -p tcp -b 1024
```

Execute the above examples in `rules_example.json`(or `rules_example.yaml`)

```shell
portbridge -f rules_example.json
```