package resingo

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

const daemonPort = 2375

//Compose is the configuration struct for a dockercompose based project.
type Compose struct {
	Projectname   string
	ComposeFiles  []string
	Host          string
	DeviceIP      string
	ClientOptions *client.Options
}

//Up does a docker-compose like up. It realies on the Compose struct to achieve
//this.
//
// This works in resin-os only.
func Up(c *Compose) error {
	co := client.Options{}
	if c.ClientOptions != nil {
		co = *c.ClientOptions
	} else {
		u, err := url.Parse(c.Host)
		if err != nil {
			if c.DeviceIP != "" {
				co.Host = fmt.Sprintf("http://%s:%d", c.DeviceIP, daemonPort)
			}
			return err
		}
		s := strings.Split(u.Host, ":")
		if len(s) == 1 {
			co.Host = fmt.Sprintf("%s:%d", c.Host, daemonPort)
		} else {
			co.Host = c.Host
		}
	}
	dc, err := client.NewDefaultFactory(co)
	if err != nil {
		return err
	}
	projectCtx := project.Context{
		ComposeFiles: c.ComposeFiles,
		ProjectName:  c.Projectname,
	}
	context := &ctx.Context{Context: projectCtx, ClientFactory: dc}
	return up(context)

}

func up(ctx *ctx.Context) error {
	project, err := docker.NewProject(ctx, nil)
	if err != nil {
		return err
	}
	o := options.Create{
		ForceRecreate: true,
		ForceBuild:    true,
	}
	return project.Up(context.Background(), options.Up{Create: o})
}
