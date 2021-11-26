package dola_test

import (
	"context"

	"github.com/numus-digital/dola"
)

func ExampleKeep() {
	keep, _ := dola.NewKeepBuilder().Build(context.Background())
	keep.Root.Add("verbose", dola.VerboseStrategy{}) //nolint:exhaustivestruct
	keep.Run(context.Background())
}
