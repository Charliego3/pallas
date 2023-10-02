package configx

type App struct {
	Network string `json:"network,omitempty" yaml:"network,omitempty"`
	Address string `json:"address,omitempty" yaml:"address"`
}

type standardAppFetcher struct{}

func (f *standardAppFetcher) Fetch() (App, error) {
	return *standard.App, nil
}
