all:
	go build -o bin/ci-operator-result-aggregation-server ./cmd/ci-operator-result-aggregation-server


test:
	./hack/test.sh

image:
	imagebuilder -t ci-operator-result-aggregation-server:latest .
