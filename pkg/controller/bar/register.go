package bar

import "k8s.io/sample-controller/pkg/controller"

func init() {
	controller.Register(controllerAgentName, NewController)
}
