package forwarder

type DataForwarder interface {
	Start() error
	forward() error
	reply() error
}
