package hcloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Hcloud struct {
	Client *hcloud.Client
}

var (
	ErrHetznerTokenIsNotExist                     = errors.New("hetzner token has to be setted")
	ErrServerCouldNotBeCreated                    = errors.New("an error occurred during server creation process")
	ErrServerCouldNotBeDeleted                    = errors.New("an error occurred during server deletion process")
	ErrServerTypesCouldNotBeRetrieved             = errors.New("server type could not be retrieved")
	ErrUnexpectedErrorDuringServerCreationProcess = errors.New("unexpected error during server creation process")
	ErrUnexpectedErrorDuringServerDeletionProcess = errors.New("unexpected error during server deletion process")
)

func (h *Hcloud) NewClient() (*Hcloud, error) {
	hcloudToken := os.Getenv("HCLOUD_TOKEN")
	if hcloudToken == "" {
		return nil, ErrHetznerTokenIsNotExist
	}

	client := hcloud.NewClient(hcloud.WithToken(hcloudToken))

	h.Client = client
	return h, nil
}

func (h *Hcloud) CreateNewInstance(ctx context.Context, opts hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, error) {
	client := h.Client
	createdServer, response, err := client.Server.Create(ctx, opts)
	if err != nil {
		return hcloud.ServerCreateResult{}, errors.Join(ErrServerCouldNotBeCreated, err)
	}

	if response.StatusCode != http.StatusCreated {
		defer func() {
			if err := response.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()
		rawBody, err := io.ReadAll(response.Body)

		if err != nil {
			return hcloud.ServerCreateResult{}, err
		}

		body := string(rawBody)
		return hcloud.ServerCreateResult{}, fmt.Errorf("%w: %s", ErrUnexpectedErrorDuringServerCreationProcess, body)
	}

	return createdServer, nil
}

func (h *Hcloud) DeleteInstance(ctx context.Context, server *hcloud.Server) (hcloud.ServerDeleteResult, error) {
	client := h.Client
	deletedServer, response, err := client.Server.DeleteWithResult(ctx, server)
	if err != nil {
		return hcloud.ServerDeleteResult{}, errors.Join(ErrServerCouldNotBeDeleted, err)
	}

	if response.StatusCode != http.StatusOK {
		defer func() {
			if err := response.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()

		rawBody, err := io.ReadAll(response.Body)

		if err != nil {
			return hcloud.ServerDeleteResult{}, err
		}

		body := string(rawBody)
		return hcloud.ServerDeleteResult{}, fmt.Errorf("%w: %s", ErrUnexpectedErrorDuringServerDeletionProcess, body)
	}

	return *deletedServer, nil
}

func (h *Hcloud) GetAllServerTypes(ctx context.Context) ([]*hcloud.ServerType, error) {
	client := h.Client
	allServerTypes, err := client.ServerType.All(ctx)
	if err != nil {
		return nil, errors.Join(ErrServerTypesCouldNotBeRetrieved, err)
	}

	return allServerTypes, nil
}
