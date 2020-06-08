package resource

import (
	"github.com/persistentsys/postgresql-go-operator/pkg/apis/postgresql/v1alpha1"
	"github.com/persistentsys/postgresql-go-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var volLog = logf.Log.WithName("resource_volumes")

// GetPostgresqlVolumeName - return name of PV used in PostgreSQL
func GetPostgresqlVolumeName(v *v1alpha1.PostgreSQL) string {
	return v.Name + "-" + v.Namespace + "-pv"
}

// GetPostgresqlVolumeClaimName - return name of PVC used in PostgreSQL
func GetPostgresqlVolumeClaimName(v *v1alpha1.PostgreSQL) string {
	return v.Name + "-pv-claim"
}

// GetPostgresqlBkpVolumeName - return name of PV used in DB Backup
// func GetPostgresqlBkpVolumeName(bkp *v1alpha1.Backup) string {
// 	return bkp.Name + "-" + bkp.Namespace + "-pv"
// }

// GetPostgresqlBkpVolumeClaimName - return name of PVC used in DB Backup
// func GetPostgresqlBkpVolumeClaimName(bkp *v1alpha1.Backup) string {
// 	return bkp.Name + "-pv-claim"
// }

// NewDbBackupPV Create a new PV object for Database Backup
// func NewDbBackupPV(bkp *v1alpha1.Backup, v *v1alpha1.PostgreSQL, scheme *runtime.Scheme) *corev1.PersistentVolume {
// 	volLog.Info("Creating new PV for Database Backup")
// 	labels := utils.PostgreSQLBkpLabels(bkp, "postgresql-backup")
// 	pv := &corev1.PersistentVolume{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: GetPostgresqlBkpVolumeName(bkp),
// 			// Namespace: v.Namespace,
// 			Labels: labels,
// 		},
// 		Spec: corev1.PersistentVolumeSpec{
// 			StorageClassName: "manual",
// 			Capacity: corev1.ResourceList{
// 				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(bkp.Spec.BackupSize),
// 			},
// 			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
// 			PersistentVolumeSource: corev1.PersistentVolumeSource{
// 				HostPath: &corev1.HostPathVolumeSource{
// 					Path: bkp.Spec.BackupPath},
// 			},
// 		},
// 	}

// 	volLog.Info("PV created for Database Backup ")
// 	controllerutil.SetControllerReference(bkp, pv, scheme)
// 	return pv
// }

// NewDbBackupPVC Create a new PV Claim object for Database Backup
// func NewDbBackupPVC(bkp *v1alpha1.Backup, v *v1alpha1.PostgreSQL, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
// 	volLog.Info("Creating new PVC for Database Backup")
// 	labels := utils.PostgreSQLBkpLabels(bkp, "postgresql-backup")
// 	storageClassName := "manual"
// 	pvc := &corev1.PersistentVolumeClaim{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      GetPostgresqlBkpVolumeClaimName(bkp),
// 			Namespace: v.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: corev1.PersistentVolumeClaimSpec{
// 			StorageClassName: &storageClassName,
// 			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
// 			Resources: corev1.ResourceRequirements{
// 				Requests: corev1.ResourceList{
// 					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(bkp.Spec.BackupSize),
// 				},
// 			},
// 			VolumeName: GetPostgresqlBkpVolumeName(bkp),
// 		},
// 	}

// 	volLog.Info("PVC created for Database Backup ")
// 	controllerutil.SetControllerReference(bkp, pvc, scheme)
// 	return pvc
// }

// NewPostgreSqlPV Create a new PV object for PostgreSQL
func NewPostgreSqlPV(v *v1alpha1.PostgreSQL, scheme *runtime.Scheme) *corev1.PersistentVolume {
	volLog.Info("Creating new PV for PostgreSQL")
	labels := utils.Labels(v, "postgresql")
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: GetPostgresqlVolumeName(v),
			// Namespace: v.Namespace,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: "manual",
			Capacity: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(v.Spec.DataStorageSize),
			},
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: v.Spec.DataStoragePath},
			},
		},
	}

	volLog.Info("PV created for PostgreSQL ")
	controllerutil.SetControllerReference(v, pv, scheme)
	return pv
}

// NewPostgreSqlPVC Create a new PV Claim object for PostgreSQL
func NewPostgreSqlPVC(v *v1alpha1.PostgreSQL, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	volLog.Info("Creating new PVC for PostgreSQL")
	labels := utils.Labels(v, "postgresql")
	storageClassName := "manual"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetPostgresqlVolumeClaimName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(v.Spec.DataStorageSize),
				},
			},
			VolumeName: GetPostgresqlVolumeName(v),
		},
	}

	volLog.Info("PVC created for PostgreSQL ")
	controllerutil.SetControllerReference(v, pvc, scheme)
	return pvc
}
