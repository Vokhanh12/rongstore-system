package domain

type TitleGl interface {
	CheckHealth() error
	GetBaseURL() string
}
