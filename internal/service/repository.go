package service

type ConversionRepository interface {
	GetOne() (string, error)
	Save() error
}
