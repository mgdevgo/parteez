package artwork

import "parteez/internal/domain/shared/repository"

type ArtworkRepository interface {
	repository.Repository[Artwork]
}
