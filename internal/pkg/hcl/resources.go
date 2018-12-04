package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

func getResourceContent(block *hcl.Block) *hcl.BodyContent {
	pluginName := block.Labels[0]
	resourceBlockSchema := GetResourceBlockSchema(pluginName)

	// content, remain, diags
	content, _, diags := block.Body.PartialContent(resourceBlockSchema)

	if diags.HasErrors() {
		fmt.Println("Error getting resource body partial content:", diags)
	}

	return content
}

// GetExprVars takes a block and a list of attribute names and will return a
// map of all the expression variables for those attributes.
func GetExprVars(block *hcl.Block) map[string][]hcl.Traversal {
	content := getResourceContent(block)

	var exprVars = map[string][]hcl.Traversal{}

	for attr := range content.Attributes {
		exprVars[attr] = content.Attributes[attr].Expr.Variables()

		// support for_each by deleting "each" line from the exprVars
		// otherwise this will cause "panic: simple: adding self edge"
		// in graph.go - addEdgesToGraph
		for _, contentVars := range exprVars[attr] {
			rootName := contentVars.RootName()

			// this might delete a full expr causing a bug if the line is something
			// like: "${each.key} ${some_module.key}" where "some_module.key" also
			// gets deleted?
			if rootName == "each" {
				delete(exprVars, attr)
			}
		}
	}

	return exprVars
}
