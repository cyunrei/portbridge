package options

type Options struct {
	SourceAddr      string `short:"s" long:"source" description:"Source address and port to bind locally" required:"true"`
	DestinationAddr string `short:"d" long:"destination" description:"Destination address and port to connect remotely" required:"true"`
	Protocol        string `short:"p" long:"protocol" description:"Source protocol type (e.g., tcp, udp)" required:"true"`
	BandwidthLimit  uint64 `short:"b" long:"bandwidth-limit" description:"Bandwidth limit in KiB" default:"0"`
	UDPBufferSize   uint64 `long:"udp-buffer-size" description:"UDP data forwarding buffer size in bytes" default:"1024"`
	UDPTimeout      uint64 `long:"udp-timeout" description:"UDP data forwarding time out in second" default:"5"`
	RulesFile       string `short:"f" long:"rules-file" description:"Batch port forwarding rules file path"`
	Help            bool   `short:"h" long:"help" description:"Print help message"`
	Version         bool   `short:"v" long:"version" description:"Print the version number"`
}
