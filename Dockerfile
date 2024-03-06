# Build httpie-go on a big golang image
FROM golang:1.13 as builder
RUN mkdir /build 
RUN git clone https://github.com/nojima/httpie-go.git /build/httpie-go
WORKDIR /build/httpie-go
RUN CGO_ENABLED=0 go build -o http ./cmd/ht

# Copy the executable to the tiny busybox image
FROM busybox
COPY --from=builder /build/httpie-go/http /bin/