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
	BackupDefaultTag = "v1.0.0-bkup"
)

type RegistryOptions struct {
	UserName     string
	Token        string
	Registry     string
	Organization string
}

func New(userName, token, registry, org string) *RegistryOptions {
	return &RegistryOptions{
		UserName:     userName,
		Token:        token,
		Registry:     registry,
		Organization: org,
	}
}

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

func BackupImage(image string, r *RegistryOptions) (string, error) {
	imageReference, err := name.ParseReference(image)
	if err != nil {
		return "", err
	}

	if imageReference.Context().Registry.RegistryStr() == r.Registry &&
		strings.HasPrefix(imageReference.Context().RepositoryStr(), fmt.Sprintf("%s/", r.Organization)) {
		return image, nil
	}

	// Some images have both sha and tag, for those images we can use default tag
	// Need to find a better solution
	tag := imageReference.Identifier()
	if strings.Contains(tag, "sha256") {
		tag = BackupDefaultTag
	}
	backupImageName := fmt.Sprintf("%s/%s:%s",
		r.Organization,
		strings.Replace(imageReference.Context().RepositoryStr(), "/", "-", -1),
		tag)

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
		Username: r.UserName,
		Password: r.Token,
	}), crane.WithTransport(tr))
	if err != nil {
		return "", err
	}

	return backupRef.String(), nil
}
