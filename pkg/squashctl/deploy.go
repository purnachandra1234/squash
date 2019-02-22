package squashctl

import (
	"fmt"
	"strings"

	"github.com/solo-io/squash/pkg/demo"
	"github.com/solo-io/squash/pkg/install"
	"github.com/spf13/cobra"
)

var defaultDemoNamespace = "default"

func (top *Options) DeployCmd(o *Options) *cobra.Command {
	dOpts := &o.DeployOptions
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the squash agent or a demo microservice",
	}

	cmd.AddCommand(
		top.deployDemoCmd(&dOpts.DemoOptions),
		top.deployAgentCmd(&dOpts.AgentOptions),
	)

	return cmd
}

func (top *Options) deployDemoCmd(demoOpts *DemoOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "deploy a demo microservice",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := top.ensureDemoDeployOpts(demoOpts); err != nil {
				return err
			}
			switch demoOpts.DemoId {
			case demo.DemoGoGo:
				return demo.DeployGoGo(top.KubeClient, demoOpts.Namespace1, demoOpts.Namespace2)
			case demo.DemoGoJava:
				return demo.DeployGoJava(top.KubeClient, demoOpts.Namespace1, demoOpts.Namespace2)
			default:
				return fmt.Errorf("Please choose a valid demo option: %v", strings.Join(demo.DemoIds, ", "))
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&demoOpts.Namespace1, "demo-namespace1", "", "namespace in which to install the sample app")
	f.StringVar(&demoOpts.Namespace2, "demo-namespace2", "", "(optional) ns for second app - defaults to 'namespace' flag's value")
	f.StringVar(&demoOpts.DemoId, "demo-id", "", "which sample microservice to deploy. Options: go-go, go-java")
	return cmd
}

func (top *Options) ensureDemoDeployOpts(dOpts *DemoOptions) error {
	if dOpts.Namespace1 == "" {
		if top.Squash.Machine {
			dOpts.Namespace1 = defaultDemoNamespace
		} else {
			top.chooseAllowedNamespace(&dOpts.Namespace1, "Select a namespace for service 1.")
		}
	}
	if dOpts.Namespace2 == "" {
		if top.Squash.Machine {
			dOpts.Namespace2 = dOpts.Namespace1
		} else {
			top.chooseAllowedNamespace(&dOpts.Namespace2, "Select a namespace for service 2.")
		}
	}
	if dOpts.DemoId == "" {
		if top.Squash.Machine {
			dOpts.DemoId = defaultDemoNamespace
		} else {
			top.chooseString("Choose a demo microservice to deploy", &dOpts.DemoId, demo.DemoIds)
		}
	}
	return nil
}

func (top *Options) deployAgentCmd(agentOpts *AgentOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "deploy a squash agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := top.ensureAgentDeployOpts(agentOpts); err != nil {
				return err
			}
			return install.InstallAgent(top.KubeClient, agentOpts.Namespace, agentOpts.Preview)
		},
	}
	f := cmd.Flags()
	f.StringVar(&agentOpts.Namespace, "agentNamespace", install.DefaultNamespace, "namespace in which to install Squash")
	f.BoolVar(&agentOpts.Preview, "preview", false, "If set, prints Squash installation yaml without installing Squash.")
	return cmd
}
func (top *Options) ensureAgentDeployOpts(dOpts *AgentOptions) error {
	// TODO(mitchdraft) - interactive mode
	return nil
}
