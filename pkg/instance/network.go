package instance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yangjinheng/huaweiyun-inventory/pkg/config"
	"io/ioutil"
	"net/http"
)

// Subnets Securitygroups 表示子网信息和安全组信息的运行时缓存
var (
	Subnets        SubnetsResult
	Securitygroups SecuritygroupsResult
)

// SubnetsResult 访问子网接口时候返回的对象
type SubnetsResult struct {
	Subnets []Subnet `json:"subnets"`
}

// Subnet 单个子网对象
type Subnet struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	CIDR             string          `json:"cidr"`
	DNSList          []string        `json:"dnsList"`
	Status           string          `json:"status"`
	VpcID            string          `json:"vpc_id"`
	Ipv6Enable       bool            `json:"ipv6_enable"`
	GatewayIP        string          `json:"gateway_ip"`
	DHCPEnable       bool            `json:"dhcp_enable"`
	PrimaryDNS       string          `json:"primary_dns"`
	SecondaryDNS     string          `json:"secondary_dns"`
	AvailabilityZone string          `json:"availability_zone"`
	NeutronNetworkID string          `json:"neutron_network_id"`
	NeutronSubnetID  string          `json:"neutron_subnet_id"`
	ExtraDHCPOpts    []ExtraDHCPOpts `json:"extra_dhcp_opts"`
}

// ExtraDHCPOpts DHCP 选项
type ExtraDHCPOpts struct {
	Optname  string `json:"opt_name"`
	Optvalue string `json:"opt_value"`
}

// GetSubnet 根据子网名称返回子网 ID
func GetSubnet(name string) (Subnet, error) {
	for _, subnet := range Subnets.Subnets {
		if subnet.Name == name {
			// fmt.Println("从缓存中读取子网ID")
			return subnet, nil
		}
	}
	request, _ := http.NewRequest("GET", "https://vpc." + config.Region +".myhuaweicloud.com/v1/" + config.ProjectID + "/subnets", bytes.NewBuffer([]byte("")))
	request.Header.Add("content-type", "application/json;charset=utf8")
	request.Header.Add("X-Project-Id", config.ProjectID)
	config.Signature.Sign(request)
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return Subnet{}, fmt.Errorf("查询子网接口访问失败 %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Subnet{}, fmt.Errorf("查询子网时没有得到响应 %s", err)
	}
	json.Unmarshal(body, &Subnets)
	for _, subnet := range Subnets.Subnets {
		if subnet.Name == name {
			// fmt.Println("从在线中读取子网ID")
			return subnet, nil
		}
	}
	return Subnet{}, fmt.Errorf("找不 %s 到这个子网", name)
}

// SecuritygroupsResult 查询安全组接口返回的对象
type SecuritygroupsResult struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

// SecurityGroup 单个安全组信息
type SecurityGroup struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	EnterpriseProjectID string              `json:"enterprise_project_id"`
	SecurityGroupRules  []SecurityGroupRule `json:"security_group_rules"`
	VpcID               *string             `json:"vpc_id,omitempty"`
}

// SecurityGroupRule 安全组规则具体的条目
type SecurityGroupRule struct {
	Description          *string     `json:"description"`
	Direction            string      `json:"direction"`
	Ethertype            string      `json:"ethertype"`
	ID                   string      `json:"id"`
	PortRangeMax         *int64      `json:"port_range_max"`
	PortRangeMin         *int64      `json:"port_range_min"`
	Protocol             *string     `json:"protocol"`
	RemoteGroupID        *string     `json:"remote_group_id"`
	RemoteIPPrefix       *string     `json:"remote_ip_prefix"`
	SecurityGroupID      string      `json:"security_group_id"`
	TenantID             string      `json:"tenant_id"`
	RemoteAddressGroupID interface{} `json:"remote_address_group_id"`
}

// GetSecurityGroups 根据安全组名字返回安全组 ID
func GetSecurityGroups(name string) (SecurityGroup, error) {
	for _, securitysroups := range Securitygroups.SecurityGroups {
		if securitysroups.Name == name {
			// fmt.Println("从缓存中读取安全组ID")
			return securitysroups, nil
		}
	}
	request, _ := http.NewRequest("GET", "https://vpc." + config.Region +".myhuaweicloud.com/v1/" + config.ProjectID +"/security-groups", bytes.NewBuffer([]byte("")))
	request.Header.Add("content-type", "application/json;charset=utf8")
	request.Header.Add("X-Project-Id", config.ProjectID)
	config.Signature.Sign(request)
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return SecurityGroup{}, fmt.Errorf("查询安全组接口访问失败 %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SecurityGroup{}, fmt.Errorf("查询安全组时没有得到响应 %s", err)
	}
	json.Unmarshal(body, &Securitygroups)
	for _, securitysroups := range Securitygroups.SecurityGroups {
		if securitysroups.Name == name {
			// fmt.Println("从在线中读取安全组ID")
			return securitysroups, nil
		}
	}
	return SecurityGroup{}, fmt.Errorf("找不 %s 到这个安全组", name)
}
