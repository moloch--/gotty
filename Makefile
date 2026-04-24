OUTPUT_DIR = ./builds
GIT_COMMIT = `git rev-parse HEAD | cut -c1-7`
VERSION = 2.0.0-alpha.3
BUILD_OPTIONS = -ldflags "-X main.Version=$(VERSION) -X main.CommitID=$(GIT_COMMIT)"

gotty: main.go server/*.go webtty/*.go backend/*.go Makefile
	go build ${BUILD_OPTIONS}

.PHONY: asset
asset: bindata/static/js/gotty-bundle.js bindata/static/index.html bindata/static/favicon.png bindata/static/css/index.css bindata/static/css/xterm.css bindata/static/css/xterm_customize.css
	go-bindata -prefix bindata -pkg server -ignore=\\.gitkeep -o server/asset.go bindata/...
	gofmt -w server/asset.go

.PHONY: all
all: asset gotty

.PHONY: clean
clean:
	rm -rf gotty bindata ./builds

bindata:
	mkdir bindata

bindata/static: bindata
	mkdir bindata/static

bindata/static/index.html: bindata/static resources/index.html
	cp resources/index.html bindata/static/index.html

bindata/static/favicon.png: bindata/static resources/favicon.png
	cp resources/favicon.png bindata/static/favicon.png

bindata/static/js: bindata/static
	mkdir -p bindata/static/js


bindata/static/js/gotty-bundle.js: bindata/static/js js/dist/gotty-bundle.js
	cp js/dist/gotty-bundle.js bindata/static/js/gotty-bundle.js

bindata/static/css: bindata/static
	mkdir -p bindata/static/css

bindata/static/css/index.css: bindata/static/css resources/index.css
	cp resources/index.css bindata/static/css/index.css

bindata/static/css/xterm_customize.css: bindata/static/css resources/xterm_customize.css
	cp resources/xterm_customize.css bindata/static/css/xterm_customize.css

bindata/static/css/xterm.css: bindata/static/css js/node_modules/@xterm/xterm/css/xterm.css
	cp js/node_modules/@xterm/xterm/css/xterm.css bindata/static/css/xterm.css

js/node_modules/@xterm/xterm/css/xterm.css: js/package-lock.json
	cd js && \
	npm install

js/dist/gotty-bundle.js: js/src/* js/package-lock.json js/webpack.config.js js/tsconfig.json js/node_modules/webpack
	cd js && \
	npx webpack

js/node_modules/webpack: js/package-lock.json
	cd js && \
	npm install

tools:
	go install github.com/go-bindata/go-bindata/go-bindata@latest
	go install github.com/mitchellh/gox@latest
	go install github.com/tcnksm/ghr@latest

test:
	if [ `gofmt -l $$(find . -name '*.go' -not -path './vendor/*') | wc -l` -gt 0 ]; then echo "go fmt error"; exit 1; fi
	go test ./...

cross_compile:
	GOARM=5 gox -os="darwin linux freebsd netbsd openbsd" -arch="386 amd64 arm" -osarch="!darwin/arm" -output "${OUTPUT_DIR}/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

targz:
	mkdir -p ${OUTPUT_DIR}/dist
	cd ${OUTPUT_DIR}/pkg/; for osarch in *; do (cd $$osarch; tar zcvf ../../dist/gotty_${VERSION}_$$osarch.tar.gz ./*); done;

shasums:
	cd ${OUTPUT_DIR}/dist; sha256sum * > ./SHA256SUMS

release:
	ghr -c ${GIT_COMMIT} --delete --prerelease -u yudai -r gotty pre-release ${OUTPUT_DIR}/dist
