package runner

import (
	"github.com/jaesbit/ffuf/pkg/ffuf"
)

func NewRunnerByName(name string, conf *ffuf.Config) ffuf.RunnerProvider {
	// We have only one Runner at the moment
	return NewSimpleRunner(conf)
}
