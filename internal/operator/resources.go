package operator

import (
	examplev1beta1 "github.com/Youngpig1998/hotelreservation-operator/api/v1beta1"
	"github.com/Youngpig1998/hotelreservation-operator/iaw-shared-helpers/pkg/resources"
	"github.com/Youngpig1998/hotelreservation-operator/iaw-shared-helpers/pkg/resources/deployments"
	"github.com/Youngpig1998/hotelreservation-operator/iaw-shared-helpers/pkg/resources/services"
	"github.com/Youngpig1998/hotelreservation-operator/iaw-shared-helpers/pkg/resources/statefulsets"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"strconv"
)

const (
	StorageRequest = "1Gi"
)

func Service(serviceName string, port int32, targetPort int32, nodePort int32) resources.Reconcileable {

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
			Labels: map[string]string{
				"io.kompose.service": serviceName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Protocol: corev1.ProtocolTCP,
				Port:     port,
				TargetPort: intstr.IntOrString{
					IntVal: targetPort,
					StrVal: strconv.Itoa(int(targetPort)),
				},
				NodePort: nodePort,
			},
			},
			Selector: map[string]string{
				"io.kompose.service": serviceName,
			},
			Type: "NodePort",
		},
	}

	return services.From(service)
}

func StatefulSet(servicesName string, app *examplev1beta1.HotelReservationApp) resources.Reconcileable {

	statefulSetName := "mongodb-" + servicesName

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: statefulSetName,
			Labels: map[string]string{
				"io.kompose.service": statefulSetName,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"io.kompose.service": statefulSetName,
				},
			},
			ServiceName: statefulSetName,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.kompose.service": statefulSetName,
					},
				},
				Spec: corev1.PodSpec{
					Volumes: nil,

					Containers: []corev1.Container{{
						Name:            "hotelreservation-" + statefulSetName,
						Image:           "mongo",
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 27017,
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      servicesName,
							MountPath: "/data/db",
						},
						},
					},
					},

					NodeSelector: map[string]string{
						"kubernetes.io/hostname": app.Spec.DataNodeName,
					},
					ServiceAccountName: "nfs-provisioner",
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: servicesName,
					Annotations: map[string]string{
						"volume.beta.kubernetes.io/storage-class": "managed-nfs-storage",
					},
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": resource.MustParse(StorageRequest),
						},
					},
				},
			},
			},
		},
	}

	return statefulsets.From(statefulSet)
}

func DeploymentForMem(servicesName string, app *examplev1beta1.HotelReservationApp) resources.Reconcileable {

	//imageName := "cp.icr.io/cp/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//if len(strings.TrimSpace(webHook.Spec.DockerRegistryPrefix)) > 0 {
	//	imageName = webHook.Spec.DockerRegistryPrefix + "/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//}

	deployName := "memcached-" + servicesName

	// Instantialize the data structure
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			//Namespace: webHook.Namespace,
			Name: deployName,
			Labels: map[string]string{
				"io.kompose.service": deployName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			// The replica is computed
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"io.kompose.service": deployName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.kompose.service": deployName,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": app.Spec.DataNodeName,
					},
					Containers: []corev1.Container{{
						Image:           "memcached",
						ImagePullPolicy: "IfNotPresent",
						Name:            "hotelreservation-" + deployName,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MEMCACHED_CACHE_SIZE",
								Value: "2048",
							},
							{
								Name:  "MEMCACHED_THREADS",
								Value: "8",
							},
						},
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	return deployments.From(deployment)
}

func DeploymentForLogic(deployName string, port int32, app *examplev1beta1.HotelReservationApp) resources.Reconcileable {

	isRunAsRoot := true
	pIsRunAsRoot := &isRunAsRoot //bool pointer

	var runAsUser int64 = 1000321000

	//imageName := "cp.icr.io/cp/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//if len(strings.TrimSpace(webHook.Spec.DockerRegistryPrefix)) > 0 {
	//	imageName = webHook.Spec.DockerRegistryPrefix + "/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//}

	// Instantialize the data structure
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deployName,
			Labels: map[string]string{
				"io.kompose.service": deployName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			// The replica should be computed
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"io.kompose.service": deployName,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.kompose.service": deployName,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": app.Spec.LogicNodeName,
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser:    &runAsUser,
						RunAsNonRoot: pIsRunAsRoot,
					},
					InitContainers: []corev1.Container{{
						Image:           "youngpig/configwriter:latest",
						ImagePullPolicy: "IfNotPresent",
						Name:            "configwriter",

						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot: pIsRunAsRoot,
						},
						Env: []corev1.EnvVar{
							{
								Name:  "LOGICNODEIP",
								Value: app.Spec.LogicNodeIp,
							},
							{
								Name:  "DATANODEIP",
								Value: app.Spec.DataNodeIp,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/var/configFiles",
								Name:      "varconfig",
							},
						},
					}},
					Containers: []corev1.Container{{
						Image:           "youngpig/hotel_reservation",
						ImagePullPolicy: "IfNotPresent",
						Name:            "hotelreservation-" + deployName,
						Command:         []string{deployName},
						Ports: []corev1.ContainerPort{{
							HostPort:      port,
							ContainerPort: port,
						}},
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.LifecycleHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"/bin/sh", "-c", "sleep 5"},
								},
							},
						},
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot: pIsRunAsRoot,
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/go/src/github.com/harlow/go-micro-services/config",
								Name:      "varconfig",
							},
						},
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						{
							Name: "varconfig",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	return deployments.From(deployment)
}

func DeploymentForConsul(app *examplev1beta1.HotelReservationApp) resources.Reconcileable {

	//imageName := "cp.icr.io/cp/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//if len(strings.TrimSpace(webHook.Spec.DockerRegistryPrefix)) > 0 {
	//	imageName = webHook.Spec.DockerRegistryPrefix + "/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//}

	// Instantialize the data structure
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			//Namespace: webHook.Namespace,
			Name: "consul",
			Labels: map[string]string{
				"io.kompose.service": "consul",
			},
		},
		Spec: appsv1.DeploymentSpec{
			// The replica is computed
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"io.kompose.service": "consul",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.kompose.service": "consul",
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": app.Spec.LogicNodeName,
					},
					Containers: []corev1.Container{{
						Image:           "consul",
						ImagePullPolicy: "IfNotPresent",
						Name:            "consul",
						Ports: []corev1.ContainerPort{{
							HostPort:      8300,
							ContainerPort: 8300,
						}, {
							HostPort:      8400,
							ContainerPort: 8400,
						}, {
							HostPort:      8500,
							ContainerPort: 8500,
						}, {
							HostPort:      8600,
							ContainerPort: 53,
							Protocol:      "UDP",
						}},
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	return deployments.From(deployment)
}

func DeploymentForJaeger(app *examplev1beta1.HotelReservationApp) resources.Reconcileable {

	//imageName := "cp.icr.io/cp/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//if len(strings.TrimSpace(webHook.Spec.DockerRegistryPrefix)) > 0 {
	//	imageName = webHook.Spec.DockerRegistryPrefix + "/opencontent-audit-webhook@sha256:f4935b3a1687aeb23922fd144f880cc5a4f00404e794a4e30cccd6392cbe29f5"
	//}

	// Instantialize the data structure
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			//Namespace: webHook.Namespace,
			Name: "jaeger",
			Labels: map[string]string{
				"io.kompose.service": "jaeger",
			},
		},
		Spec: appsv1.DeploymentSpec{
			// The replica is computed
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"io.kompose.service": "jaeger",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.kompose.service": "jaeger",
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": app.Spec.LogicNodeName,
					},
					Containers: []corev1.Container{{
						Image:           "jaegertracing/all-in-one:latest",
						ImagePullPolicy: "IfNotPresent",
						Name:            "jaeger",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 14269,
						}, {
							HostPort:      5778,
							ContainerPort: 5778,
						}, {
							HostPort:      14268,
							ContainerPort: 14268,
						}, {
							HostPort:      5775,
							ContainerPort: 5775,
							Protocol:      "UDP",
						}, {
							ContainerPort: 14267,
						}, {
							HostPort:      16686,
							ContainerPort: 16686,
						}, {
							HostPort:      6831,
							ContainerPort: 6831,
							Protocol:      "UDP",
						}, {
							HostPort:      6832,
							ContainerPort: 6832,
							Protocol:      "UDP",
						}},
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	return deployments.From(deployment)
}
