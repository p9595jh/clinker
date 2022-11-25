package service_test

import (
	"clinker-backend/internal/domain/service"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/p9595jh/transform"
)

func TestEthAddr(t *testing.T) {
	s := service.NewProcessService(validator.New(), transform.New())
	s.Initializer()

	type Dater struct {
		Date string `validate:"required,date"`
	}
	type Address struct {
		Addr string `validate:"required,ethAddr"`
	}
	type Hash struct {
		TxHash string `validate:"required,txHash"`
	}

	dateSample := Dater{
		Date: "2022-11-24",
	}
	err := s.Validate(&dateSample)
	t.Log(err)

	addrSample := Address{
		Addr: "0x4d943a7C1f2AF858BfEe8aB499fbE76B1D046eC7",
	}
	err = s.Validate(&addrSample)
	t.Log(err)

	hashSample := Hash{
		TxHash: "0x9cd2e21195e2548f3d6bb6977a9b5302f6cee3985f7d1af635f93b523a7d25f4",
	}
	err = s.Validate(&hashSample)
	t.Log(err)
}
