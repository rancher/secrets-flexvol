package secrets

import (
	"os"

	flexvol "github.com/cloudnautique/rancher-flexvol"
	"github.com/docker/docker/pkg/mount"
)

const (
	volRoot = "/var/lib/rancher/volumes"
)

type FlexVolume struct{}

func (sv *FlexVolume) Init() error {
	return nil
}

func (sv *FlexVolume) Attach(params map[string]interface{}) (string, error) {
	return "", flexvol.ErrNotSupported
}

func (sv *FlexVolume) Detach(device string) error {
	return flexvol.ErrNotSupported
}

func (sv *FlexVolume) Mount(dir, device string, params map[string]interface{}) error {
	//Default volume mode
	mode := int(0700)
	mountOpts := "size=10m"

	if uMode, ok := params["mode"].(int); ok {
		mode = int(uMode)
	}

	if mOpts, ok := params["mountOpts"].(string); ok {
		mountOpts = mOpts
	}

	if err := os.MkdirAll(dir, os.FileMode(mode)); err != nil {
		return err
	}

	if err := mount.Mount("tmpfs", dir, "tmpfs", mountOpts); err != nil {
		return err
	}

	return nil
}

func (sv *FlexVolume) Unmount(dir string) error {
	return mount.Unmount(dir)
}
