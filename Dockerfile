# syntax=docker/dockerfile:1.0.0-experimental
# https://blog.oddbit.com/post/2019-02-24-docker-build-learns-about-secr/

FROM golang:1.20-alpine as builder
  # git+ssh required to access private GitHub repositories
  # ca-certificates required so we have a set of CA certs to
  # copy into our final container
  RUN apk add alpine-sdk openssh git ca-certificates

  RUN mkdir -m 700 /root/.ssh; \
    touch -m 600 /root/.ssh/known_hosts; \
    ssh-keyscan github.com > /root/.ssh/known_hosts

  RUN git config --global \
    url."git@github.com:jtyers".insteadOf "https://github.com/jtyers" 

  WORKDIR /usr/local/src/app
  
  # copy go module files first, and run go get - this is to speed
  # up container builds when deps have not changed
  COPY go.sum .
  COPY go.mod .

  ENV GOPRIVATE=github.com/jtyers/*
  # FIXME needed to let gin-jwt/v2 download
  #ENV GOPRIVATE=github.com/jtyers/tmaas-*
   
  # build it somewhere outside of the GOPATH (/go for golang images)
  COPY . .

  RUN --mount=type=ssh \
       go test -cover -v \
    && CGO_ENABLED=0 go build -o app
 
FROM scratch
  COPY --from=builder /usr/local/src/app/app /app

  # scratch contains nothing, which means no default CA certs, which are
  # needed for us to make outbound HTTPS requests (eg to GCP APIs), so
  # we copy them in here from Alpine's included set
  COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

  ENTRYPOINT [ "/app" ]
