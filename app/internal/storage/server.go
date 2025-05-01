package storage

type Information struct {
	Role string
}

func NewInformation() *Information {
	return &Information{
		Role: "master",
	}
}
