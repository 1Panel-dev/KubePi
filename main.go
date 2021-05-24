package main

import (
	"embed"
	"github.com/KubeOperator/ekko/cmd"
	"github.com/KubeOperator/ekko/pkg/route"
	"github.com/spf13/cobra"
)

//go:embed web
var embedWeb embed.FS

func main() {
	route.EmbedWeb = embedWeb
	cobra.CheckErr(cmd.RootCmd.Execute())

}
