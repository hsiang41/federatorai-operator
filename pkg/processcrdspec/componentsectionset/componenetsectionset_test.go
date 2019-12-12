package componentsectionset

import (
	"github.com/stretchr/testify/assert"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestReplaceOrAppendEnvVar(t *testing.T) {

	type testCaseHave struct {
		source []corev1.EnvVar
		target []corev1.EnvVar
	}

	type testCase struct {
		have testCaseHave
		want []corev1.EnvVar
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{
				source: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "e1",
						Value: "e1v",
					},
					corev1.EnvVar{
						Name:  "e2",
						Value: "e2v",
					},
				},
				target: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "e3",
						Value: "e3v",
					},
				},
			},
			want: []corev1.EnvVar{
				corev1.EnvVar{
					Name:  "e1",
					Value: "e1v",
				},
				corev1.EnvVar{
					Name:  "e2",
					Value: "e2v",
				},
				corev1.EnvVar{
					Name:  "e3",
					Value: "e3v",
				},
			},
		},
		testCase{
			have: testCaseHave{
				source: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "e1",
						Value: "newe1v",
					},
					corev1.EnvVar{
						Name:  "e2",
						Value: "e2v",
					},
				},
				target: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "e1",
						Value: "e1v",
					},
				},
			},
			want: []corev1.EnvVar{
				corev1.EnvVar{
					Name:  "e1",
					Value: "newe1v",
				},
				corev1.EnvVar{
					Name:  "e2",
					Value: "e2v",
				},
			},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		actual := replaceOrAppendEnvVar(tc.have.target, tc.have.source)
		assert.ElementsMatch(tc.want, actual)
	}

}
