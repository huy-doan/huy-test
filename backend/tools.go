//go:build tools
// +build tools

package tools

import (
	_ "github.com/go-playground/validator/v10"
	_ "github.com/graph-gophers/dataloader"
	_ "github.com/graph-gophers/dataloader/v7"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/urfave/cli/v2"
)
