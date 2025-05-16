package v1alpha1

type UpstreamServer struct {
	HostHeader string `json:"hostHeader,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
	Weight     int    `json:"weight"`
	TlsVerify  bool   `json:"tlsVerify,omitempty"`
}

type Upstream struct {
	Name       string           `json:"name"`
	HostHeader string           `json:"hostHeader,omitempty"`
	Servers    []UpstreamServer `json:"servers,omitempty"`
}
