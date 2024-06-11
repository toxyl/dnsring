# dnsring
Ring of DNS servers for spreading IP lookups over multiple DNS servers. Uses an in-memory cache to avoid repeated lookups. Only returns the first IP if matches and prefers IPv4 over IPv6.

# Usage Example
```go
    DNSRing := dnsring.New(dns.NewServer("8.8.8.8", 53, 2000), dns.NewServer("8.8.4.4", 53, 2000))
    host := "google.com"
	ip := DNSRing.Next().HostToIP(host)
	if ip == "" {
        fmt.Printf("%s did not resolve to an IP\n", host)
	} else {
        fmt.Printf("%s has the IP: %s\n", host, ip)
    }
```
