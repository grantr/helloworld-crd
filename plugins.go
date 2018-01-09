package main

// Add a line here for each controller plugin we want to register.
import (
	_ "k8s.io/sample-controller/pkg/controller/foo"
)
