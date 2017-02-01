package secrets

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func (s *secret) setDefaults() error {
	if s.Mode == "" {
		s.Mode = DefaultMode
	}

	if s.UID == "" {
		s.UID = DefaultUID
	}

	if s.GID == "" {
		s.GID = DefaultGID
	}

	return nil
}

func (s *secret) writeFile(basedir string, content []byte) error {
	fullPath := path.Join(basedir, s.Name)
	output := make([]byte, base64.StdEncoding.DecodedLen(len(content)))
	// Make sure defaults are set otherwise things could fail silently.

	l, err := base64.StdEncoding.Decode(output, content)
	if err != nil {
		return err
	}

	if err := s.setDefaults(); err != nil {
		return err
	}

	mode, err := strconv.ParseUint(s.Mode, 8, 32)
	if err != nil {
		return err
	}

	// Create the file and always truncate
	err = ioutil.WriteFile(fullPath, output[:l], os.FileMode(mode))
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(s.UID)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(s.GID)
	if err != nil {
		return err
	}

	if err = os.Chown(fullPath, uid, gid); err != nil {
		return err
	}

	if err = os.Chmod(fullPath, os.FileMode(mode)); err != nil {
		return err
	}

	return nil
}
