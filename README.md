# kubectl-tks

A kubectl plugin for working with Twilio flavored clusters.

## Commands

### kubectl tks check versions

This will check the core service versions and what version they require to upgrade to
another version of Kubernetes

```
❯ kubectl tks check versions --kube-version v1.19
SERVICE              OUT OF DATE   CURRENT VERSION      REQUIRED VERSION
coredns              true          v1.6.6               1.8.4
kube-state-metrics   false         v2.1.0               2.1.0
metrics-server       true          v0.3.6               0.5.0
kube-proxy           true          v1.17.9-eksbuild.1   1.19.6-eksbuild.2
```

### kubectl tks get nodes

This will return the `nodes` in a cluster but organized by role and also with SID, instance-id
and instances-type added to the output.

```
❯ kubectl tks get nodes
ROLE             NAME                            SID                                  INSTANCE ID           INSTANCE TYPE   STATUS
anchore-engine   ip-172-24-1-211.ec2.internal    HOe7ab5047bca0edcefab63aa8abb48668   i-0bccdc2726708b734   c5.4xlarge      Ready
anchore-engine   ip-172-24-1-66.ec2.internal     HO57a9110a666baa55a3f450805d908cac   i-035823b7e00fe2181   c5.4xlarge      Ready
anchore-engine   ip-172-24-26-2.ec2.internal     HOc8361ef31ae5df77905d463d3325e4b8   i-00dcf621752e0c6e6   c5.4xlarge      Ready
anchore-engine   ip-172-24-30-207.ec2.internal   HOb95dd70c4019169c14750914dcbcc682   i-0706c9d8d345e49be   c5.4xlarge      Ready
```