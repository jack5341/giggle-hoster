package handler

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jack5341/giggle-hoster/internal/cf"
	hcloudi "github.com/jack5341/giggle-hoster/internal/hcloud"
	"github.com/pocketbase/pocketbase/core"
)

func CreateNode(e *core.RecordCreateEvent) error {
	game := strings.ToLower(e.Record.GetString("game"))
	size := strings.ToLower(e.Record.GetString("size"))

	h, err := hcloudi.NewClient()
	if err != nil {
		return err
	}

	serverName := fmt.Sprintf("%s-%s-%s", game, size, uuid.New())
	serverOpts := hcloud.ServerCreateOpts{
		Name: strings.ToLower(serverName),
		ServerType: &hcloud.ServerType{
			Name: "cax11",
		},
		Image: &hcloud.Image{
			Name: "ubuntu-20.04",
		},
		Location: &hcloud.Location{Name: "fsn1"},
		Labels: map[string]string{
			"game": game,
			"size": size,
		},
	}

	server, err := h.CreateNewInstance(e.HttpContext.Request().Context(), serverOpts)
	if err != nil {
		return err
	}

	c, err := cf.NewClient()
	if err != nil {
		return err
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s.%s.nedim-akar.cloud", uid.String(), game)
	dnsRecordParam := map[string]interface{}{
		"Type":    "A",
		"Name":    name,
		"Content": server.Server.PublicNet.IPv4.IP.String(),
		"TTL":     120,
	}

	_, err = c.CreateRecord(e.HttpContext.Request().Context(), "ae60ae0f1bed3c9369bd296e50b501aa", dnsRecordParam)
	if err != nil {
		return err
	}

	e.Record.Set("url", name)
	return nil
}
