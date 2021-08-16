package v1

import "time"

type MkyAuditLog struct {
	CreateTime time.Time `json:"create_time" description:"创建时间"`
	UserID     int       `json:"user_id" description:"用户ID"`
	StaffID    int       `json:"staff_id" description:"涉及用户ID"`
	IP         string    `json:"ip" description:"ip地址"`
	OPType     int       `json:"op_type" description:"操作类型"`
	OPName     string    `json:"op_name" description:"操作名称"`
	OPDesc     string    `json:"op_desc" description:"操作描述"`
	LogLevel   string    `json:"log_level" description:"日志等级"`
	LogType    int       `json:"log_type" description:"日志类型, 1: 审计日志"`
}

type ListMkyAuditLogResp struct {
	ListResponseMeta
	MkyAuditLogs []*MkyAuditLog `json:"logs"`
}

type Response struct {
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

type BatchLogsReq struct {
	Datas []*MkyAuditLog `json:"datas"`
}
