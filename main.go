package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yangjinheng/huaweiyun-inventory/pkg/config"
	"github.com/yangjinheng/huaweiyun-inventory/pkg/instance"
	"github.com/yangjinheng/huaweiyun-inventory/pkg/output"
)

// 为 init 设定初始变量
var (
	Configuration config.Config
	Command       string
)

// 利用 init 加载解析配置文件
func init() {
	// 读取命令行参数
	if len(os.Args) < 2 {
		Command = ""
	} else {
		Command = os.Args[1]
	}

	// 读取配置文件
	configuration, err := config.LoadParamer("./config.yaml")
	Configuration = configuration
	if err != nil {
		fmt.Println("读取配置文件错误", err)
		os.Exit(1)
	}
}

// 单元测试
func testing() {
	// // 创建主机，并且等待 JOB 的完成，返回主机信息
	// server1 := config.Host{
	// 	HostName:       "jh-operator",
	// 	AdminPass:      "qloud@123",
	// 	System:         "CentOS 7.6 64bit",
	// 	FlavorRef:      "s2.large.2",
	// 	Subnetname:     "vpc-omd-bak",
	// 	RootVolume:     config.RootVolume{Size: 80, Volumetype: "SAS"},
	// 	DataVolume:    []config.DataVolume{config.DataVolume{Size: 100, Volumetype: "SAS"}},
	// 	SecurityGroups: []string{"Sys-default", "nomad-cluster", "kube"},
	// 	ServerTags: []config.ServerTags{
	// 		config.ServerTags{Key: "Key1", Value: "Value2"},
	// 		config.ServerTags{Key: "Key2", Value: "Value2"},
	// 	},
	// }
	// job, err := instance.CreateServer(server1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(job)
	// query, _ := instance.WaitJobSucess(job)
	// fmt.Println(query)

	// 根据字符串在线查询子网 ID，带运行时缓存
	// subnet, err := instance.GetSubnet("vpc-omd-bak")
	// fmt.Println(subnet.ID, err)

	// 根据字符串在线查询安全组 ID，带运行时缓存
	// x, err := instance.GetSecurityGroups("nomad-cluster")
	// fmt.Println(x.ID, err)

	// 根据镜像名称和规格型号在线查询镜像 ID，带运行时缓存
	// image, err := instance.GetImage("CentOS 7.4 64bit", "s2.large.2")
	// fmt.Println(image.ID, err)

	// 根据服务器 ID 在线查询服务器详细信息
	// server, err := instance.ShowServer("4d59e03e-be4c-499a-96b2-98d501ae47f5")
	// fmt.Println(server, err)

	// 模拟创建主机返回的数据，然后根据配置文件和创建主机的返回结果组装 Inventory
	test := `[{"server":{"fault":null,"id":"5c023952-214f-48d3-b2ea-5105ac42c4c5","name":"instance2","addresses":{"338b5e89-e2ee-4d58-9dca-e792e7c49fc7":[{"version":"4","addr":"10.8.2.13","OS-EXT-IPS-MAC:mac_addr":"fa:16:3e:91:da:51","OS-EXT-IPS:type":"fixed","OS-EXT-IPS:port_id":"a61d5dc3-dcdf-437e-98d1-7e44cf0fbc26"}]},"flavor":{"disk":"0","vcpus":"2","ram":"4096","id":"s2.large.2","name":"s2.large.2"},"accessIPv4":"","accessIPv6":"","status":"ACTIVE","progress":0,"hostId":"b5abb43c3e917bfe7cae9a5696f6099114f83af5784b58c324331cc7","updated":"2020-01-06T13:48:33Z","created":"2020-01-06T13:47:48Z","metadata":{"metering.image_id":"3a64bd37-955e-40cd-ab9e-129db56bc05d","metering.imagetype":"gold","metering.resourcespeccode":"s2.large.2.linux","image_name":"CentOS 7.6 64bit","os_bit":"64","cascaded.instance_extrainfo":"pcibridge:1","metering.resourcetype":"1","vpc_id":"338b5e89-e2ee-4d58-9dca-e792e7c49fc7","os_type":"Linux","charging_mode":"0"},"tags":["key1=value1","key2=value2"],"description":"","locked":false,"serverNovaLock":false,"config_drive":"","tenant_id":"6de764f622104bae93d5bbb45c6e2ec9","user_id":"051d10718a0026961fd5c011adf485b5","key_name":null,"os-extended-volumes:volumes_attached":[{"device":"/dev/vda","bootIndex":"0","id":"9bd883f0-57cd-4c7f-9992-8939ce19f1f7","delete_on_termination":"true"},{"device":"/dev/vdb","bootIndex":null,"id":"f056eea2-c0e3-4d7b-bc62-2021d36b7a30","delete_on_termination":"false"}],"OS-EXT-STS:task_state":null,"OS-EXT-STS:power_state":1,"OS-EXT-STS:vm_state":"active","OS-EXT-SRV-ATTR:host":"pod04.cnnorth1a","OS-EXT-SRV-ATTR:instance_name":"instance-005cc6ec","OS-EXT-SRV-ATTR:hypervisor_hostname":"nova010@42","OS-DCF:diskConfig":"MANUAL","OS-EXT-AZ:availability_zone":"cn-north-1a","os:scheduler_hints":{},"OS-EXT-SRV-ATTR:root_device_name":"/dev/vda","OS-EXT-SRV-ATTR:ramdisk_id":"","enterprise_project_id":"0","OS-EXT-SRV-ATTR:user_data":"IyEvYmluL2Jhc2gKZWNobyAncm9vdDpxbG91ZEAxMjMnIHwgY2hwYXNzd2QgOw==","OS-SRV-USG:launched_at":"2020-01-06T13:48:08.000000","OS-EXT-SRV-ATTR:kernel_id":"","OS-EXT-SRV-ATTR:launch_index":0,"host_status":"UP","OS-EXT-SRV-ATTR:reservation_id":"r-hplnzlyn","OS-EXT-SRV-ATTR:hostname":"instance2","OS-SRV-USG:terminated_at":null,"sys_tags":[{"key":"_sys_enterprise_project_id","value":"0"}],"security_groups":[{"id":"2be5680c-a600-43ed-a702-2573b0370744","name":"default"},{"id":"2eee56f1-d206-4a4d-9ac5-23d95d784e09","name":"nomad-cluster"}],"image":{"id":"3a64bd37-955e-40cd-ab9e-129db56bc05d"}}},{"server":{"fault":null,"id":"535fd02e-c854-4f32-82ef-ea3c40b00797","name":"instance1","addresses":{"338b5e89-e2ee-4d58-9dca-e792e7c49fc7":[{"version":"4","addr":"10.8.2.117","OS-EXT-IPS-MAC:mac_addr":"fa:16:3e:92:4a:ed","OS-EXT-IPS:type":"fixed","OS-EXT-IPS:port_id":"1698d03c-f925-4918-9a89-26d34e2c08c8"}]},"flavor":{"disk":"0","vcpus":"2","ram":"4096","id":"s2.large.2","name":"s2.large.2"},"accessIPv4":"","accessIPv6":"","status":"ACTIVE","progress":0,"hostId":"b5abb43c3e917bfe7cae9a5696f6099114f83af5784b58c324331cc7","updated":"2020-01-06T13:48:29Z","created":"2020-01-06T13:47:48Z","metadata":{"metering.image_id":"3a64bd37-955e-40cd-ab9e-129db56bc05d","metering.imagetype":"gold","metering.resourcespeccode":"s2.large.2.linux","image_name":"CentOS 7.6 64bit","os_bit":"64","cascaded.instance_extrainfo":"pcibridge:1","metering.resourcetype":"1","vpc_id":"338b5e89-e2ee-4d58-9dca-e792e7c49fc7","os_type":"Linux","charging_mode":"0"},"tags":["key1=value1","key2=value2"],"description":"","locked":false,"serverNovaLock":false,"config_drive":"","tenant_id":"6de764f622104bae93d5bbb45c6e2ec9","user_id":"051d10718a0026961fd5c011adf485b5","key_name":null,"os-extended-volumes:volumes_attached":[{"device":"/dev/vda","bootIndex":"0","id":"b8479df6-48ac-41b3-b42c-c4cc05173e3b","delete_on_termination":"true"},{"device":"/dev/vdb","bootIndex":null,"id":"ab7201dc-b040-4d57-96a6-7593051f16f6","delete_on_termination":"false"}],"OS-EXT-STS:task_state":null,"OS-EXT-STS:power_state":1,"OS-EXT-STS:vm_state":"active","OS-EXT-SRV-ATTR:host":"pod04.cnnorth1a","OS-EXT-SRV-ATTR:instance_name":"instance-005cc6eb","OS-EXT-SRV-ATTR:hypervisor_hostname":"nova010@4","OS-DCF:diskConfig":"MANUAL","OS-EXT-AZ:availability_zone":"cn-north-1a","os:scheduler_hints":{},"OS-EXT-SRV-ATTR:root_device_name":"/dev/vda","OS-EXT-SRV-ATTR:ramdisk_id":"","enterprise_project_id":"0","OS-EXT-SRV-ATTR:user_data":"IyEvYmluL2Jhc2gKZWNobyAncm9vdDpxbG91ZEAxMjMnIHwgY2hwYXNzd2QgOw==","OS-SRV-USG:launched_at":"2020-01-06T13:48:07.000000","OS-EXT-SRV-ATTR:kernel_id":"","OS-EXT-SRV-ATTR:launch_index":0,"host_status":"UP","OS-EXT-SRV-ATTR:reservation_id":"r-gage1q4p","OS-EXT-SRV-ATTR:hostname":"instance1","OS-SRV-USG:terminated_at":null,"sys_tags":[{"key":"_sys_enterprise_project_id","value":"0"}],"security_groups":[{"id":"2be5680c-a600-43ed-a702-2573b0370744","name":"default"},{"id":"2eee56f1-d206-4a4d-9ac5-23d95d784e09","name":"nomad-cluster"}],"image":{"id":"3a64bd37-955e-40cd-ab9e-129db56bc05d"}}}]`
	sucesshosts := []instance.ServerElement{}
	json.Unmarshal([]byte(test), &sucesshosts)
	inv := output.GenInventory(Configuration.Instances, sucesshosts)
	hosts, _ := json.Marshal(inv)
	fmt.Println(string(hosts))
}

func main() {
	// 单元测试
	// testing()

	// 命令行传参为 Inventory 则根据配置文件创建主机，完成后返回 Inventory
	if Command == "--list" {
		sucesshosts := instance.QuicklyCreate(Configuration.Instances)
		inv := output.GenInventory(Configuration.Instances, sucesshosts)
		hosts, _ := json.Marshal(inv)
		fmt.Println(string(hosts))
	}

	if Command == "" || Command == "--help" || Command == "-h" {
		fmt.Println("Usage:", os.Args[0], "[apply|--list]")
		fmt.Println("\tapply: 创建主机并返回主机名和 IP 地址")
		fmt.Println("\t--list: 创建主机并返回 Ansible Inventory")
	}

	// 命令行传参为 apply 则根据配置文件创建主机，完成后返回裸格式
	if Command == "apply" {
		fmt.Println("apply")
	}
}

// 1. 将创建结果保存在本地作为缓存，如果本地状态文件的修改时间没有超出 30 分钟则直接解析为 Inventory 返回，否则在线查询
// 2. 提供一个根据状态文件销毁实例的功能，根据状态文件对比，变更配置？
// 3. 提供一个关机降配的功能，降低消费
// 4. 根据组在线查询所有符合条件的主机，返回 Inventory
