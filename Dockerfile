FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .
# Build a small image
FROM scratch
COPY --from=builder /dist/main /
ENTRYPOINT ["/main"]

EXPOSE 8080
#ENV herokuDB1='host=ec2-34-243-101-244.eu-west-1.compute.amazonaws.com port=5432 user=hvbofdxjbkkdgq password=ff9c8195d4fa5205036cb92a384e142c9ca7bfbbc5f7639f038b4925bacdfea9 dbname=d62omvefcmhpmq'
#RUN go test ./dist/main