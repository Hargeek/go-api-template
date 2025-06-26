package service

type HelloServiceImpl struct{}

func NewHelloServiceImpl() *HelloServiceImpl {
	return &HelloServiceImpl{}
}

func (h *HelloServiceImpl) Hello() string {
	return "Hello, World!"
}
