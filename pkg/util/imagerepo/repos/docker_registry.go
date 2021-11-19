package repos

import (
	"context"
	"crypto/tls"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"io"
	"net/http"
)

func NewDockerRegistryClient(endpoint, username, password string) *dockerRegistryClient {
	return &dockerRegistryClient{
		Username: username,
		Password: password,
		EndPoint: endpoint,
	}
}

type dockerRegistryClient struct {
	Username string
	Password string
	EndPoint string
}

func (c *dockerRegistryClient) ListRepos() (names []string, err error) {
	return
}

func (c *dockerRegistryClient) ListImages(project string) (images []string, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		},
	}

	basicTransport := &BasicTransport{
		Transport: transport,
		Username:  c.Username,
		Password:  c.Password,
		URL:       c.EndPoint,
	}
	registry, err1 := client.NewRegistry(c.EndPoint, basicTransport)
	if err1 != nil {
		err = err1
		return
	}
	allImages := make([]string, 0)
	var last string
	for {
		imageNoTags := make([]string, 10)
		count, err2 := registry.Repositories(context.Background(), imageNoTags, last)
		if err2 == io.EOF {
			allImages = append(allImages, imageNoTags[:count]...)
			break
		} else if err2 != nil {
			err = err2
			return
		}
		last = imageNoTags[count-1]
		allImages = append(allImages, imageNoTags...)
	}
	for _, image := range allImages {
		tags, err3 := c.listImageTags(image, basicTransport)
		if err3 != nil {
			err = err3
			return
		}

		if len(tags) != 0 {
			for _, tag := range tags {
				images = append(images, image+":"+tag)
			}
		}
	}
	return
}

func (c *dockerRegistryClient) listImageTags(imageName string, transport http.RoundTripper) (tags []string, err error) {
	ref, err1 := reference.Parse(imageName)
	if err1 != nil {
		err = err1
		return
	}
	repo, err2 := client.NewRepository(ref.(reference.Named), c.EndPoint, transport)
	if err2 != nil {
		err = err2
		return
	}
	ctx := context.Background()
	tagService := repo.Tags(ctx)
	return tagService.All(ctx)
}
