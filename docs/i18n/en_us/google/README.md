# Google Provider

## Introduction

This article provides users with the instructions to create and launch a K3s cluster on a Google Compute Engine instance, and to add nodes for an existing K3s cluster on GCE. In addition, this article provides guidance of advanced usages of running K3s on GCE, such as setting up private registry, and enabling UI components.

## Prerequisites

To ensure that GCE instances can be created and accessed successfully, please follow the instructions below.

### Setting up Environment

Configure the following environment variables for the host on which you are running `autok3s`.

```bash
export GOOGLE_SERVICE_ACCOUNT_FILE='<service-account-file-path>'
export GOOGLE_SERVICE_ACCOUNT='<service-account-name>'
```

### Setting up Service Account

Please refer [here](https://cloud.google.com/iam/docs/service-accounts?_ga=2.117192902.-1539515613.1620703671) for more Service Account settings.

Please make sure your service account has permission to specified project and compute resource permission.

### Setting up Security Group

The GCE instances need to apply the following **minimum** Security Group Rules:

<details>

```bash
Rule        Protocol    Port      Source             Description
InBound     TCP         22        ALL                SSH Connect Port
InBound     TCP         6443      K3s agent nodes    Kubernetes API
InBound     TCP         10250     K3s server & agent Kubelet
InBound     UDP         8472      K3s server & agent (Optional) Required only for Flannel VXLAN
InBound     TCP         2379,2380 K3s server nodes   (Optional) Required only for embedded ETCD
OutBound    ALL         ALL       ALL                Allow All
```

</details>

## Creating a K3s cluster

Please use `autok3s create` command to create a cluster in your GCE instance.

### Normal Cluster

The following command uses google as cloud provider, creates a K3s cluster named "myk3s", and assign it with 1 master node and 1 worker node:

```bash
autok3s -d create -p google --name myk3s --master 1 --worker 1 --project <your-project>
```

### HA Cluster

Please use one of the following commands to create an HA cluster.

#### Embedded etcd

The following command uses google as cloud provider, creates an HA K3s cluster named "myk3s", and assigns it with 3 master nodes.

```bash
autok3s -d create -p google --name myk3s --master 3 --cluster --project <your-project>
```

#### External Database

The following requirements must be met before creating an HA K3s cluster with an external database:

- The number of master nodes in this cluster must be greater or equal to 1.
- The external database information must be specified within `--datastore "PATH"` parameter.

In the example below, `--master 2` specifies the number of master nodes to be 2, `--datastore "PATH"` specifies the external database information. As a result, requirements listed above are met.

Run the command below and create an HA K3s cluster with an external database:

```bash
autok3s -d create -p google --name myk3s --master 2 --datastore "mysql://<user>:<password>@tcp(<ip>:<port>)/<db>"
```

## Join K3s Nodes

Please use `autok3s join` command to add one or more nodes for an existing K3s cluster.

### Normal Cluster

The command below shows how to add a worker node for an existing K3s cluster named "myk3s".

```bash
autok3s -d join -p google --name myk3s --worker 1
```

### HA Cluster

The commands to add one or more nodes for an existing HA K3s cluster varies based on the types of HA cluster. Please choose one of the following commands to run.

```bash
autok3s -d join -p google --name myk3s --master 2 --worker 1
```

## Delete K3s Cluster

This command will delete a k3s cluster named "myk3s".

```bash
autok3s -d delete -p google --name myk3s
```

## List K3s Clusters

This command will list the clusters that you have created on this machine.

```bash
autok3s list
```

```bash
NAME         REGION      PROVIDER  STATUS   MASTERS  WORKERS    VERSION
myk3s    asia-northeast1  google   Running  1        0        v1.20.2+k3s1
```

## Describe k3s cluster

This command will show detail information of a specified cluster, such as instance status, node IP, kubelet version, etc.

```bash
autok3s describe -n <clusterName> -p google
```

> Note：There will be multiple results if using the same name to create with different providers, please use `-p <provider>` to choose a specified cluster. i.e. `autok3s describe cluster myk3s -p google`

```bash
Name: myk3s
Provider: google
Region: asia-northeast1
Zone: asia-northeast1-b
Master: 1
Worker: 0
Status: Running
Version: v1.20.2+k3s1
Nodes:
  - internal-ip: [x.x.x.x]
    external-ip: [x.x.x.x]
    instance-status: RUNNING
    instance-id: xxxxxxxx
    roles: control-plane,master
    status: Ready
    hostname: xxxxxxxx
    container-runtime: containerd://1.4.3-k3s1
    version: v1.20.2+k3s1
```

## Access K3s Cluster

After the cluster is created, `autok3s` will automatically merge the `kubeconfig` so that you can access the cluster.

```bash
autok3s kubectl config use-context myk3s.asia-northeast1.google
autok3s kubectl <sub-commands> <flags>
```

In the scenario of multiple clusters, the access to different clusters can be completed by switching context.

```bash
autok3s kubectl config get-contexts
autok3s kubectl config use-context <context>
```

## SSH K3s Cluster's Node

Login to a specific k3s cluster node via ssh, i.e. myk3s.

```bash
autok3s ssh --provider google --name myk3s
```

## Other Usages

More usage details please running `autok3s <sub-command> --provider google --help` commands.

## Advanced Usages

We integrate some advanced components such as private registries and UI, related to the current provider.

### Setting up Private Registry

Below are examples showing how you may configure `/etc/autok3s/registries.yaml` on your current node when using TLS, and make it take effect on k3s cluster by `autok3s`.

```bash
mirrors:
  docker.io:
    endpoint:
      - "https://mycustomreg.com:5000"
configs:
  "mycustomreg:5000":
    auth:
      username: xxxxxx # this is the registry username
      password: xxxxxx # this is the registry password
    tls:
      cert_file: # path to the cert file used in the registry
      key_file:  # path to the key file used in the registry
      ca_file:   # path to the ca file used in the registry
```

When running `autok3s create` or `autok3s join` command, it will take effect with the`--registry /etc/autok3s/registries.yaml` flag, i.e:

```bash
autok3s -d create \
    --provider google \
    --name myk3s \
    --master 1 \
    --worker 1 \
    --registry /etc/autok3s/registries.yaml
```

### Enabling GCP Cloud Controller Manager(CCM)

Will enable [gcp-cloud-provider](https://github.com/kubernetes/cloud-provider-gcp) for K3s

```bash
autok3s -d create -p google \
    ... \
    --cloud-controller-manager 
```

### Enable UI Component

AutoK3s support 2 kinds of UI Component, including [kubernetes/dashboard](https://github.com/kubernetes/dashboard) and [cnrancher/kube-explorer](https://github.com/cnrancher/kube-explorer).

#### Enable Kubernetes dashboard

You can enable Kubernetes dashboard using following command.

```bash
autok3s -d create -p google \
    ... \
    --enable dashboard
```
If you want to create user token to access dashboard, please following this [docs](https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md).

#### Enable kube-explorer dashboard

You can enable kube-explorer using following command.

```bash
autok3s explorer --context myk3s.asia-northeast1.google --port 9999
```

You can enable kube-explorer when creating K3s Cluster by UI.

![](../../../assets/enable-kube-explorer-by-create-cluster.png)

You can also enable/disable kube-explorer any time from UI, and access kube-explorer dashboard by `dashboard` button.

![](../../../assets/enable-kube-explorer-by-button.png)