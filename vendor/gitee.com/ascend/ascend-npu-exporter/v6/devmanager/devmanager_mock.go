/* Copyright(C) 2021-2023. Huawei Technologies Co.,Ltd. All rights reserved.
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

// Package devmanager this for device driver manager mock
package devmanager

import (
	"huawei.com/npu-exporter/v6/devmanager/common"
	"huawei.com/npu-exporter/v6/devmanager/dcmi"
)

// DeviceManagerMock common device manager mock for Ascend910/310P/310
type DeviceManagerMock struct {
}

// Init load symbol and initialize dcmi
func (d *DeviceManagerMock) Init() error {
	return nil
}

// ShutDown clean the dynamically loaded resource
func (d *DeviceManagerMock) ShutDown() error {
	return nil
}

// GetDevType return mock type
func (d *DeviceManagerMock) GetDevType() string {
	return common.Ascend910
}

// GetDeviceCount get npu device count
func (d *DeviceManagerMock) GetDeviceCount() (int32, error) {
	return 1, nil
}

// GetCardList  get all card list
func (d *DeviceManagerMock) GetCardList() (int32, []int32, error) {
	return 1, []int32{0}, nil
}

// GetDeviceNumInCard  get all device list in one card
func (d *DeviceManagerMock) GetDeviceNumInCard(cardID int32) (int32, error) {
	return 1, nil
}

// GetDeviceList get all device logicID list
func (d *DeviceManagerMock) GetDeviceList() (int32, []int32, error) {
	return 1, []int32{0}, nil
}

// GetDeviceHealth query npu device health status
func (d *DeviceManagerMock) GetDeviceHealth(logicID int32) (uint32, error) {
	return 0, nil
}

// GetDeviceNetWorkHealth query npu device network health status
func (d *DeviceManagerMock) GetDeviceNetWorkHealth(logicID int32) (uint32, error) {
	return 0, nil
}

// GetDeviceUtilizationRate get npu device utilization
func (d *DeviceManagerMock) GetDeviceUtilizationRate(logicID int32, deviceType common.DeviceType) (uint32, error) {
	return 1, nil
}

// GetDeviceTemperature get npu device temperature
func (d *DeviceManagerMock) GetDeviceTemperature(logicID int32) (int32, error) {
	return 1, nil
}

// GetDeviceVoltage get npu device voltage
func (d *DeviceManagerMock) GetDeviceVoltage(logicID int32) (float32, error) {
	return 1, nil
}

// GetDevicePowerInfo get npu device power info
func (d *DeviceManagerMock) GetDevicePowerInfo(logicID int32) (float32, error) {
	return 1, nil
}

// GetDeviceFrequency get npu device work frequency
func (d *DeviceManagerMock) GetDeviceFrequency(logicID int32, deviceType common.DeviceType) (uint32, error) {
	return 1, nil
}

// GetDeviceMemoryInfo get npu memory information
func (d *DeviceManagerMock) GetDeviceMemoryInfo(logicID int32) (*common.MemoryInfo, error) {
	return &common.MemoryInfo{
		MemorySize:      1,
		MemoryAvailable: 1,
		Frequency:       1,
		Utilization:     1,
	}, nil
}

// GetDeviceHbmInfo get npu HBM module memory and frequency information
func (d *DeviceManagerMock) GetDeviceHbmInfo(logicID int32) (*common.HbmInfo, error) {
	return &common.HbmInfo{
		MemorySize:        1,
		Frequency:         1,
		Usage:             1,
		Temp:              1,
		BandWidthUtilRate: 1,
	}, nil
}

// GetDeviceErrorCode get npu device error code
func (d *DeviceManagerMock) GetDeviceErrorCode(logicID int32) (int32, int64, error) {
	return int32(0), int64(0), nil
}

// GetChipInfo get npu device error code
func (d *DeviceManagerMock) GetChipInfo(logicID int32) (*common.ChipInfo, error) {
	chip := &common.ChipInfo{
		Type:    "ascend",
		Name:    common.Chip910,
		Version: "v1",
	}
	return chip, nil
}

// GetPhysicIDFromLogicID get device physic id from logic id
func (d *DeviceManagerMock) GetPhysicIDFromLogicID(logicID int32) (int32, error) {
	return 1, nil
}

// GetLogicIDFromPhysicID get device logic id from physic id
func (d *DeviceManagerMock) GetLogicIDFromPhysicID(physicID int32) (int32, error) {
	return 1, nil
}

// GetDeviceLogicID get device logic id from card id and device id
func (d *DeviceManagerMock) GetDeviceLogicID(cardID, deviceID int32) (int32, error) {
	return 1, nil
}

// GetDeviceIPAddress get device ip address
func (d *DeviceManagerMock) GetDeviceIPAddress(logicID, ipType int32) (string, error) {
	if ipType == 0 {
		return "127.0.0.1", nil
	}
	return "::1", nil
}

// CreateVirtualDevice create virtual device
func (d *DeviceManagerMock) CreateVirtualDevice(logicID int32, vDevInfo common.CgoCreateVDevRes) (common.
	CgoCreateVDevOut, error) {
	return common.CgoCreateVDevOut{}, nil
}

// GetVirtualDeviceInfo get virtual device info
func (d *DeviceManagerMock) GetVirtualDeviceInfo(logicID int32) (common.VirtualDevInfo, error) {
	return common.VirtualDevInfo{}, nil
}

// DestroyVirtualDevice destroy virtual device
func (d *DeviceManagerMock) DestroyVirtualDevice(logicID int32, vDevID uint32) error {
	return nil
}

// GetMcuPowerInfo get mcu power info for cardID
func (d *DeviceManagerMock) GetMcuPowerInfo(cardID int32) (float32, error) {
	return 1, nil
}

// GetCardIDDeviceID get cardID and deviceID by logicID
func (d *DeviceManagerMock) GetCardIDDeviceID(logicID int32) (int32, int32, error) {
	return 0, 0, nil
}

// GetProductType get product type success
func (d *DeviceManagerMock) GetProductType(cardID, deviceID int32) (string, error) {
	return "", nil
}

// GetAllProductType get all product type success
func (d *DeviceManagerMock) GetAllProductType() ([]string, error) {
	return []string{}, nil
}

// GetNpuWorkMode get npu chip work mode SMP success
func (d *DeviceManagerMock) GetNpuWorkMode() string {
	return common.SMPMode
}

// SetDeviceReset set device reset success
func (d *DeviceManagerMock) SetDeviceReset(cardID, deviceID int32) error {
	return nil
}

// GetDeviceBootStatus get device boot status success
func (d *DeviceManagerMock) GetDeviceBootStatus(logicID int32) (int, error) {
	return common.BootStartFinish, nil
}

// GetDeviceAllErrorCode get device all error code success
func (d *DeviceManagerMock) GetDeviceAllErrorCode(logicID int32) (int32, []int64, error) {
	return 0, []int64{}, nil
}

// SubscribeDeviceFaultEvent subscribe device fault event success
func (d *DeviceManagerMock) SubscribeDeviceFaultEvent(logicID int32) error {
	return nil
}

// SetFaultEventCallFunc set fault event call func success
func (d *DeviceManagerMock) SetFaultEventCallFunc(businessFunc func(common.DevFaultInfo)) error {
	return nil
}

// GetDieID get die id success
func (d *DeviceManagerMock) GetDieID(logicID int32, dcmiDieType dcmi.DcmiDieType) (string, error) {
	return "ABCDEFGHIGKLMNOPQRSTUVWXYZ01234567890123", nil
}

// GetDevProcessInfo get process info
func (d *DeviceManagerMock) GetDevProcessInfo(logicID int32) (*common.DevProcessInfo, error) {
	return &common.DevProcessInfo{}, nil
}

func (d *DeviceManagerMock) GetPCIeBusInfo(logicID int32) (string, error) {
	return "0000:61:00.0", nil
}

func (d *DeviceManagerMock) GetBoardInfo(logicID int32) (common.BoardInfo, error) {
	return common.BoardInfo{}, nil
}

// GetProductTypeArray test for get product type array
func (d *DeviceManagerMock) GetProductTypeArray() []string {
	return []string{common.Atlas200ISoc}
}

// GetPCIEBandwidth get pcie bandwidth
func (d *DeviceManagerMock) GetPCIEBandwidth(logicID int32, _ int) (common.PCIEBwStat, error) {
	return common.PCIEBwStat{}, nil
}

// SetIsTrainingCard set IsTrainingCard
func (d *DeviceManagerMock) SetIsTrainingCard() error {
	return nil
}

// IsTrainingCard get IsTrainingCard
func (d *DeviceManagerMock) IsTrainingCard() bool {
	return true
}

// GetDcmiVersion get dcmi version
func (d *DeviceManagerMock) GetDcmiVersion() string {
	return "v1"
}

// GetValidChipInfo get valid chip info from all npu
func (d *DeviceManagerMock) GetValidChipInfo() (common.ChipInfo, error) {
	return common.ChipInfo{}, nil
}
