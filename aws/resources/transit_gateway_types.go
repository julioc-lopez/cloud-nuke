package resources

import (
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/go-commons/errors"
)

// TransitGatewaysVpcAttachment - represents all transit gateways vpc attachments
type TransitGatewaysVpcAttachment struct {
	Client ec2iface.EC2API
	Region string
	Ids    []string
}

func (tgw *TransitGatewaysVpcAttachment) Init(session *session.Session) {
	tgw.Client = ec2.New(session)
}

// ResourceName - the simple name of the aws resource
func (tgw *TransitGatewaysVpcAttachment) ResourceName() string {
	return "transit-gateway-attachment"
}

// MaxBatchSize - Tentative batch size to ensure AWS doesn't throttle
func (tgw *TransitGatewaysVpcAttachment) MaxBatchSize() int {
	return maxBatchSize
}

// ResourceIdentifiers - The Ids of the transit gateways
func (tgw *TransitGatewaysVpcAttachment) ResourceIdentifiers() []string {
	return tgw.Ids
}

func (tgw *TransitGatewaysVpcAttachment) GetAndSetIdentifiers(configObj config.Config) ([]string, error) {
	identifiers, err := tgw.getAll(configObj)
	if err != nil {
		return nil, err
	}

	tgw.Ids = awsgo.StringValueSlice(identifiers)
	return tgw.Ids, nil
}

// Nuke - nuke 'em all!!!
func (tgw *TransitGatewaysVpcAttachment) Nuke(identifiers []string) error {
	if err := tgw.nukeAll(awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

// TransitGatewaysRouteTables - represents all transit gateways route tables
type TransitGatewaysRouteTables struct {
	Client ec2iface.EC2API
	Region string
	Ids    []string
}

func (tgw *TransitGatewaysRouteTables) Init(session *session.Session) {
	tgw.Client = ec2.New(session)
}

// ResourceName - the simple name of the aws resource
func (tgw *TransitGatewaysRouteTables) ResourceName() string {
	return "transit-gateway-route-table"
}

// MaxBatchSize - Tentative batch size to ensure AWS doesn't throttle
func (tgw *TransitGatewaysRouteTables) MaxBatchSize() int {
	return maxBatchSize
}

// ResourceIdentifiers - The arns of the transit gateways route tables
func (tgw *TransitGatewaysRouteTables) ResourceIdentifiers() []string {
	return tgw.Ids
}

func (tgw *TransitGatewaysRouteTables) GetAndSetIdentifiers(configObj config.Config) ([]string, error) {
	identifiers, err := tgw.getAll(configObj)
	if err != nil {
		return nil, err
	}

	tgw.Ids = awsgo.StringValueSlice(identifiers)
	return tgw.Ids, nil
}

// Nuke - nuke 'em all!!!
func (tgw *TransitGatewaysRouteTables) Nuke(identifiers []string) error {
	if err := tgw.nukeAll(awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

// TransitGateways - represents all transit gateways
type TransitGateways struct {
	Client ec2iface.EC2API
	Region string
	Ids    []string
}

func (tgw *TransitGateways) Init(session *session.Session) {
	tgw.Client = ec2.New(session)
}

// ResourceName - the simple name of the aws resource
func (tgw *TransitGateways) ResourceName() string {
	return "transit-gateway"
}

// MaxBatchSize - Tentative batch size to ensure AWS doesn't throttle
func (tgw *TransitGateways) MaxBatchSize() int {
	return maxBatchSize
}

// ResourceIdentifiers - The Ids of the transit gateways
func (tgw *TransitGateways) ResourceIdentifiers() []string {
	return tgw.Ids
}

func (tgw *TransitGateways) GetAndSetIdentifiers(configObj config.Config) ([]string, error) {
	identifiers, err := tgw.getAll(configObj)
	if err != nil {
		return nil, err
	}

	tgw.Ids = awsgo.StringValueSlice(identifiers)
	return tgw.Ids, nil
}

// Nuke - nuke 'em all!!!
func (tgw *TransitGateways) Nuke(identifiers []string) error {
	if err := tgw.nukeAll(awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
