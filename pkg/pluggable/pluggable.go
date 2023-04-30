package pluggable

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/commentcov/commentcov/pkg/common"
	"github.com/commentcov/commentcov/pkg/filepath"
	"github.com/commentcov/commentcov/proto"
)

var (
	// PluginHandshakeConfig holds configuration for plugin.GRPCPlugin.
	PluginHandshakeConfig = plugin.HandshakeConfig{
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "Hello",
	}
)

// GetPluginFromClient returns Pluggable object.
func GetPluginFromClient(client *plugin.Client, pluginName string) (Pluggable, error) {
	grpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	raw, err := grpcClient.Dispense(pluginName)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	p := raw.(Pluggable)
	return p, nil
}

// Consume aggregates CoverageItems from multi publishers.
func Consume(_ hclog.Logger, queue <-chan common.Pair[[]*proto.CoverageItem, error]) ([]*proto.CoverageItem, error) {
	items := make([]*proto.CoverageItem, 0)

	var err error
	for pair := range queue {
		if pair.V2 != nil {
			err = pair.V2
		}

		items = append(items, pair.V1...)
	}

	if err != nil {
		return []*proto.CoverageItem{}, fmt.Errorf("%w", err)
	}

	return items, nil
}

// Publish receive a list of target files and call the plugin MeasureCoverage logic.
func Publish(
	wg *sync.WaitGroup,
	_ hclog.Logger,
	p Pluggable,
	filenames []string,
	excluded filepath.ExcludeFileSet,
	queue chan<- common.Pair[[]*proto.CoverageItem, error],
) {
	defer wg.Done()

	items := make([]*proto.CoverageItem, 0)
	cis, err := p.MeasureCoverage(filenames)

	// Double check against the exclusion list.
	for _, ci := range cis {
		if _, ok := excluded[ci.File]; ok {
			continue
		}

		items = append(items, ci)
	}

	queue <- common.Pair[[]*proto.CoverageItem, error]{
		V1: items,
		V2: err,
	}
}
