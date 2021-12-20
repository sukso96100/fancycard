# build stage
FROM golang:1.17-bullseye AS build

RUN mkdir /build
COPY ./ /build
WORKDIR /build
RUN go build -o fancycard .

# product stage
FROM chromedp/headless-shell:98.0.4758.9

RUN apt-get update && \
    apt-get install -y dumb-init ca-certificates procps fonts-noto-cjk fonts-noto-color-emoji && \
    apt-get autoclean -y && apt-get clean -y && \
    apt-get autoremove -y && rm -rf /var/lib/{apt,dpkg,cache,log} && \
    mkdir /app
    
WORKDIR /app
COPY --from=build /build/fancycard .

ENV PATH=$PATH:/headless-shell

EXPOSE 8080
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./fancycard"]