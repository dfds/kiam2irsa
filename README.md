# :warning: Repository not maintained :warning:

Please note that this repository is currently archived, and is no longer being maintained.

- It may contain code, or reference dependencies, with known vulnerabilities
- It may contain out-dated advice, how-to's or other forms of documentation
- 
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
./kiam2irsa pods --help
```

### Find ServiceAccount status using default kubeconfig ~/.kube/config

```bash
./kiam2irsa sa
```

### Find ServiceAccount status using custom kubeconfig through environment variable

```bash
KUBECONFIG=~/.kube/my-cluster.config
./kiam2irsa sa
```

### Find ServiceAccount status using kubeconfig through argument passing

```bash
./kiam2irsa sa --kubeconfig ~/.kube/my-cluster.config
```

### Find pods only using KIAM
```bash
./kiam2irsa pods --status KIAM
```

### Find pods only migrated to IRSA, but that still supports KIAM
```bash
./kiam2irsa pods --status BOTH
```

### Find pods fully migrated to IRSA
```bash
./kiam2irsa pods --status IRSA
```

## Build instructions

```bash
go build .
```
