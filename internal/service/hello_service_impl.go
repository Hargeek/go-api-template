package service

type HelloServiceImpl struct{}

func (h *HelloServiceImpl) Hello() string {
	return "Hello, World!"
}
