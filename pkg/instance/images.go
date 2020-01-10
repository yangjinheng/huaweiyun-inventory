package instance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yangjinheng/huaweiyun-inventory/pkg/config"
	"io/ioutil"
	"net/http"
)

// Images 作为运行时缓存
var Images ImagesResult

// ImagesResult 从接口中返回的镜像列表
type ImagesResult struct {
	Images []Image `json:"images"`
}

// Image 单个镜像对象
type Image struct {
	Schema                  string   `json:"schema"`
	IsOffshelved            *string  `json:"__is_offshelved,omitempty"`
	RootOrigin              *string  `json:"__root_origin,omitempty"`
	MinDisk                 int64    `json:"min_disk"`
	CreatedAt               string   `json:"created_at"`
	Originalimagename       *string  `json:"__originalimagename,omitempty"`
	ImageSourceType         string   `json:"__image_source_type"`
	ContainerFormat         string   `json:"container_format"`
	ImageSize               *string  `json:"__image_size,omitempty"`
	SystemSupportMarket     *bool    `json:"__system_support_market,omitempty"`
	File                    string   `json:"file"`
	UpdatedAt               string   `json:"updated_at"`
	Protected               bool     `json:"protected"`
	Productcode             *string  `json:"__productcode,omitempty"`
	MaxRAM                  *string  `json:"max_ram,omitempty"`
	Checksum                string   `json:"checksum"`
	ID                      string   `json:"id"`
	Description             *string  `json:"__description,omitempty"`
	Isregistered            string   `json:"__isregistered"`
	MinRAM                  int64    `json:"min_ram"`
	Lazyloading             *string  `json:"__lazyloading,omitempty"`
	Owner                   string   `json:"owner"`
	DataOrigin              *string  `json:"__data_origin,omitempty"`
	HwFirmwareType          *string  `json:"hw_firmware_type,omitempty"`
	HwVifMultiqueueEnabled  *string  `json:"hw_vif_multiqueue_enabled,omitempty"`
	OSType                  string   `json:"__os_type"`
	Imagetype               string   `json:"__imagetype"`
	Visibility              string   `json:"visibility"`
	SourceOwner             *string  `json:"__source_owner,omitempty"`
	VirtualEnvType          string   `json:"virtual_env_type"`
	AccountCode             *string  `json:"__account_code,omitempty"`
	Tags                    []string `json:"tags"`
	Platform                string   `json:"__platform"`
	Size                    int64    `json:"size"`
	SourceImageID           *string  `json:"__source_image_id,omitempty"`
	OSBit                   string   `json:"__os_bit"`
	OSVersion               string   `json:"__os_version"`
	Name                    string   `json:"name"`
	Self                    string   `json:"self"`
	DiskFormat              string   `json:"disk_format"`
	VirtualSize             int64    `json:"virtual_size"`
	SystemSupportExport     *bool    `json:"__system_support_export,omitempty"`
	Status                  string   `json:"status"`
	SupportKVM              *string  `json:"__support_kvm,omitempty"`
	SupportKVMHi1822Hiovs   *string  `json:"__support_kvm_hi1822_hiovs,omitempty"`
	SupportAgentList        *string  `json:"__support_agent_list,omitempty"`
	SupportXen              *string  `json:"__support_xen,omitempty"`
	OSFeatureList           *string  `json:"__os_feature_list,omitempty"`
	SequenceNum             *string  `json:"__sequence_num,omitempty"`
	ImageDisplayname        *string  `json:"__image_displayname,omitempty"`
	EnterpriseProjectID     *string  `json:"enterprise_project_id,omitempty"`
	SupportLiveResize       *string  `json:"__support_live_resize,omitempty"`
	SupportKVMNvmeSpdk      *string  `json:"__support_kvm_nvme_spdk,omitempty"`
	SupportTest             *string  `json:"__support_test,omitempty"`
	SupportLargememory      *string  `json:"__support_largememory,omitempty"`
	SupportKVMGPUType       *string  `json:"__support_kvm_gpu_type,omitempty"`
	BackupID                *string  `json:"__backup_id,omitempty"`
	SupportHighperformance  *string  `json:"__support_highperformance,omitempty"`
	SupportDiskintensive    *string  `json:"__support_diskintensive,omitempty"`
	SupportArm              *string  `json:"__support_arm,omitempty"`
	IsConfigInit            *string  `json:"__is_config_init,omitempty"`
	SupportKVMFPGAType      *string  `json:"__support_kvm_fpga_type,omitempty"`
	SupportKVMAscend310     *string  `json:"__support_kvm_ascend_310,omitempty"`
	SupportFcInject         *string  `json:"__support_fc_inject,omitempty"`
	SupportP3L              *string  `json:"__support_p3l,omitempty"`
	SupportH2LSell          *string  `json:"__support_h2l_sell,omitempty"`
	SupportH2LTest          *string  `json:"__support_h2l_test,omitempty"`
	SupportS43XlTest        *string  `json:"__support_s43xl_test,omitempty"`
	SupportH23Xl            *string  `json:"__support_h23xl,omitempty"`
	SupportP3LTest          *string  `json:"__support_p3l_test,omitempty"`
	SupportH2L              *string  `json:"__support_h2l,omitempty"`
	SupportP3LPod18         *string  `json:"__support_p3l_pod18,omitempty"`
	SupportH2LSellJfWyd     *string  `json:"__support_h2l_sell_jf_wyd,omitempty"`
	SupportD1T              *string  `json:"__support_d1t,omitempty"`
	SupportKVMInfiniband    *string  `json:"__support_kvm_infiniband,omitempty"`
	WholeImage              *string  `json:"__whole_image,omitempty"`
	SupportKVMHi1822Hisriov *string  `json:"__support_kvm_hi1822_hisriov,omitempty"`
}

// GetImage 根据镜像名称和服务器规格查找出镜像对象
func GetImage(name string, flavor string) (Image, error) {
	for _, image := range Images.Images {
		if image.Name == name {
			// fmt.Println("从缓存中读取镜像ID")
			return image, nil
		}
	}
	request, _ := http.NewRequest("GET", "https://ims." + config.Region +".myhuaweicloud.com/v2/cloudimages?__os_bit=64&flavor_id="+flavor+"&name="+name+"&status=active", bytes.NewBuffer([]byte("")))
	request.Header.Add("content-type", "application/json;charset=utf8")
	request.Header.Add("X-Project-Id", "6de764f622104bae93d5bbb45c6e2ec9")
	config.Signature.Sign(request)
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return Image{}, fmt.Errorf("查询镜像接口访问失败 %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Image{}, fmt.Errorf("查询镜像时没有得到响应 %s", err)
	}
	json.Unmarshal(body, &Images)
	for _, image := range Images.Images {
		if image.Name == name {
			// fmt.Println("从线上中读取镜像ID")
			return image, nil
		}
	}
	return Image{}, fmt.Errorf("找不 %s 到这个镜像", name)
}
