FROM golang:1.19.2

RUN apt-get update && apt-get install -y ca-certificates git-core ssh

RUN go install github.com/cespare/reflex@latest && \
    go install github.com/google/wire/cmd/wire@latest

ADD ./cicd/id_rsa_shared /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/vireocloud.insteadOf https://github.com/vireocloud

COPY . /app
WORKDIR /app
RUN go clean -modcache
RUN go get -u github.com/vireocloud/property-pros-sdk
RUN go install github.com/vireocloud/property-pros-sdk
RUN go mod tidy 
RUN go mod download

COPY reflex.conf /usr/local/etc/

VOLUME /go

ENTRYPOINT [ "reflex", "-d", "none", "-c", "/usr/local/etc/reflex.conf"]