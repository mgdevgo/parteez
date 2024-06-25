package services

import (
	"iditusi/internal/services/local"
)

type AccessService interface {
}

var _ AccessService = (*local.AccessService)(nil)
