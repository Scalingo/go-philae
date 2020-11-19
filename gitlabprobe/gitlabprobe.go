package gitlabprobe

import (
	"github.com/Scalingo/go-philae/v4/statusioprobe"
)

func NewGitlabProbe(name string) statusioprobe.StatusIOProbe {
	return statusioprobe.NewStatusIOProbe(name, "5b36dc6502d06804c08349f7")
}
