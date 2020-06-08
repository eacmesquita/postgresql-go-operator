package postgresql

import (
	"context"

	postgresqlv1alpha1 "github.com/persistentsys/postgresql-go-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/persistentsys/postgresql-go-operator/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const postgresqlPort = 80
const pvStorageName = "postgresql-pv-storage"
const pvClaimName = "postgresql-pv-claim"

func postgresqlDeploymentName(v *postgresqlv1alpha1.PostgreSQL) string {
	return v.Name + "-server"
}

func postgresqlServiceName(v *postgresqlv1alpha1.PostgreSQL) string {
	return v.Name + "-service"
}

func postgresqlAuthName() string {
	return "postgresql-auth"
}

func (r *ReconcilePostgreSQL) postgresqlDeployment(v *postgresqlv1alpha1.PostgreSQL) *appsv1.Deployment {
	labels := utils.Labels(v, "postgresql")
	size := v.Spec.Size
	image := v.Spec.Image

	dbname := v.Spec.Database
	//rootpwd := v.Spec.Rootpwd

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: postgresqlAuthName()},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: postgresqlAuthName()},
			Key:                  "password",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgresqlDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: pvStorageName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvClaimName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "postgresql-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: postgresqlPort,
							Name:          "postgresql",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      pvStorageName,
								MountPath: "/var/lib/postgresql/data",
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "POSTGRES_DB",
								Value: dbname,
							},
							{
								Name:      "POSTGRES_USER",
								ValueFrom: userSecret,
							},
							{
								Name:      "POSTGRES_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcilePostgreSQL) postgresqlService(v *postgresqlv1alpha1.PostgreSQL) *corev1.Service {
	labels := utils.Labels(v, "postgresql")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgresqlServiceName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       postgresqlPort,
				TargetPort: intstr.FromInt(3306),
				NodePort:   v.Spec.Port,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcilePostgreSQL) updatePostgresqlStatus(v *postgresqlv1alpha1.PostgreSQL) error {
	//v.Status.BackendImage = postgresqlImage
	err := r.client.Status().Update(context.TODO(), v)
	return err
}

func (r *ReconcilePostgreSQL) postgresqlAuthSecret(v *postgresqlv1alpha1.PostgreSQL) *corev1.Secret {

	username := v.Spec.Username
	password := v.Spec.Password

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgresqlAuthName(),
			Namespace: v.Namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"username": []byte(username),
			"password": []byte(password),
		},
	}
	controllerutil.SetControllerReference(v, secret, r.scheme)
	return secret
}
