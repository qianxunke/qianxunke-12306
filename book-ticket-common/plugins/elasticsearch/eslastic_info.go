package elasticsearch

type ElasticConfigInfo struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Index           string `yaml:"index"`
	Type            string `yaml:"type"`
	isOpeneSetSniff bool   `yaml:"isOpeneSetSniff"`
}
