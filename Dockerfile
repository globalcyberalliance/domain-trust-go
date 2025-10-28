FROM golang:1.25.3-alpine3.22 AS build

RUN apk add git make \
    && apk cache clean

WORKDIR /src

# Install setup dependencies (e.g. betteralign).
COPY Makefile ./
RUN make setup

# Cache Go modules and local vendored dependencies needed for replacements.
COPY go.mod go.sum ./
RUN go mod download

COPY . /src

RUN make

FROM alpine:3.22

RUN apk add bash \
    && apk cache clean

COPY --from=build /src/bin/client /usr/bin/client

WORKDIR /app

ENTRYPOINT ["client"]

CMD ["client"]
