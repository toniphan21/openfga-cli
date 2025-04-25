package generate

import (
	"fmt"
	"testing"

	"github.com/openfga/cli/internal/authorizationmodel"
)

func TestCollectForListObject(t *testing.T) {
	inputModel := `
model
schema 1.1

type user

type tenant
	relations
		define admin: [user]
		define member: [user]

type project
	relations
		define tenant: [tenant]
		define owner: [user]
		define can_read: owner or (admin from tenant)
		define can_write: owner or (admin from tenant)
	`

	authModel := authorizationmodel.AuthzModel{}

	err := authModel.ReadModelFromString(inputModel, authorizationmodel.ModelFormatFGA)
	if err != nil {
		panic(err)
	}

	g := &PklGenerator{
		Model:      authModel.TypeDefinitions,
		Convention: &PklConvention{Config: make(map[string]PklConventionConfig)},
	}

	cAssignments := g.collectGenAssignments()
	for _, assignment := range cAssignments {
		fmt.Printf("name: %v\n", assignment.name)
		for _, relation := range assignment.relations {
			fmt.Printf("  relation: %v\n", relation.name)
			for on, object := range relation.objects {
				fmt.Printf("   object[%v]: %v\n", on, object.name)
			}
		}
		fmt.Println("----")
	}

	fmt.Println("==========")

	cAssertions := g.collectGenAssertions()
	for _, assertion := range cAssertions {
		fmt.Println(assertion.name)
		fmt.Println(assertion.relations)
		fmt.Println("----")
	}

	fmt.Println("==========")

	result := g.collectGenListObjects(cAssignments, cAssertions)
	for _, item := range result {
		fmt.Println(item.name)
		fmt.Println(item.relations)
		fmt.Println(item.accesses)
	}
}
