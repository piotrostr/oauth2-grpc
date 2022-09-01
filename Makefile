all: install compile-proto generate-cert

install:
	./hack/install.sh

compile-proto:
	./hack/compile-proto.sh

generate-cert:
	./hack/generate-cert.sh

run:
	./hack/run.sh

