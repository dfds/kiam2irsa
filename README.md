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

## Requirements

A kubeconfig file that exists in `~/.kube` directory. By default it will attempt to read the `~/.kube/config` file.
If you want it to read a different kubeconfig file, specify the filename with full path as shown below.

## Usage

```bash
go run main.go ~/.kube/my-cluster.config
```
