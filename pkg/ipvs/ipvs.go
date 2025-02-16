// Copyright © 2021 Alibaba Group Holding Ltd.
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

package ipvs

import (
	"fmt"
	"net"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	LvsCareStaticPodName = "kube-lvscare"
	LvsCareCommand       = "/usr/bin/lvscare"
)

// LvsStaticPodYaml return lvs care static pod yaml
func LvsStaticPodYaml(vip net.IP, masters []net.IP, image string) (string, error) {
	if vip == nil || len(masters) == 0 || image == "" {
		return "", fmt.Errorf("invalid args to create Lvs static pod")
	}

	args := []string{"care", "--vs", net.JoinHostPort(vip.String(), "6443"), "--health-path", "/healthz", "--health-schem", "https"}
	for _, m := range masters {
		args = append(args, "--rs")
		args = append(args, net.JoinHostPort(m.String(), "6443"))
	}
	flag := true
	pod := componentPod(v1.Container{
		Name:            LvsCareStaticPodName,
		Image:           image,
		Command:         []string{LvsCareCommand},
		Args:            args,
		ImagePullPolicy: v1.PullIfNotPresent,
		SecurityContext: &v1.SecurityContext{Privileged: &flag},
	})

	yml, err := yaml.Marshal(pod)
	if err != nil {
		return "", fmt.Errorf("failed to decode lvs care static pod yaml: %s", err)
	}

	return string(yml), nil
}

// componentPod returns a Pod object from the container and volume specifications
func componentPod(container v1.Container) v1.Pod {
	hostPathType := v1.HostPathUnset
	mountName := "lib-modules"
	volumes := []v1.Volume{
		{Name: mountName, VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: "/lib/modules",
				Type: &hostPathType,
			},
		}},
	}
	container.VolumeMounts = []v1.VolumeMount{
		{Name: mountName, ReadOnly: true, MountPath: "/lib/modules"},
	}

	return v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      container.Name,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: v1.PodSpec{
			Containers:  []v1.Container{container},
			HostNetwork: true,
			Volumes:     volumes,
		},
	}
}
