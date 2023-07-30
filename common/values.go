package common

import (
	"fmt"
	"github.com/weplanx/go/values"
	"strings"
)

type Values struct {
	Mode      string `env:"MODE" envDefault:"debug"`
	Address   string `env:"ADDRESS" envDefault:":3000"`
	Namespace string `env:"NAMESPACE,required"`
	Key       string `env:"KEY,required"`

	Database struct {
		Url   string `env:"URL,required"`
		Name  string `env:"NAME,required"`
		Redis string `env:"REDIS,required"`
	} `envPrefix:"DATABASE_"`

	Nats struct {
		Hosts []string `env:"HOSTS,required" envSeparator:","`
		Nkey  string   `env:"NKEY,required"`
	} `envPrefix:"NATS_"`

	*Extra
}

type Extra struct {
	BaseUrl     string
	IpAddress   string
	IpSecretId  string
	IpSecretKey string
	*values.DynamicValues
}

func (x Values) IsRelease() bool {
	return x.Mode == "release"
}

func (x Values) Name(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}
