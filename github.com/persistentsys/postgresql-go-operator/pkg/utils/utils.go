package utils

import (
	"github.com/persistentsys/postgresql-go-operator/pkg/apis/postgresql/v1alpha1"
)

func Labels(v *v1alpha1.PostgreSQL, tier string) map[string]string {
	return map[string]string{
		"app":        "PostgreSQL",
		"PostgreSQL_cr": v.Name,
		"tier":       tier,
	}
}

