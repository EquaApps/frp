package client

import (
	"fmt"

	"github.com/johncoker233/frpaaa/dao"
	"github.com/johncoker233/frpaaa/models"
)

type ValidateableClientRequest interface {
	GetClientSecret() string
	GetClientId() string
}

func ValidateClientRequest(req ValidateableClientRequest) (*models.ClientEntity, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	if req.GetClientId() == "" || req.GetClientSecret() == "" {
		return nil, fmt.Errorf("invalid request")
	}

	var (
		cli *models.ClientEntity
		err error
	)

	if cli, err = dao.ValidateClientSecret(req.GetClientId(), req.GetClientSecret()); err != nil {
		return nil, err
	}

	return cli, nil
}
