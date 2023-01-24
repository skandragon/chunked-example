#
# Copyright 2023 Michael Graff.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#
# Install the latest versions of our mods.  This is done as a separate step
# so it will pull from an image cache if possible, unless there are changes.
#
FROM --platform=${BUILDPLATFORM} golang:1.19-alpine AS buildmod
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

#
# Compile the code.
#
FROM buildmod AS build-binaries
COPY . .
ARG GIT_BRANCH
ARG GIT_HASH
ARG BUILD_TYPE
ARG TARGETOS
ARG TARGETARCH
ENV GIT_BRANCH=${GIT_BRANCH} GIT_HASH=${GIT_HASH} BUILD_TYPE=${BUILD_TYPE}
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
RUN mkdir /out
RUN go build -o /out/chunked-client app/chunked-client/*.go
RUN go build -o /out/chunked-client-bufio app/chunked-client-bufio/*.go
RUN go build -o /out/chunked-server app/chunked-server/*.go

#
# Build the chunked-client image.  This should be a --target on docker build.
#
FROM scratch AS chunked-client-image
WORKDIR /app
COPY --from=build-binaries /out/chunked-client /app
EXPOSE 8090
CMD ["/app/chunked-client"]

#
# Build the chunked-client-bufio image.  This should be a --target on docker build.
#
FROM scratch AS chunked-client-bufio-image
WORKDIR /app
COPY --from=build-binaries /out/chunked-client-bufio /app
EXPOSE 8090
CMD ["/app/chunked-client-bufio"]

#
# Build the chunked-server image.  This should be a --target on docker build.
#
FROM scratch AS chunked-server-image
WORKDIR /app
COPY --from=build-binaries /out/chunked-server /app
EXPOSE 8090
CMD ["/app/chunked-server"]
