package registry

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultDockerOrg = "backuphub"

	DefaultRegistry = "index.docker.io"

	UserName = "imgbackupcontroller"
	PassKey  = "6bb4afd8-f760-462d-a4ec-34a568e29ab3"
)

var (
	tr = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

func PullAndUploadImage(image string) (string, error) {
	// get the original image's reference
	originalRef, err := name.ParseReference(image)
	if err != nil {
		return "", err
	}

	// check if the image already belongs to the backup repository
	if originalRef.Context().Registry.RegistryStr() == DefaultRegistry &&
		strings.HasPrefix(originalRef.Context().RepositoryStr(), fmt.Sprintf("%s/", DefaultDockerOrg)){
		return image, nil
	}

	// build the backup image reference
	backupImageName := fmt.Sprintf("%s/%s:%s",
		DefaultDockerOrg,
		strings.Replace(originalRef.Context().RepositoryStr(), "/", "-", -1),
		originalRef.Identifier())

	backupRef, err := name.ParseReference(backupImageName)
	if err != nil {
		return "", err
	}

	// Supports only schema 2 of OCI specs
	// https://github.com/google/go-containerregistry/issues/377
	originalImage, err := crane.Pull(image)
	if err != nil {
		return "", err
	}

	err = crane.Push(originalImage, backupRef.String(), crane.WithAuth(&authn.Basic{
		Username: UserName,
		Password: PassKey,
	}), crane.WithTransport(tr))
	if err != nil {
		return "", err
	}

	return backupRef.String(), nil
}
