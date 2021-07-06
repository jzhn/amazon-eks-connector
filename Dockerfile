FROM public.ecr.aws/amazonlinux/amazonlinux:2

COPY dist/AWSWesleyExternalClusterConnector_linux_amd64/eks-connector /var/eks/connector

ENTRYPOINT ["/var/eks/connector"]