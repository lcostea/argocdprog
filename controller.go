/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	apv1alpha1 "lcostea.io/knapps/pkg/apis/knapps/v1alpha1"
	clientset "lcostea.io/knapps/pkg/generated/clientset/versioned"
	samplescheme "lcostea.io/knapps/pkg/generated/clientset/versioned/scheme"
	informers "lcostea.io/knapps/pkg/generated/informers/externalversions/knapps/v1alpha1"
	listers "lcostea.io/knapps/pkg/generated/listers/knapps/v1alpha1"
)

const controllerAgentName = "argocdprog-controller"

// Controller is the controller implementation for Argocdprog resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	argocdprogLister listers.ArgocdprogLister
	argocdprogSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	argocdprogInformer informers.ArgocdprogInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		sampleclientset:  sampleclientset,
		argocdprogLister: argocdprogInformer.Lister(),
		argocdprogSynced: argocdprogInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Argocdprogs"),
		recorder:         recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Argocdprog resources change
	argocdprogInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueArgocdprog,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueArgocdprog(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Argocdprog controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.argocdprogSynced, c.argocdprogSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Argocdprog resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Argocdprog resource to be synced.
		if when, err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		} else if when != time.Duration(0) {
			c.workqueue.AddAfter(key, when)
		} else {
			// Finally we Forget this item so it does not
			// get queued again until another change happens.
			c.workqueue.Forget(obj)
		}

		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Argocdprog resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) (time.Duration, error) {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return time.Duration(0), err
	}

	// Get the Argocdprog resource with this namespace/name
	ap, err := c.argocdprogLister.Argocdprogs(namespace).Get(name)
	if err != nil {
		// The Argocdprog resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("Argocdprog '%s' in work queue no longer exists", key))
			return time.Duration(0), err
		}

		return time.Duration(0), err
	}
	if until, _ := newSchedule(ap.Spec.Schedule, ap.Status.ScheduleStatus); until > 0 {
		ap.Status.SyncStatus = ""
		klog.Infof("New schedule found for %s, status is reseting", ap.Name)
	}
	if ap.Status.SyncStatus == "" {
		reset := ""
		c.updateArgocdprogStatus(ap, apv1alpha1.SyncPending, &reset)
		klog.Infof("%s status is set to Pending", ap.Name)
	}

	switch ap.Status.SyncStatus {
	case apv1alpha1.SyncPending:
		at, err := timeUntilSchedule(ap.Spec.Schedule)
		if err != nil {
			return time.Duration(0), err
		}
		if at > 0 {
			return at, nil
		}
		c.updateArgocdprogStatus(ap, apv1alpha1.SyncRunning, nil)
		klog.Infof("%s status is set to Running", ap.Name)
	case apv1alpha1.SyncRunning:
		statusCode := executeSyncRequest(ap.Spec.SyncApp, ap.Spec.ApiServer)
		at := time.Now().UTC().Format("2006-01-02T15:04:05")
		syncStatus := apv1alpha1.SyncDone
		if statusCode == 200 {
			klog.Infof("Sync call to %s was successfull", ap.Spec.SyncApp)
			klog.Infof("%s status is set to Done", ap.Name)
		} else {
			klog.Infof("%s status is set to Error", ap.Name)
			syncStatus = apv1alpha1.SyncError
		}
		c.updateArgocdprogStatus(ap, syncStatus, &at)
	default:
		return time.Duration(0), nil
	}

	return time.Duration(0), nil
}

// enqueueArgocdprog takes a Argocdprog resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Argocdprog.
func (c *Controller) enqueueArgocdprog(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Argocdprog resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Argocdprog resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Argocdprog, we should not do anything more
		// with it.
		if ownerRef.Kind != "Argocdprog" {
			return
		}

		Argocdprog, err := c.argocdprogLister.Argocdprogs(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.Infof("ignoring orphaned object '%s' of Argocdprog '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueArgocdprog(Argocdprog)
		return
	}
}

func timeUntilSchedule(schedule string) (time.Duration, error) {
	now := time.Now()
	layout := "2006-01-02T15:04:05"
	at, err := time.Parse(layout, schedule)
	if err != nil {
		return time.Duration(0), err
	}
	return at.Sub(now), nil
}

func newSchedule(schedule string, last string) (time.Duration, error) {
	layout := "2006-01-02T15:04:05"
	prev, err := time.Parse(layout, last)
	if err != nil {
		return time.Duration(0), err
	}
	next, err := time.Parse(layout, schedule)
	if err != nil {
		return time.Duration(0), err
	}
	return next.Sub(prev), nil
}

func executeSyncRequest(appName string, apiServer string) int {
	syncTkn := os.Getenv("ARGOCDPROG_SYNC_TOKEN")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //because the server is using a self generated cert
	}
	client := http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}
	jsonBody := []byte("{\"name\": \"" + appName + "\",\"prune\": false}'")

	url := apiServer + "/api/v1/applications/" + appName + "/sync"
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+syncTkn)

	resp, err := client.Do(request)
	if err != nil {
		klog.Errorln(err)
		return 0
	}

	defer resp.Body.Close()
	return resp.StatusCode
}

func (c *Controller) updateArgocdprogStatus(ap *apv1alpha1.Argocdprog, status string, at *string) {
	apCopy := ap.DeepCopy()
	apCopy.Status.SyncStatus = status
	if at != nil {
		apCopy.Status.ScheduleStatus = *at
	}

	_, err := c.sampleclientset.KnappsV1alpha1().Argocdprogs(ap.Namespace).UpdateStatus(apCopy)
	if err != nil {
		klog.Error(err)
	}
}
