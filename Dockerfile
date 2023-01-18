FROM --platform=$BUILDPLATFORM golang:1.19 as builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY *.go ./

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -ldflags="-s -w" -o server .

FROM gcr.io/distroless/static as final

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/server .

CMD ["/server"]