# build golang backend application
FROM amd64/golang:1.25.3-alpine AS backend
ARG COMMIT_SHA
ARG VERSION_TAG
ARG GIT_BRANCH
ENV CGO_ENABLED=0 GOARCH=amd64 GOOS=linux
WORKDIR /app
COPY src/backend .
RUN go build -ldflags "-s -w \
    -X 'opendataaggregator/config.CommitHash=${COMMIT_SHA}' \
    -X 'opendataaggregator/config.VersionTag=${VERSION_TAG}' \
    -X 'opendataaggregator/config.Branch=${GIT_BRANCH}'" \
    -trimpath -mod=vendor -o build/opendata_aggregator main.go

# build vue frontend static files
FROM amd64/node:24-alpine AS frontend
ENV PNPM_HOME="/pnpm" PATH+=":$PNPM_HOME"
RUN corepack enable && corepack use pnpm@latest
WORKDIR /app
COPY src/frontend/package.json src/frontend/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY src/frontend/vite.config.ts src/frontend/tsconfig.json src/frontend/index.html ./
COPY src/frontend/public public
COPY src/frontend/src src
RUN pnpm build

# final stage, full build application
FROM amd64/alpine:latest
ARG COMMIT_SHA
ARG VERSION_TAG
ARG GIT_BRANCH
LABEL git.remote.origin.branch="${GIT_BRANCH}"
LABEL git.remote.origin.tag="${VERSION_TAG}"
LABEL git.remote.origin.commit_sha="${COMMIT_SHA}"
WORKDIR /opt
COPY --from=frontend app/dist /opt/static
COPY --from=backend app/build/opendata_aggregator /usr/bin/opendata_aggregator
CMD [ "opendata_aggregator", "server", "--serve", "--config=config.toml" ]
