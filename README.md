### 子网IP管理方案

1、利用dhcpd服务器

步骤

（1）用户创建网络，网络有网段信息，路由信息，start、end ip段，这些信息当作subnet填入dhcpd的配置，使用固定的地址启动dhcp服务器。

（2）用户创建虚拟机配置网络时，有网络选择信息，网络选择后，用户可以配置该网络网段的静态ip转（3）

或者dhcp方式获取IP转（4）。

（3）生成mac地址(mac地址管理模块)，cloudinit配置vm sriov网卡指定的静态ip，虚拟机配置的主机名(hostname)写入到dhcpd的subnet 的statiIP配置。

（4）生成mac地址(mac地址管理模块)，cloudinit配置vm sriov网卡指定的dhcp，并且指定dhcp服务器地址。



mac地址管理模块采用开源组件kubemacpool（[k8snetworkplumbingwg/kubemacpool (github.com)](https://github.com/k8snetworkplumbingwg/kubemacpool)）。

2、利用k8s webhook加go-ipma

使用k8s webhook管理macaddress的分配（kubemacpool）。

使用go-ipma（[metal-stack/go-ipam:](https://github.com/metal-stack/go-ipam)）管理ip address的分配。

go-ipam(backend redis/etcd)->cloudinit secret->kubevirt vm

运行go-ipam服务。

```
docker run -it -d --name go-ipam  --net host --env GOIPAM_LOG_LEVEL=debug --env GOIPAM_REDIS_HOST=172.28.240.178 --env GOIPAM_REDIS_PORT=6379 ghcr.io/metal-stack/go-ipam:custom redis
```



3、利用k8s crd crd-controller

使用k8s webhook管理macaddress的分配（kubemacpool）。

kube-ipam（[routerd/kube-ipam](https://github.com/routerd/kube-ipam)）,没有ip冲突检测机制，需要根据crd创建后的状态判断IP是否已经分配。

kube-ipam->cloudinit secret->kubevirt vm

方案对比

|             | DHCP                                                 | kube-ipam                  | go-ipam                           |
| ----------- | ---------------------------------------------------- | -------------------------- | --------------------------------- |
| ipv4/v6支持 | 支持                                                 | 支持                       | 支持                              |
| 多子网支持  | 配置复杂易出错                                       | 配置简单                   | 配置简单                          |
| ip冲突检测  | 自动检测                                             | 根据返回状态手动判断       | 自动检测                          |
| ip信息存储  | 不需要，内存保存，服务重启通过广播收集同网段已使用ip | 不需要，crd对象保存        | 需要，后端存储(etcd/redis)保存    |
| api调用方式 | 通过下发命令修改dhcp配置文件                         | 通过调用k8s api创建crd对象 | 通过开发和构建restful组件提供服务 |
| 网络限制    | vm的网络接口必须与dhcp服务监听的网络接口连通         | 无，由业务获取ip填入vm中   | 无，由业务获取ip填入vm中          |

三种方案的部署图

白色框绿色字体的是需要自定义开发的组件。

![子网IP管理架构](.\子网IP管理架构.png)



### api设计

| url   | 操作            | 输入参数                                | 输出参数                                                |
| ----- | --------------- | --------------------------------------- | ------------------------------------------------------- |
| /cidr | get/post/delete | /cidr?cidr=192.168.0.0/16               | {httpstatus:200, cidr: 192.168.0.0/16}                  |
| /ip   | post/delete     | /ip?cidr=192.168.0.0/16&ip=192.168.0.24 | {httpstatus:200, cidr: 192.168.0.0/16,ip: 192.168.0.24} |
|       |                 |                                         |                                                         |

