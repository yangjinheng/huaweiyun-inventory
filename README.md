## 主要作用

从配置文件中读取主机的配置列表，调用华为云的 API 接口，使用 Goroutine 完成主机的申请，根据返回的主机信息列表为 Ansible 提供一个 Inventory 源。

*   应用

1.  从申请主机到使用 ansible-playbook 部署服务，全自动化，可用在在线扩容服务的场景

## 文件组织

~~~bash
.
├── README.md                         # 本文档
├── example.config.yaml               # 示例配置文件，需要改名为 config.yaml 才能被识别
├── go.mod                            # go mod
├── go.sum                            # go mod
├── inventory.json                    # inventory 示例输出
├── main.go                           # main 入口
└── pkg
    ├── config
    │   └── config.go                 # 解析 config.yaml
    ├── instance
    │   ├── images.go                 # 根据镜像名称在线查询镜像 ID
    │   ├── instance.go               # 创建主机和并发创建
    │   ├── network.go                # 根据子网名称在线查询子网 ID
    │   └── queryjob.go               # 阻塞等待创建主机 JOB 完成
    ├── output
    │   ├── inventory.go              # 华为云返回的主机信息转换为 inventory
    │   └── row.go                    # 未完成
    └── signer
        ├── escape.go
        └── signer.go                 # 请求签名工具，在 config 中初始化
~~~

## 元信息配置

~~~yaml
- hostname: node2
  adminpass: 123456
  system: CentOS 7.6 64bit
  flavorref: s2.large.2
  subnetname: VPC 下具体的子网名称
  rootvolume: { size: 80, volumetype: SAS }
  datavolume:
    - { size: 100, volumetype: SAS }
  securitgroups: ["Sys-default", "nomad-cluster"]
  servertags:
    - { key: "ansible_group", value: "kubernetes ceph registry" }   // value 使用空格分隔该主机所属的 Ansible 组，除了这个键其他 tag 都被视作主机的环境变量
    - { key: "ansible_ssh_port", value: "22" }                      // 环境变量
    - { key: "prodction", value: "true" }                           // 环境变量
~~~

## 使用方法

1.  将 example.config.yaml 改名为 config.yaml，修改这个文件，填入华为云账号信息
2.  修改 config.yaml 中的主机申请信息，注意 servertags 它决定了主机的组和主机的变量信息，这一步需要首先理解 Inventory 规范
3.  编译本项目 go build
4.  执行 ansible-playbook -i <编译出来的文件> <playbook.yaml> 完成从申请主机到部署服务的全部过程

