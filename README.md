<div align="center">
    <img src="https://gitlab.com/lu-ka/goldig/-/raw/main/goldig.png">
</div>

# GOLDIG
Execute commands (sneaky) on Windows via DNS TXT records

This small utility will create a DNS server and serve your malicious TXT record. Just copy & paste the command from the console output into your Windows shell. 

<img src="https://gitlab.com/lu-ka/goldig/-/raw/main/run.png">

## why
This was made for all the people that don't want to buy a domain to test this attack or just want to run this internally.

But be aware when accessing this externally, that larger corporations will often block external DNS servers or external DNS traffic altogether. In such cases you have to buy a domain and setup a public TXT record.

## Kudos
+ [Alh4zr3d](https://twitter.com/Alh4zr3d) for the [tweet](https://twitter.com/Alh4zr3d/status/1566489367232651264) that started it this
+ [phuslu](https://github.com/phuslu/) for the [fastdns](https://github.com/phuslu/fastdns) library

## Meaning
+ from GO + dig
+ with the 'l' means golden in German

## Installation
```bash
go install gitlab.com/lu-ka/goldig@latest
```
## Build from source
```bash
git clone https://gitlab.com/lu-ka/goldig
cd goldig
go build

# or via makefile
make linuxgit re
make windows
make all
```

## Usage
```
usage: (sudo) ./goldig
-i	Interface to bind to (default: eth0)
-p	Port to use (default: 53)
-c	Command to serve (default: calc.exe)
-d	Domain to answer to (default: example.com)
-h	Shows this text

# port 53 usually requires higher privileges
```

## License
This code is released under the [MIT License](https://gitlab.com/lu-ka/goldig/blob/main/LICENSE).

## Legal
This code is provided for educational use only. If you engage in any illegal activity the author does not take any responsibility for it. By using this code, you agree with these terms.
