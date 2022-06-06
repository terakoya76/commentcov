package cmd

import (
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"

	"github.com/commentcov/commentcov/pkg/common"
	"github.com/commentcov/commentcov/pkg/filepath"
	"github.com/commentcov/commentcov/pkg/pluggable"
	"github.com/commentcov/commentcov/pkg/report"
	"github.com/commentcov/commentcov/proto"
)

const (
	// pluginFileBatch determines how many files each plugin process will process.
	pluginFileBatch int = 100
)

// coverageCmd returns coverage info.
var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "generate coverage reports",
	Long:  "generate coverage reports",
	Run: func(cmd *cobra.Command, args []string) {
		logger := hclog.New(&hclog.LoggerOptions{
			Output:     hclog.DefaultOutput,
			Level:      hclog.Trace,
			Name:       "commentcov",
			JSONFormat: true,
		})

		fileset, err := filepath.Extract(cfg.TargetPath, cfg.ExcludePaths)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		pWG := &sync.WaitGroup{}
		queue := make(chan common.Pair[[]*proto.CoverageItem, error])
		for _, ext := range fileset.Extensions() {
			pc := pluggable.FindPluginConfig(cfg.Plugins, ext)
			if pc == nil {
				continue
			}

			pluginName := pc.Name()
			logger.Info("Install Plugin", "plugin", pluginName)
			if err = pc.Install(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			cmd := pc.GetExecuteCommand()
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: pluggable.PluginHandshakeConfig,
				VersionedPlugins: map[int]plugin.PluginSet{
					1: {
						pluginName: &pluggable.CommentcovPlugin{},
					},
				},
				Cmd: cmd,
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolGRPC,
				},
				Logger: logger,
			})
			defer client.Kill()

			p, err := pluggable.GetPluginFromClient(client, pluginName)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			batched := common.Batched(fileset.Files(ext), pluginFileBatch)
			pWG.Add(len(batched))

			for _, files := range batched {
				go pluggable.Publish(pWG, logger, p, files, queue)
			}
		}

		go func() {
			pWG.Wait()
			close(queue)
		}()

		items, err := pluggable.Consume(logger, queue)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		mode, err := report.StringToMode(cfg.Mode)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		report.Report(mode, items)
		os.Exit(0)
	},
}
