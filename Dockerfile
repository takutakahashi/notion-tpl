FROM golang as builder

WORKDIR /src
COPY go.mod /src/
COPY go.sum /src/
RUN go mod download

COPY . /src/
RUN make build

FROM ubuntu

COPY --from=builder /src/dist/cmd /bin/notion-tpl
COPY --from=builder /src/src/hugo.md.tpl /etc/notion-tpl/
