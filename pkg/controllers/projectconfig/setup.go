package projectconfig

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"

	"github.com/openshift/hypershift-toolkit/pkg/cmd/cpoperator"
	"github.com/openshift/hypershift-toolkit/pkg/controllers"
)

const (
	ConfigNamespace = "openshift-config"
)

func Setup(cfg *cpoperator.ControlPlaneOperatorConfig) error {
	if err := setupProjectConfigObserver(cfg); err != nil {
		return err
	}

	if err := setupControllerManagerCAUpdater(cfg); err != nil {
		return err
	}

	return nil
}

func setupProjectConfigObserver(cfg *cpoperator.ControlPlaneOperatorConfig) error {
	openshiftClient, err := configclient.NewForConfig(cfg.TargetConfig())
	if err != nil {
		return err
	}
	informerFactory := configinformers.NewSharedInformerFactory(openshiftClient, controllers.DefaultResync)
	cfg.Manager().Add(manager.RunnableFunc(func(stopCh <-chan struct{}) error {
		informerFactory.Start(stopCh)
		return nil
	}))
	configMaps := informerFactory.Core().V1().ConfigMaps()
	reconciler := &ProjectConfigObserver{
		Client:         cfg.Manager().GetClient(),
		TargetCMLister: projects.Lister(),
		Namespace:      cfg.Namespace(),
		Log:            cfg.Logger().WithName("ManagedCAObserver"),
	}
	c, err := controller.New("ca-configmap-observer", cfg.Manager(), controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}
	if err := c.Watch(&source.Informer{Informer: configMaps.Informer()}, controllers.NamedResourceHandler(RouterCAConfigMap, ServiceCAConfigMap)); err != nil {
		return err
	}
	return nil
}

func setupControllerManagerCAUpdater(cfg *cpoperator.ControlPlaneOperatorConfig) error {
	reconciler := &ControllerManagerCAUpdater{
		Client:    cfg.Manager().GetClient(),
		Namespace: cfg.Namespace(),
		Log:       cfg.Logger().WithName("ControllerManagerCAUpdater"),
		InitialCA: cfg.InitialCA(),
	}
	c, err := controller.New("controller-manager-ca-updater", cfg.Manager(), controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}
	if err := c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, controllers.NamedResourceHandler(ControllerManagerAdditionalCAConfigMap)); err != nil {
		return err
	}
	return nil
}
