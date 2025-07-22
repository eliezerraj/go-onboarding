package model

import (
	"time"
	go_core_pg "github.com/eliezerraj/go-core/database/pg"
	go_core_observ "github.com/eliezerraj/go-core/observability" 
)

type AppServer struct {
	InfoPod 		*InfoPod 					`json:"info_pod"`
	Server     		*Server     				`json:"server"`
	ConfigOTEL		*go_core_observ.ConfigOTEL	`json:"otel_config"`
	DatabaseConfig	*go_core_pg.DatabaseConfig  `json:"database"`
	AwsService		*AwsService					`json:"aws_services"`
	Cert			*Cert						`json:"cert_tls_server"`
}

type InfoPod struct {
	PodName				string 	`json:"pod_name"`
	ApiVersion			string 	`json:"version"`
	OSPID				string 	`json:"os_pid"`
	IPAddress			string 	`json:"ip_address"`
	AvailabilityZone 	string 	`json:"availabilityZone"`
	IsAZ				bool   	`json:"is_az"`
	Env					string `json:"enviroment,omitempty"`
	AccountID			string `json:"account_id,omitempty"`
}

type Server struct {
	Port 			int `json:"port"`
	ReadTimeout		int `json:"readTimeout"`
	WriteTimeout	int `json:"writeTimeout"`
	IdleTimeout		int `json:"idleTimeout"`
	CtxTimeout		int `json:"ctxTimeout"`
}

type AwsService struct {
	AwsRegion			string `json:"aws_region"`
	BucketName			string `json:"bucket_name"`
	FilePath			string `json:"file_path"`
}

type MessageRouter struct {
	Message			string `json:"message"`
}

type Onboarding struct {
	Person 			*Person `json:"person"`
}

type Cert struct {
	IsTLS				bool	`json:"server_tls"`	
	CertPEM 			[]byte 	`json:"cert_pen"`
	CertPEMStr 			string 	`json:"cert_pen_str"`		
	CertPrivKeyPEM		[]byte  `json:"private_key"`
	CertPrivKeyPEMStr	string  `json:"private_key_str"`	 
}

type Person struct {
	ID			int 		`json:"id,omitempty"`
	PersonID	string		`json:"person_id,omitempty"`
	Name 		string 		`json:"name,omitempty"`
	CreatedAt	time.Time 	`json:"created_at,omitempty"`
	UpdatedAt	*time.Time 	`json:"updated_at,omitempty"`
	TenantID	string 		`json:"tenant_id,omitempty"`
}

type OnboardingFile struct {
	BucketName	string	`json:"bucket_name,omitempty"`
	FilePath	string 	`json:"file_path"`
	FileName	string	`json:"file_name,omitempty"`
	File		[]byte	`json:"file,omitempty"`
}