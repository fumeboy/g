package example

import (
	"g"
	"g/gen/errgen"
	"g/gen/recovergen"
	"g/gen/wiregen"
	"os"
	"testing"
)

func TestGen(_ *testing.T) {
	wd, _ := os.Getwd()
	gg := g.G(wd, os.Environ(), "", nil)
	gg.Gen(errgen.Gen, recovergen.Gen, wiregen.Gen)
}
