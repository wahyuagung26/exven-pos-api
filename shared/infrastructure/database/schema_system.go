package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net"
	"time"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type RetentionType string
type RetentionStatus string

const (
	RetentionTypeTransactionArchive RetentionType = "transaction_archive"
	RetentionTypeTransactionDelete  RetentionType = "transaction_delete"
	RetentionTypeAuditCleanup       RetentionType = "audit_cleanup"

	RetentionStatusSuccess RetentionStatus = "success"
	RetentionStatusFailed  RetentionStatus = "failed"
	RetentionStatusPartial RetentionStatus = "partial"
)

type AuditLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	TenantID  uint64    `gorm:"not null;index:idx_audit_logs_tenant_date"`
	UserID    *uint64   `gorm:"index:idx_audit_logs_user_date;constraint:OnDelete:SET NULL"`
	Action    string    `gorm:"size:100;not null"`
	TableName string    `gorm:"size:100;index:idx_audit_logs_table_record"`
	RecordID  *uint64   `gorm:"index:idx_audit_logs_table_record"`
	OldValues JSONMap   `gorm:"type:jsonb"`
	NewValues JSONMap   `gorm:"type:jsonb"`
	IPAddress net.IP    `gorm:"type:inet"`
	UserAgent string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_audit_logs_tenant_date;index:idx_audit_logs_user_date"`

	Tenant Tenant `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

type DataRetentionLog struct {
	ID              uint64          `gorm:"primaryKey;autoIncrement"`
	TenantID        uint64          `gorm:"not null;index:idx_data_retention_logs_tenant_type_date"`
	RetentionType   RetentionType   `gorm:"not null;index:idx_data_retention_logs_tenant_type_date"`
	RecordsAffected int             `gorm:"not null"`
	DateFrom        time.Time       `gorm:"type:date;not null"`
	DateTo          time.Time       `gorm:"type:date;not null"`
	ExecutionTime   *float64        `gorm:"type:decimal(8,3)"`
	Status          RetentionStatus `gorm:"default:'success'"`
	ErrorMessage    string          `gorm:"type:text"`
	CreatedAt       time.Time       `gorm:"autoCreateTime;index:idx_data_retention_logs_tenant_type_date"`

	Tenant Tenant `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
}
