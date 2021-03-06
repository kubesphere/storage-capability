package main

import (
	"flag"
	snapclientset "github.com/kubernetes-csi/external-snapshotter/v2/pkg/client/clientset/versioned"
	snapinformers "github.com/kubernetes-csi/external-snapshotter/v2/pkg/client/informers/externalversions"
	"github.com/kubesphere/storage-capability/pkg/controller"
	crdclientset "github.com/kubesphere/storage-capability/pkg/generated/clientset/versioned"
	crdinformers "github.com/kubesphere/storage-capability/pkg/generated/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"time"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := controller.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	// Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	crdClient, err := crdclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building storage capability clientset: %s", err.Error())
	}
	list, err := crdClient.Discovery().ServerResourcesForGroupVersion("storage.kubesphere.io/v1alpha1")
	if err != nil {
		klog.Error(err)
	}
	klog.Info("list: ", list.String())

	snapClient, err := snapclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building snapshot clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	crdInformerFactory := crdinformers.NewSharedInformerFactory(crdClient, time.Second*30)
	snapInformerFactory := snapinformers.NewSharedInformerFactory(snapClient, time.Second*30)

	controller := controller.NewController(kubeClient, crdClient,
		kubeInformerFactory.Storage().V1().StorageClasses(),
		snapInformerFactory.Snapshot().V1beta1().VolumeSnapshotClasses(),
		crdInformerFactory.Storage().V1alpha1().ProvisionerCapabilities(),
		crdInformerFactory.Storage().V1alpha1().StorageClassCapabilities(),
	)

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	crdInformerFactory.Start(stopCh)
	snapInformerFactory.Start(stopCh)

	if err = controller.Start(stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
