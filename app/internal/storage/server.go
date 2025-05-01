package storage

type Information struct {
	Role string
}

func NewInformation(role string) *Information {
	return &Information{
		Role: role,
	}
}
