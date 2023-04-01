# goldig
Execute commands on Windows via DNS records.

This utility will create a DNS server and serve a malicious TXT record. Just copy & paste the command from the console output into your Windows shell. 

## Usage
```
usage: goldig [OPTIONS]

Options:
    -i 	Interface to bind to (default: eth0)
    -c 	Command to serve (default: calc.exe)
    -h 	Shows this text
```

## Installation
```bash
go install gitlab.com/rtfmkiesel/goldig@latest
```

## Build from source
```bash
git clone https://gitlab.com/rtfmkiesel/goldig
cd goldig
# to build binary in the current directory
go build -ldflags="-s -w" .
# to build & install binary into GOPATH/bin
go install .
```

## Why
This was made for all the people who don't want to buy a domain to test this attack or just want to run this internally. Be aware when accessing this externally that larger corporations will often block external DNS servers or external DNS traffic altogether. In such cases you have to buy a domain and setup a public DNS record.

## Kudos
+ [Alh4zr3d](https://twitter.com/Alh4zr3d) for the [tweet](https://twitter.com/Alh4zr3d/status/1566489367232651264) that started it this
+ [phuslu](https://github.com/phuslu/) for the [fastdns](https://github.com/phuslu/fastdns) library
+ [I_Am_Jakoby](https://twitter.com/I_Am_Jakoby) for the "> 255 char" [payload](https://twitter.com/I_Am_Jakoby/status/1570770517589917697)

## Legal
This code is provided for educational use only. If you engage in any illegal activity the author does not take any responsibility for it. By using this code, you agree with these terms.
