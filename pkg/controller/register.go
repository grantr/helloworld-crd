package controller

import (
	"sync"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	clientset "k8s.io/sample-controller/pkg/client/clientset/versioned"
	informers "k8s.io/sample-controller/pkg/client/informers/externalversions"
)

type Interface interface {
	Run(threadiness int, stopCh <-chan struct{}) error
}

type Constructor func(kubernetes.Interface, clientset.Interface, kubeinformers.SharedInformerFactory, informers.SharedInformerFactory) Interface

var constructorsMutex sync.Mutex
var constructors = make(map[string]Constructor)

func Register(name string, ctor Constructor) {
	constructorsMutex.Lock()
	defer constructorsMutex.Unlock()
	_, found := constructors[name]
	if found {
		glog.Fatalf("Controller %q was registered twice", name)
	}
	glog.V(4).Infof("Registered controller %q", name)
	constructors[name] = ctor
}

type multiController struct {
	controllers []Interface
}

func New(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	sampleInformerFactory informers.SharedInformerFactory) Interface {
	var mc multiController
	constructorsMutex.Lock()
	defer constructorsMutex.Unlock()
	for _, ctor := range constructors {
		mc.controllers = append(mc.controllers,
			ctor(kubeclientset, sampleclientset, kubeInformerFactory,
				sampleInformerFactory))
	}
	return &mc
}

func (mc *multiController) Run(threadiness int, stopCh <-chan struct{}) error {
	// Spin up a go routine for each controller.
	errors := make(chan error)
	for _, c := range mc.controllers {
		// We need a copy for the go routine.
		ctrlr := c
		go func() {
			// We don't expect this to return until stop is called,
			// but if it does, propagate it back.
			errors <- ctrlr.Run(threadiness, stopCh)
		}()
	}

	// Wait for a response from each go routine.
	for _ = range mc.controllers {
		if err := <-errors; err != nil {
			return err
		}
	}

	// TODO(mattmoor): Should we listen for stopCh?
	return nil
}
