package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/stretchr/testify/assert"
)

func TestEffectiveBlueprint_GetWantedDogus(t *testing.T) {
	t.Run("should get only present dogus", func(t *testing.T) {
		ldapDogu := Dogu{
			Name: cescommons.QualifiedName{
				SimpleName: "ldap",
				Namespace:  "official",
			},
			Version: &version3213,
		}
		absentMysqlDogu := Dogu{
			Name: cescommons.QualifiedName{
				SimpleName: "mysql",
				Namespace:  "official",
			},
			Absent: true,
		}
		postgresqlDogu := Dogu{
			Name: cescommons.QualifiedName{
				SimpleName: "postgresql",
				Namespace:  "official",
			},
			Version: &version3213,
		}
		effectiveBlueprint := &EffectiveBlueprint{
			Dogus: []Dogu{ldapDogu, absentMysqlDogu, postgresqlDogu},
		}

		result := effectiveBlueprint.GetWantedDogus()

		assert.Len(t, result, 2)
		assert.Contains(t, result, ldapDogu)
		assert.Contains(t, result, postgresqlDogu)
	})
}
