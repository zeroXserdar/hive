package clients

import (
	"errors"
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/hive-taiko-clients/clients"
	"net"
)

var _ clients.ManagedClient = &HiveManagedClient{}

type HiveOptionsGenerator func() ([]hivesim.StartOption, error)

type HiveManagedClient struct {
	T                    *hivesim.T
	OptionsGenerator     HiveOptionsGenerator
	HiveClientDefinition *hivesim.ClientDefinition

	HiveClient        *hivesim.Client
	extraStartOptions []hivesim.StartOption
}

func (h *HiveManagedClient) IsRunning() bool {
	return h.HiveClient != nil
}

func (h *HiveManagedClient) Start() error {
	h.T.Logf("Starting client %s", h.ClientType())
	opts, err := h.OptionsGenerator()
	//h.T.Logf("With first Option %v", opts[0])
	h.T.Logf("With Options from Generator %v", opts)
	h.T.Logf("With Extra Start Options %v", h.extraStartOptions)
	if err != nil {
		return fmt.Errorf("unable to get start options: %v", err)
	}

	if opts == nil {
		opts = make([]hivesim.StartOption, 0)
	}

	if h.extraStartOptions != nil {
		opts = append(opts, h.extraStartOptions...)
	}

	h.T.Logf("With Name %s", h.HiveClientDefinition.Name)
	h.T.Logf("With Final Options %v", opts)
	h.HiveClient = h.T.StartClient(h.HiveClientDefinition.Name, opts...)
	if h.HiveClient == nil {
		return fmt.Errorf("unable to launch client")
	}
	h.T.Logf(
		"Started client %s, container %s",
		h.ClientType(),
		h.HiveClient.Container,
	)
	return nil
}

func (h *HiveManagedClient) AddStartOption(opts ...interface{}) {
	if h.extraStartOptions == nil {
		h.extraStartOptions = make([]hivesim.StartOption, 0)
	}
	for _, o := range opts {
		if o, ok := o.(hivesim.StartOption); ok {
			h.extraStartOptions = append(h.extraStartOptions, o)
		}
	}
}

func (h *HiveManagedClient) GetIP() net.IP {
	if h.HiveClient == nil {
		return net.IP{}
	}
	return h.HiveClient.IP
}

func (h *HiveManagedClient) Shutdown() error {
	if err := h.T.Sim.StopClient(h.T.SuiteID, h.T.TestID, h.HiveClient.Container); err != nil {
		return err
	}
	h.HiveClient = nil
	return nil
}

func (h *HiveManagedClient) GetEnodeURL() (string, error) {
	return h.HiveClient.EnodeURL()
}

func (h *HiveManagedClient) ClientType() string {
	return h.HiveClientDefinition.Name
}

func (h *HiveManagedClient) GetHost() string {
	if h.HiveClient == nil {
		return ""
	}
	return h.HiveClient.IP.String()
}

func (h *HiveManagedClient) GetAddress() string {
	if h.HiveClient == nil {
		return ""
	}
	return h.HiveClient.IP.String()
}

func (h *HiveManagedClient) GetEnvVar(t *hivesim.T, testSuite hivesim.SuiteID, test hivesim.TestID, node string, network string, varName string) (string, error) {
	resp, err := t.Sim.ClientExec(testSuite, test, node, []string{fmt.Sprintf("cat /saved_env.txt | grep %s | cut -d '=' -f2-", varName)})
	if err != nil {
		return "", err
	}
	if resp.ExitCode != 0 {
		return "", errors.New("unexpected exit code for getting Env Var")
	}

	output := resp.Stdout

	return output, nil
}
