package cmd

import "github.com/sirupsen/logrus"

var (
	log           *logrus.Logger
	cfgFile       string
	homeConfigDir string
	Version       string
)
