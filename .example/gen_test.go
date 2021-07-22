package example

import (
	"github.com/fumeboy/g"
	"github.com/fumeboy/g/gen/errgen"
	"github.com/fumeboy/g/gen/recovergen"
	"github.com/fumeboy/g/gen/wiregen"
	"os"
	"testing"
)

func TestGen(_ *testing.T) {
	wd, _ := os.Getwd()
	gg := g.G(wd, os.Environ(), "", nil)
	gg.Gen(errgen.Gen, recovergen.Gen, wiregen.Gen)
}
