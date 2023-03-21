when compiling go code be sure to export CGO_ENABLED=0 
if CGO_ENABLED=1 the webhook pod will fail on startup with the following error:

/usr/local/bin/webhook: /lib64/libc.so.6: version `GLIBC_2.32' not found (required by /usr/local/bin/webhook)
/usr/local/bin/webhook: /lib64/libc.so.6: version `GLIBC_2.34' not found (required by /usr/local/bin/webhook)


to build the go module:

go build -o bin/webhook main.go


