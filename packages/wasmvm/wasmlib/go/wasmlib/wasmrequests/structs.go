// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmrequests

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"

type CallRequest struct {
	// caller assets that the call is allowed to access
	Allowance []byte
	Contract  wasmtypes.ScHname
	Function  wasmtypes.ScHname
	Params    []byte
}

func NewCallRequestFromBytes(buf []byte) *CallRequest {
	dec := wasmtypes.NewWasmDecoder(buf)
	data := &CallRequest{}
	data.Allowance = wasmtypes.BytesDecode(dec)
	data.Contract  = wasmtypes.HnameDecode(dec)
	data.Function  = wasmtypes.HnameDecode(dec)
	data.Params    = wasmtypes.BytesDecode(dec)
	dec.Close()
	return data
}

func (o *CallRequest) Bytes() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.BytesEncode(enc, o.Allowance)
	wasmtypes.HnameEncode(enc, o.Contract)
	wasmtypes.HnameEncode(enc, o.Function)
	wasmtypes.BytesEncode(enc, o.Params)
	return enc.Buf()
}

type ImmutableCallRequest struct {
	Proxy wasmtypes.Proxy
}

func (o ImmutableCallRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o ImmutableCallRequest) Value() *CallRequest {
	return NewCallRequestFromBytes(o.Proxy.Get())
}

type MutableCallRequest struct {
	Proxy wasmtypes.Proxy
}

func (o MutableCallRequest) Delete() {
	o.Proxy.Delete()
}

func (o MutableCallRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o MutableCallRequest) SetValue(value *CallRequest) {
	o.Proxy.Set(value.Bytes())
}

func (o MutableCallRequest) Value() *CallRequest {
	return NewCallRequestFromBytes(o.Proxy.Get())
}

type DeployRequest struct {
	Description string
	Name        string
	Params      []byte
	ProgHash    wasmtypes.ScHash
}

func NewDeployRequestFromBytes(buf []byte) *DeployRequest {
	dec := wasmtypes.NewWasmDecoder(buf)
	data := &DeployRequest{}
	data.Description = wasmtypes.StringDecode(dec)
	data.Name        = wasmtypes.StringDecode(dec)
	data.Params      = wasmtypes.BytesDecode(dec)
	data.ProgHash    = wasmtypes.HashDecode(dec)
	dec.Close()
	return data
}

func (o *DeployRequest) Bytes() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.StringEncode(enc, o.Description)
	wasmtypes.StringEncode(enc, o.Name)
	wasmtypes.BytesEncode(enc, o.Params)
	wasmtypes.HashEncode(enc, o.ProgHash)
	return enc.Buf()
}

type ImmutableDeployRequest struct {
	Proxy wasmtypes.Proxy
}

func (o ImmutableDeployRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o ImmutableDeployRequest) Value() *DeployRequest {
	return NewDeployRequestFromBytes(o.Proxy.Get())
}

type MutableDeployRequest struct {
	Proxy wasmtypes.Proxy
}

func (o MutableDeployRequest) Delete() {
	o.Proxy.Delete()
}

func (o MutableDeployRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o MutableDeployRequest) SetValue(value *DeployRequest) {
	o.Proxy.Set(value.Bytes())
}

func (o MutableDeployRequest) Value() *DeployRequest {
	return NewDeployRequestFromBytes(o.Proxy.Get())
}

type PostRequest struct {
	// caller assets that the call is allowed to access
	Allowance []byte
	ChainID   wasmtypes.ScChainID
	Contract  wasmtypes.ScHname
	Delay     uint32
	Function  wasmtypes.ScHname
	Params    []byte
	// assets that are transferred into caller account
	Transfer  []byte
}

func NewPostRequestFromBytes(buf []byte) *PostRequest {
	dec := wasmtypes.NewWasmDecoder(buf)
	data := &PostRequest{}
	data.Allowance = wasmtypes.BytesDecode(dec)
	data.ChainID   = wasmtypes.ChainIDDecode(dec)
	data.Contract  = wasmtypes.HnameDecode(dec)
	data.Delay     = wasmtypes.Uint32Decode(dec)
	data.Function  = wasmtypes.HnameDecode(dec)
	data.Params    = wasmtypes.BytesDecode(dec)
	data.Transfer  = wasmtypes.BytesDecode(dec)
	dec.Close()
	return data
}

func (o *PostRequest) Bytes() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.BytesEncode(enc, o.Allowance)
	wasmtypes.ChainIDEncode(enc, o.ChainID)
	wasmtypes.HnameEncode(enc, o.Contract)
	wasmtypes.Uint32Encode(enc, o.Delay)
	wasmtypes.HnameEncode(enc, o.Function)
	wasmtypes.BytesEncode(enc, o.Params)
	wasmtypes.BytesEncode(enc, o.Transfer)
	return enc.Buf()
}

type ImmutablePostRequest struct {
	Proxy wasmtypes.Proxy
}

func (o ImmutablePostRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o ImmutablePostRequest) Value() *PostRequest {
	return NewPostRequestFromBytes(o.Proxy.Get())
}

type MutablePostRequest struct {
	Proxy wasmtypes.Proxy
}

func (o MutablePostRequest) Delete() {
	o.Proxy.Delete()
}

func (o MutablePostRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o MutablePostRequest) SetValue(value *PostRequest) {
	o.Proxy.Set(value.Bytes())
}

func (o MutablePostRequest) Value() *PostRequest {
	return NewPostRequestFromBytes(o.Proxy.Get())
}

type SendRequest struct {
	Address  wasmtypes.ScAddress
	Transfer []byte
}

func NewSendRequestFromBytes(buf []byte) *SendRequest {
	dec := wasmtypes.NewWasmDecoder(buf)
	data := &SendRequest{}
	data.Address  = wasmtypes.AddressDecode(dec)
	data.Transfer = wasmtypes.BytesDecode(dec)
	dec.Close()
	return data
}

func (o *SendRequest) Bytes() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.AddressEncode(enc, o.Address)
	wasmtypes.BytesEncode(enc, o.Transfer)
	return enc.Buf()
}

type ImmutableSendRequest struct {
	Proxy wasmtypes.Proxy
}

func (o ImmutableSendRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o ImmutableSendRequest) Value() *SendRequest {
	return NewSendRequestFromBytes(o.Proxy.Get())
}

type MutableSendRequest struct {
	Proxy wasmtypes.Proxy
}

func (o MutableSendRequest) Delete() {
	o.Proxy.Delete()
}

func (o MutableSendRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o MutableSendRequest) SetValue(value *SendRequest) {
	o.Proxy.Set(value.Bytes())
}

func (o MutableSendRequest) Value() *SendRequest {
	return NewSendRequestFromBytes(o.Proxy.Get())
}

type TransferRequest struct {
	AgentID  wasmtypes.ScAgentID
	Transfer []byte
}

func NewTransferRequestFromBytes(buf []byte) *TransferRequest {
	dec := wasmtypes.NewWasmDecoder(buf)
	data := &TransferRequest{}
	data.AgentID  = wasmtypes.AgentIDDecode(dec)
	data.Transfer = wasmtypes.BytesDecode(dec)
	dec.Close()
	return data
}

func (o *TransferRequest) Bytes() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.AgentIDEncode(enc, o.AgentID)
	wasmtypes.BytesEncode(enc, o.Transfer)
	return enc.Buf()
}

type ImmutableTransferRequest struct {
	Proxy wasmtypes.Proxy
}

func (o ImmutableTransferRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o ImmutableTransferRequest) Value() *TransferRequest {
	return NewTransferRequestFromBytes(o.Proxy.Get())
}

type MutableTransferRequest struct {
	Proxy wasmtypes.Proxy
}

func (o MutableTransferRequest) Delete() {
	o.Proxy.Delete()
}

func (o MutableTransferRequest) Exists() bool {
	return o.Proxy.Exists()
}

func (o MutableTransferRequest) SetValue(value *TransferRequest) {
	o.Proxy.Set(value.Bytes())
}

func (o MutableTransferRequest) Value() *TransferRequest {
	return NewTransferRequestFromBytes(o.Proxy.Get())
}
