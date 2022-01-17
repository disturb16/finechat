# project argument is set outside
# to be able to use it on both phases, builder and final
ARG project=finechatbot

# ================================
# ======= BUILD STAGE ============
FROM golang:1.16-alpine AS builder
ARG project

RUN apk update
RUN apk add --no-cache \
    git \
    alpine-sdk \
    upx

# Install go pacakge dependencies.
WORKDIR /${project}
COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x

# Copy the rest of the project
WORKDIR /${project}
COPY . .

RUN go build -o ./${project} \
    && upx ./${project}

# ============================
# ======= FINAL STAGE ========
FROM alpine:3.10
ARG project
ENV entry=${project}

RUN apk update
RUN apk add --no-cache ca-certificates curl

WORKDIR /${project}

# Create non-root user for app
RUN adduser -D -g 'appuser' appuser && \
    chown -R appuser:appuser /${project}

# Copy files from builder and bundler
COPY --from=builder /${project}/${project} .

USER appuser

HEALTHCHECK --interval=30s \
    --timeout=5s \
    --retries=3 \
    --start-period=20s \
    CMD curl \
    -f http://localhost:8080/healthcheck || exit 1

ENTRYPOINT ./${entry}