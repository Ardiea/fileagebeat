// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"github.com/aws/aws-sdk-go-v2/private/protocol/query"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/iam-2010-05-08/DeleteAccountAliasRequest
type DeleteAccountAliasInput struct {
	_ struct{} `type:"structure"`

	// The name of the account alias to delete.
	//
	// This parameter allows (through its regex pattern (http://wikipedia.org/wiki/regex))
	// a string of characters consisting of lowercase letters, digits, and dashes.
	// You cannot start or finish with a dash, nor can you have two dashes in a
	// row.
	//
	// AccountAlias is a required field
	AccountAlias *string `min:"3" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteAccountAliasInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteAccountAliasInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "DeleteAccountAliasInput"}

	if s.AccountAlias == nil {
		invalidParams.Add(aws.NewErrParamRequired("AccountAlias"))
	}
	if s.AccountAlias != nil && len(*s.AccountAlias) < 3 {
		invalidParams.Add(aws.NewErrParamMinLen("AccountAlias", 3))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/iam-2010-05-08/DeleteAccountAliasOutput
type DeleteAccountAliasOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteAccountAliasOutput) String() string {
	return awsutil.Prettify(s)
}

const opDeleteAccountAlias = "DeleteAccountAlias"

// DeleteAccountAliasRequest returns a request value for making API operation for
// AWS Identity and Access Management.
//
// Deletes the specified AWS account alias. For information about using an AWS
// account alias, see Using an Alias for Your AWS Account ID (https://docs.aws.amazon.com/IAM/latest/UserGuide/AccountAlias.html)
// in the IAM User Guide.
//
//    // Example sending a request using DeleteAccountAliasRequest.
//    req := client.DeleteAccountAliasRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/iam-2010-05-08/DeleteAccountAlias
func (c *Client) DeleteAccountAliasRequest(input *DeleteAccountAliasInput) DeleteAccountAliasRequest {
	op := &aws.Operation{
		Name:       opDeleteAccountAlias,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DeleteAccountAliasInput{}
	}

	req := c.newRequest(op, input, &DeleteAccountAliasOutput{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return DeleteAccountAliasRequest{Request: req, Input: input, Copy: c.DeleteAccountAliasRequest}
}

// DeleteAccountAliasRequest is the request type for the
// DeleteAccountAlias API operation.
type DeleteAccountAliasRequest struct {
	*aws.Request
	Input *DeleteAccountAliasInput
	Copy  func(*DeleteAccountAliasInput) DeleteAccountAliasRequest
}

// Send marshals and sends the DeleteAccountAlias API request.
func (r DeleteAccountAliasRequest) Send(ctx context.Context) (*DeleteAccountAliasResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &DeleteAccountAliasResponse{
		DeleteAccountAliasOutput: r.Request.Data.(*DeleteAccountAliasOutput),
		response:                 &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// DeleteAccountAliasResponse is the response type for the
// DeleteAccountAlias API operation.
type DeleteAccountAliasResponse struct {
	*DeleteAccountAliasOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// DeleteAccountAlias request.
func (r *DeleteAccountAliasResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
