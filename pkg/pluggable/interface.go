package pluggable

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/commentcov/commentcov/proto"
)

// Pluggable is the interface.
// Each plugin implements Pluggable interface via gRPC.
type Pluggable interface {
	MeasureCoverage(files []string) ([]*proto.CoverageItem, error)
}

// CommentcovPlugin implements plugin.GRPCPlugin.
type CommentcovPlugin struct {
	plugin.Plugin

	Impl Pluggable
}

// GRPCServer is a part of plugin.GRPCPlugin interface.
func (p *CommentcovPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCommentcovPluginServer(
		s,
		&GRPCServer{
			Impl: p.Impl,
		},
	)

	return nil
}

// GRPCClient is a part of plugin.GRPCPlugin interface.
func (p *CommentcovPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	c := &GRPCClient{
		client: proto.NewCommentcovPluginClient(conn),
	}

	return c, nil
}

// GRPCServer Dependency(Plugin Implementation) Injected Object.
type GRPCServer struct {
	proto.UnimplementedCommentcovPluginServer

	Impl Pluggable
}

// MeasureCoverage calls the Pluggable implementation via gRPC.
func (s *GRPCServer) MeasureCoverage(_ context.Context, in *proto.MeasureCoverageIn) (*proto.MeasureCoverageOut, error) {
	cis, err := s.Impl.MeasureCoverage(in.Files)
	if err != nil {
		out := &proto.MeasureCoverageOut{
			CoverageItems: []*proto.CoverageItem{},
		}

		return out, fmt.Errorf("%w", err)
	}

	out := &proto.MeasureCoverageOut{
		CoverageItems: cis,
	}

	return out, nil
}

// GRPCClient is the interface to be called from the host.
type GRPCClient struct {
	client proto.CommentcovPluginClient
}

// MeasureCoverage implements Pluggable.
func (c *GRPCClient) MeasureCoverage(files []string) ([]*proto.CoverageItem, error) {
	res, err := c.client.MeasureCoverage(context.Background(), &proto.MeasureCoverageIn{
		Files: files,
	})
	if err != nil {
		ci := []*proto.CoverageItem{}
		return ci, fmt.Errorf("%w", err)
	}

	return res.CoverageItems, nil
}
