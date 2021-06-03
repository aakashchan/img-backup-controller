module github.com/lostbrain101/img-backup-controller

go 1.13

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/go-containerregistry v0.4.1
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/klog/v2 v2.8.0 // indirect
	sigs.k8s.io/controller-runtime v0.8.3
)
