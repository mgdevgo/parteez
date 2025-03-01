package account

import "time"

type AccountID int

type Account struct {
	ID         AccountID
	Email      string
	Password   string
	TelegramID int64
	Profile    []Profile
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProfileID int

type Profile struct {
	ID          ProfileID
	Type        ProfileType
	DisplayName string
	IsActive    bool
}

type ProfileType string

const (
	ProfileUser    ProfileType = "USER"
	ProfileAdmin   ProfileType = "ADMIN"
	ProfileArtist  ProfileType = "ARTIST"
	ProfileGuest   ProfileType = "GUEST"
	ProfilePremium ProfileType = "PREMIUM"
)

type Guest interface {
	Login() error
}

type User interface {
	Guest
	LoginWithPassword(password string) error
}

type Artist interface {
	User
	IsVerified() bool
}

type Admin interface {
	User
}
