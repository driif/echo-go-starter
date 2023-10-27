### -----------------------
# --- Stage: development
# --- Purpose: Local development environment
# --- https://hub.docker.com/_/golang
# --- https://github.com/microsoft/vscode-remote-try-go/blob/master/.devcontainer/Dockerfile
### -----------------------
FROM golang:1.21.3-bullseye as dev

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive
ENV CGO_ENABLED=0

# postgresql-support: Add the official postgres repo to install the matching postgresql-client tools of your stack
# https://wiki.postgresql.org/wiki/Apt
# run lsb_release -c inside the container to pick the proper repository flavor
# e.g. stretch=>stretch-pgdg, buster=>buster-pgdg, bullseye=>bullseye-pgdg
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ bullseye-pgdg main" \
    | tee /etc/apt/sources.list.d/pgdg.list \
    && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc \
    | apt-key add -

# Install required system dependencies
RUN apt-get update \
    && apt-get install -y \
    #
    # Mandadory minimal linux packages
    # Installed at development stage and app stage
    # Do not forget to add mandadory linux packages to the final app Dockerfile stage below!
    #
    # -- START MANDADORY --
    ca-certificates \
    # --- END MANDADORY ---
    #
    # Development specific packages
    # Only installed at development stage and NOT available in the final Docker stage
    # based upon
    # https://github.com/microsoft/vscode-remote-try-go/blob/master/.devcontainer/Dockerfile
    # https://raw.githubusercontent.com/microsoft/vscode-dev-containers/master/script-library/common-debian.sh
    #
    # icu-devtools: https://stackoverflow.com/questions/58736399/how-to-get-vscode-liveshare-extension-working-when-running-inside-vscode-remote
    # graphviz: https://github.com/google/pprof#building-pprof
    # -- START DEVELOPMENT --
    apt-utils \
    dialog \
    openssh-client \
    less \
    iproute2 \
    procps \
    lsb-release \
    locales \
    sudo \
    bash-completion \
    bsdmainutils \
    graphviz \
    xz-utils \
    postgresql-client-15 \
    icu-devtools \
    tmux \
    rsync \
    git \
    tzdata \
    # --- END DEVELOPMENT ---
    #
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN go install github.com/volatiletech/sqlboiler/v4@latest && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest && \
    go install github.com/rubenv/sql-migrate/sql-migrate@latest && \
    go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["air", "-c", ".air.toml"]

# Build Image
FROM golang:1.21.3-bullseye as build

ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Production Image (distroless) with only the binary
FROM gcr.io/distroless/static-debian11 as app
COPY --from=build /app/main /app/main
COPY --from=build /app/embedding /app/embedding

WORKDIR /app

CMD ["./main", "run"]