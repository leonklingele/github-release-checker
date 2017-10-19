APP_DIR = $(HOME)/.github-release-checker
CONFIG_PATH = "$(APP_DIR)/config.toml"

INSTALL = install

.PHONY: all
all: build config

.PHONY: build
build:
	go build .

.PHONY: install
install: config
	go install

.PHONY: config
config:
	@if [ ! -f "${CONFIG_PATH}" ]; then \
		OLD_UMASK=$(shell echo `umask`) ; \
		umask 077 > /dev/null ; \
		mkdir -p $(APP_DIR) ; \
		$(INSTALL) -c -m 0600 config.toml.in $(APP_DIR)/config.toml ; \
		umask $(OLD_UMASK) > /dev/null ; \
		echo "Config file installed at ${CONFIG_PATH}" ; \
	fi
