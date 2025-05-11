package cognitarium

import (
	"context"
	"errors"
	"fmt"
	"strings"

	schema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6"
	"google.golang.org/grpc"
)

const W3IDPrefix = "https://w3id.org/axone/ontology/v4"

var (
	VcBodySubject = schema.IRI_Full("dataverse:credential:body#subject")
	VcBodyType    = schema.IRI_Full("dataverse:credential:body#type")
	VcBodyClaim   = schema.IRI_Full("dataverse:credential:body#claim")
)

var (
	ErrNoResult    = errors.New("no result")
	ErrVarNotFound = errors.New("variable not found")
	ErrTypeMimatch = errors.New("variable type mismatch")
)

func ref[T any](v T) *T {
	return &v
}

//nolint:funlen
func GetResourceGovAddrQuery(resource string) schema.SelectQuery {
	return schema.SelectQuery{
		Limit: ref(1),
		Prefixes: []schema.Prefix{
			{
				Prefix:    "gov",
				Namespace: fmt.Sprintf("%s/schema/credential/governance/text/", W3IDPrefix),
			},
		},
		Select: []schema.SelectItem{
			{
				Variable: ref(schema.SelectItem_Variable("code")),
			},
		},
		Where: schema.WhereClause{
			Bgp: &schema.WhereClause_Bgp{
				Patterns: []schema.TriplePattern{
					{
						Subject: schema.VarOrNode{Variable: ref(schema.VarOrNode_Variable("credId"))},
						Predicate: schema.VarOrNamedNode{
							NamedNode: &schema.VarOrNamedNode_NamedNode{Full: &VcBodySubject},
						},
						Object: schema.VarOrNodeOrLiteral{
							Node: &schema.VarOrNodeOrLiteral_Node{
								NamedNode: &schema.Node_NamedNode{Full: ref(schema.IRI_Full(resource))},
							},
						},
					},
					{
						Subject: schema.VarOrNode{Variable: ref(schema.VarOrNode_Variable("credId"))},
						Predicate: schema.VarOrNamedNode{
							NamedNode: &schema.VarOrNamedNode_NamedNode{Full: &VcBodyType},
						},
						Object: schema.VarOrNodeOrLiteral{
							Node: &schema.VarOrNodeOrLiteral_Node{
								NamedNode: &schema.Node_NamedNode{Prefixed: ref(schema.IRI_Prefixed("gov:GovernanceTextCredential"))},
							},
						},
					},
					{
						Subject: schema.VarOrNode{Variable: ref(schema.VarOrNode_Variable("credId"))},
						Predicate: schema.VarOrNamedNode{
							NamedNode: &schema.VarOrNamedNode_NamedNode{Full: &VcBodyClaim},
						},
						Object: schema.VarOrNodeOrLiteral{Variable: ref(schema.VarOrNodeOrLiteral_Variable("claim"))},
					},
					{
						Subject: schema.VarOrNode{Variable: ref(schema.VarOrNode_Variable("claim"))},
						Predicate: schema.VarOrNamedNode{
							NamedNode: &schema.VarOrNamedNode_NamedNode{Prefixed: ref(schema.IRI_Prefixed("gov:isGovernedBy"))},
						},
						Object: schema.VarOrNodeOrLiteral{Variable: ref(schema.VarOrNodeOrLiteral_Variable("gov"))},
					},
					{
						Subject: schema.VarOrNode{Variable: ref(schema.VarOrNode_Variable("gov"))},
						Predicate: schema.VarOrNamedNode{
							NamedNode: &schema.VarOrNamedNode_NamedNode{Prefixed: ref(schema.IRI_Prefixed("gov:fromGovernance"))},
						},
						Object: schema.VarOrNodeOrLiteral{Variable: ref(schema.VarOrNodeOrLiteral_Variable("code"))},
					},
				},
			},
		},
	}
}

// GetGovernanceAddressForResource queries the governance address for a given resource DID.
func GetGovernanceAddressForResource(
	ctx context.Context, cc grpc.ClientConnInterface, address string, resourceDID string,
) (string, error) {
	query := GetResourceGovAddrQuery(resourceDID)
	response, err := Select(ctx, cc, address, &schema.QueryMsg_Select{Query: query})
	if err != nil {
		return "", err
	}

	if len(response.Results.Bindings) != 1 {
		return "", ErrNoResult
	}

	codeBinding, ok := response.Results.Bindings[0]["code"]
	if !ok {
		return "", fmt.Errorf("variable %q: %w", "code", ErrVarNotFound)
	}
	code, ok := codeBinding.ValueType.(schema.URI)
	if !ok {
		return "", fmt.Errorf("expected URI, got %T: %w", codeBinding.ValueType, ErrTypeMimatch)
	}

	if code.Value.Full == nil {
		return "", fmt.Errorf("code URI is nil")
	}
	codeURI := string(*code.Value.Full)
	addr := codeURI
	if i := strings.LastIndex(codeURI, ":"); i != -1 {
		addr = codeURI[i+1:]
	}

	return addr, nil
}
