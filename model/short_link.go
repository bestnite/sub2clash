package model

type ShortLink struct {
	ID              string        `gorm:"unique"`
	Config          ConvertConfig `gorm:"serializer:json"`
	Password        string
	LastRequestTime int64
}
