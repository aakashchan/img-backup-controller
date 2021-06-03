package registry

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DockerOrg = "backuphub"
	Tag = "bkup"

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

func PullAndUploadImage(image string) (string,error) {
	img, err := crane.Pull(image)

	imageDetails := strings.Split(image, ":")


	imageRepo := strings.Split(imageDetails[0], "/")[1]
	imageTag := imageDetails[1]

	if err != nil {
		return "",err
	}

	destination := fmt.Sprintf("%s/%s:%s", DockerOrg, imageRepo, imageTag)

	err = crane.Push(img,destination, crane.WithAuth(&authn.Basic{
		Username:      UserName,
		Password:      PassKey,
	}), crane.WithTransport(tr))
	if err != nil {
		return "",err
	}

	return destination,nil
}