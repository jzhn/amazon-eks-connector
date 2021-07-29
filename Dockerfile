FROM public.ecr.aws/amazonlinux/amazonlinux:2 as build

FROM scratch

# copy ca-bundle.crt from AmazonLinux2...
COPY --from=build /etc/ssl/certs/ca-bundle.crt /etc/ssl/certs/
COPY dist/AWSWesleyExternalClusterConnector_linux_amd64/eks-connector /var/eks/connector

ENTRYPOINT ["/var/eks/connector"]