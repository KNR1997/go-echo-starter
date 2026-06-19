package domain

type AuditLog struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   int    `gorm:"index;not null" json:"user_id"`
	Username string `gorm:"size:64;index;default:''" json:"username"`
	Module   string `gorm:"size:64;index;default:''" json:"module"`
}
