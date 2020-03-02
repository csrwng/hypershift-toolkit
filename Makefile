SRC_DIRS = cmd pkg


.PHONY: default
default: build

.PHONY: build
build:  bindata control-plane-operator
	go build -o bin/hypershift github.com/openshift/hypershift-toolkit/cmd/hypershift

.PHONY: hypershift-installer
hypershift-installer: bindata
	go build -o bin/hypershift-installer github.com/openshift/hypershift-toolkit/contrib/cmd/hypershift-installer

.PHONY: bindata
bindata:
	hack/update-generated-bindata.sh

.PHONY: verify-bindata
verify-bindata:
	hack/verify-generated-bindata.sh

.PHONY: verify-gofmt
verify-gofmt:
	@echo Verifying gofmt
	@gofmt -l -s $(SRC_DIRS)>.out 2>&1 || true
	@[ ! -s .out ] || \
	  (echo && echo "*** Please run 'make fmt' in order to fix the following:" && \
	  cat .out && echo && rm .out && false)
	@rm .out

.PHONY: verify
verify: verify-gofmt verify-bindata

# Build manager binary
.PHONY: control-plane-operator
control-plane-operator:
	go build -o bin/control-plane-operator github.com/openshift/hypershift-toolkit/cmd/control-plane-operator
