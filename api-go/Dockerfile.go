#Pickin a base image that is light 
FROM golang:1.24-bullseye AS build

#Crearing a directory and cd into the directory
WORKDIR /app

#Copy only files required to install dependencies(better layer caching)
COPY go.mod go.sum ./

# Use cache mount to speed up install of existing dependencies
RUN --mount=type=cache,target=/go/pkg/mode \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM build AS dev

#install air for hot reload & delve for debugging
RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest


COPY . .

CMD [ "air", "-c", ".air.toml" ]

FROM build AS build-production

#adding a nonroot user for security purpose
RUN useradd -u 1001 nonroot

COPY . .

#compile healthcheck
RUN go build \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o healthcheck \
    ./healthcheck/healthcheck.go

#compile aplication during build rather than at runtime
RUN go build \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o api-golang
#use seperate dtage for deployable image
FROM scratch

#set gin mode
ENV GIN_MODE=release

WORKDIR /

#Copy the password file
COPY --from=build-production /etc/passwd /etc/passwd

#Copy the healthecheck file from the build stage
COPY --from=build-production /app/healthcheck/healthcheck healthcheck

#Copy the app binary from the build stage
COPY --from=build-production /app/api-golang api-golang

#Use nonroot user
USER nonroot

#Exposing the port specified in the src code
EXPOSE 8080

CMD [ "/api-golang" ]