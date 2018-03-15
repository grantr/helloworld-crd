package bar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	fakekubeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1"
	clientset "k8s.io/sample-controller/pkg/client/clientset/versioned"
	fakeclientset "k8s.io/sample-controller/pkg/client/clientset/versioned/fake"
	informers "k8s.io/sample-controller/pkg/client/informers/externalversions"
)

func getTestBar() *v1alpha1.Bar {
	return &v1alpha1.Bar{
		ObjectMeta: metav1.ObjectMeta{
			// SelfLink is required for an event to be created
			SelfLink:  "/apis/samplecontroller.k8s.io/v1alpha1/namespaces/test/bars/test-bar",
			Name:      "test-bar",
			Namespace: "test",
		},
	}
}

func newRunningTestController(t *testing.T) (
	kubeClient kubernetes.Interface,
	sampleClient clientset.Interface,
	controller *Controller,
	kubeInformer kubeinformers.SharedInformerFactory,
	sampleInformer informers.SharedInformerFactory,
	stopCh chan struct{}) {

	// Create fake clients
	kubeClient = fakekubeclientset.NewSimpleClientset()
	sampleClient = fakeclientset.NewSimpleClientset()

	// Create informer factories with fake clients. The second parameter sets the
	// resync period to zero, disabling it.
	kubeInformer = kubeinformers.NewSharedInformerFactory(kubeClient, 0)
	sampleInformer = informers.NewSharedInformerFactory(sampleClient, 0)

	// Create a controller and safe cast it to the proper type. This is necessary
	// because NewController returns controller.Interface.
	controller, ok := NewController(
		kubeClient,
		sampleClient,
		kubeInformer,
		sampleInformer,
	).(*Controller)
	if !ok {
		t.Fatal("cast to *Controller failed")
	}

	// Start the informers. This must happen after the call to NewController,
	// otherwise there are no informers to be started.
	stopCh = make(chan struct{})
	kubeInformer.Start(stopCh)
	sampleInformer.Start(stopCh)

	// Run the controller.
	go func() {
		if err := controller.Run(2, stopCh); err != nil {
			t.Fatalf("Error running controller: %v", err)
		}
	}()

	return
}

// Verify that an event is generated when a Bar is created.
func TestCreateGeneratesEvent(t *testing.T) {
	_, sampleClient, controller, _, _, stopCh := newRunningTestController(t)
	testBar := getTestBar()

	// Create an event watcher on the controller's broadcaster. If an event is
	// created, this will run and close stopCh, ending the test.
	controller.broadcaster.StartEventWatcher(func(e *corev1.Event) {
		assert.Equal(t, "Bar", e.InvolvedObject.Kind)
		assert.Equal(t, "samplecontroller.k8s.io", e.InvolvedObject.APIVersion)
		assert.Equal(t, testBar.Name, e.InvolvedObject.Name)
		assert.Equal(t, MessageResourceSynced, e.Message)
		assert.Equal(t, SuccessSynced, e.Reason)
		assert.Equal(t, corev1.EventTypeNormal, e.Type)
		close(stopCh)
	})

	// Create a testBar. This should cause an event to be created.
	sampleClient.SamplecontrollerV1alpha1().Bars("test").Create(testBar)

	// Wait up to 3 seconds for stopCh to close, otherwise fail the test.
	select {
	case <-stopCh:
		return
	case <-time.After(time.Second * 3):
		t.Fatal("timed out waiting for event")
	}
}

// Verify that an event is generated when a Bar is updated.
func TestUpdateGeneratesEvent(t *testing.T) {
	_, sampleClient, controller, _, _, stopCh := newRunningTestController(t)

	testBar := getTestBar()
	assert.Empty(t, testBar.Spec.DeploymentName)

	// Create an event to update.
	sampleClient.SamplecontrollerV1alpha1().Bars("test").Create(testBar)

	// Create an event watcher on the controller's broadcaster. If an event is
	// created, this will run and close stopCh, ending the test.
	controller.broadcaster.StartEventWatcher(func(e *corev1.Event) {
		assert.Equal(t, "Bar", e.InvolvedObject.Kind)
		assert.Equal(t, "samplecontroller.k8s.io", e.InvolvedObject.APIVersion)
		assert.Equal(t, testBar.Name, e.InvolvedObject.Name)
		assert.Equal(t, "test-deployment", testBar.Spec.DeploymentName)
		assert.Equal(t, MessageResourceSynced, e.Message)
		assert.Equal(t, SuccessSynced, e.Reason)
		assert.Equal(t, corev1.EventTypeNormal, e.Type)
		close(stopCh)
	})

	// Modify the deploymentName of the test Bar
	testBar.Spec.DeploymentName = "test-deployment"

	// Create a testBar. This should cause an event to be created.
	sampleClient.SamplecontrollerV1alpha1().Bars("test").Update(testBar)

	// Wait up to 3 seconds for stopCh to close, otherwise fail the test.
	select {
	case <-stopCh:
		return
	case <-time.After(time.Second * 3):
		t.Fatal("timed out waiting for event")
	}
}