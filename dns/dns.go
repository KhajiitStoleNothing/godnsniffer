package dns

import (
	//Third Party
	"github.com/miekg/dns"
	"github.com/fatih/color"
	tw "github.com/olekukonko/tablewriter"
	//
	"io/ioutil"
	"net"
	"strings"
	"config"
	"context"
	"time"
	"os"
	"sync"
)
//shitty hacks =(
var internalConfig config.DnsConfig
//
var fileMutex = &sync.Mutex{}

type DnsServer struct {
	tcpserver *dns.Server
	udpserver *dns.Server
	config config.DnsConfig
}

func NewDnsServer (config config.DnsConfig) *DnsServer{
	internalConfig = config
	dns.HandleFunc(".",DnsRequestHandler)
	tcpDnsServer := &dns.Server{
		Addr: config.ListenIP+":53",
		Net: "tcp",
	}
	udpDnsServer := &dns.Server{
		Addr: config.ListenIP+":53",
		Net: "udp",
	}
	dnserver := &DnsServer {
		tcpserver: tcpDnsServer,
		udpserver: udpDnsServer,
		config: config,
	}
	return dnserver
}

func DnsRequestHandler (w dns.ResponseWriter, r *dns.Msg){
	domain := r.Question[0].Name
	color.Green("[+] Got DNS request for %s",domain)
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	rr1 := new(dns.A)
	rr1.Hdr = dns.RR_Header {
		Name: domain,
		Rrtype: dns.TypeA,
		Class: dns.ClassINET,
		Ttl: uint32(internalConfig.Ttl),
	}
	if strings.HasSuffix(domain,internalConfig.Zone) {
		rr1.A = net.ParseIP(internalConfig.ListenIP)
		fileMutex.Lock()
		//os.Truncate("./logs/log.txt", 1024*1024)
		f,_:=os.OpenFile("./logs/log.txt",os.O_RDWR|os.O_CREATE,0600)
		f.Seek(142,0)
		remainder,_ := ioutil.ReadAll(f)
		f.Seek(142,0)
		table := tw.NewWriter(f)
		table.SetBorders(tw.Border{Left: true, Top: true, Right: true, Bottom: false})
		data := []string{tw.Pad(time.Now().Format("2006.01.02 15:04")," ",20),tw.Pad(w.RemoteAddr().String()," ",20) ,tw.Pad(domain," ",20)}
		table.Append(data)
		table.Render()
		f.Write(remainder)
		f.Close()
		fileMutex.Unlock()
	}
	m.Answer = []dns.RR{rr1}
	w.WriteMsg(m)
}

func (dnserver *DnsServer) RunTCP() {
	color.Cyan("[!] Starting TCP DNS server at %s",dnserver.config.ListenIP)
	dnserver.tcpserver.ListenAndServe()
}

func (dnserver *DnsServer) RunUDP() {
	color.Cyan("[!] Starting UDP DNS server at %s",dnserver.config.ListenIP)
	dnserver.udpserver.ListenAndServe()
}

func (dnserver *DnsServer) Shutdown() {
	_, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	dnserver.tcpserver.Shutdown()
	dnserver.udpserver.Shutdown()
}
