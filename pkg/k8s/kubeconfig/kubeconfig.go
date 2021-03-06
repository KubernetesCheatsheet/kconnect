/*
Copyright 2020 The kconnect Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubeconfig

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	logger = logrus.WithField("package", "kubeconfig")
)

// Write will write the kubeconfig to the specified file. If there
// is an existing kubeconfig it will be merged
func Write(path string, clusterConfig *api.Config, setCurrent bool) error {
	logger.Debugf("writing kubeconfig: %s", path)

	pathOptions := clientcmd.NewDefaultPathOptions()
	if path != "" {
		pathOptions.LoadingRules.ExplicitPath = path
	}

	existingConfig, err := pathOptions.GetStartingConfig()
	if err != nil {
		return fmt.Errorf("getting existing kubeconfig: %w", err)
	}

	logger.Debug("merging kubeconfig files")
	for k, v := range clusterConfig.Clusters {
		existingConfig.Clusters[k] = v
	}
	for k, v := range clusterConfig.AuthInfos {
		existingConfig.AuthInfos[k] = v
	}
	for k, v := range clusterConfig.Contexts {
		existingConfig.Contexts[k] = v
	}

	if setCurrent {
		logger.Infof("setting current context to: %s", clusterConfig.CurrentContext)
		existingConfig.CurrentContext = clusterConfig.CurrentContext
	}

	if err := clientcmd.ModifyConfig(pathOptions, *existingConfig, true); err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}
	logger.Infof("kubeconfig updated: %s", path)

	return nil
}
