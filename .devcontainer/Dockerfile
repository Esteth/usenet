FROM golang:1.13

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# This Dockerfile adds a non-root user with sudo access. Use the "remoteUser"
# property in devcontainer.json to use it. On Linux, the container user's GID/UIDs
# will be updated to match your local UID/GID (when using the dockerFile property).
# See https://aka.ms/vscode-remote/containers/non-root-user for details.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Configure apt, install packages and tools
RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils dialog 2>&1 \
    #
    # Verify git, process tools, lsb-release (common in install instructions for CLIs) installed
    && apt-get -y install git iproute2 procps lsb-release \
    #
    # Install Go tools w/module support
    && mkdir -p /tmp/gotools \
    && cd /tmp/gotools
RUN GO111MODULE=on go get -v golang.org/x/tools/gopls 2>&1
RUN GO111MODULE=on go get -v \
        honnef.co/go/tools \
        golang.org/x/tools/cmd/gorename \
        golang.org/x/tools/cmd/goimports \
        golang.org/x/tools/cmd/guru \
        golang.org/x/lint/golint \
        github.com/mdempsky/gocode \
        github.com/cweill/gotests \
        github.com/haya14busa/goplay/cmd/goplay \
        github.com/sqs/goreturns \
        github.com/josharian/impl \
        github.com/davidrjenni/reftools/cmd/fillstruct \
        github.com/ramya-rao-a/go-outline  \
        github.com/acroca/go-symbols  \
        github.com/godoctor/godoctor \
        github.com/rogpeppe/godef  \
        github.com/zmb3/gogetdoc \
        github.com/fatih/gomodifytags \
        github.com/mgechev/revive  \
        github.com/go-delve/delve/cmd/dlv 2>&1 \
    #
    # Install Go tools w/o module support
    && go get -v github.com/alecthomas/gometalinter 2>&1 \
    #
    # Install gocode-gomod
    && go get -x -d github.com/stamblerre/gocode 2>&1 \
    && go build -o gocode-gomod github.com/stamblerre/gocode \
    && mv gocode-gomod $GOPATH/bin/ \
    #
    # Install golangci-lint
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 2>&1 \
    #
    # Create a non-root user to use if preferred - see https://aka.ms/vscode-remote/containers/non-root-user.
    && groupadd --gid $USER_GID $USERNAME \
    && useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME \
    # [Optional] Add sudo support
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    #
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* /go/src /tmp/gotools

# Update this to "on" or "off" as appropriate
ENV GO111MODULE=auto

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog