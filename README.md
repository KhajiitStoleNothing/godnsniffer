Go DNS Sniffer
=======

Go DNS Sniffer: Simple sniffer and logger of DNS requestes, written on GO.

### Building From Source and Installation

To build godnsniffer from source, simply run ```go get github.com/khajiitstolenothing/godnsniffer``` and ```cd``` into the project source directory. Then, run ```go build```. After this, you should have a binary called ```godnsffier``` in the current directory.

### Better way to install and use

Just grub binary from release for Linux =)

### Usage

By default you only need to specify ```-zone``` flag with domain name you have. For security purposes additionally change ```-username``` and ```-password``` used to login to the web interface with logs. Then run the binary and access logs of DNS requests through web interface (by default ```http://<IP>:8080/secure_logs```).

### Command line args

```
Usage of ./godnsniffer:
  -dns-ip string
    	IPv4 address for DNS to listen on (default "0.0.0.0")
  -log-dir string
    	Default web directory to server logs (default "/secure_logs")
  -password string
    	Default password for web interface (default "GiveMeLogz")
  -port int
    	Port to run webserver with logs on (default 8080)
  -ttl int
    	TTL of DNS records (default 86400)
  -username string
    	Default username for web interface (default "admin")
  -web-ip string
    	IPv4 address for webserver to listen on. Default 0.0.0.0 (default "0.0.0.0")
  -zone string
    	DNS zone to serve (default "test.pw")
```

### Notice

I know that code is crap, but it works, so whatever. If you want to do it better feel free to fork.

