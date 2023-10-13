FROM golang:1.20.10

RUN apt-get update && apt-get install -y ca-certificates git-core ssh

RUN go install github.com/cespare/reflex@latest && \
    go install github.com/google/wire/cmd/wire@latest

ADD ./cicd/id_rsa_shared /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa && \
    echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config && \
    git config --global url.ssh://git@github.com/vireocloud.insteadOf https://github.com/vireocloud

COPY . /app
WORKDIR /app
RUN go clean -modcache && \
    export GONOSUMDB="github.com/vireocloud" && \
    go get github.com/vireocloud/property-pros-sdk && \
    go mod tidy && \
    go mod download

COPY reflex.conf /usr/local/etc/

VOLUME /go

ENTRYPOINT [ "reflex", "-d", "none", "-c", "/usr/local/etc/reflex.conf"]