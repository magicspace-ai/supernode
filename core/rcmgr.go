package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	network "github.com/libp2p/go-libp2p/core/network"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
)

func GetResourceMgr() (network.ResourceManager, error) {

	curDir, _ := os.Getwd()
	limitsFile := path.Join(curDir, "config", "limits.json")
	 
	var lconfig io.Reader
	
	if _, err := os.Stat(limitsFile); os.IsNotExist(err) {
		
		ilimit, err := rcmgr.InfiniteLimits.MarshalJSON()
		
		if err != nil {
			return nil, fmt.Errorf("failed to convert infinite limits to json, err=%w", err)
		}

		lconfig =  bytes.NewReader(ilimit)

	}  else {

		lconfig, err = os.Open(limitsFile)

		if err != nil {
			return nil, fmt.Errorf("failed to open ./config/limits.json, err=%w", err)
		}
	}
	
	limiter, err := rcmgr.NewDefaultLimiterFromJSON(lconfig)	

	if err != nil {
		return nil, fmt.Errorf("failed to create new limits from json, err=%w", err)
	}
	
	rcm, err := rcmgr.NewResourceManager(limiter)

	if err != nil {
		return nil, fmt.Errorf("resource manager failed to initialize, err=%w", err)
	}

	return rcm, nil 
}