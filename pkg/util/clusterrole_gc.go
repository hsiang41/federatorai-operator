package util

import (
	"context"

	"github.com/spf13/viper"
	rbacv1 "k8s.io/api/rbac/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetGCSecret(ctx context.Context, clnt client.Client) (rbacv1.ClusterRole, error) {
	name := "alameda-gc"
	if viper.IsSet("gcClusterRole") {
		name = viper.GetString("gcClusterRole")
	}
	clusterRole := rbacv1.ClusterRole{}

	if err := clnt.Get(ctx, client.ObjectKey{Name: name}, &clusterRole); err != nil {
		return clusterRole, err
	}
	return clusterRole, nil
}

func GetOrCreateGCCluster(clnt client.Client) (*rbacv1.ClusterRole, error) {
	gcClusterRoleName := "alameda-gc"
	if viper.IsSet("gcClusterRole") {
		gcClusterRoleName = viper.GetString("gcClusterRole")
	}
	gcClusterRole := &rbacv1.ClusterRole{}

	err := clnt.Get(context.TODO(), client.ObjectKey{
		Name: gcClusterRoleName}, gcClusterRole)
	if err != nil && k8sErrors.IsNotFound(err) {
		newGCClusterRole := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: gcClusterRoleName,
			},
		}
		err := clnt.Create(context.TODO(), newGCClusterRole)
		if err == nil {
			return newGCClusterRole, nil
		}
	} else if err != nil {
		return nil, err
	}
	return gcClusterRole, nil
}
