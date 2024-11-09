
#!/bin/bash
RUN_NAME=packer.exe
mkdir -p tool/resource
mkdir -p tool/outcome
go build -o tool/${RUN_NAME} main.go