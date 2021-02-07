package hanako

import "github.com/t-oki/pollen-api/internal/domain/entity"

type PollenRepositoryImpl struct{}

func NewPollenRepositoryImpl() entity.PollenRepository {
	return &PollenRepositoryImpl{}
}
