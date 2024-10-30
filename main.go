package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	config2 "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg"
	k8sv1 "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/config"
	// +kubebuilder:scaffold:imports
)

var (
	// Version of the application
	Version = "0.0.0"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	// These variables are here to avoid errors during leader election.
	leaseDuration = time.Second * 60
	renewDeadline = time.Second * 40
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(k8sv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	ctx := ctrl.SetupSignalHandler()
	restConfig := config2.GetConfigOrDie()
	operatorConfig, err := config.NewOperatorConfig(Version)
	if err != nil {
		setupLog.Error(err, "unable to create operator config")
		os.Exit(1)
	}
	err = startOperator(ctx, restConfig, operatorConfig, flag.CommandLine, os.Args)
	if err != nil {
		setupLog.Error(err, "unable to start operator")
		os.Exit(1)
	}
}

func startOperator(
	ctx context.Context,
	restConfig *rest.Config,
	operatorConfig *config.OperatorConfig,
	flags *flag.FlagSet,
	args []string,
) error {
	k8sManager, err := NewK8sManager(restConfig, operatorConfig, flags, args)
	if err != nil {
		return fmt.Errorf("unable to start manager: %w", err)
	}

	var recorder eventRecorder = k8sManager.GetEventRecorderFor("k8s-blueprint-operator")
	bootstrap, err := pkg.Bootstrap(restConfig, recorder, operatorConfig.Namespace)
	if err != nil {
		return fmt.Errorf("unable to bootstrap application context: %w", err)
	}

	err = configureManager(k8sManager, bootstrap.Reconciler)
	if err != nil {
		return fmt.Errorf("unable to configure manager: %w", err)
	}

	err = startK8sManager(ctx, k8sManager)
	if err != nil {
		return fmt.Errorf("unable to start operator: %w", err)
	}
	return err
}

func NewK8sManager(
	restConfig *rest.Config,
	operatorConfig *config.OperatorConfig,
	flags *flag.FlagSet, args []string,
) (manager.Manager, error) {
	options := getK8sManagerOptions(flags, args, operatorConfig)
	return ctrl.NewManager(restConfig, options)
}

func configureManager(k8sManager controllerManager, blueprintReconciler *reconciler.BlueprintReconciler) error {
	err := blueprintReconciler.SetupWithManager(k8sManager)
	if err != nil {
		return fmt.Errorf("unable to configure reconciler: %w", err)
	}

	err = addChecks(k8sManager)
	if err != nil {
		return fmt.Errorf("unable to add checks to the manager: %w", err)
	}

	return nil
}

func getK8sManagerOptions(flags *flag.FlagSet, args []string, operatorConfig *config.OperatorConfig) ctrl.Options {
	controllerOpts := ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{DefaultNamespaces: map[string]cache.Config{
			operatorConfig.Namespace: {},
		}},
		WebhookServer:    webhook.NewServer(webhook.Options{Port: 9443}),
		LeaderElectionID: "ae48821c.cloudogu.com",
		LeaseDuration:    &leaseDuration,
		RenewDeadline:    &renewDeadline,
	}
	controllerOpts, zapOpts := parseManagerFlags(flags, args, controllerOpts)

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zapOpts)))

	return controllerOpts
}

func parseManagerFlags(flags *flag.FlagSet, args []string, ctrlOpts ctrl.Options) (ctrl.Options, zap.Options) {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flags.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flags.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flags.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	zapOpts := zap.Options{
		Development: config.IsStageDevelopment(),
	}
	zapOpts.BindFlags(flags)
	// Ignore errors; flags is set to exit on errors
	_ = flags.Parse(args)

	ctrlOpts.Metrics = metricsserver.Options{BindAddress: metricsAddr}
	ctrlOpts.HealthProbeBindAddress = probeAddr
	ctrlOpts.LeaderElection = enableLeaderElection

	return ctrlOpts, zapOpts
}

func addChecks(k8sManager controllerManager) error {
	if err := k8sManager.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}
	if err := k8sManager.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %w", err)
	}

	return nil
}

func startK8sManager(ctx context.Context, k8sManager controllerManager) error {
	logger := log.FromContext(ctx).WithName("k8s-manager-start")
	logger.Info("starting manager")
	if err := k8sManager.Start(ctx); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}
