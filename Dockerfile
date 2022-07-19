FROM node:14.18.1 as stage-web-build

LABEL stage=stage-web-build

WORKDIR /build/kubepi/web

COPY . .

RUN make build_web

RUN rm -fr web

FROM golang:1.16 as stage-bin-build

ENV GOPROXY="https://goproxy.cn,direct"

ENV CGO_ENABLED=0

ENV GO111MODULE=on

LABEL stage=stage-bin-build

WORKDIR /build/kubepi/bin

COPY --from=stage-web-build /build/kubepi/web .

RUN go mod download

RUN make build_gotty
RUN make build_bin

FROM alpine:3.16

WORKDIR /

COPY --from=stage-bin-build /build/kubepi/bin/dist/usr /usr

RUN ARCH=$(uname -m) \
    && case $ARCH in aarch64) ARCH="arm64";; x86_64) ARCH="amd64";; esac \
    && echo "ARCH: " $ARCH \
    && apk add --update --no-cache bash bash-completion curl wget openssl iputils busybox-extras vim \
    && sed -i "s/nobody:\//nobody:\/nonexistent/g" /etc/passwd \
    && curl -sLf https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/kubectl/v1.22.1/${ARCH}/kubectl > /usr/bin/kubectl \
    && chmod +x /usr/bin/kubectl \
    && cd /opt/ \
    && wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/kubectl-aliases/kubectl-aliases.tar.gz \
    && tar zxvf kubectl-aliases.tar.gz \
    && rm -rf kubectl-aliases.tar.gz \
    && chmod -R 755 kubectl-aliases \
    && wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/fzf/0.21.0/fzf.tar.gz \
    && tar zxvf fzf.tar.gz \
    && rm -rf fzf.tar.gz \
    && chmod -R 755 fzf \
    && yes | fzf/install \
    && ln -s fzf/bin/fzf /usr/local/bin/fzf \
    && cd /tmp/ \
    && wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/k9s/v0.24.14/k9s_Linux_${ARCH}.tar.gz \
    && tar -xvf k9s_Linux_${ARCH}.tar.gz \
    && chmod +x k9s \
    && mv k9s /usr/bin \
    && KUBECTX_VERSION=v0.9.4 \
    && wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/kubens/${KUBECTX_VERSION}/kubens_${KUBECTX_VERSION}_linux_${ARCH}.tar.gz \
    && tar -xvf kubens_${KUBECTX_VERSION}_linux_${ARCH}.tar.gz \
    && chmod +x kubens \
    && mv kubens /usr/bin \
    && wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/kubectx/${KUBECTX_VERSION}/kubectx_${KUBECTX_VERSION}_linux_${ARCH}.tar.gz \
    && tar -xvf kubectx_${KUBECTX_VERSION}_linux_${ARCH}.tar.gz \
    && chmod +x kubectx \
    && mv kubectx /usr/bin \
    && curl -L https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/get-helm-3 | bash \
    && chmod +x /usr/local/bin/gotty \
    && chmod 555 /bin/busybox \
    && rm -rf /tmp/* /var/tmp/* /var/cache/apk/* \
    && chmod -R 755 /tmp \
    && mkdir -p /opt/webkubectl

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY vimrc.local /etc/vim

EXPOSE 80

USER root

CMD ["kubepi-server"]
