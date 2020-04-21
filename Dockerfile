FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/github.com/ingvagabund/ci-operator-result-aggregation-server
COPY . .

RUN go build -o ci-operator-result-aggregation-server ./cmd/ci-operator-result-aggregation-server

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /go/src/github.com/ingvagabund/ci-operator-result-aggregation-server/ci-operator-result-aggregation-server /usr/bin/

LABEL io.k8s.display-name="CI Operator Result Aggregation Server" \
      io.k8s.description="This is a component of OpenShift for collecting and reporting CI operator errors" \
      io.openshift.tags="openshift,ci-operator-result-aggregation-server" \
      maintainer="AOS ???, <jchaloup@redhat.com>"
