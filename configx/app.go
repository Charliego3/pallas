package configx

type App struct {
	Network string `json:"network,omitempty" yaml:"network,omitempty"`
	Address string `json:"address,omitempty" yaml:"address"`
}

type embeddedAppFetcher struct{}

func (f *embeddedAppFetcher) Fetch() (App, error) {
	return *embedded.App, nil
}
