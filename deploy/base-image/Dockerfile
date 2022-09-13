FROM cosmtrek/air:v1.40.4

ENV GO111MODULE on

ENV GOSUMDB off

# Precompile the entire go standard library into the first Docker cache layer: useful for other projects too!
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go install -v -installsuffix cgo -a std

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin \
    && golangci-lint --version

RUN go install github.com/smartystreets/goconvey@latest

RUN go install github.com/mailru/easyjson/...@latest

RUN go install go.k6.io/xk6/cmd/xk6@latest

RUN go install go.k6.io/k6@latest

RUN go install  github.com/pascaldekloe/colfer/cmd/colf@latest

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user \
    && echo 'user:x:65534:65534::/:' > /user/passwd \
    && echo 'user:x:65534:' > /user/group