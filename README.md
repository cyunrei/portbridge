# Portbridge

Portbridge is a port-forwarding tool that supports IPv4, IPv6, as well as both TCP and UDP protocols. It supports
multiple platforms.

# Example

Visit 1.1.1.1 ipv6 DNS port via 127.0.0.2:53

```shell
portbridge -s 127.0.0.2:53 -d [2606:4700:4700::1111]:53 -p udp
```

Solve terraria not supporting game join via ipv6 address issue

```shell
portbridge -s 127.0.0.1:7777 -d [ipv6 addr]:7777 -p tcp
```

Expose local TCP port 8080 to 8081

```shell
portbridge -s :8081 -d 127.0.0.1:8080 -p tcp
```
