# KIAM to IRSA migration check

## Purpose

This is a tool that can be useful when migrating from KIAM to IAM Roles for Service Accounts (IRSA).
The tool will find all Kubernetes service accounts that has the annotation:

```yaml
metadata:
    annotations:
        eks.amazonaws.com/role-arn: "<ANY ARN>"
```

and doesn't have the annotation:

```yaml
metadata:
    annotations:
        eks.amazonaws.com/sts-regional-endpoints: "true"
```

## Usage

```bash
./kiam2irsa sa --kubeconfig ~/.kube/my-cluster.config
```
