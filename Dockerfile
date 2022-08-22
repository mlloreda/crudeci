FROM golang:1.19.0-alpine3.16
MAINTAINER Miguel Lloreda <mig.lloreda@gmail.com>

WORKDIR /crudeci

COPY . .
RUN apk update && apk add --no-cache make build-base
RUN make build

CMD ["make", "test"]
