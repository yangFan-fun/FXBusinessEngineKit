package tool

import "FXBusinessEngineKit_Server/configuration"

type OrganizationModel struct {
	Id               int    `gorm:"column:id"`
	Organization     string `gorm:"column:organization"`
	RegistrationDate int64  `gorm:"column:registrationDate"`
	ProductName      string `gorm:"column:productName"`
	ProductId        string `gorm:"column:productId"`
}

func (receiver OrganizationModel) TableName() string {
	return configuration.FXDatabaseProduct
}
