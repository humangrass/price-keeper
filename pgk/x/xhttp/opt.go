package xhttp

import (
	"fmt"
	"time"
)

type Opt struct {
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	EnableHealthMethod bool          `yaml:"enable_health_method"`
}

func (o *Opt) Validate() error {
	if len(o.Host) == 0 {
		return fmt.Errorf("incorrect host: %v", o.Host)
	}
	if o.Port <= 0 {
		return fmt.Errorf("incorrect port value: %v", o.Port)
	}
	if o.ReadTimeout == 0 || o.WriteTimeout == 0 {
		return fmt.Errorf("timeouts can't equals zero")
	}
	return nil
}
