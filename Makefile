# Build emyt
default: linux

.PHONY: emyt_linux
linux:
	@echo "Building emyt binary to './builds/emyt'"
	@(cd cmd/; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-s -w" -o ../builds/emyt)

.PHONY: emyt_osx
osx:
	@echo "Building emyt(emyt_osx) binary to './builds/emyt_osx'"
	@(cd cmd/; CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build --ldflags "-s -w" -o ../builds/emyt_osx)

.PHONY: emyt_win
windows:
	@echo "Building emyt(emyt_windows) binary to './builds/emyt_win.exe'"
	@(cd cmd/; CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build --ldflags "-s -w" -o ../builds/emyt_win.exe)

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf emyt_win.exe emyt_osx emyt
 
install: install_linux

.PHONY: emyt_install_linux
install_linux:
	@echo "Installing emyt Proxy to /usr/sbin/emyt directory"
	@cp -f builds/emyt /usr/sbin/emyt
	@mkdir -p /etc/emyt
	@mkdir -p /etc/emyt/ssl
	@mkdir -p /var/log/emyt
	@mkdir -p /var/www/html
	@cp -n app.yaml.example /etc/emyt
	@cp -n certgen.sh /etc/emyt/ssl
	@printf "[Unit]\nDescription=emyt the Reverse Proxy\n\n[Service]\nType=simple\nRestart=always\n\RestartSec=5s\nExecStart=/usr/sbin/emyt\n\n[Install]\nWantedBy=multi-user.target\n" > /lib/systemd/system/emyt-proxy.service
	@echo "Start emyt service using"
	@echo "-----------------------------------"
	@echo " $> sudo service emyt-proxy start"
	@echo "-----------------------------------"
	@echo "Enabling Service"
	@systemctl enable emyt-proxy.service


