ARG RUST_VERSION=latest

FROM rust:${RUST_VERSION}

WORKDIR /app

ARG CURRENT_GROUP=1000
ARG CURRENT_USERID=1000

RUN set -eu && \
    addgroup --system --gid ${CURRENT_GROUP} builder; \
    adduser --system --no-create-home --ingroup builder --uid ${CURRENT_USERID} builder;

USER builder

RUN set -eu && \
    rustup component add \
        rustfmt \
        clippy;

COPY Cargo.toml Cargo.lock /app/

COPY src/main.rs /app/src/

RUN cargo fetch

