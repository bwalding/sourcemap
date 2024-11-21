FROM golang:latest

COPY . /go/sourcemap
WORKDIR /go/sourcemap
RUN go install .
CMD [ "/go/bin/sourcemap", "--output",  "/sources", "--debug"]
