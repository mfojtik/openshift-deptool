package show

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mfojtik/openshift-deptool/pkg/repository"
)

type ShowOptions struct {
	Repository       string
	UpstreamTag      string
	DownstreamBranch string
}

func NewShow() *cobra.Command {
	options := &ShowOptions{}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current dependency levels and list all cherry picks for repository",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().UTC().UnixNano())
			if err := options.Complete(); err != nil {
				klog.Error(err)
				os.Exit(1)
			}
			if err := options.Validate(); err != nil {
				klog.Error(err)
				os.Exit(1)
			}
			if err := options.Run(); err != nil {
				klog.Error(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&options.Repository, "repository", options.Repository, "Repository to analyze (eg. 'openshift/kubernetes-apimachinery').")
	cmd.Flags().StringVar(&options.UpstreamTag, "upstream-tag", options.UpstreamTag, "Upstream version tag name (eg. 'kubernetes-1.14.0').")
	cmd.Flags().StringVar(&options.DownstreamBranch, "downstream-branch", options.DownstreamBranch, "Downstream branch name (eg. 'oc-4.2-kubernetes-1.14.0').")
	return cmd
}

func (o *ShowOptions) Complete() error {
	return nil
}

func (o *ShowOptions) Validate() error {
	if len(o.Repository) == 0 {
		return fmt.Errorf("repository must be specified")
	}
	if len(o.UpstreamTag) == 0 {
		return fmt.Errorf("upstream tag must be specified")
	}
	if len(o.DownstreamBranch) == 0 {
		return fmt.Errorf("downstream branch must be specified")
	}
	return nil
}

func (o ShowOptions) Run() error {
	repo, err := repository.New("https://github.com/" + o.Repository)
	if err != nil {
		return err
	}
	commits, err := repo.ListUpstreamCommits(o.UpstreamTag, o.DownstreamBranch)
	if err != nil {
		return err
	}
	for _, commit := range commits {
		fmt.Fprintf(os.Stdout, "%s: %s\n", commit.Hash.String()[0:8], strings.Split(commit.Message, "\n")[0])
	}

	return nil
}
