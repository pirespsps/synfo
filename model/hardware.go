package model

type Hardware interface {
	Overall() ([]byte, error)
	Extensive() ([]byte, error)
}
