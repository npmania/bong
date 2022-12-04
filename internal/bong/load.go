package bong

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func LoadBongs(filename string) (BongMap, error) {
	logrus.WithFields(logrus.Fields{
		"filename": filename,
	}).Info("loading bong from file")

	if _, err := os.Stat(filename); err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": filename,
			"error":    err,
		}).Error("failed file lookup")
		return nil, err
	}

	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	bongMap := make(BongMap)
	if err = yaml.Unmarshal(raw, &bongMap); err != nil {
		return nil, err
	}

	for bg := range bongMap {
		b := bongMap[bg]
		b.MainUrl = strings.ReplaceAll(bongMap[bg].BongUrl, "%s", "%[1]s")
		b.BongUrl = strings.ReplaceAll(bongMap[bg].BongUrl, "%s", "%[1]s")
		bongMap[bg] = b
	}

	if err = bongMap.validate(); err != nil {
		return nil, err
	}

	return bongMap, nil
}
