package secrets

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	flexvol "github.com/cloudnautique/rancher-flexvol"
	"github.com/docker/docker/pkg/mount"
)

const (
	volRoot = "/var/lib/rancher/volumes"
)

type FlexVolume struct{}

func (sv *FlexVolume) Name() string {
	return "rancher-secrets"
}

func (sv *FlexVolume) Init() (*flexvol.DriverOutput, error) {
	return &flexvol.DriverOutput{Status: "Success"}, nil
}

func (sv *FlexVolume) Attach(params map[string]interface{}) (*flexvol.DriverOutput, error) {
	var volName string
	var ok bool

	//Default volume mode
	mode := int(0700)
	status := "Failure"
	mountOpts := "size=10m"
	output := &flexvol.DriverOutput{}

	if volName, ok = params["volumeName"].(string); !ok {
		output.Message = "No volumeName set"
	}

	if uMode, ok := params["mode"].(int); ok {
		mode = int(uMode)
	}

	if mOpts, ok := params["mountOpts"].(string); ok {
		mountOpts = mOpts
	}

	path := fmt.Sprintf("%s/%s", volRoot, volName)

	os.Mkdir(path, os.FileMode(mode))

	if err := mount.Mount("tmpfs", path, "tmpfs", mountOpts); err == nil {
		status = "Success"
	}

	output.Status = status
	output.Device = path

	return output, nil
}

func (sv *FlexVolume) Detach(device string) (*flexvol.DriverOutput, error) {
	status := "Failure"
	output := &flexvol.DriverOutput{}

	logrus.Errorf("DEVICE: %s", device)

	output.Status = status

	return output, nil
}

func (sv *FlexVolume) Mount(dir, device, params string) (*flexvol.DriverOutput, error) {
	status := "Failure"
	output := &flexvol.DriverOutput{}

	output.Status = status
	logrus.Errorf("Dir: %s, DEV: %s, PARAMS: %s", dir, device, params)

	return output, nil
}

func (sv *FlexVolume) Unmount(dir string) (*flexvol.DriverOutput, error) {
	status := "Failure"
	output := &flexvol.DriverOutput{}

	logrus.Errorf("Dir: %s", dir)
	output.Status = status

	return output, nil
}
