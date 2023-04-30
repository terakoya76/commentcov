package execute

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/commentcov/commentcov/pkg/common"
	"github.com/commentcov/commentcov/pkg/config"
	"github.com/commentcov/commentcov/pkg/filepath"
	"github.com/commentcov/commentcov/pkg/pluggable"
	"github.com/commentcov/commentcov/pkg/report"
	"github.com/commentcov/commentcov/proto"
)

const (
	// pluginFileBatch determines how many files each plugin process will process.
	pluginFileBatch int = 100
)

func Run(cfg *config.ViperConfig) error {
	logger := hclog.New(&hclog.LoggerOptions{
		Output:     hclog.DefaultOutput,
		Level:      hclog.Trace,
		Name:       common.ProjectName,
		JSONFormat: true,
	})

	excluded, err := filepath.NewExcludeFileSet(cfg.ExcludePaths)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fileset, err := filepath.Extract(cfg.TargetPath, excluded)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	pWG := &sync.WaitGroup{}
	queue := make(chan common.Pair[[]*proto.CoverageItem, error])
	for _, ext := range fileset.Extensions() {
		// cf. https://stackoverflow.com/questions/45617758/defer-in-the-loop-what-will-be-better
		err = func() error {
			pc := pluggable.FindPluginConfig(cfg.Plugins, ext)
			if pc == nil {
				return nil
			}

			pluginName := pc.Name()
			logger.Info("Install Plugin", "plugin", pluginName)
			if err = pc.Install(); err != nil {
				return fmt.Errorf("%w", err)
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

			var p pluggable.Pluggable
			p, err = pluggable.GetPluginFromClient(client, pluginName)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			batched := common.Batched(fileset.Files(ext), pluginFileBatch)
			pWG.Add(len(batched))

			for _, files := range batched {
				go pluggable.Publish(pWG, logger, p, files, excluded, queue)
			}

			return nil
		}()

		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	go func() {
		pWG.Wait()
		close(queue)
	}()

	items, err := pluggable.Consume(logger, queue)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	mode, err := report.StringToMode(cfg.Mode)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	report.Report(mode, items)
	return nil
}
