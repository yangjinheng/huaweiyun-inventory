package instance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/yangjinheng/huaweiyun-inventory/pkg/config"
)

// --------------------------------------------------------- 创建实例传递的 JSON ---------------------------------------------------------

// Instance 访问创建主机接口时候传递的 JSON
type Instance struct {
	Server Server `json:"server"`
}

// Server 表示服务器的具体配置
type Server struct {
	Name             string               `json:"name"`
	AdminPass        string               `json:"adminPass"`
	AvailabilityZone string               `json:"availability_zone"`
	FlavorRef        string               `json:"flavorRef"`
	ImageRef         string               `json:"imageRef"`
	RootVolume       config.RootVolume    `json:"root_volume"`
	DataVolumes      *[]config.DataVolume `json:"data_volumes,omitempty"`
	Vpcid            string               `json:"vpcid"`
	Nics             []Nics               `json:"nics"`
	SecurityGroups   []SecurityGroupID    `json:"security_groups"`
	Publicip         *Publicip            `json:"publicip,omitempty"`
	Extendparam      Extendparam          `json:"extendparam"`
	ServerTags       []config.ServerTags  `json:"server_tags"`
}

// Nics 子网信息
type Nics struct {
	SubnetID string `json:"subnet_id"`
}

// SecurityGroupID 安全组对象
type SecurityGroupID struct {
	ID string `json:"id"`
}

// Publicip 公网 IP 对象
type Publicip struct {
	Eip Eip `json:"eip"`
}

// Eip 公网 IP
type Eip struct {
	Iptype      string         `json:"iptype"`
	Bandwidth   Bandwidth      `json:"bandwidth"`
	Extendparam EipExtendparam `json:"extendparam"`
}

// Bandwidth 公网 IP
type Bandwidth struct {
	Size       int    `json:"size"`
	Sharetype  string `json:"sharetype"`
	Chargemode string `json:"chargemode"`
}

// EipExtendparam 公网 IP
type EipExtendparam struct {
	ChargingMode string `json:"chargingMode"`
}

// Extendparam 服务器扩展信息
type Extendparam struct {
	ChargingMode string `json:"chargingMode"`
	IsAutoPay    string `json:"isAutoPay"`
}

// ServerTags 服务器标签
type ServerTags struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// --------------------------------------------------------- 创建实例 ---------------------------------------------------------

// CreateServer 接收主要参数，返回一个服务器的结构体定义
func CreateServer(paramer config.Host) (Job, error) {
	// 在线取得子网ID
	subnet, err := GetSubnet(paramer.Subnetname)
	if err != nil {
		return Job{}, err
	}

	// 在线取得安全组ID
	securitygroups := []SecurityGroupID{}
	for _, sgname := range paramer.SecurityGroups {
		sgobj, err := GetSecurityGroups(sgname)
		if err != nil {
			return Job{}, err
		}
		securitygroup := SecurityGroupID{ID: sgobj.ID}
		securitygroups = append(securitygroups, securitygroup)
	}

	// 在线取得镜像ID
	image, err := GetImage(paramer.System, paramer.FlavorRef)
	if err != nil {
		return Job{}, err
	}

	// 创建实例对象结构体
	server := Instance{}
	server.Server = Server{
		AvailabilityZone: "cn-north-1a",
		Name:             paramer.HostName,
		AdminPass:        paramer.AdminPass,
		ImageRef:         image.ID,
		FlavorRef:        paramer.FlavorRef,
		Vpcid:            subnet.VpcID,
		RootVolume:       paramer.RootVolume,
		DataVolumes:      &paramer.DataVolume,
		Nics:             []Nics{{SubnetID: subnet.ID}},
		SecurityGroups:   securitygroups,
		// Publicip:         Publicip{Eip: Eip{Iptype: "5_bgp", Bandwidth: Bandwidth{Size: 100, Sharetype: "PER", Chargemode: "traffic"}, Extendparam: EipExtendparam{ChargingMode: "postPaid"}}},
		Extendparam: Extendparam{ChargingMode: "postPaid", IsAutoPay: "IsAutoPay"},
		ServerTags:  paramer.ServerTags,
	}

	// 创建服务器请求对象
	serverjson, _ := json.Marshal(server)
	request, _ := http.NewRequest("POST", "https://ecs."+config.Region+".myhuaweicloud.com/v1/"+config.ProjectID+"/cloudservers", bytes.NewBuffer(serverjson))
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-Project-Id", config.ProjectID)

	// 请求签名
	config.Signature.Sign(request)

	// 发出请求
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return Job{}, fmt.Errorf("请求地址错误 %s", request.URL)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Job{}, fmt.Errorf("创建主机失败 %s", request.URL)
	}
	// 读取结果并返回
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Job{}, fmt.Errorf("无法读取创建主机的响应报文")
	}
	job := Job{}
	json.Unmarshal(body, &job)
	return job, nil
}

// QuicklyCreate 接收一个主机信息切片，返回一个主机切片
func QuicklyCreate(hosts []config.Host) (sucesshosts []ServerElement) {
	channel := make(chan ServerElement)
	var waitgroup sync.WaitGroup
	for _, host := range hosts {
		waitgroup.Add(1)
		go func(host config.Host) {
			// 创建主机
			job, err := CreateServer(host)
			if err != nil {
				fmt.Println("创建主机时候出现问题", err)
				waitgroup.Done()
				return
			}
			// 等待 JOB 完成
			sucessjob, err := WaitJobSucess(job)
			if err != nil {
				fmt.Println("获取 job id 出错", err)
				waitgroup.Done()
				return
			}
			// 根据 serverid 查询主机的详细信息并返回
			var serverid string
			for _, subjob := range sucessjob.Entities.SubJobs {
				if subjob.JobType == "createSingleServer" {
					serverid = subjob.Entities.ServerID
				}
			}
			server, err := ShowServer(serverid)
			if err != nil {
				waitgroup.Done()
				return
			}
			// 写入通道
			channel <- server
		}(host)
	}
	// 从通道中取得创建好的主机
	go func() {
		for host := range channel {
			sucesshosts = append(sucesshosts, host)
			waitgroup.Done()
		}
	}()
	waitgroup.Wait()
	return sucesshosts
}

// --------------------------------------------------------- 查询实例返回的 JSON ---------------------------------------------------------

// ServerElement 是调用查询 Server 时候返回的单个 Server 的详细信息
type ServerElement struct {
	Server ServerClass `json:"server"`
}

// ServerClass 主题信息
type ServerClass struct {
	Fault                            *Fault                             `json:"fault"`
	ID                               string                             `json:"id"`
	Name                             string                             `json:"name"`
	Addresses                        map[string][]Address               `json:"addresses"`
	Flavor                           Flavor                             `json:"flavor"`
	AccessIPv4                       string                             `json:"accessIPv4"`
	AccessIPv6                       string                             `json:"accessIPv6"`
	Status                           string                             `json:"status"`
	Progress                         *int64                             `json:"progress"`
	HostID                           string                             `json:"hostId"`
	Updated                          string                             `json:"updated"`
	Created                          string                             `json:"created"`
	Metadata                         Metadata                           `json:"metadata"`
	Tags                             []string                           `json:"tags"`
	Description                      string                             `json:"description"`
	Locked                           bool                               `json:"locked"`
	ServerNovaLock                   bool                               `json:"serverNovaLock"`
	ConfigDrive                      string                             `json:"config_drive"`
	TenantID                         string                             `json:"tenant_id"`
	UserID                           string                             `json:"user_id"`
	KeyName                          *string                            `json:"key_name"`
	OSExtendedVolumesVolumesAttached []OSExtendedVolumesVolumesAttached `json:"os-extended-volumes:volumes_attached"`
	OSEXTSTSTaskState                *string                            `json:"OS-EXT-STS:task_state"`
	OSEXTSTSPowerState               int64                              `json:"OS-EXT-STS:power_state"`
	OSEXTSTSVMState                  string                             `json:"OS-EXT-STS:vm_state"`
	OSEXTSRVATTRHost                 string                             `json:"OS-EXT-SRV-ATTR:host"`
	OSEXTSRVATTRInstanceName         string                             `json:"OS-EXT-SRV-ATTR:instance_name"`
	OSEXTSRVATTRHypervisorHostname   string                             `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	OSDCFDiskConfig                  string                             `json:"OS-DCF:diskConfig"`
	OSEXTAZAvailabilityZone          string                             `json:"OS-EXT-AZ:availability_zone"`
	OSSchedulerHints                 OSSchedulerHints                   `json:"os:scheduler_hints"`
	OSEXTSRVATTRRootDeviceName       string                             `json:"OS-EXT-SRV-ATTR:root_device_name"`
	OSEXTSRVATTRRamdiskID            string                             `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	EnterpriseProjectID              string                             `json:"enterprise_project_id"`
	OSEXTSRVATTRUserData             string                             `json:"OS-EXT-SRV-ATTR:user_data"`
	OSSRVUSGLaunchedAt               string                             `json:"OS-SRV-USG:launched_at"`
	OSEXTSRVATTRKernelID             string                             `json:"OS-EXT-SRV-ATTR:kernel_id"`
	OSEXTSRVATTRLaunchIndex          int64                              `json:"OS-EXT-SRV-ATTR:launch_index"`
	HostStatus                       string                             `json:"host_status"`
	OSEXTSRVATTRReservationID        string                             `json:"OS-EXT-SRV-ATTR:reservation_id"`
	OSEXTSRVATTRHostname             string                             `json:"OS-EXT-SRV-ATTR:hostname"`
	OSSRVUSGTerminatedAt             *string                            `json:"OS-SRV-USG:terminated_at"`
	SysTags                          []SysTag                           `json:"sys_tags"`
	SecurityGroups                   []ServerSecurityGroup              `json:"security_groups"`
	Image                            ServerImage                        `json:"image"`
}

// Address 地址信息
type Address struct {
	Version            string `json:"version"`
	Addr               string `json:"addr"`
	OSEXTIPSMACMACAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
	OSEXTIPSType       string `json:"OS-EXT-IPS:type"`
	OSEXTIPSPortID     string `json:"OS-EXT-IPS:port_id"`
}

// Fault 故障信息
type Fault struct {
	Code    int64  `json:"code"`
	Created string `json:"created"`
	Message string `json:"message"`
}

// Flavor 规格信息
type Flavor struct {
	Disk  string `json:"disk"`
	Vcpus string `json:"vcpus"`
	RAM   string `json:"ram"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

// ServerImage 使用的镜像 ID
type ServerImage struct {
	ID string `json:"id"`
}

// Metadata 元数据信息
type Metadata struct {
	MeteringImageID           string  `json:"metering.image_id"`
	MeteringImagetype         string  `json:"metering.imagetype"`
	MeteringResourcespeccode  string  `json:"metering.resourcespeccode"`
	ImageName                 string  `json:"image_name"`
	OSBit                     string  `json:"os_bit"`
	CascadedInstanceExtrainfo string  `json:"cascaded.instance_extrainfo"`
	MeteringResourcetype      string  `json:"metering.resourcetype"`
	VpcID                     string  `json:"vpc_id"`
	OSType                    string  `json:"os_type"`
	ChargingMode              string  `json:"charging_mode"`
	SupportAgentList          *string `json:"__support_agent_list,omitempty"`
	MeteringOrderID           *string `json:"metering.order_id,omitempty"`
	LockCheckEndpoint         *string `json:"lockCheckEndpoint,omitempty"`
	LockSource                *string `json:"lockSource,omitempty"`
	EcmResStatus              *string `json:"EcmResStatus,omitempty"`
	MeteringProductID         *string `json:"metering.product_id,omitempty"`
	LockSourceID              *string `json:"lockSourceId,omitempty"`
	LockScene                 *string `json:"lockScene,omitempty"`
}

// OSExtendedVolumesVolumesAttached 挂载的卷
type OSExtendedVolumesVolumesAttached struct {
	Device              string  `json:"device"`
	BootIndex           *string `json:"bootIndex"`
	ID                  string  `json:"id"`
	DeleteOnTermination string  `json:"delete_on_termination"`
}

// OSSchedulerHints 未知
type OSSchedulerHints struct {
	BuildNearHostIP *string `json:"build_near_host_ip,omitempty"`
	CIDR            *string `json:"cidr,omitempty"`
	DifferentHost   *string `json:"different_host,omitempty"`
	Group           *string `json:"group,omitempty"`
	SameHost        *string `json:"same_host,omitempty"`
}

// ServerSecurityGroup 安全组信息
type ServerSecurityGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SysTag TAG
type SysTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// --------------------------------------------------------- 查询实例 ---------------------------------------------------------

// ShowServer 根据 ServerID 返回实例的详细信息
func ShowServer(serverid string) (ServerElement, error) {
	request, _ := http.NewRequest("GET", "https://ecs."+config.Region+".myhuaweicloud.com/v1/"+config.ProjectID+"/cloudservers/"+serverid, bytes.NewBuffer([]byte("")))
	request.Header.Add("content-type", "application/json;charset=utf8")
	request.Header.Add("X-Project-Id", config.ProjectID)
	config.Signature.Sign(request)
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return ServerElement{}, fmt.Errorf("访问查询实例接口失败 %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ServerElement{}, fmt.Errorf("查询实例接口的响应读取失败 %s", err)
	}
	server := ServerElement{}
	json.Unmarshal(body, &server)
	return server, nil
}
