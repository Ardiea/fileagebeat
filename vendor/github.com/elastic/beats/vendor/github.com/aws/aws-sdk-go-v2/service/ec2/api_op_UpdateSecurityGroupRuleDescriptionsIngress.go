// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/UpdateSecurityGroupRuleDescriptionsIngressRequest
type UpdateSecurityGroupRuleDescriptionsIngressInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// The ID of the security group. You must specify either the security group
	// ID or the security group name in the request. For security groups in a nondefault
	// VPC, you must specify the security group ID.
	GroupId *string `type:"string"`

	// [EC2-Classic, default VPC] The name of the security group. You must specify
	// either the security group ID or the security group name in the request.
	GroupName *string `type:"string"`

	// The IP permissions for the security group rule.
	//
	// IpPermissions is a required field
	IpPermissions []IpPermission `locationNameList:"item" type:"list" required:"true"`
}

// String returns the string representation
func (s UpdateSecurityGroupRuleDescriptionsIngressInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *UpdateSecurityGroupRuleDescriptionsIngressInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "UpdateSecurityGroupRuleDescriptionsIngressInput"}

	if s.IpPermissions == nil {
		invalidParams.Add(aws.NewErrParamRequired("IpPermissions"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/UpdateSecurityGroupRuleDescriptionsIngressResult
type UpdateSecurityGroupRuleDescriptionsIngressOutput struct {
	_ struct{} `type:"structure"`

	// Returns true if the request succeeds; otherwise, returns an error.
	Return *bool `locationName:"return" type:"boolean"`
}

// String returns the string representation
func (s UpdateSecurityGroupRuleDescriptionsIngressOutput) String() string {
	return awsutil.Prettify(s)
}

const opUpdateSecurityGroupRuleDescriptionsIngress = "UpdateSecurityGroupRuleDescriptionsIngress"

// UpdateSecurityGroupRuleDescriptionsIngressRequest returns a request value for making API operation for
// Amazon Elastic Compute Cloud.
//
// Updates the description of an ingress (inbound) security group rule. You
// can replace an existing description, or add a description to a rule that
// did not have one previously.
//
// You specify the description as part of the IP permissions structure. You
// can remove a description for a security group rule by omitting the description
// parameter in the request.
//
//    // Example sending a request using UpdateSecurityGroupRuleDescriptionsIngressRequest.
//    req := client.UpdateSecurityGroupRuleDescriptionsIngressRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/UpdateSecurityGroupRuleDescriptionsIngress
func (c *Client) UpdateSecurityGroupRuleDescriptionsIngressRequest(input *UpdateSecurityGroupRuleDescriptionsIngressInput) UpdateSecurityGroupRuleDescriptionsIngressRequest {
	op := &aws.Operation{
		Name:       opUpdateSecurityGroupRuleDescriptionsIngress,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &UpdateSecurityGroupRuleDescriptionsIngressInput{}
	}

	req := c.newRequest(op, input, &UpdateSecurityGroupRuleDescriptionsIngressOutput{})
	return UpdateSecurityGroupRuleDescriptionsIngressRequest{Request: req, Input: input, Copy: c.UpdateSecurityGroupRuleDescriptionsIngressRequest}
}

// UpdateSecurityGroupRuleDescriptionsIngressRequest is the request type for the
// UpdateSecurityGroupRuleDescriptionsIngress API operation.
type UpdateSecurityGroupRuleDescriptionsIngressRequest struct {
	*aws.Request
	Input *UpdateSecurityGroupRuleDescriptionsIngressInput
	Copy  func(*UpdateSecurityGroupRuleDescriptionsIngressInput) UpdateSecurityGroupRuleDescriptionsIngressRequest
}

// Send marshals and sends the UpdateSecurityGroupRuleDescriptionsIngress API request.
func (r UpdateSecurityGroupRuleDescriptionsIngressRequest) Send(ctx context.Context) (*UpdateSecurityGroupRuleDescriptionsIngressResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &UpdateSecurityGroupRuleDescriptionsIngressResponse{
		UpdateSecurityGroupRuleDescriptionsIngressOutput: r.Request.Data.(*UpdateSecurityGroupRuleDescriptionsIngressOutput),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// UpdateSecurityGroupRuleDescriptionsIngressResponse is the response type for the
// UpdateSecurityGroupRuleDescriptionsIngress API operation.
type UpdateSecurityGroupRuleDescriptionsIngressResponse struct {
	*UpdateSecurityGroupRuleDescriptionsIngressOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// UpdateSecurityGroupRuleDescriptionsIngress request.
func (r *UpdateSecurityGroupRuleDescriptionsIngressResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
