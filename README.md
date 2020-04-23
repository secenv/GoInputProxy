# GoInputProxy
goInputProxy - log and pass inputs (STDIN, arguments, environment variables) to /htdocs/cgibin_, log outputs to /tmp/ and system logs.

This go program was used to debug the "/htdocs/cgibin" binary executable file in the D-Link DIR-859L and other SOHO routers. It helped us discover and debug multiple vulnerabilities:

https://medium.com/@s1kr10s/d-link-dir-859-rce-unautenticated-cve-2019-17621-en-d94b47a15104

https://medium.com/@s1kr10s/d-link-dir-859-unauthenticated-rce-in-ssdpcgi-http-st-cve-2019-20215-en-2e799acb8a73

https://medium.com/@s1kr10s/d-link-dir-859-rce-unauthenticated-cve-2019-20216-cve-2019-20217-en-6bca043500ae

## Usage

To use goInputProxy, you have to compile it for the target architecture, upload the resulting executable file to the device, give it execution permissions, rename /htdocs/cgibin to /htdocs/cgibin_ and then rename goInputProxy as /htdocs/cgibin (be careful and keep a copy of the original /htdocs/cgibin somewhere!!).

In short, run the following commands:

In your local machine, run

`$ GOARCH=mips go build -ldflags '-s -w' -o /tmp/inputProxy inputProxy.go; cd /tmp; python3 -m http.server`

In the (emulated) router, modify the $MYIP variable and then run

`$ cd /tmp/; rm inputProxy; MYIP=192.168.0.101; wget http://$MYIP:8000/inputProxy; cp inputProxy /htdocs/cgibin; chmod 777 /htdocs/cgibin`

## What is it for?
In the aforementioned routers, "/htdocs/cgibin" receives inputs from the Mathopd HTTP server through STDIN, as arguments and as environment variables. goInputProxy will:
- send the inputs to system logs (in this case, /dev/ttyS0) and to /tmp/LOG_CGIBIN
- pass these inputs to the target program (/htdocs/cgibin_, originally called /htdocs/cgibin)
- receive and log its outputs
- quit
