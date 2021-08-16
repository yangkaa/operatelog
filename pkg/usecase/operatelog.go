package usecase

import (
	v1 "goodrain.com/operatelog/pkg/controller/v1"
	"goodrain.com/operatelog/pkg/models"
	"gorm.io/gorm"
)

type MkyAuditLogUsecase interface {
	Create(req *v1.MkyAuditLog) error
	List(startTime, endTime, query string, page, pageSize int) (*v1.ListMkyAuditLogResp, error)
	BatchCreateLogs(batchLogs *v1.BatchLogsReq) error
}

func NewMkyAuditLogUcase() MkyAuditLogUsecase {
	return &mkyAuditLogUcase{}
}

type mkyAuditLogUcase struct {
	DB              *gorm.DB                     `inject:""`
	MkyAuditLogRepo models.MkyAuditLogRepository `inject:""`
}

func (m *mkyAuditLogUcase) Create(req *v1.MkyAuditLog) error {
	mkyAuditLog := &models.MkyAuditLog{
		UserID:   req.UserID,
		StaffID:  req.StaffID,
		IP:       req.IP,
		OPType:   req.OPType,
		OPName:   req.OPName,
		OPDesc:   req.OPDesc,
		LogLevel: req.LogLevel,
		LogType:  req.LogType,
	}
	return m.MkyAuditLogRepo.Create(mkyAuditLog)
}

func (m *mkyAuditLogUcase) List(startTime, endTime, query string, page, pageSize int) (*v1.ListMkyAuditLogResp, error) {
	logs, total, err := m.MkyAuditLogRepo.List(startTime, endTime, query, page, pageSize)
	if err != nil {
		return nil, err
	}
	var mkyAuditLogs []*v1.MkyAuditLog
	for _, log := range logs {
		mkyAuditLogs = append(mkyAuditLogs, &v1.MkyAuditLog{
			CreateTime: log.CreatedAt,
			UserID:     log.UserID,
			StaffID:    log.StaffID,
			IP:         log.IP,
			OPType:     log.OPType,
			OPName:     log.OPName,
			OPDesc:     log.OPDesc,
			LogLevel:   log.LogLevel,
			LogType:    log.LogType,
		})
	}
	resp := &v1.ListMkyAuditLogResp{
		MkyAuditLogs: mkyAuditLogs,
	}
	resp.PageNo = page
	resp.PageSize = pageSize
	resp.TotalCount = total
	return resp, nil
}

func (m *mkyAuditLogUcase) BatchCreateLogs(batchLogs *v1.BatchLogsReq) error {
	var mkyAuditLogs []*models.MkyAuditLog
	for _, req := range batchLogs.Datas {
		mkyAuditLog := &models.MkyAuditLog{
			UserID:   req.UserID,
			StaffID:  req.StaffID,
			IP:       req.IP,
			OPType:   req.OPType,
			OPName:   req.OPName,
			OPDesc:   req.OPDesc,
			LogLevel: req.LogLevel,
			LogType:  req.LogType,
		}
		mkyAuditLogs = append(mkyAuditLogs, mkyAuditLog)
	}
	return m.MkyAuditLogRepo.BatchCreate(mkyAuditLogs)
}
