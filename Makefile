all: install compile-proto

install:
	./hack/install.sh

compile-proto: 
	./hack/compile-proto.sh
