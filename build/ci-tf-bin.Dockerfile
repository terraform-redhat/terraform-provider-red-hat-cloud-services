FROM registry.access.redhat.com/ubi8/ubi-minimal:latest AS builder
WORKDIR /root

# Update the base image and install necessary packages
RUN microdnf update -y && \
    microdnf install -y git make go-toolset && \
    microdnf clean all

# terraform-provider-red-hat-cloud-services repo
COPY . ./terraform-provider-red-hat-cloud-services


RUN cd terraform-provider-red-hat-cloud-services && go mod tidy && go mod vendor && make build &&\
    echo 'RUN done'

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /root
COPY --from=builder /root/terraform-provider-red-hat-cloud-services* /root/
