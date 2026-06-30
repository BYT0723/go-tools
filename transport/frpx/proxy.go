package frpx

type ProxyConfig struct {
	Name          string
	Type          string   // "tcp", "http", "https"
	LocalAddr     string   // local service address (host:port)
	RemotePort    int      // remote port for tcp
	CustomDomains []string // for http/https
	SubDomain     string
	HTTPUser      string
	HTTPPwd       string
	HostHeaderRw  string
}

func (cfg ProxyConfig) toNewProxy(user string) NewProxy {
	name := cfg.Name
	if user != "" {
		name = user + "." + name
	}
	return NewProxy{
		ProxyName:     name,
		ProxyType:     cfg.Type,
		RemotePort:    cfg.RemotePort,
		CustomDomains: cfg.CustomDomains,
		SubDomain:     cfg.SubDomain,
		HTTPUser:      cfg.HTTPUser,
		HTTPPwd:       cfg.HTTPPwd,
		HostHeaderRw:  cfg.HostHeaderRw,
	}
}
