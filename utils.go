package urknall

import (
	"crypto/sha256"
	"fmt"

	"github.com/dynport/urknall/cmd"
)

func renderTemplate(builder Template) (*packageImpl, error) {
	p := &packageImpl{reference: builder}
	e := validateTemplate(builder)
	if e != nil {
		return nil, e
	}
	builder.Render(p)
	return p, nil
}

func executeCommand(cmd cmd.Command, build *Build, checksumDir, taskName string) (e error) {
	sCmd := cmd.Shell()
	for _, env := range build.Env {
		sCmd = env + " " + sCmd
	}
	r := &remoteTaskRunner{build: build, cmd: sCmd, command: cmd, dir: checksumDir, taskName: taskName}
	return r.run()
}

func commandChecksum(c cmd.Command) (string, error) {
	if c, ok := c.(interface {
		Checksum() string
	}); ok {
		return c.Checksum(), nil
	}
	s := sha256.New()
	if _, e := s.Write([]byte(c.Shell())); e != nil {
		return "", e
	}
	return fmt.Sprintf("%x", s.Sum(nil)), nil
}
