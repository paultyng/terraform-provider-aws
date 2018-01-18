package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type DataSourceAvailabilityZone struct {
	Name       string `tf:",optional,computed"`
	Region     string `tf:",computed"`
	NameSuffix string `tf:",computed"`
	State      string `tf:",optional,computed"`
}

func (ds *DataSourceAvailabilityZone) ID() string {
	return ds.Name
}

func (ds *DataSourceAvailabilityZone) Read(meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	req := &ec2.DescribeAvailabilityZonesInput{}

	if ds.Name != "" {
		req.ZoneNames = []*string{aws.String(ds.Name)}
	}

	req.Filters = buildEC2AttributeFilterList(
		map[string]string{
			"state": ds.State,
		},
	)
	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		req.Filters = nil
	}

	log.Printf("[DEBUG] Reading Availability Zone: %s", req)
	resp, err := conn.DescribeAvailabilityZones(req)
	if err != nil {
		return err
	}
	if resp == nil || len(resp.AvailabilityZones) == 0 {
		return fmt.Errorf("no matching AZ found")
	}
	if len(resp.AvailabilityZones) > 1 {
		return fmt.Errorf("multiple AZs matched; use additional constraints to reduce matches to a single AZ")
	}

	az := resp.AvailabilityZones[0]

	ds.Name = *az.ZoneName
	ds.Region = *az.RegionName
	ds.State = *az.State

	// As a convenience when working with AZs generically, we expose
	// the AZ suffix alone, without the region name.
	// This can be used e.g. to create lookup tables by AZ letter that
	// work regardless of region.
	ds.NameSuffix = ds.Name[len(ds.Region):]

	return nil
}
