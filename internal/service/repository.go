package service

type OrdersRepository interface {
	GetOne() (string, error)
	Save() error
}
