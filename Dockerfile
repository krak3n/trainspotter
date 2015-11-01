FROM alpine:3.2
MAINTAINER Chris Reeves <hello@chris.reeves.io>

# Install OS Dependencies
RUN apk update \
    && apk add \
        build-base \
        go \
        go-tools \
        ca-certificates \
        make \
        bash \
        git \
    && rm -rf /var/cache/apk/* \
    && wget https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm \
    && chmod +x gpm \
    && mv gpm /usr/local/bin

# Environment Variables
ENV GOPATH /trainspotter
ENV PATH $PATH:$GOPATH/bin

# Insttal Go Package Manager and Install Dependencies
COPY ./Godeps /Godeps
RUN  gpm install

# Working Directory
WORKDIR /trainspotter/src/github.com/krak3n/Trainspotter

# Set our Application Entrypoint
ENTRYPOINT ["trainspotter"]

# Copy Application Source to expected directory
COPY . /trainspotter/src/github.com/krak3n/Trainspotter

# Install Application
RUN make install
