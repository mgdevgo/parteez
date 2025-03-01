package artist

type ArtistID int

type Artist struct {
	ID          ArtistID
	Name        string
	SocialLinks []SocialLink
}

type SocialLink string
