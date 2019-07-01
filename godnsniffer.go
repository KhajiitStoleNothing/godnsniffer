package main

import (
	"flag"
	"web"
	"config"
	"os"
	"os/signal"
	"dns"
	"strings"
//Third party
	"github.com/fatih/color"
	tw "github.com/olekukonko/tablewriter"
)

const Banner = `
   ___   ___         ___    __  __      __       _  __  __           
  / _ \ /___\       /   \/\ \ \/ _\    / _\_ __ (_)/ _|/ _| ___ _ __ 
 / /_\///  //      / /\ /  \/ /\ \     \ \| '_ \| | |_| |_ / _ \ '__|
/ /_\\/ \_//      / /_// /\  / _\ \    _\ \ | | | |  _|  _|  __/ |   
\____/\___/      /___,'\_\ \/  \__/    \__/_| |_|_|_| |_|  \___|_|   
                                                                     
`

func main() {
	var config config.Config
	color.Yellow(Banner)
	flag.StringVar(&config.WebConfig.ListenIP, "web-ip", "0.0.0.0", "IPv4 address for webserver to listen on. Default 0.0.0.0")
	flag.IntVar(&config.WebConfig.ListenWebPort, "port", 8080, "Port to run webserver with logs on")
	flag.StringVar(&config.WebConfig.Username, "username", "admin", "Default username for web interface")
	flag.StringVar(&config.WebConfig.Password, "password", "GiveMeLogz", "Default password for web interface")
	flag.StringVar(&config.WebConfig.LogDirectory,"log-dir","/secure_logs","Default web directory to server logs")
	flag.StringVar(&config.DnsConfig.ListenIP, "dns-ip","0.0.0.0","IPv4 address for DNS to listen on")
	flag.StringVar(&config.DnsConfig.Zone,"zone","test.pw","DNS zone to serve")
	flag.IntVar(&config.DnsConfig.Ttl,"ttl",86400,"TTL of DNS records")
	flag.Parse()
	if _, err := os.Stat("./logs"); os.IsNotExist(err){
		os.Mkdir("./logs", os.ModePerm)
		f, _ := os.OpenFile("./logs/log.txt",os.O_WRONLY|os.O_CREATE,0600)
		table := tw.NewWriter(f)
		data := []string{tw.Pad("DATE"," ",20),tw.Pad("IP"," ",20),tw.Pad("DOMAIN"," ",20)}
		table.Append(data)
		table.Render()

	}
	if !strings.HasSuffix(config.DnsConfig.Zone,"."){
		config.DnsConfig.Zone += "."
	}
	if !strings.HasPrefix(config.WebConfig.LogDirectory,"/"){
		config.WebConfig.LogDirectory = "/" + config.WebConfig.LogDirectory
	}
	//Starting and serving web server
	webConfig := config.WebConfig
	webServer := web.NewWebServer(webConfig)
	go webServer.Run()
	//Starting and serving dns server
	dnsConfig := config.DnsConfig
	dnsServer := dns.NewDnsServer(dnsConfig)
	go dnsServer.RunTCP()
	go dnsServer.RunUDP()
	//Detecting command-line interrupt
	c := make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt)
	<-c
	color.Red("\n[!!!!]CTRL+C Recevied. Exiting gracefully")
	//Gracefull shutdown
	webServer.Shutdown()
	dnsServer.Shutdown()
}
