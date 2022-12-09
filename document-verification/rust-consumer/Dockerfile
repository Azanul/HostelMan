FROM rust:1.65-alpine3.15

RUN apk add --no-cache gcc musl-dev linux-headers libressl-dev

COPY . .
RUN cargo build --release
CMD ["./target/release/docker"]