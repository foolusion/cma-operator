package sds_package_manager

import (
	"fmt"
	"github.com/samsung-cnct/cma-operator/pkg/util/cma"
	"github.com/samsung-cnct/cma-operator/pkg/util/cmagrpc"
	"github.com/samsung-cnct/cma-operator/pkg/util/sds/callback"
	"github.com/spf13/viper"
	"time"

	"k8s.io/apimachinery/pkg/fields"

	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/juju/loggo"
	api "github.com/samsung-cnct/cma-operator/pkg/apis/cma/v1alpha1"
	"github.com/samsung-cnct/cma-operator/pkg/generated/cma/client/clientset/versioned"
	"github.com/samsung-cnct/cma-operator/pkg/util"
	"github.com/samsung-cnct/cma-operator/pkg/util/helmutil"
	"github.com/samsung-cnct/cma-operator/pkg/util/k8sutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

var (
	logger loggo.Logger
)

const (
	WaitForClusterChangeMaxTries         = 3
	WaitForClusterChangeTimeInterval     = 5 * time.Second
	KubernetesNamespaceViperVariableName = "kubernetes-namespace"
	ClusterRequestIDAnnotation           = "requestID"
	ClusterCallbackURLAnnotation         = "callbackURL"
)

type SDSPackageManagerController struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller

	client        *versioned.Clientset
	cmaGRPCClient cmagrpc.ClientInterface
}

func NewSDSPackageManagerController(config *rest.Config) (output *SDSPackageManagerController, err error) {
	cmaGRPCClient, err := cmagrpc.CreateNewDefaultClient()
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = k8sutil.DefaultConfig
	}
	client := versioned.NewForConfigOrDie(config)

	// create sdscluster list/watcher
	sdsPackageManagerListWatcher := cache.NewListWatchFromClient(
		client.CmaV1alpha1().RESTClient(),
		api.SDSPackageManagerResourcePlural,
		viper.GetString(KubernetesNamespaceViperVariableName),
		fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the SDSCluster key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the SDSPackageManager than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(sdsPackageManagerListWatcher, &api.SDSPackageManager{}, 30*time.Second, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	output = &SDSPackageManagerController{
		informer:      informer,
		indexer:       indexer,
		queue:         queue,
		client:        client,
		cmaGRPCClient: cmaGRPCClient,
	}
	output.SetLogger()
	return
}

func (c *SDSPackageManagerController) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two SDSClusters with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.processItem(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}

// processItem is the business logic of the controller.
func (c *SDSPackageManagerController) processItem(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a SDSPackageManager, so that we will see a delete for one SDSPackageManager
		fmt.Printf("SDSPackageManager -->%s<-- does not exist anymore\n", key)
	} else {
		SDSPackageManager := obj.(*api.SDSPackageManager)
		clusterName := SDSPackageManager.GetName()
		fmt.Printf("SDSPackageManager -->%s<-- does exist (name=%s)!\n", key, clusterName)

		switch SDSPackageManager.Status.Phase {
		case api.PackageManagerPhaseNone, api.PackageManagerPhasePending:
			c.deployTiller(SDSPackageManager)
			break
		case api.PackageManagerPhaseInstalling:
			c.waitForTiller(SDSPackageManager)
			break
		}
	}
	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *SDSPackageManagerController) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		glog.Infof("Error syncing packageManager %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping packageManager %q out of the queue: %v", key, err)
}

func (c *SDSPackageManagerController) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	glog.Info("Starting SDSCluster controller")

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping SDSCluster controller")
}

func (c *SDSPackageManagerController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *SDSPackageManagerController) deployTiller(packageManager *api.SDSPackageManager) (bool, error) {
	sdsCluster, err := cma.GetSDSCluster(packageManager.Spec.Cluster.Name, viper.GetString(KubernetesNamespaceViperVariableName), nil)
	if err != nil {
		logger.Infof("Error trying to package manager -->%s<-- could not get sdscluster becuase of %s", packageManager.Spec.Cluster.Name, err)
		return false, err
	}
	if sdsCluster.Status.Phase != api.ClusterPhaseReady {
		logger.Infof("Could not deploy package manager -->%s<-- as cluster is not ready, in phase -->%s<--", sdsCluster.Name, sdsCluster.Status.Phase)
		return false, err
	}

	config, err := c.getRestConfigForRemoteCluster(packageManager.Spec.Cluster.Name, packageManager.Namespace, nil)
	if err != nil {
		return false, err
	}

	k8sutil.CreateNamespace(k8sutil.GenerateNamespace(packageManager.Spec.Namespace), config)
	k8sutil.CreateServiceAccount(k8sutil.GenerateServiceAccount("tiller-sa"), packageManager.Spec.Namespace, config)
	if packageManager.Spec.Permissions.ClusterWide {
		k8sutil.CreateClusterRole(helmutil.GenerateClusterAdminRole("tiller-"+packageManager.Spec.Namespace), config)
		k8sutil.CreateClusterRoleBinding(k8sutil.GenerateSingleClusterRolebinding("tiller-"+packageManager.Spec.Namespace, "tiller-sa", packageManager.Spec.Namespace, "tiller-"+packageManager.Spec.Namespace), config)
	} else {
		logger.Infof("Not cluster wide")
		namespaces := append(packageManager.Spec.Permissions.Namespaces, packageManager.Spec.Namespace)
		for _, namespace := range namespaces {
			logger.Infof("Creating namespace %s", namespace)
			k8sutil.CreateNamespace(k8sutil.GenerateNamespace(namespace), config)
			k8sutil.CreateRole(helmutil.GenerateAdminRole(packageManager.Spec.Namespace+"-tiller"), namespace, config)
			k8sutil.CreateRoleBinding(k8sutil.GenerateSingleRolebinding(packageManager.Spec.Namespace+"-tiller", "tiller-sa", packageManager.Spec.Namespace, packageManager.Spec.Namespace+"-tiller"), namespace, config)
		}
	}
	k8sutil.CreateJob(helmutil.GenerateTillerInitJob(
		helmutil.TillerInitOptions{
			BackoffLimit:   4,
			Name:           "tiller-install-job",
			Namespace:      packageManager.Spec.Namespace,
			ServiceAccount: "tiller-sa",
			Version:        packageManager.Spec.Version}), packageManager.Spec.Namespace, config)

	if packageManager.Annotations[ClusterCallbackURLAnnotation] != "" {
		// We need to notify someone that the package manager is being deployed(again)
		message := &sdscallback.ClusterMessage{
			State:        sdscallback.ClusterMessageStateInProgress,
			StateText:    api.PackageManagerPhaseInstalling,
			ProgressRate: 0,
		}
		sdscallback.DoCallback(packageManager.Annotations[ClusterCallbackURLAnnotation], message)
	}

	packageManager.Status.Phase = api.PackageManagerPhaseInstalling
	_, err = c.client.CmaV1alpha1().SDSPackageManagers(packageManager.Namespace).Update(packageManager)
	if err == nil {
		logger.Infof("Deployed tiller on -->%s<--", packageManager.Spec.Name)
	} else {
		logger.Infof("Could not update the status error was %s", err)
	}

	return true, nil
}

func retrieveClusterRestConfig(name string, kubeconfig string) (*rest.Config, error) {
	// Let's create a tempfile and line it up for removal
	file, err := ioutil.TempFile(os.TempDir(), "kraken-kubeconfig")
	defer os.Remove(file.Name())
	file.WriteString(kubeconfig)

	clusterConfig, err := clientcmd.BuildConfigFromFlags("", file.Name())
	if os.Getenv("CLUSTERMANAGERAPI_INSECURE_TLS") == "true" {
		clusterConfig.TLSClientConfig = rest.TLSClientConfig{Insecure: true}
	}

	if err != nil {
		logger.Errorf("Could not load kubeconfig for cluster -->%s<--", name)
		return nil, err
	}
	return clusterConfig, nil
}

func (c *SDSPackageManagerController) getRestConfigForRemoteCluster(clusterName string, namespace string, config *rest.Config) (*rest.Config, error) {
	sdscluster, err := c.client.CmaV1alpha1().SDSClusters(viper.GetString(KubernetesNamespaceViperVariableName)).Get(clusterName, v1.GetOptions{})
	if err != nil {
		glog.Errorf("Failed to retrieve SDSCluster -->%s<--, error was: %s", clusterName, err)
		return nil, err
	}
	cluster, err := c.cmaGRPCClient.GetCluster(cmagrpc.GetClusterInput{Name: clusterName, Provider: sdscluster.Spec.Provider})
	if err != nil {
		glog.Errorf("Failed to retrieve Cluster Status -->%s<--, error was: %s", clusterName, err)
		return nil, err
	}
	if cluster.Kubeconfig == "" {
		glog.Errorf("Could not install tiller yet for cluster -->%s<-- cluster is not ready, status is -->%s<--", cluster.Name, cluster.Status)
		return nil, err
	}

	remoteConfig, err := retrieveClusterRestConfig(clusterName, cluster.Kubeconfig)
	if err != nil {
		glog.Errorf("Could not install tiller yet for cluster -->%s<-- cluster is not ready, error is: %v", clusterName, err)
		return nil, err
	}

	return remoteConfig, nil
}

func (c *SDSPackageManagerController) SetLogger() {
	logger = util.GetModuleLogger("pkg.controllers.sds_package_manager", loggo.INFO)
}

func (c *SDSPackageManagerController) waitForTiller(packageManager *api.SDSPackageManager) (result bool, err error) {
	config, err := c.getRestConfigForRemoteCluster(packageManager.Spec.Cluster.Name, packageManager.Namespace, nil)
	if err != nil {
		return false, err
	}

	clientset, _ := kubernetes.NewForConfig(config)
	timeout := 0
	for timeout < 2000 {
		job, err := clientset.BatchV1().Jobs(packageManager.Spec.Namespace).Get("tiller-install-job", v1.GetOptions{})
		if err == nil {
			if job.Status.Succeeded > 0 {
				packageManager.Status.Phase = api.PackageManagerPhaseImplemented
				packageManager.Status.Ready = true
				_, err = c.client.CmaV1alpha1().SDSPackageManagers(packageManager.Namespace).Update(packageManager)
				if err == nil {
					logger.Infof("Tiller running on -->%s<--", packageManager.Spec.Name)
					c.handleHavingPackageManager(packageManager)
				} else {
					logger.Infof("Could not update the status error was %s", err)
				}
				return true, nil
			}
		}
		time.Sleep(5 * time.Second)
		timeout++
	}
	return false, nil
}

func (c *SDSPackageManagerController) handleHavingPackageManager(packageManager *api.SDSPackageManager) {
	if packageManager.Annotations[ClusterCallbackURLAnnotation] != "" {
		// We need to notify someone that the package manager is being deployed(again)
		message := &sdscallback.ClusterMessage{
			State:        sdscallback.ClusterMessageStateCompleted,
			StateText:    api.PackageManagerPhaseImplemented,
			ProgressRate: 100,
		}
		sdscallback.DoCallback(packageManager.Annotations[ClusterCallbackURLAnnotation], message)
	}

}
