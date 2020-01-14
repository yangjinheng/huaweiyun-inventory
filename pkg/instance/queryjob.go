package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/yangjinheng/huaweiyun-inventory/pkg/config"
)

// QueryJob 查询 job 正确返回
type QueryJob struct {
	JobID      string           `json:"job_id"`
	JobType    string           `json:"job_type"`
	BeginTime  string           `json:"begin_time"`
	EndTime    string           `json:"end_time"`
	Status     string           `json:"status"`
	ErrorCode  string           `json:"error_code"`
	FailReason string           `json:"fail_reason"`
	Entities   QueryJobEntities `json:"entities"`
}

// QueryJobEntities 的实体
type QueryJobEntities struct {
	SubJobsTotal int64    `json:"sub_jobs_total"`
	SubJobs      []SubJob `json:"sub_jobs"`
}

// SubJob 某一个 JOB
type SubJob struct {
	JobID      string         `json:"job_id"`
	JobType    string         `json:"job_type"`
	BeginTime  string         `json:"begin_time"`
	EndTime    string         `json:"end_time"`
	Status     string         `json:"status"`
	ErrorCode  string         `json:"error_code"`
	FailReason string         `json:"fail_reason"`
	Entities   SubJobEntities `json:"entities"`
}

// SubJobEntities 某一个 JOB 任务返回的实体
type SubJobEntities struct {
	ServerID string `json:"server_id"`
}

// Job 是创建主机接口返回的 jobid 结构体
type Job struct {
	JobID string `json:"job_id"`
}

// Getjob 是根据 jobid 查询 job 的状态并返回
func Getjob(job Job) (QueryJob, error) {
	request, _ := http.NewRequest("GET", "https://ecs."+config.Region+".myhuaweicloud.com/v1/"+config.ProjectID+"/jobs/"+job.JobID, nil)
	request.Header.Add("content-type", "application/json;charset=utf8")
	request.Header.Add("X-Project-Id", config.ProjectID)
	config.Signature.Sign(request)
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return QueryJob{}, fmt.Errorf("请求地址错误 %s", request.URL)
	}
	defer resp.Body.Close()
	// 如果收到 200 的响应码
	if resp.StatusCode != 200 {
		return QueryJob{}, fmt.Errorf("查询JOB失败 %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return QueryJob{}, fmt.Errorf("查询 JOB 没有得到正确响应")
	}
	result := QueryJob{}
	json.Unmarshal(body, &result)
	return result, nil
}

// WaitJobSucess 等待创建主机的 JOB 运行成功
func WaitJobSucess(job Job) (QueryJob, error) {
	for {
		job, err := Getjob(job)
		time.Sleep(2e9)
		if err != nil {
			continue
		}
		if job.Status == "FAIL" {
			return QueryJob{}, fmt.Errorf("JOB执行失败")
		}
		if job.Status == "SUCCESS" {
			for _, subjob := range job.Entities.SubJobs {
				if subjob.JobType == "createSingleServer" && subjob.Status == "SUCCESS" {
					return job, nil
				}
			}
		}
	}
}
