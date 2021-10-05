# kubectl-check

![example](preview.png)

A kubectl plugin that checks your clusters for component compatibility and Kubernetes version end of life. This plugin is meant to assist with the upgrade process for these needs:

- Assessing what components need to be upgraded before upgrading to a specific Kubernetes version
- Assessing how much time is available before the Kubernetes version is sunsetted in different platforms like EKS

This cluster is _not_ responsible for

- Finding deprecated APIs, other tools can do this well, like `kubent`

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

This command will check the end of life date for a version using `endoflife.data` and print how much time is left as well as the end of life data. The default option is to check upstream Kubernetes only:

```shell
❯ kubectl check endoflife
TYPE         VERSION   EOL DATE     DAYS LEFT
Kubernetes   1.19      2021-10-28   29
```

To change what product is evaluated, use the `--product` flag

```shell
❯ kubectl check endoflife --product amazon-eks
TYPE         VERSION   EOL DATE     DAYS LEFT
EKS          1.19      2022-04-01   184
```

This command also will exit `0` or `1` depending on the flags set. Exiting `1` means the cluster is within a threshold that the user deems expired. Here are some examples below

```shell
❯ kubectl check endoflife --product amazon-eks --expiry-range 200
TYPE         VERSION   EOL DATE     DAYS LEFT
amazon-eks   1.19      2022-04-01   177
exit status 1

❯ kubectl check endoflife --product amazon-eks --expiry-range 150
TYPE         VERSION   EOL DATE     DAYS LEFT
amazon-eks   1.19      2022-04-01   177
```

In the first example we have set anything less than 200 days to be expired and exit with a `1`. In the second example, we have shortened this timeframe and the command exits with a `0`. This becomes more valuable when we start to script with this command like below

```shell
if ! kubectl check endoflife --product amazon-eks --silent --expiry-range 30 ; then 
    echo "starting upgrade now..."
fi
```

### kubectl check versions

This will check the component versions and what version they require to upgrade to
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

---

## Version Config

This configuration is highly configurable, here is an example of a simple configuration that has version compatbility for `coredns` across multiple Kubernetes releases.

```json
{
    "v1.19": [
        {
            "Name": "coredns",
            "Version": "1.7.0"
        }
    ],
    "v1.20": [
        {
            "Name": "coredns",
            "Version": "1.7.0"
        }
    ],
    "v1.21": [
        {
            "Name": "coredns",
            "Version": "1.8.0"
        }
    ],
    "v1.22": [
        {
            "Name": "coredns",
            "Version": "1.8.4"
        }
    ]
}
```