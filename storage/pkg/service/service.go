package service

// Empty by design.

type Service interface {
}

type service struct {
}

func NewService() Service {
	return &service{}
}
