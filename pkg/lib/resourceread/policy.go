package resourceread

import (
	"github.com/pkg/errors"

	v1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	policyScheme = runtime.NewScheme()
	policyCodecs = serializer.NewCodecFactory(policyScheme)
)

func init() {
	if err := v1beta1.AddToScheme(policyScheme); err != nil {
		log.Error(err, "Fail AddToScheme")
	}
}

func ReadPolicyV1beta1PodSecurityPolicy(objBytes []byte) (v1beta1.PodSecurityPolicy, error) {
	requiredObj, err := runtime.Decode(policyCodecs.UniversalDecoder(v1beta1.SchemeGroupVersion), objBytes)
	if err != nil {
		return v1beta1.PodSecurityPolicy{}, errors.Wrap(err, "decode PodSecurityPolicy failed")
	}
	return *requiredObj.(*v1beta1.PodSecurityPolicy), nil
}
