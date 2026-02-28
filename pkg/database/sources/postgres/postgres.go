package postgres

const SourceKind = "postgres"

type Source struct {
}

func New() *Source {
	return &Source{}
}
