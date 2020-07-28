package model

type Campaign struct {
	Base

	Account string `json:"account" sql:"index"`

	PromoType string  `json:"promoType" validate:"required"`
	LowValue  float32 `json:"lowValue,string"`
	MedValue  float32 `json:"medValue,string"`
	HighValue float32 `json:"highValue,string"`
}
