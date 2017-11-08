/*Copyright IT People Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

******************************************************************/

// Date Created:
// Author: David Medvedev + Leon Matthews
// Organization: IT People Corporation
// Last Update: Nov 08 2017
package main

type Part struct {
	PartNumber string `json:"partnumber"`
	SupplierID string `json:"supplierid"`
	AuditInfo  AuditInfo
}

type PartDetail struct {
	PartNumber  string `json:"partnumber"`
	Description string `json:"description"`
}

type PurchaseOrder struct {
	PONumber    string   `json:"poNumber"`
	PartNumbers []string `json:"parts"`
	AuditInfo   AuditInfo
}

type OrderStatus struct {
	PONumber string `json:"poNumber"`
	Status   string `json:"status"`
}

type AuditInfo struct {
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
	CreatedTS string `json:"createdTS"`
	UpdatedTS string `json:"updatedTS"`
}
