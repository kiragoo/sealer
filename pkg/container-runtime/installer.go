// Copyright © 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package containerruntime

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/sealerio/sealer/common"

	"github.com/sealerio/sealer/pkg/infradriver"
)

const (
	DefaultDockerCRISocket     = "/var/run/dockershim.sock"
	DefaultContainerdCRISocket = "/run/containerd/containerd.sock"
	DefaultSystemdCgroupDriver = "systemd"
	DefaultCgroupDriver        = "cgroupfs"
	DockerDockerCertsDir       = "/etc/docker/certs.d"
	DockerConfigFileName       = "config.json"
)

type Installer interface {
	InstallOn(hosts []net.IP) error

	GetInfo() (Info, error)

	UnInstallFrom(hosts []net.IP) error

	//Upgrade() (ContainerRuntimeInfo, error)
	//Rollback() (ContainerRuntimeInfo, error)
}

type Config struct {
	Type         string
	LimitNofile  string `json:"limitNofile,omitempty" yaml:"limitNofile,omitempty"`
	CgroupDriver string `json:"cgroupDriver,omitempty" yaml:"cgroupDriver,omitempty"`
}

type Info struct {
	Config
	CRISocket      string
	CertsDir       string
	ConfigFilePath string
}

func NewInstaller(conf Config, driver infradriver.InfraDriver) (Installer, error) {
	if conf.Type == "docker" {
		return &DockerInstaller{
			rootfs: driver.GetClusterRootfsPath(),
			driver: driver,
			Info: Info{
				CertsDir:       DockerDockerCertsDir,
				CRISocket:      DefaultDockerCRISocket,
				Config:         conf,
				ConfigFilePath: filepath.Join(common.GetHomeDir(), ".docker", DockerConfigFileName),
			},
		}, nil
	}

	if conf.Type == "containerd" {
		return &ContainerdInstaller{
			rootfs: driver.GetClusterRootfsPath(),
			driver: driver,
		}, nil
	}

	return nil, fmt.Errorf("please enter the correct container type")
}
