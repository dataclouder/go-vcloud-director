package types

import "encoding/xml"

// VAppV2 is functionally identical to a vApp. This type is the basis for a more capable vApp
type VAppV2 VApp

// VmType is an alias to CreateItem. This is used here to make easier the relationship with the API docs, where it is
// always referred to as VmType
type VmType CreateItem

// ComposeVAppParamsV2 is an enhanced and revised version of ComposeVAppParams
// This structure can handle multiple VMs at the same time, while the old one can't
type ComposeVAppParamsV2 struct {
	XMLName xml.Name `xml:"ComposeVAppParams"`
	Ovf     string   `xml:"xmlns:ovf,attr"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	// Attributes
	Name        string `xml:"name,attr,omitempty"`        // Typically used to name or identify the subject of the request. For example, the name of the object being created or modified.
	Deploy      bool   `xml:"deploy,attr"`                // True if the vApp should be deployed at instantiation. Defaults to true.
	PowerOn     bool   `xml:"powerOn,attr"`               // True if the vApp should be powered-on at instantiation. Defaults to true.
	LinkedClone bool   `xml:"linkedClone,attr,omitempty"` // Reserved. Unimplemented.
	// Elements
	Description         string                         `xml:"Description,omitempty"`         // Optional description.
	VAppParent          *Reference                     `xml:"VAppParent,omitempty"`          // Reserved. Unimplemented.
	InstantiationParams *InstantiationParams           `xml:"InstantiationParams,omitempty"` // Instantiation parameters for the composed vApp.
	SourcedItem         []*SourcedCompositionItemParam `xml:"SourcedItem,omitempty"`         // Composition items. One of: vApp vAppTemplate VM.
	CreateItem          []*VmType                      `xml:"CreateItem,omitempty"`          // Composition items. One of: vApp vAppTemplate VM.
	AllEULAsAccepted    bool                           `xml:"AllEULAsAccepted,omitempty"`    // True confirms acceptance of all EULAs in a vApp template. Instantiation fails if this element is missing, empty, or set to false and one or more EulaSection elements are present.
}

// ReComposeVAppParamsV2 is an enhanced and revised version of ReComposeVAppParams
// This structure can handle multiple VMs at the same time, while the old one can't
type ReComposeVAppParamsV2 struct {
	XMLName xml.Name `xml:"RecomposeVAppParams"`
	Ovf     string   `xml:"xmlns:ovf,attr"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	// Attributes
	Name        string `xml:"name,attr,omitempty"`        // Typically used to name or identify the subject of the request. For example, the name of the object being created or modified.
	Deploy      bool   `xml:"deploy,attr"`                // True if the vApp should be deployed at instantiation. Defaults to true.
	PowerOn     bool   `xml:"powerOn,attr"`               // True if the vApp should be powered-on at instantiation. Defaults to true.
	LinkedClone bool   `xml:"linkedClone,attr,omitempty"` // Reserved. Unimplemented.
	// Elements
	Description         string                         `xml:"Description,omitempty"`         // Optional description.
	VAppParent          *Reference                     `xml:"VAppParent,omitempty"`          // Reserved. Unimplemented.
	InstantiationParams *InstantiationParams           `xml:"InstantiationParams,omitempty"` // Instantiation parameters for the composed vApp.
	SourcedItem         []*SourcedCompositionItemParam `xml:"SourcedItem,omitempty"`         // Composition item. One of: vApp vAppTemplate VM.
	CreateItem          []*VmType                      `xml:"CreateItem,omitempty"`          // Composition items. One of: vApp vAppTemplate VM.
	ReconfigureItem     []*VmType                      `xml:"ReconfigureItem,omitempty"`     // Existing Vm to be reconfigured during recomposition.
	DeleteItem          *Reference                     `xml:"DeleteItem,omitempty"`          // Reference to a Vm to be deleted during recomposition.
	AllEULAsAccepted    bool                           `xml:"AllEULAsAccepted,omitempty"`    // True confirms acceptance of all EULAs in a vApp template. Instantiation fails if this element is missing, empty, or set to false and one or more EulaSection elements are present.
}