# kubectl-check

![example](preview.png)

A kubectl plugin for working with Twilio flavored clusters.

## Quickstart

Building and installing this plugin can be done by cloning this repo and running `make install` like this

```
make install
```

Once installed you can now run 

```
kubectl check
```

## Commands

### kubectl check endoflife

This command will check the end of life date for a version using `endoflife.data`. The default option is to check upstream Kubernetes first:

```shell
❯ kubectl check endoflife
TYPE         VERSION   EOL DATE     DAYS LEFT
Kubernetes   1.19      2021-10-28   29
```

To add EKS output, you can add an `--eks` flag

```shell
❯ go run main.go endoflife --eks
TYPE         VERSION   EOL DATE     DAYS LEFT
EKS          1.19      2022-04-01   184
Kubernetes   1.19      2021-10-28   29
```

### kubectl check versions

This will check the core service versions and what version they require to upgrade to
another version of Kubernetes, this uses the config located at `cmd/config` and can support
multiple Kubernetes versions

```
❯ kubectl check versions --kube-version v1.19
SERVICE                 OUT OF DATE   CURRENT VERSION      REQUIRED VERSION
cluster-autoscaler      false         v1.19.1              1.19.0
coredns                 false         latest               1.8.4
kube-state-metrics      false         v2.1.0               2.1.0
metrics-server          false         v0.5.0               0.5.0
kube-proxy              false         v1.19.6-eksbuild.2   1.19.6-eksbuild.2
node-problem-detector   false         v0.8.9               0.8.9
cert-manager            false         v1.4.0               1.4.0
```