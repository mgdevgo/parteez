package artwork

import "parteez/internal/repository"

type ArtworkRepository interface {
	repository.Repository[Artwork]
}
