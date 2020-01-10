package output

import (
	"strings"

	"github.com/yangjinheng/huaweiyun-inventory/pkg/config"
	"github.com/yangjinheng/huaweiyun-inventory/pkg/instance"
)

// Inventory 表示 Inventory 记录
type Inventory map[string]Inner

// Inner 利用 omitempty 特性描述内部三种数据结构
type Inner struct {
	Children *[]string                          `json:"children,omitempty"`
	Hosts    *[]string                          `json:"hosts,omitempty"`
	Hostvars *map[string]map[string]interface{} `json:"hostvars,omitempty"`
}

// InSlice 是 Slice 成员检查函数
func InSlice(s string, array []string) bool {
	for _, item := range array {
		if item == s {
			return true
		}
	}
	return false
}

// GenInventory 申请生成 Inventory
func GenInventory(HostsRequest []config.Host, HostsResponse []instance.ServerElement) Inventory {
	// 创建一个 Inventory 对象
	inventory := Inventory{"all": Inner{Children: &[]string{}}, "_meta": Inner{Hostvars: &map[string]map[string]interface{}{}}}

	// 遍历申请主机列表
	for _, hostrequest := range HostsRequest {
		// 得到申请主机时候用的子网，有缓存秒得到
		subnet, _ := instance.GetSubnet(hostrequest.Subnetname)

		// 遍历子网信息
		for _, host := range HostsResponse {
			// 找到申请的主机的网卡们
			if host.Server.Name == hostrequest.HostName {
				netcard := host.Server.Addresses[subnet.VpcID]
				// 从网卡列表中取出这个子网的 IP 地址
				for _, card := range netcard {
					if card.OSEXTIPSType == "fixed" {
						HostName := hostrequest.HostName
						AccessIPV4 := card.Addr

						// 取得 ServerTags 保存的主机组信息、并收集所有变量到 map
						AnsibleGroup, Variables := "", map[string]interface{}{}
						for _, tag := range hostrequest.ServerTags {
							if tag.Key == "ansible_group" {
								AnsibleGroup = tag.Value
							} else {
								Variables[tag.Key] = tag.Value
							}
						}
						Variables["ansible_ssh_host"] = AccessIPV4
						Variables["ansible_ssh_pass"] = hostrequest.AdminPass

						// 将收集所有变量放入到 _meta.hostvar 中
						hostvar := *inventory["_meta"].Hostvars
						hostvar[HostName] = Variables

						// 将主机组放入 all.Children 中，将主机放入所属的组的切片内
						for _, group := range strings.Split(AnsibleGroup, " ") {
							// 如果 all.children 没有这个组就追加
							groups := inventory["all"].Children
							if !InSlice(group, *groups) {
								*groups = append(*groups, group)
							}
							// 将主机名放入相应的主机组内
							if (inventory[group] == Inner{}) {
								inventory[group] = Inner{Hosts: &[]string{}}
							}
							hosts := inventory[group].Hosts
							*hosts = append(*hosts, HostName)
						}
					}
				}
			}
		}
	}
	return inventory
}
