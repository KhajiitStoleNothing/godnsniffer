package config

type WebConfig struct {
	ListenIP string
	ListenWebPort int
	Username string
	Password string
	LogDirectory string
}

type DnsConfig struct {
	ListenIP string
	Zone string
	Ttl int
}
type Config struct {
	WebConfig WebConfig
	DnsConfig DnsConfig
	Banner string
}
