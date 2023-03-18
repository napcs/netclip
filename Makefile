.PHONY: all

VERSION = $(shell go run bin/netclip.go -v | cut -c10- )

all: windows_64 mac_silicon mac_intel linux_64 release

windows_64:
	mkdir -p dist/windows_64
	env GOOS=windows GOARCH=amd64 go build -o dist/windows_64/netclip.exe bin/netclip.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows_64/netclipsvc.exe service/netclipsvc.go


mac_intel:
	mkdir -p dist/mac_intel
	env GOOS=darwin GOARCH=amd64 go build -o dist/mac_intel/netclip bin/netclip.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/mac_intel/netclipsvc service/netclipsvc.go

mac_silicon:
	mkdir -p dist/mac_silicon
	env GOOS=darwin GOARCH=arm64 go build -o dist/mac_silicon/netclip bin/netclip.go
	env GOOS=darwin GOARCH=arm64 go build -o dist/mac_silicon/netclipsvc service/netclipsvc.go

linux_64:
	mkdir -p dist/linux_64
	env GOOS=linux GOARCH=amd64 go build -o dist/linux_64/netclip bin/netclip.go
	env GOOS=linux GOARCH=amd64 go build -o dist/linux_64/netclipsvc service/netclipsvc.go

release:
	cd dist/windows_64 && zip netclip_${VERSION}_windows64.zip netclip.exe netclipsvc.exe && mv netclip_${VERSION}_windows64.zip ../
	cd dist/mac_intel && zip netclip_${VERSION}_mac_intel.zip netclip netclipsvc && mv netclip_${VERSION}_mac_intel.zip ../
	cd dist/mac_silicon && zip netclip_${VERSION}_mac_silicon.zip netclip netclipsvc && mv netclip_${VERSION}_mac_silicon.zip ../
	cd dist/linux_64 && zip netclip_${VERSION}_linux64.zip netclip netclipsvc && mv netclip_${VERSION}_linux64.zip ../
	cd dist/linux_64 && tar -czvf netclip_${VERSION}_linux64.tar.gz netclip netclipsvc && mv netclip_${VERSION}_linux64.tar.gz ../

.PHONY: clean
clean:
	rm -r dist/


