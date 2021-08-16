package models

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type MkyAuditLog struct {
	Model
	UserID   int    `gorm:"column:user_id" comment:"用户ID"`
	StaffID  int    `gorm:"column:staff_id" comment:"涉及用户ID"`
	IP       string `gorm:"column:ip;size:50" comment:"ip地址"`
	OPType   int    `gorm:"column:op_type" comment:"操作类型"`
	OPName   string `gorm:"column:op_name;size:50" comment:"操作名称"`
	OPDesc   string `gorm:"column:op_desc;size:3000" comment:"操作描述"`
	LogLevel string `gorm:"column:log_level;size:50" comment:"日志等级"`
	LogType  int    `gorm:"column:log_type" comment:"日志类型, 1: 审计日志"`
}

type MkyAuditLogRepository interface {
	Create(log *MkyAuditLog) error
	List(startTime, endTime, query string, page, pageSize int) ([]*MkyAuditLog, int64, error)
	BatchCreate(logs []*MkyAuditLog) error
}

type MkyAuditLogRepo struct {
	DB *gorm.DB `inject:""`
}

func NewMkyAuditLogRepo() MkyAuditLogRepository {
	return &MkyAuditLogRepo{}
}

func (m *MkyAuditLogRepo) WithTransaction(db *gorm.DB) MkyAuditLogRepository {
	return &MkyAuditLogRepo{DB: db}
}

func (m *MkyAuditLogRepo) Create(log *MkyAuditLog) error {
	if err := m.DB.Create(log).Error; err != nil {
		return errors.Wrap(err, "create log failed")
	}
	return nil
}

func (m *MkyAuditLogRepo) List(startTime, endTime, query string, page, pageSize int) ([]*MkyAuditLog, int64, error) {
	db := m.DB
	if startTime != "" {
		db = db.Where("create_time > ?", startTime)
	}
	if endTime != "" {
		db = db.Where("create_time < ?", endTime)
	}
	if query != "" {
		db = db.Where("op_desc like ?", "%"+query+"%")
	}

	var total int64
	if err := db.Model(&MkyAuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count log failed")
	}

	if pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	var logs []*MkyAuditLog
	if err := db.Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(err, "list logs failed")
	}
	return logs, total, nil
}

func (m *MkyAuditLogRepo) BatchCreate(logs []*MkyAuditLog) error {
	if err := m.DB.CreateInBatches(logs, 100).Error; err != nil {
		return errors.Wrap(err, "batch create log failed")
	}
	return nil
}
