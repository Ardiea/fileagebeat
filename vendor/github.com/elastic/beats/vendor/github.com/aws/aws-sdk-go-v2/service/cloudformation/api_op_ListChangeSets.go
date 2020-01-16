// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package cloudformation

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

// The input for the ListChangeSets action.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15/ListChangeSetsInput
type ListChangeSetsInput struct {
	_ struct{} `type:"structure"`

	// A string (provided by the ListChangeSets response output) that identifies
	// the next page of change sets that you want to retrieve.
	NextToken *string `min:"1" type:"string"`

	// The name or the Amazon Resource Name (ARN) of the stack for which you want
	// to list change sets.
	//
	// StackName is a required field
	StackName *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s ListChangeSetsInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ListChangeSetsInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "ListChangeSetsInput"}
	if s.NextToken != nil && len(*s.NextToken) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("NextToken", 1))
	}

	if s.StackName == nil {
		invalidParams.Add(aws.NewErrParamRequired("StackName"))
	}
	if s.StackName != nil && len(*s.StackName) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("StackName", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// The output for the ListChangeSets action.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15/ListChangeSetsOutput
type ListChangeSetsOutput struct {
	_ struct{} `type:"structure"`

	// If the output exceeds 1 MB, a string that identifies the next page of change
	// sets. If there is no additional page, this value is null.
	NextToken *string `min:"1" type:"string"`

	// A list of ChangeSetSummary structures that provides the ID and status of
	// each change set for the specified stack.
	Summaries []ChangeSetSummary `type:"list"`
}

// String returns the string representation
func (s ListChangeSetsOutput) String() string {
	return awsutil.Prettify(s)
}

const opListChangeSets = "ListChangeSets"

// ListChangeSetsRequest returns a request value for making API operation for
// AWS CloudFormation.
//
// Returns the ID and status of each active change set for a stack. For example,
// AWS CloudFormation lists change sets that are in the CREATE_IN_PROGRESS or
// CREATE_PENDING state.
//
//    // Example sending a request using ListChangeSetsRequest.
//    req := client.ListChangeSetsRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15/ListChangeSets
func (c *Client) ListChangeSetsRequest(input *ListChangeSetsInput) ListChangeSetsRequest {
	op := &aws.Operation{
		Name:       opListChangeSets,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &ListChangeSetsInput{}
	}

	req := c.newRequest(op, input, &ListChangeSetsOutput{})
	return ListChangeSetsRequest{Request: req, Input: input, Copy: c.ListChangeSetsRequest}
}

// ListChangeSetsRequest is the request type for the
// ListChangeSets API operation.
type ListChangeSetsRequest struct {
	*aws.Request
	Input *ListChangeSetsInput
	Copy  func(*ListChangeSetsInput) ListChangeSetsRequest
}

// Send marshals and sends the ListChangeSets API request.
func (r ListChangeSetsRequest) Send(ctx context.Context) (*ListChangeSetsResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ListChangeSetsResponse{
		ListChangeSetsOutput: r.Request.Data.(*ListChangeSetsOutput),
		response:             &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ListChangeSetsResponse is the response type for the
// ListChangeSets API operation.
type ListChangeSetsResponse struct {
	*ListChangeSetsOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// ListChangeSets request.
func (r *ListChangeSetsResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
