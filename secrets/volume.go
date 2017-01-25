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
	hostKeyPath = "/var/lib/rancher/etc/ssl/host.key"
)

// FlexVolume is a struct to implement the Rancher Volume interface
type FlexVolume struct {
	secretWriter SecretWriter
}

// Init implements the flex volume interface and is a no-op at this time
func (sv *FlexVolume) Init() error {
	return nil
}

// Create is implemented for Docker volume plugin API
func (sv *FlexVolume) Create(options map[string]interface{}) (map[string]interface{}, error) {
	resp := map[string]interface{}{}
	if name, ok := options["name"].(string); ok {
		volPath := path.Join(volRoot, "staging", name)

		if err := createTmpfs(volPath, options); err != nil {
			logrus.Error(err)
			return resp, err
		}

		resp["device"] = volPath
		resp["name"] = name

		return resp, nil
	}

	logrus.Error("Name not given as a string")
	return resp, errors.New("Name not given")
}

// Delete is implemented for Docker volume plugin API it detaches the
// volume and removes its content.
func (sv *FlexVolume) Delete(options map[string]interface{}) error {
	if device, ok := options["device"].(string); ok {
		return sv.Detach(device)
	}
	return nil
}

// Attach is implemeneted as a no-op for the flexvolume API
func (sv *FlexVolume) Attach(params map[string]interface{}) (string, error) {
	options, err := newOptions(params)

	// func (sv *FlexVolume) Mount(params map[string]interface{}) (string, error) {
	if options.Name == "" {
		return "", errors.New("Volume Name not given")
	}

	volumeDevice := path.Join(volRoot, "staging", options.Name)

	if err := createTmpfs(volumeDevice, params); err != nil {
		logrus.Error(err)
		return "", err
	}

	secretGetter, err := NewRancherSecretGetter(options)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	secrets, err := secretGetter.GetSecrets(options)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	decryptor, err := NewRSADecryptor(hostKeyPath)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	secretWriter, err := NewRSASecretFileWriter(decryptor)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return volumeDevice, secretWriter.Write(secrets, volumeDevice)
}

// Detach effectively erases the volume.
func (sv *FlexVolume) Detach(device string) error {
	if err := mount.Unmount(device); err != nil {
		return err
	}

	return os.RemoveAll(device)
}

// Mount implements does a bind mount of the volume to the target directory
func (sv *FlexVolume) Mount(dir, device string, params map[string]interface{}) error {
	//Default volume mode
	return mount.Mount(device, dir, "none", "bind,rw")

}

// Unmount undoes the bind mount, and removes the target directory
func (sv *FlexVolume) Unmount(dir string) error {
	// This will be a bind mount
	if err := mount.Unmount(dir); err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func createTmpfs(dir string, options map[string]interface{}) error {
	mounted, err := mount.Mounted(dir)
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

	return mount.Mount("tmpfs", dir, "tmpfs", mountOpts)
}
