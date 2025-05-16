package cmd

/*
Copyright © 2020 Peter Howe <pnhowe@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

var configSetName, configSetValue, configDeleteName string
var configFull, detailIsPrimary bool
var detailHostname, detailSite, detailBlueprint, detailFoundation, detailInterfaceName string
var detailPrimary int
var detailSecondary string
var detailSubnet, detailReason, detailPXE string
var detailPrefix, detailGatewayOffset, detailOffset int
var scriptFile string
var detailAddParent, detailDeleteParent, detailAddFoundationBluePrint, detailDeleteFoundationBluePrint, detailAddType, detailDeleteType, detailAddIfaceName, detailDeleteIfaceName string
var detailName, detailDescription, detailParent string
var detailZone int
var detailDatacenter, detailCluster string
var detailHost int
var detailBuiltPercentage int
var detailMembers []int
var detailUsername, detailPassword string
var detailAMTUsername, detailAMTPassword, detailAMTIP string
var detailIPMIUsername, detailIPMIPassword, detailIPMIIP, detailIPMISOL string
var detailRedFishUsername, detailRedFishPassword, detailRedFishIP, detailRedFishSOL string
var detailLocator, detailPlot, detailComplex, detailPhysicalLocation, detailLinkName, detailMac, detailPxeName string
var detailIsProvisioning bool
var detailNetwork int
var detailAddressBlock, detailVlan, detailMTU int
var detailFailLikelihood, detailDelayVariance int
