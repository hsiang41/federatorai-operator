package updateenvvar

import (
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/util"
	securityv1 "github.com/openshift/api/security/v1"
)

func AssignServiceAccountsToSecurityContextConstraints(scc *securityv1.SecurityContextConstraints, ns string) {
	serviceAccount := "serviceaccount:" + ns
	for index, value := range scc.Users {
		if strings.Contains(value, util.NamespaceServiceAccount) {
			newUser := strings.Replace(scc.Users[index], util.NamespaceServiceAccount, serviceAccount, -1)
			scc.Users[index] = newUser
		}
	}
}
