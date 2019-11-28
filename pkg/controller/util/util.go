package util

import (
	"context"
	"fmt"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"

	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	alamedaServiceLockName          = "alamedaservice-lock"
	alamedaServiceLockAnnotationKey = "alamedaservices.federatorai.containers.ai/name"
)

func GetAlamedaServiceLockName() string {
	return alamedaServiceLockName
}

func GetAlamedaServiceLockAnnotationKey() string {
	return alamedaServiceLockAnnotationKey
}

func GetAlamedaServiceLock(ctx context.Context, k8sClient client.Client) (rbacv1.ClusterRole, error) {
	lock := rbacv1.ClusterRole{}
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: alamedaServiceLockName}, &lock); err != nil {
		return lock, err
	}
	return lock, nil
}

func IsAlamedaServiceLockOwnedByAlamedaService(lock rbacv1.ClusterRole, alamedaService federatoraiv1alpha1.AlamedaService) bool {

	// For backward compatibility, keep previous logic that descides wheather lock is owned by AlamedaSerivce
	old := false
	for _, ownerReference := range lock.OwnerReferences {
		if ownerReference.UID == alamedaService.UID {
			old = true
			break
		}
	}

	new := lock.ObjectMeta.Annotations != nil &&
		lock.ObjectMeta.Annotations[alamedaServiceLockAnnotationKey] == fmt.Sprintf("%s/%s", alamedaService.Namespace, alamedaService.Name)

	return old || new
}
