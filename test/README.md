# AWSWesleyExternalClusterConnector Testing

## What

This is a __temporary__ doc for setting up EKS connector pod for CCS testing.

## How

### Prepare
See [EKS connector beta onboarding doc](https://quip-amazon.com/lyErAIMd6Se1/How-to-onboard-to-EKS-connector-beta) for steps.

### Deploy

To deploy it we need to create an SSM hybrid activation first.
__For testing, put a high number of activation instance__ so that we don't need to create activation often when SSM
agent restarts.

```bash
# Fill in the activation ID and activation code.
export SSM_ACTIVATION_ID=""
export SSM_ACTIVATION_CODE=""

# Apply the manifest
sed "s~%SSM_ACTIVATION_ID%~$SSM_ACTIVATION_ID~g; s~%SSM_ACTIVATION_CODE%~$(echo -n $SSM_ACTIVATION_CODE | base64)~g" \
    eks-connector.yaml | kubectl apply -f -
# After a few seconds the connector pod should be healthy in kubernetes.

# Now get the managed instance at SSM.
aws ssm describe-instance-information --filters Key=ActivationIds,Values=$SSM_ACTIVATION_ID
# If you are lucky you should see exactly one managed instance.
# Alternatively, grep the logs at init container, which should print out the instance id.

# Now execute non interactive command
# NOTE: fill in TARGET with your own managed instance id like `mi-069f7e4b6ce64c0ce`
aws ssm start-session \
    --target TARGET \
    --document-name AWS-StartNonInteractiveCommand \
    --parameters '{"command": ["curl --unix-socket /var/eks/shared/connector.sock -H \"x-aws-eks-identity-arn: arn:aws:iam::123456789012:user/srajakum\" http://localhost/api/v1/pods"]}'
```

### Cleanup

Just delete the manifest

```bash
sed "s~%SSM_ACTIVATION_ID%~$SSM_ACTIVATION_ID~g; s~%SSM_ACTIVATION_CODE%~$(echo -n $SSM_ACTIVATION_CODE | base64)~g" \
    eks-connector.yaml | kubectl delete -f -
```