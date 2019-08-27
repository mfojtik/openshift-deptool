module github.com/mfojtik/openshift-deptool

go 1.12

require (
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	gopkg.in/src-d/go-git.v4 v4.13.1
	k8s.io/client-go v0.0.0-20190822053941-f4e58ce6093c
	k8s.io/component-base v0.0.0-20190823013255-e3d4ac5c99fb
	k8s.io/klog v0.4.0
)

replace golang.org/x/net => golang.org/x/net v0.0.0-20190822053941-7f726cade0ab
