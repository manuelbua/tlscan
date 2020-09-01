# GNU Makefile

BIN_DIR := "bin"
CHK_SUMS := "$(BIN_DIR)/checksums.txt"
BASE_OUT := "$(BIN_DIR)/tlscan"
MAIN_GO := "cmd/tlscan/main.go"
LDFLAGS_SH := $(shell go-version-ldflags.sh)
UPX_CMD := "upx --brute"

default: clean darwin linux windows compress package integrity

clean:
	echo "Cleaning.."
	test -n "$(BASE_OUT)*" && rm -f $(BASE_OUT)*
	test -n "$(CHK_SUMS)" && rm -f $(CHK_SUMS)

darwin:
	echo "Building for Darwin/64.."
	GOOS=darwin GOARCH=amd64 go build -o $(BASE_OUT)_darwin64 -ldflags="$(LDFLAGS_SH)" $(MAIN_GO)

linux:
	echo "Building for Linux/64.."
	GOOS=linux GOARCH=amd64 go build -o $(BASE_OUT)_linux64 -ldflags="$(LDFLAGS_SH)" $(MAIN_GO)

windows:
	echo "Building for Windows/64.."
	GOOS=windows GOARCH=amd64 go build -o $(BASE_OUT)_win64 -ldflags="$(LDFLAGS_SH)" $(MAIN_GO)

compress:
	echo "Compressing executables.."
	$(shell $$( \
	"$(UPX_CMD)" "$(BASE_OUT)_darwin64" & \
	"$(UPX_CMD)" "$(BASE_OUT)_linux64" & \
	"$(UPX_CMD)" "$(BASE_OUT)_win64"))

package:
	echo "Packaging executables.."
	zip -j "$(BASE_OUT)_darwin.zip" "$(BASE_OUT)_darwin64"
	zip -j "$(BASE_OUT)_linux64.zip" "$(BASE_OUT)_linux64"
	zip -j "$(BASE_OUT)_win64.zip" "$(BASE_OUT)_win64"

integrity:
	echo "Computing integrity.."
	shasum -a 256 $(BASE_OUT)*.zip > $(CHK_SUMS)

