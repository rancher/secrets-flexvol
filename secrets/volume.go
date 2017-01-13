package secrets

import (
	"errors"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/mount"
)

const (
	volRoot     = "/var/lib/rancher/volumes/rancher-secrets"
	hostKeyPath = "/var/lib/rancher/etc/ssl/server.key"
)

// FlexVolume is an empty struct to implement the interface
type FlexVolume struct {
	secretWriter SecretWriter
}

// Init implements the flex volume interface
func (sv *FlexVolume) Init() error {
	return nil
}

// Create is implemented for Docker volume plugin API but is a no-op
func (sv *FlexVolume) Create(options map[string]interface{}) (map[string]interface{}, error) {
	resp := map[string]interface{}{}
	if name, ok := options["name"].(string); ok {
		logrus.Infof("Create Called for volume: %s", name)
		volPath := path.Join(volRoot, "staging", name)

		if err := createTmpfs(volPath, options); err != nil {
			logrus.Error(err)
			return resp, err
		}

		resp["device"] = volPath
		resp["name"] = name
		logrus.Infof("Returning: %#v", resp)
		return resp, nil
	}

	logrus.Error("Name not given as a string")
	return resp, errors.New("Name not given")
}

// Delete is implemented for Docker volume plugin API but is a no-op
func (sv *FlexVolume) Delete(options map[string]interface{}) error {
	logrus.Infof("Delete called: %#v", options)
	if device, ok := options["device"].(string); ok {
		return sv.Detach(device)
	}
	return nil
}

// Attach is implemeneted as a no-op for the flexvolume API
func (sv *FlexVolume) Attach(params map[string]interface{}) (string, error) {
	// func (sv *FlexVolume) Mount(params map[string]interface{}) (string, error) {
	logrus.Infof("Attach Params: %#v", params)
	name, ok := params["name"].(string)
	if !ok {
		return "", errors.New("Volume Name not given")
	}

	volumeDevice := path.Join(volRoot, "staging", name)

	if err := createTmpfs(volumeDevice, params); err != nil {
		return "", err
	}

	//secretGetter, err := NewRancherSecretGetter(params)
	//if err != nil {
	//return "", err
	//}

	//secrets, err := secretGetter.GetSecrets(params)
	//if err != nil {
	//return "", err
	//}

	//logrus.Debugf("Secrets: %#v", secrets)

	//decryptor, err := NewRSADecryptor(hostKeyPath)
	// if err != nil {
	//return "", err
	//}

	// secretWriter, err := NewRSASecretFileWriter(decryptor, params)
	// if err != nil {
	// return "", err
	// }

	// err = secretWriter.Write(secrets, dir)

	return volumeDevice, nil
}

// Detach effectively erases the volume.
func (sv *FlexVolume) Detach(device string) error {
	logrus.Infof("Detach: Device %s", device)
	if err := mount.Unmount(device); err != nil {
		return err
	}

	return os.RemoveAll(device)
}

// Mount implements does a bind mount of the volume
func (sv *FlexVolume) Mount(dir, device string, params map[string]interface{}) error {
	//Default volume mode
	logrus.Infof("Mounting: %s dev %s with params %#v", dir, device, params)
	return mount.Mount(device, dir, "none", "bind,rw")

}

// Unmount is a no-op
func (sv *FlexVolume) Unmount(dir string) error {
	logrus.Infof("Unmounting: %s", dir)
	// This will be a bind mount
	if err := mount.Unmount(dir); err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func createTmpfs(dir string, options map[string]interface{}) error {
	mounted, err := mount.Mounted(dir)
	logrus.Infof("Mounted: %v Err: %s", mounted, err)
	if mounted || err != nil {
		return err
	}

	mode := int(0755)
	mountOpts := "size=10m"

	if uMode, ok := options["mode"].(int); ok {
		mode = int(uMode)
	}

	if mOpts, ok := options["mountOpts"].(string); ok {
		mountOpts = mOpts
	}

	if err := os.MkdirAll(dir, os.FileMode(mode)); err != nil {
		return err
	}

	logrus.Infof("mounting: %s with opts: %s", dir, mountOpts)

	return mount.Mount("tmpfs", dir, "tmpfs", mountOpts)
}
