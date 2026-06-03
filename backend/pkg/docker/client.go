package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient(host string) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) GetNodes(ctx context.Context) ([]types.Node, error) {
	nodes, err := c.cli.NodeList(ctx, types.NodeListOptions{})
	return nodes, err
}

func (c *Client) GetContainers(ctx context.Context) ([]types.Container, error) {
	containers, err := c.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	return containers, err
}

func (c *Client) GetServices(ctx context.Context) ([]types.Service, error) {
	services, err := c.cli.ServiceList(ctx, types.ServiceListOptions{})
	return services, err
}

func (c *Client) GetContainerStats(ctx context.Context, containerID string) (*types.StatsJSON, error) {
	stats, err := c.cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()
	var statsJSON := &types.StatsJSON{}
	return statsJSON, nil
}

func (c *Client) RestartContainer(ctx context.Context, containerID string) error {
	timeout := 10
	return c.cli.ContainerRestart(ctx, containerID, &timeout)
}

func (c *Client) DrainNode(ctx context.Context, nodeID string) error {
	nodeSpec := types.NodeSpec{
		Availability: types.NodeAvailabilityDrain,
	}
	_, err := c.cli.NodeUpdate(ctx, nodeID, types.NodeUpdateOptions{Version: types.Version{}, Spec: nodeSpec})
	return err
}

func (c *Client) Close() error {
	return c.cli.Close()
}
