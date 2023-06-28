## development #################################################################

FROM docker.ddollar.dev/go:1.20 AS development

ENV DEVELOPMENT=true
ENV PORT=8000
ENV VITE_CLIENT_PORT=443

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

EXPOSE 8000

WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm install
RUN npm rebuild

WORKDIR /src
COPY go.mod go.sum ./
COPY vendor vendor
RUN go list -mod=vendor ./vendor/... | xargs go build -mod=vendor --ldflags="-s -w"

COPY . .

ENTRYPOINT ["go", "run", "."]

## package #####################################################################

FROM docker.ddollar.dev/go:1.20 AS package

ENV CGO_ENABLED=0
ENV VITE_PREFIX=/

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm install
RUN npm rebuild

WORKDIR /src
COPY go.mod go.sum ./
COPY vendor vendor
RUN go list -mod=vendor ./vendor/... | xargs go build -mod=vendor --ldflags="-s -w"

COPY . .

RUN make build

## production ##################################################################

FROM docker.ddollar.dev/scratch AS production

ENV PATH=/opt/bin:$PATH
ENV PORT=8000

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

EXPOSE 8000

WORKDIR /

COPY --from=package /src/dist/* /opt/bin/
