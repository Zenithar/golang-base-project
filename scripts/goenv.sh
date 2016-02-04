goroot="$1"
gopath="$2"

go_version_min="1.5"
go_version_install="1.5.3"

vernum() {
	printf "%03d%03d%03d" $(echo "$1" | tr '.' ' ')
}

if command -v "go" >/dev/null; then
    go_version=$(go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
fi

# If we satisfy the version requirement, there is nothing to do. Otherwise
# proceed downloading and installing a go environment.
if [ $(vernum ${go_version}) -ge $(vernum ${go_version_min}) ]; then
	return
fi

export GOPATH="${gopath}"
export GOROOT="${goroot}/${go_version_install}"

export PATH="$PATH:$GOROOT/bin"

if [ ! -x "${GOROOT}/bin/go" ]; then

	mkdir -p "${GOROOT}"

	os=$(uname | tr A-Z a-z)
	arch=$(uname -m | sed -e 's/x86_64/amd64/' | sed -e 's/i.86/386/')

	url="https://golang.org/dl"
	tarball="go${go_version_install}.${os}-${arch}.tar.gz"

	wget -qO- "${url}/${tarball}" | tar -C "${GOROOT}" --strip 1 -xz
fi
