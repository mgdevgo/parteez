package artwork

import "parteez/internal/domain/shared"

type ArtworkRepository interface {
	shared.Repository[Artwork]
}
