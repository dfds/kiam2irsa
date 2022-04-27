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

### Getting CLI help

```bash
./kiam2irsa --help

./kiam2irsa sa --help
```

### With default kubeconfig ~/.kube/config

```bash
./kiam2irsa sa
```

### With custom kubeconfig through environment variable

```bash
KUBECONFIG=~/.kube/my-cluster.config
./kiam2irsa sa
```

### With kubeconfig through argument passing

```bash
./kiam2irsa sa --kubeconfig ~/.kube/my-cluster.config
```

## Build instructions

```bash
go build .
```
