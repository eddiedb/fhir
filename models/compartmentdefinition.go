// Copyright (c) 2011-2017, HL7, Inc & The MITRE Corporation
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
//     * Redistributions of source code must retain the above copyright notice, this
//       list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright notice,
//       this list of conditions and the following disclaimer in the documentation
//       and/or other materials provided with the distribution.
//     * Neither the name of HL7 nor the names of its contributors may be used to
//       endorse or promote products derived from this software without specific
//       prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type CompartmentDefinition struct {
	DomainResource `bson:",inline"`
	Url            string                                   `bson:"url,omitempty" json:"url,omitempty"`
	Name           string                                   `bson:"name,omitempty" json:"name,omitempty"`
	Title          string                                   `bson:"title,omitempty" json:"title,omitempty"`
	Status         string                                   `bson:"status,omitempty" json:"status,omitempty"`
	Experimental   *bool                                    `bson:"experimental,omitempty" json:"experimental,omitempty"`
	Date           *FHIRDateTime                            `bson:"date,omitempty" json:"date,omitempty"`
	Publisher      string                                   `bson:"publisher,omitempty" json:"publisher,omitempty"`
	Contact        []ContactDetail                          `bson:"contact,omitempty" json:"contact,omitempty"`
	Description    string                                   `bson:"description,omitempty" json:"description,omitempty"`
	Purpose        string                                   `bson:"purpose,omitempty" json:"purpose,omitempty"`
	UseContext     []UsageContext                           `bson:"useContext,omitempty" json:"useContext,omitempty"`
	Jurisdiction   []CodeableConcept                        `bson:"jurisdiction,omitempty" json:"jurisdiction,omitempty"`
	Code           string                                   `bson:"code,omitempty" json:"code,omitempty"`
	Search         *bool                                    `bson:"search,omitempty" json:"search,omitempty"`
	Resource       []CompartmentDefinitionResourceComponent `bson:"resource,omitempty" json:"resource,omitempty"`
}

// Custom marshaller to add the resourceType property, as required by the specification
func (resource *CompartmentDefinition) MarshalJSON() ([]byte, error) {
	resource.ResourceType = "CompartmentDefinition"
	// Dereferencing the pointer to avoid infinite recursion.
	// Passing in plain old x (a pointer to CompartmentDefinition), would cause this same
	// MarshallJSON function to be called again
	return json.Marshal(*resource)
}

func (x *CompartmentDefinition) GetBSON() (interface{}, error) {
	x.ResourceType = "CompartmentDefinition"
	// See comment in MarshallJSON to see why we dereference
	return *x, nil
}

// The "compartmentDefinition" sub-type is needed to avoid infinite recursion in UnmarshalJSON
type compartmentDefinition CompartmentDefinition

// Custom unmarshaller to properly unmarshal embedded resources (represented as interface{})
func (x *CompartmentDefinition) UnmarshalJSON(data []byte) (err error) {
	x2 := compartmentDefinition{}
	if err = json.Unmarshal(data, &x2); err == nil {
		if x2.Contained != nil {
			for i := range x2.Contained {
				x2.Contained[i], err = MapToResource(x2.Contained[i], true)
				if err != nil {
					return err
				}
			}
		}
		*x = CompartmentDefinition(x2)
		return x.checkResourceType()
	}
	return
}

func (x *CompartmentDefinition) checkResourceType() error {
	if x.ResourceType == "" {
		x.ResourceType = "CompartmentDefinition"
	} else if x.ResourceType != "CompartmentDefinition" {
		return errors.New(fmt.Sprintf("Expected resourceType to be CompartmentDefinition, instead received %s", x.ResourceType))
	}
	return nil
}

type CompartmentDefinitionResourceComponent struct {
	BackboneElement `bson:",inline"`
	Code            string   `bson:"code,omitempty" json:"code,omitempty"`
	Param           []string `bson:"param,omitempty" json:"param,omitempty"`
	Documentation   string   `bson:"documentation,omitempty" json:"documentation,omitempty"`
}
