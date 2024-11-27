type User struct {
	ID            uint   `gorm:"primaryKey"`
	NationalID    string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
	PasswordHash  string `gorm:"not null"`
	FirstName     string `gorm:"size:100"`
	LastName      string
	DateOfBirth   time.Time
	City          string
	WalletBalance float64 `gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}