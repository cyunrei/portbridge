package options

type Options struct {
	SourceAddr      string `short:"s" long:"source" description:"Source address and port to bind locally" required:"true"`
	DestinationAddr string `short:"d" long:"destination" description:"Destination address and port to connect remotely" required:"true"`
	Protocol        string `short:"p" long:"protocol" description:"Specify the source protocol type" required:"true"`
	BandwidthLimit  uint64 `short:"b" long:"bandwidth-limit" description:"TCP Bandwidth limit in KiB" default:"0"`
	UDPBufferSize   uint64 `long:"udp-buffer-size" description:"UDP data forwarding buffer size in bytes" default:"1024"`
	RuleFile        string `short:"f" long:"rule-file" description:"Batch port forwarding file path"`
	Help            bool   `short:"h" long:"help" description:"Show help message"`
	Version         bool   `short:"v" long:"version" description:"Print the version number"`
}
