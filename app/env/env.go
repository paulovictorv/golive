package env

type Env struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	CdnId  string `yaml:"cdnId"`
	Bucket string `yaml:"bucket"`
}
