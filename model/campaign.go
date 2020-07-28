package model

type Campaign struct {
	Base

	Account string `json:"account" sql:"index"`

	PromoType string  `json:"promoType" validate:"required"`
	LowValue  float32 `json:"lowValue"`
	MedValue  float32 `json:"medValue"`
	HighValue float32 `json:"highValue"`
}
