# Native Provider

## Introduction

This article provides users with the instructions to create and launch a edge cloud on a virtual machine(VM) or edge device(nvidia Jetson or Raspberry Pi), and to add nodes for an existing cloud on edge device. In addition, this article provides guidance of advanced usages of running cloud on VM, such as setting up private registry, and enabling UI components.

## Prerequisites

### Operating System on VM

You will need a VM that is capable of running popular Linux distributions such as **Ubuntu, Debian and Raspbian**, and register or set `SSH key/password` for them.

### Setting up Security Group

The VM needs to apply the following **minimum** Security Group Rules:

<details>

```bash
Rule        Protocol    Port      Source             Description
InBound     TCP         22        ALL                SSH Connect Port
InBound     TCP         6443      K3s agent nodes    Kubernetes API
InBound     TCP         10250     K3s server & agent Kubelet
InBound     TCP         8999      K3s dashboard      (Optional) Required only for Dashboard UI
InBound     UDP         8472      K3s server & agent (Optional) Required only for Flannel VXLAN
InBound     TCP         2379,2380 K3s server nodes   (Optional) Required only for embedded ETCD
OutBound    ALL         ALL       ALL                Allow All
```

</details>

## Creating a edge cloud

Please use `ecm create` command to create a cluster in your VM.

### Normal Cluster

The following command creates a edge cloud named "mycloud", and assign it with 2 master nodes and 2 worker nodes:

```bash
ecm -d create \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --master-ips <master-ip-1,master-ip-2> \
    --worker-ips <worker-ip-1,worker-ip-2>
```

### HA Cloud

Please use one of the following commands to create an HA cloud.

#### Embedded etcd

The following command creates an HA edge cloud named "mycloud", and assigns it with 3 master nodes.

```bash
ecm -d create \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --master-ips <master-ip-1,master-ip-2,master-ip-3> \
    --cluster
```

#### External Database

The following requirements must be met before creating an HA edge cloud with an external database:

- The number of master nodes in this cluster must be greater or equal to 1.
- The external database information must be specified within `--datastore "PATH"` parameter.

In the example below, `--master-ips <master-ip-1,master-ip-2>` specifies the number of master nodes to be 2, `--datastore "PATH"` specifies the external database information. As a result, requirements listed above are met.

Run the command below and create an HA K3s cluster with an external database:

```bash
ecm -d create \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --master-ips <master-ip-1,master-ip-2> \
    --datastore "mysql://<user>:<password>@tcp(<ip>:<port>)/<db>"
```

## Join edge cloud Nodes

Please use `ecm join` command to add one or more nodes for an existing edge cloud.

### Normal edge cloud

The command below shows how to add 2 worker nodes for an existing edge cloud named "mycloud".

```bash
ecm -d join \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --worker-ips <worker-ip-2,worker-ip-3>
```

If you want to join a worker node to an existing edge cloud which is not handled by ecm, please use the following command.

> PS: The existing cluster is not handled by ecm, so it's better to use the same ssh connect information for both master node and worker node so that we can access both VM with the same ssh config.

```bash
ecm -d join \
    --name mycloud \
    --ip <master-ip> \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --worker-ips <worker-ip>
```

### HA cloud

The commands to add one or more nodes for an existing HAedge cloud varies based on the types of HA cluster. Please choose one of the following commands to run.

#### Embedded etcd

Run the command below, to add 2 master nodes for an Embedded etcd HA cluster(embedded etcd: >= 1.19.1-k3s1).

```bash
ecm -d join \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --master-ips <master-ip-2,master-ip-3>
```

#### External Database

Run the command below, to add 2 master nodes for an HA edge cloud with external database, you will need to fill in `--datastore "PATH"` as well.

```bash
ecm -d join \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --master-ips <master-ip-2,master-ip-3> \
    --datastore "mysql://<user>:<password>@tcp(<ip>:<port>)/<db>"
```

## Delete edge cloud

This command will delete a edge cloud named "mycloud".

```bash
ecm -d delete --provider native --name mycloud
```

> PS: If the cluster is an existing K3s cluster which is not handled by ecm, we won't uninstall it when delete edge cloud from ecm.

## List edge clouds 

This command will list the clusters that you have created on this machine.

```bash
ecm list
```

```bash
   NAME     REGION  PROVIDER  STATUS   MASTERS  WORKERS    VERSION
  mycloud             native    Running  1        0        v1.22.6+k3s1
```

## Describe edge cloud

This command will show detail information of a specified cluster, such as instance status, node IP, kubelet version, etc.

```bash
ecm  describe -n <clusterName> 
```

> Note：There will be multiple results if using the same name to create with different providers, please use `-p <provider>` to choose a specified cluster. i.e. `ecm describe cluster mycloud`

```bash
Name: mycloud
Provider: native
Region:
Zone:
Master: 1
Worker: 0
Status: Running
Version: v1.22.6+k3s1
Nodes:
  - internal-ip: [x.x.x.x]
    external-ip: [x.x.x.x]
    instance-status: -
    instance-id: xxxxxxxxxx
    roles: control-plane,master
    status: Ready
    hostname: test
    container-runtime: containerd://1.5.9-k3s1
    version: v1.22.6+k3s1
```

## Access edge cloud

After the cluster is created, `ecm` will automatically merge the `kubeconfig` so that you can access the cluster.

```bash
ecm kubectl config use-context mycloud
ecm kubectl <sub-commands> <flags>
```

In the scenario of multiple clusters, the access to different clusters can be completed by switching context.

```bash
ecm kubectl config get-contexts
ecm kubectl config use-context <context>
```

## SSH edge cloud's Node

Login to a specific k3s cluster node via ssh, i.e. mycloud.

```bash
ecm ssh --name mycloud
```

> If the edge cloud is an existing one which is not handled by ecm, you can't use Execute Shell from UI, but you can access the cluster nodes via CLI.

If the ssh config is different between the existing nodes and current nodes(joined with ecm), you can use the command below to switch the ssh config

```bash
ecm ssh --name mycloud <ip> --ssh-user ubuntu --ssh-key-path ~/.ssh/id_rsa
```

## Other Usages

More usage details please running `ecm <sub-command> --provider native --help` commands.

## Advanced Usages

We integrate some advanced components such as private registries and UI related to the current provider.

### Setting up Private Registry

When running `ecm create` or `ecm join` command, it takes effect with the`--registry /etc/ecm/registries.yaml` flag, i.e.:

```bash
ecm -d create \
    --name mycloud \
    --ssh-user <ssh-user> \
    --ssh-key-path <ssh-key-path> \
    --master-ips <master-ip-1,master-ip-2> \
    --worker-ips <worker-ip-1,worker-ip-2> \
    --registry /etc/ecm/registries.yaml
```

Below are examples showing how you may configure `/etc/ecm/registries.yaml` on your current node when using TLS, and make it take effect on edge cloud by `ecm`.

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

### Enable UI Component

ecm support UI Component, including [kubernetes/dashboard](https://github.com/kubernetes/dashboard) .

#### Enable Kubernetes dashboard

You can enable Kubernetes dashboard using following command.

```bash
ecm -d create \
    ... \
    --enable dashboard
```
If you want to create user token to access dashboard, please following this [docs](https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md).
