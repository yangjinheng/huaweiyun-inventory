package config

import (
	"io/ioutil"
	"strings"

	"github.com/yangjinheng/huaweiyun-inventory/pkg/signer"
	"gopkg.in/yaml.v3"
)

// Signature 签名工具
var (
	Signature signer.Signer
	Region string
	ProjectID string
)

// Config 描述配置文件有哪些配置段
type Config struct {
	Huaweicloud Huaweicloud `yaml:"huaweicloud"`
	Instances   []Host      `yaml:"instances"`
}

// Huaweicloud 是配置文件中华为云相关配置
type Huaweicloud struct {
	Accesskey  string `yaml:"access_key"`
	Secretkey  string `yaml:"secret_key"`
	ProjectID  string `yaml:"project_id"`
	Iamaddress string `yaml:"iam_address"`
}

// Host 是 CreateServer 创建主机函数接收的参数，它是配置文件的一部分
type Host struct {
	HostName       string       `json:"hostname" yaml:"hostname"`
	AdminPass      string       `json:"adminPass" yaml:"adminpass"`
	System         string       `json:"system" yaml:"system"`
	FlavorRef      string       `json:"flavorref" yaml:"flavorref"`
	Subnetname     string       `json:"subnetname" yaml:"subnetname"`
	RootVolume     RootVolume   `json:"rootvolume" yaml:"rootvolume"`
	DataVolume     []DataVolume `json:"datavolume" yaml:"datavolume"`
	SecurityGroups []string     `json:"securitgroups" yaml:"securitgroups"`
	ServerTags     []ServerTags `json:"servertags" yaml:"servertags"`
}

// RootVolume 系统盘
type RootVolume struct {
	Volumetype string `json:"volumetype" yaml:"volumetype"`
	Size       int    `json:"size" yaml:"size"`
}

// DataVolume 数据盘
type DataVolume struct {
	Volumetype string `json:"volumetype" yaml:"volumetype"`
	Size       int    `json:"size" yaml:"size"`
}

// ServerTags 标签
type ServerTags struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// LoadParamer 解析配置文件
func LoadParamer(config string) (Config, error) {
	paramer := Config{}
	file, err := ioutil.ReadFile(config)
	if err != nil {
		return paramer, err
	}
	err = yaml.Unmarshal(file, &paramer)
	if err != nil {
		return paramer, err
	}
	Signature = signer.Signer{
		Key: paramer.Huaweicloud.Accesskey,
		Secret: paramer.Huaweicloud.Secretkey,
	}
	Region = strings.Split(paramer.Huaweicloud.Iamaddress, ".")[1]
	ProjectID = paramer.Huaweicloud.ProjectID
	return paramer, nil
}
