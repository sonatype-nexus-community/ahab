package packages

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/go-sona-types/cyclonedx"
)

type Hasher struct {
	logLady *logrus.Logger
	Files   []cyclonedx.File
}

func New(logLady *logrus.Logger) *Hasher {
	return &Hasher{logLady: logLady}
}

func (h Hasher) PopulateListOfHashes(path string) {
	for _, v := range filepath.SplitList(path) {
		err := filepath.Walk(v, func(path string, f os.FileInfo, err error) (hashErr error) {
			hashErr = h.getHashAndAppend(path)
			if hashErr != nil {
				return hashErr
			}
			return nil
		})
		if err != nil {
			h.logLady.Error(err)
		}
	}
}

func (h *Hasher) getHashAndAppend(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	hashString, err := h.getSha1(f)
	if err != nil {
		return err
	}
	h.Files = append(h.Files, cyclonedx.File{Path: path, Extension: filepath.Ext(path), Hash: hashString})
	return
}

func (h Hasher) getSha1(f *os.File) (hashString string, err error) {
	hash := sha1.New()
	if _, err = io.Copy(hash, f); err != nil {
		h.logLady.Error(err)
		return
	}

	hashString = hex.EncodeToString(hash.Sum(nil))

	return
}
