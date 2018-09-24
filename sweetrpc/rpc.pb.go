// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

package sweetrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type WpaConnectionUpdate_WpaConnectionUpdateState int32

const (
	WpaConnectionUpdate_CONNECTING WpaConnectionUpdate_WpaConnectionUpdateState = 0
	WpaConnectionUpdate_CONNECTED  WpaConnectionUpdate_WpaConnectionUpdateState = 1
	WpaConnectionUpdate_FAILED     WpaConnectionUpdate_WpaConnectionUpdateState = 2
)

var WpaConnectionUpdate_WpaConnectionUpdateState_name = map[int32]string{
	0: "CONNECTING",
	1: "CONNECTED",
	2: "FAILED",
}

var WpaConnectionUpdate_WpaConnectionUpdateState_value = map[string]int32{
	"CONNECTING": 0,
	"CONNECTED":  1,
	"FAILED":     2,
}

func (x WpaConnectionUpdate_WpaConnectionUpdateState) String() string {
	return proto.EnumName(WpaConnectionUpdate_WpaConnectionUpdateState_name, int32(x))
}

func (WpaConnectionUpdate_WpaConnectionUpdateState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{5, 0}
}

type GetInfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInfoRequest) Reset()         { *m = GetInfoRequest{} }
func (m *GetInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetInfoRequest) ProtoMessage()    {}
func (*GetInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{0}
}
func (m *GetInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetInfoRequest.Unmarshal(m, b)
}
func (m *GetInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetInfoRequest.Marshal(b, m, deterministic)
}
func (dst *GetInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInfoRequest.Merge(dst, src)
}
func (m *GetInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetInfoRequest.Size(m)
}
func (m *GetInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetInfoRequest proto.InternalMessageInfo

type GetInfoResponse struct {
	Serial               string   `protobuf:"bytes,1,opt,name=serial,proto3" json:"serial,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Commit               string   `protobuf:"bytes,3,opt,name=commit,proto3" json:"commit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInfoResponse) Reset()         { *m = GetInfoResponse{} }
func (m *GetInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetInfoResponse) ProtoMessage()    {}
func (*GetInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{1}
}
func (m *GetInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetInfoResponse.Unmarshal(m, b)
}
func (m *GetInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetInfoResponse.Marshal(b, m, deterministic)
}
func (dst *GetInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInfoResponse.Merge(dst, src)
}
func (m *GetInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetInfoResponse.Size(m)
}
func (m *GetInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetInfoResponse proto.InternalMessageInfo

func (m *GetInfoResponse) GetSerial() string {
	if m != nil {
		return m.Serial
	}
	return ""
}

func (m *GetInfoResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *GetInfoResponse) GetCommit() string {
	if m != nil {
		return m.Commit
	}
	return ""
}

type GetWpaConnectionInfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWpaConnectionInfoRequest) Reset()         { *m = GetWpaConnectionInfoRequest{} }
func (m *GetWpaConnectionInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetWpaConnectionInfoRequest) ProtoMessage()    {}
func (*GetWpaConnectionInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{2}
}
func (m *GetWpaConnectionInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWpaConnectionInfoRequest.Unmarshal(m, b)
}
func (m *GetWpaConnectionInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWpaConnectionInfoRequest.Marshal(b, m, deterministic)
}
func (dst *GetWpaConnectionInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWpaConnectionInfoRequest.Merge(dst, src)
}
func (m *GetWpaConnectionInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetWpaConnectionInfoRequest.Size(m)
}
func (m *GetWpaConnectionInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWpaConnectionInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetWpaConnectionInfoRequest proto.InternalMessageInfo

type GetWpaConnectionInfoResponse struct {
	Ssid                 string   `protobuf:"bytes,1,opt,name=ssid,proto3" json:"ssid,omitempty"`
	State                string   `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	Ip                   string   `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	Message              string   `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWpaConnectionInfoResponse) Reset()         { *m = GetWpaConnectionInfoResponse{} }
func (m *GetWpaConnectionInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetWpaConnectionInfoResponse) ProtoMessage()    {}
func (*GetWpaConnectionInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{3}
}
func (m *GetWpaConnectionInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWpaConnectionInfoResponse.Unmarshal(m, b)
}
func (m *GetWpaConnectionInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWpaConnectionInfoResponse.Marshal(b, m, deterministic)
}
func (dst *GetWpaConnectionInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWpaConnectionInfoResponse.Merge(dst, src)
}
func (m *GetWpaConnectionInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetWpaConnectionInfoResponse.Size(m)
}
func (m *GetWpaConnectionInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWpaConnectionInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetWpaConnectionInfoResponse proto.InternalMessageInfo

func (m *GetWpaConnectionInfoResponse) GetSsid() string {
	if m != nil {
		return m.Ssid
	}
	return ""
}

func (m *GetWpaConnectionInfoResponse) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *GetWpaConnectionInfoResponse) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *GetWpaConnectionInfoResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ConnectWpaNetworkRequest struct {
	Ssid                 string   `protobuf:"bytes,1,opt,name=ssid,proto3" json:"ssid,omitempty"`
	Psk                  string   `protobuf:"bytes,2,opt,name=psk,proto3" json:"psk,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectWpaNetworkRequest) Reset()         { *m = ConnectWpaNetworkRequest{} }
func (m *ConnectWpaNetworkRequest) String() string { return proto.CompactTextString(m) }
func (*ConnectWpaNetworkRequest) ProtoMessage()    {}
func (*ConnectWpaNetworkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{4}
}
func (m *ConnectWpaNetworkRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectWpaNetworkRequest.Unmarshal(m, b)
}
func (m *ConnectWpaNetworkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectWpaNetworkRequest.Marshal(b, m, deterministic)
}
func (dst *ConnectWpaNetworkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectWpaNetworkRequest.Merge(dst, src)
}
func (m *ConnectWpaNetworkRequest) XXX_Size() int {
	return xxx_messageInfo_ConnectWpaNetworkRequest.Size(m)
}
func (m *ConnectWpaNetworkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectWpaNetworkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectWpaNetworkRequest proto.InternalMessageInfo

func (m *ConnectWpaNetworkRequest) GetSsid() string {
	if m != nil {
		return m.Ssid
	}
	return ""
}

func (m *ConnectWpaNetworkRequest) GetPsk() string {
	if m != nil {
		return m.Psk
	}
	return ""
}

type WpaConnectionUpdate struct {
	Status               WpaConnectionUpdate_WpaConnectionUpdateState `protobuf:"varint,1,opt,name=status,proto3,enum=sweetrpc.WpaConnectionUpdate_WpaConnectionUpdateState" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                     `json:"-"`
	XXX_unrecognized     []byte                                       `json:"-"`
	XXX_sizecache        int32                                        `json:"-"`
}

func (m *WpaConnectionUpdate) Reset()         { *m = WpaConnectionUpdate{} }
func (m *WpaConnectionUpdate) String() string { return proto.CompactTextString(m) }
func (*WpaConnectionUpdate) ProtoMessage()    {}
func (*WpaConnectionUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{5}
}
func (m *WpaConnectionUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WpaConnectionUpdate.Unmarshal(m, b)
}
func (m *WpaConnectionUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WpaConnectionUpdate.Marshal(b, m, deterministic)
}
func (dst *WpaConnectionUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WpaConnectionUpdate.Merge(dst, src)
}
func (m *WpaConnectionUpdate) XXX_Size() int {
	return xxx_messageInfo_WpaConnectionUpdate.Size(m)
}
func (m *WpaConnectionUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_WpaConnectionUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_WpaConnectionUpdate proto.InternalMessageInfo

func (m *WpaConnectionUpdate) GetStatus() WpaConnectionUpdate_WpaConnectionUpdateState {
	if m != nil {
		return m.Status
	}
	return WpaConnectionUpdate_CONNECTING
}

type SubscribeWpaNetworkScanUpdatesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeWpaNetworkScanUpdatesRequest) Reset()         { *m = SubscribeWpaNetworkScanUpdatesRequest{} }
func (m *SubscribeWpaNetworkScanUpdatesRequest) String() string { return proto.CompactTextString(m) }
func (*SubscribeWpaNetworkScanUpdatesRequest) ProtoMessage()    {}
func (*SubscribeWpaNetworkScanUpdatesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{6}
}
func (m *SubscribeWpaNetworkScanUpdatesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeWpaNetworkScanUpdatesRequest.Unmarshal(m, b)
}
func (m *SubscribeWpaNetworkScanUpdatesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeWpaNetworkScanUpdatesRequest.Marshal(b, m, deterministic)
}
func (dst *SubscribeWpaNetworkScanUpdatesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeWpaNetworkScanUpdatesRequest.Merge(dst, src)
}
func (m *SubscribeWpaNetworkScanUpdatesRequest) XXX_Size() int {
	return xxx_messageInfo_SubscribeWpaNetworkScanUpdatesRequest.Size(m)
}
func (m *SubscribeWpaNetworkScanUpdatesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeWpaNetworkScanUpdatesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeWpaNetworkScanUpdatesRequest proto.InternalMessageInfo

type WpaNetwork struct {
	Bssid                string   `protobuf:"bytes,1,opt,name=bssid,proto3" json:"bssid,omitempty"`
	Frequency            string   `protobuf:"bytes,2,opt,name=frequency,proto3" json:"frequency,omitempty"`
	SignalLevel          string   `protobuf:"bytes,3,opt,name=signal_level,json=signalLevel,proto3" json:"signal_level,omitempty"`
	Flags                string   `protobuf:"bytes,4,opt,name=flags,proto3" json:"flags,omitempty"`
	Ssid                 string   `protobuf:"bytes,5,opt,name=ssid,proto3" json:"ssid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WpaNetwork) Reset()         { *m = WpaNetwork{} }
func (m *WpaNetwork) String() string { return proto.CompactTextString(m) }
func (*WpaNetwork) ProtoMessage()    {}
func (*WpaNetwork) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{7}
}
func (m *WpaNetwork) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WpaNetwork.Unmarshal(m, b)
}
func (m *WpaNetwork) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WpaNetwork.Marshal(b, m, deterministic)
}
func (dst *WpaNetwork) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WpaNetwork.Merge(dst, src)
}
func (m *WpaNetwork) XXX_Size() int {
	return xxx_messageInfo_WpaNetwork.Size(m)
}
func (m *WpaNetwork) XXX_DiscardUnknown() {
	xxx_messageInfo_WpaNetwork.DiscardUnknown(m)
}

var xxx_messageInfo_WpaNetwork proto.InternalMessageInfo

func (m *WpaNetwork) GetBssid() string {
	if m != nil {
		return m.Bssid
	}
	return ""
}

func (m *WpaNetwork) GetFrequency() string {
	if m != nil {
		return m.Frequency
	}
	return ""
}

func (m *WpaNetwork) GetSignalLevel() string {
	if m != nil {
		return m.SignalLevel
	}
	return ""
}

func (m *WpaNetwork) GetFlags() string {
	if m != nil {
		return m.Flags
	}
	return ""
}

func (m *WpaNetwork) GetSsid() string {
	if m != nil {
		return m.Ssid
	}
	return ""
}

type WpaNetworkScanUpdate struct {
	// Types that are valid to be assigned to Update:
	//	*WpaNetworkScanUpdate_Appeared
	//	*WpaNetworkScanUpdate_Changed
	//	*WpaNetworkScanUpdate_Gone
	Update               isWpaNetworkScanUpdate_Update `protobuf_oneof:"update"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *WpaNetworkScanUpdate) Reset()         { *m = WpaNetworkScanUpdate{} }
func (m *WpaNetworkScanUpdate) String() string { return proto.CompactTextString(m) }
func (*WpaNetworkScanUpdate) ProtoMessage()    {}
func (*WpaNetworkScanUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{8}
}
func (m *WpaNetworkScanUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WpaNetworkScanUpdate.Unmarshal(m, b)
}
func (m *WpaNetworkScanUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WpaNetworkScanUpdate.Marshal(b, m, deterministic)
}
func (dst *WpaNetworkScanUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WpaNetworkScanUpdate.Merge(dst, src)
}
func (m *WpaNetworkScanUpdate) XXX_Size() int {
	return xxx_messageInfo_WpaNetworkScanUpdate.Size(m)
}
func (m *WpaNetworkScanUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_WpaNetworkScanUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_WpaNetworkScanUpdate proto.InternalMessageInfo

type isWpaNetworkScanUpdate_Update interface {
	isWpaNetworkScanUpdate_Update()
}

type WpaNetworkScanUpdate_Appeared struct {
	Appeared *WpaNetwork `protobuf:"bytes,1,opt,name=appeared,proto3,oneof"`
}

type WpaNetworkScanUpdate_Changed struct {
	Changed *WpaNetwork `protobuf:"bytes,2,opt,name=changed,proto3,oneof"`
}

type WpaNetworkScanUpdate_Gone struct {
	Gone *WpaNetwork `protobuf:"bytes,3,opt,name=gone,proto3,oneof"`
}

func (*WpaNetworkScanUpdate_Appeared) isWpaNetworkScanUpdate_Update() {}

func (*WpaNetworkScanUpdate_Changed) isWpaNetworkScanUpdate_Update() {}

func (*WpaNetworkScanUpdate_Gone) isWpaNetworkScanUpdate_Update() {}

func (m *WpaNetworkScanUpdate) GetUpdate() isWpaNetworkScanUpdate_Update {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *WpaNetworkScanUpdate) GetAppeared() *WpaNetwork {
	if x, ok := m.GetUpdate().(*WpaNetworkScanUpdate_Appeared); ok {
		return x.Appeared
	}
	return nil
}

func (m *WpaNetworkScanUpdate) GetChanged() *WpaNetwork {
	if x, ok := m.GetUpdate().(*WpaNetworkScanUpdate_Changed); ok {
		return x.Changed
	}
	return nil
}

func (m *WpaNetworkScanUpdate) GetGone() *WpaNetwork {
	if x, ok := m.GetUpdate().(*WpaNetworkScanUpdate_Gone); ok {
		return x.Gone
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*WpaNetworkScanUpdate) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _WpaNetworkScanUpdate_OneofMarshaler, _WpaNetworkScanUpdate_OneofUnmarshaler, _WpaNetworkScanUpdate_OneofSizer, []interface{}{
		(*WpaNetworkScanUpdate_Appeared)(nil),
		(*WpaNetworkScanUpdate_Changed)(nil),
		(*WpaNetworkScanUpdate_Gone)(nil),
	}
}

func _WpaNetworkScanUpdate_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*WpaNetworkScanUpdate)
	// update
	switch x := m.Update.(type) {
	case *WpaNetworkScanUpdate_Appeared:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Appeared); err != nil {
			return err
		}
	case *WpaNetworkScanUpdate_Changed:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Changed); err != nil {
			return err
		}
	case *WpaNetworkScanUpdate_Gone:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Gone); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("WpaNetworkScanUpdate.Update has unexpected type %T", x)
	}
	return nil
}

func _WpaNetworkScanUpdate_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*WpaNetworkScanUpdate)
	switch tag {
	case 1: // update.appeared
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(WpaNetwork)
		err := b.DecodeMessage(msg)
		m.Update = &WpaNetworkScanUpdate_Appeared{msg}
		return true, err
	case 2: // update.changed
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(WpaNetwork)
		err := b.DecodeMessage(msg)
		m.Update = &WpaNetworkScanUpdate_Changed{msg}
		return true, err
	case 3: // update.gone
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(WpaNetwork)
		err := b.DecodeMessage(msg)
		m.Update = &WpaNetworkScanUpdate_Gone{msg}
		return true, err
	default:
		return false, nil
	}
}

func _WpaNetworkScanUpdate_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*WpaNetworkScanUpdate)
	// update
	switch x := m.Update.(type) {
	case *WpaNetworkScanUpdate_Appeared:
		s := proto.Size(x.Appeared)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *WpaNetworkScanUpdate_Changed:
		s := proto.Size(x.Changed)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *WpaNetworkScanUpdate_Gone:
		s := proto.Size(x.Gone)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type UpdateRequest struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateRequest) Reset()         { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()    {}
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{9}
}
func (m *UpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRequest.Unmarshal(m, b)
}
func (m *UpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRequest.Marshal(b, m, deterministic)
}
func (dst *UpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRequest.Merge(dst, src)
}
func (m *UpdateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRequest.Size(m)
}
func (m *UpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRequest proto.InternalMessageInfo

func (m *UpdateRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type UpdateResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateResponse) Reset()         { *m = UpdateResponse{} }
func (m *UpdateResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateResponse) ProtoMessage()    {}
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{10}
}
func (m *UpdateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateResponse.Unmarshal(m, b)
}
func (m *UpdateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateResponse.Marshal(b, m, deterministic)
}
func (dst *UpdateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateResponse.Merge(dst, src)
}
func (m *UpdateResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateResponse.Size(m)
}
func (m *UpdateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetInfoRequest)(nil), "sweetrpc.GetInfoRequest")
	proto.RegisterType((*GetInfoResponse)(nil), "sweetrpc.GetInfoResponse")
	proto.RegisterType((*GetWpaConnectionInfoRequest)(nil), "sweetrpc.GetWpaConnectionInfoRequest")
	proto.RegisterType((*GetWpaConnectionInfoResponse)(nil), "sweetrpc.GetWpaConnectionInfoResponse")
	proto.RegisterType((*ConnectWpaNetworkRequest)(nil), "sweetrpc.ConnectWpaNetworkRequest")
	proto.RegisterType((*WpaConnectionUpdate)(nil), "sweetrpc.WpaConnectionUpdate")
	proto.RegisterType((*SubscribeWpaNetworkScanUpdatesRequest)(nil), "sweetrpc.SubscribeWpaNetworkScanUpdatesRequest")
	proto.RegisterType((*WpaNetwork)(nil), "sweetrpc.WpaNetwork")
	proto.RegisterType((*WpaNetworkScanUpdate)(nil), "sweetrpc.WpaNetworkScanUpdate")
	proto.RegisterType((*UpdateRequest)(nil), "sweetrpc.UpdateRequest")
	proto.RegisterType((*UpdateResponse)(nil), "sweetrpc.UpdateResponse")
	proto.RegisterEnum("sweetrpc.WpaConnectionUpdate_WpaConnectionUpdateState", WpaConnectionUpdate_WpaConnectionUpdateState_name, WpaConnectionUpdate_WpaConnectionUpdateState_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SweetClient is the client API for Sweet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SweetClient interface {
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error)
	GetWpaConnectionInfo(ctx context.Context, in *GetWpaConnectionInfoRequest, opts ...grpc.CallOption) (*GetWpaConnectionInfoResponse, error)
	ConnectWpaNetwork(ctx context.Context, in *ConnectWpaNetworkRequest, opts ...grpc.CallOption) (Sweet_ConnectWpaNetworkClient, error)
	SubscribeWpaNetworkScanUpdates(ctx context.Context, in *SubscribeWpaNetworkScanUpdatesRequest, opts ...grpc.CallOption) (Sweet_SubscribeWpaNetworkScanUpdatesClient, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
}

type sweetClient struct {
	cc *grpc.ClientConn
}

func NewSweetClient(cc *grpc.ClientConn) SweetClient {
	return &sweetClient{cc}
}

func (c *sweetClient) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error) {
	out := new(GetInfoResponse)
	err := c.cc.Invoke(ctx, "/sweetrpc.Sweet/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sweetClient) GetWpaConnectionInfo(ctx context.Context, in *GetWpaConnectionInfoRequest, opts ...grpc.CallOption) (*GetWpaConnectionInfoResponse, error) {
	out := new(GetWpaConnectionInfoResponse)
	err := c.cc.Invoke(ctx, "/sweetrpc.Sweet/GetWpaConnectionInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sweetClient) ConnectWpaNetwork(ctx context.Context, in *ConnectWpaNetworkRequest, opts ...grpc.CallOption) (Sweet_ConnectWpaNetworkClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Sweet_serviceDesc.Streams[0], "/sweetrpc.Sweet/ConnectWpaNetwork", opts...)
	if err != nil {
		return nil, err
	}
	x := &sweetConnectWpaNetworkClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Sweet_ConnectWpaNetworkClient interface {
	Recv() (*WpaConnectionUpdate, error)
	grpc.ClientStream
}

type sweetConnectWpaNetworkClient struct {
	grpc.ClientStream
}

func (x *sweetConnectWpaNetworkClient) Recv() (*WpaConnectionUpdate, error) {
	m := new(WpaConnectionUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sweetClient) SubscribeWpaNetworkScanUpdates(ctx context.Context, in *SubscribeWpaNetworkScanUpdatesRequest, opts ...grpc.CallOption) (Sweet_SubscribeWpaNetworkScanUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Sweet_serviceDesc.Streams[1], "/sweetrpc.Sweet/SubscribeWpaNetworkScanUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &sweetSubscribeWpaNetworkScanUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Sweet_SubscribeWpaNetworkScanUpdatesClient interface {
	Recv() (*WpaNetworkScanUpdate, error)
	grpc.ClientStream
}

type sweetSubscribeWpaNetworkScanUpdatesClient struct {
	grpc.ClientStream
}

func (x *sweetSubscribeWpaNetworkScanUpdatesClient) Recv() (*WpaNetworkScanUpdate, error) {
	m := new(WpaNetworkScanUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sweetClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/sweetrpc.Sweet/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SweetServer is the server API for Sweet service.
type SweetServer interface {
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	GetWpaConnectionInfo(context.Context, *GetWpaConnectionInfoRequest) (*GetWpaConnectionInfoResponse, error)
	ConnectWpaNetwork(*ConnectWpaNetworkRequest, Sweet_ConnectWpaNetworkServer) error
	SubscribeWpaNetworkScanUpdates(*SubscribeWpaNetworkScanUpdatesRequest, Sweet_SubscribeWpaNetworkScanUpdatesServer) error
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
}

func RegisterSweetServer(s *grpc.Server, srv SweetServer) {
	s.RegisterService(&_Sweet_serviceDesc, srv)
}

func _Sweet_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SweetServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sweetrpc.Sweet/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SweetServer).GetInfo(ctx, req.(*GetInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sweet_GetWpaConnectionInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWpaConnectionInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SweetServer).GetWpaConnectionInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sweetrpc.Sweet/GetWpaConnectionInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SweetServer).GetWpaConnectionInfo(ctx, req.(*GetWpaConnectionInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sweet_ConnectWpaNetwork_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectWpaNetworkRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SweetServer).ConnectWpaNetwork(m, &sweetConnectWpaNetworkServer{stream})
}

type Sweet_ConnectWpaNetworkServer interface {
	Send(*WpaConnectionUpdate) error
	grpc.ServerStream
}

type sweetConnectWpaNetworkServer struct {
	grpc.ServerStream
}

func (x *sweetConnectWpaNetworkServer) Send(m *WpaConnectionUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func _Sweet_SubscribeWpaNetworkScanUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeWpaNetworkScanUpdatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SweetServer).SubscribeWpaNetworkScanUpdates(m, &sweetSubscribeWpaNetworkScanUpdatesServer{stream})
}

type Sweet_SubscribeWpaNetworkScanUpdatesServer interface {
	Send(*WpaNetworkScanUpdate) error
	grpc.ServerStream
}

type sweetSubscribeWpaNetworkScanUpdatesServer struct {
	grpc.ServerStream
}

func (x *sweetSubscribeWpaNetworkScanUpdatesServer) Send(m *WpaNetworkScanUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func _Sweet_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SweetServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sweetrpc.Sweet/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SweetServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Sweet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sweetrpc.Sweet",
	HandlerType: (*SweetServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _Sweet_GetInfo_Handler,
		},
		{
			MethodName: "GetWpaConnectionInfo",
			Handler:    _Sweet_GetWpaConnectionInfo_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Sweet_Update_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConnectWpaNetwork",
			Handler:       _Sweet_ConnectWpaNetwork_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeWpaNetworkScanUpdates",
			Handler:       _Sweet_SubscribeWpaNetworkScanUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_77a6da22d6a3feb1) }

var fileDescriptor_77a6da22d6a3feb1 = []byte{
	// 581 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x8d, 0xf3, 0xe3, 0x24, 0xd3, 0xaf, 0xf9, 0xcc, 0x12, 0x81, 0x09, 0x6d, 0x45, 0x57, 0x2a,
	0x20, 0x2e, 0x42, 0x15, 0x24, 0xee, 0x90, 0x0a, 0x69, 0x08, 0x91, 0xaa, 0x20, 0x25, 0xa0, 0x4a,
	0x70, 0x81, 0x36, 0xce, 0xc4, 0x58, 0x4d, 0xec, 0x65, 0xd7, 0x6e, 0xc5, 0x43, 0xf0, 0x26, 0x5c,
	0xf1, 0x0c, 0x3c, 0x18, 0x5a, 0xef, 0x3a, 0x6e, 0x54, 0xa7, 0x70, 0xe7, 0x99, 0x39, 0x33, 0x73,
	0x76, 0x66, 0x8e, 0xa1, 0x29, 0xb8, 0xd7, 0xe5, 0x22, 0x8a, 0x23, 0xd2, 0x90, 0x57, 0x88, 0xb1,
	0xe0, 0x1e, 0x75, 0xa0, 0x35, 0xc4, 0x78, 0x14, 0x2e, 0xa2, 0x09, 0x7e, 0x4b, 0x50, 0xc6, 0xf4,
	0x33, 0xfc, 0xbf, 0xf6, 0x48, 0x1e, 0x85, 0x12, 0xc9, 0x3d, 0xb0, 0x25, 0x8a, 0x80, 0x2d, 0x5d,
	0xeb, 0x91, 0xf5, 0xb4, 0x39, 0x31, 0x16, 0x71, 0xa1, 0x7e, 0x89, 0x42, 0x06, 0x51, 0xe8, 0x96,
	0xd3, 0x40, 0x66, 0xaa, 0x0c, 0x2f, 0x5a, 0xad, 0x82, 0xd8, 0xad, 0xe8, 0x0c, 0x6d, 0xd1, 0x7d,
	0x78, 0x38, 0xc4, 0xf8, 0x9c, 0xb3, 0x7e, 0x14, 0x86, 0xe8, 0xc5, 0x41, 0x14, 0x5e, 0xef, 0x2d,
	0x60, 0xaf, 0x38, 0x6c, 0x88, 0x10, 0xa8, 0x4a, 0x19, 0xcc, 0x0d, 0x8d, 0xf4, 0x9b, 0xb4, 0xa1,
	0x26, 0x63, 0x16, 0xa3, 0xa1, 0xa0, 0x0d, 0xd2, 0x82, 0x72, 0xc0, 0x4d, 0xf3, 0x72, 0xc0, 0x15,
	0xd5, 0x15, 0x4a, 0xc9, 0x7c, 0x74, 0xab, 0x9a, 0xaa, 0x31, 0xe9, 0x09, 0xb8, 0xa6, 0xdb, 0x39,
	0x67, 0x63, 0x8c, 0xaf, 0x22, 0x71, 0x61, 0xf8, 0x14, 0xf6, 0x73, 0xa0, 0xc2, 0xe5, 0x85, 0xe9,
	0xa6, 0x3e, 0xe9, 0x4f, 0x0b, 0xee, 0x6e, 0x70, 0xfe, 0xc8, 0xe7, 0x8a, 0xc3, 0x18, 0x6c, 0x45,
	0x26, 0x91, 0x69, 0x7e, 0xab, 0xf7, 0xb2, 0x9b, 0x8d, 0xbd, 0x5b, 0x00, 0x2f, 0xf2, 0x4d, 0xd5,
	0x5b, 0x26, 0xa6, 0x0a, 0x1d, 0x80, 0xbb, 0x0d, 0x43, 0x5a, 0x00, 0xfd, 0xf7, 0xe3, 0xf1, 0xa0,
	0xff, 0x61, 0x34, 0x1e, 0x3a, 0x25, 0xb2, 0x0b, 0x4d, 0x63, 0x0f, 0x4e, 0x1d, 0x8b, 0x00, 0xd8,
	0x6f, 0x5f, 0x8f, 0xce, 0x06, 0xa7, 0x4e, 0x99, 0x3e, 0x81, 0xa3, 0x69, 0x32, 0x93, 0x9e, 0x08,
	0x66, 0x98, 0x3f, 0x79, 0xea, 0x31, 0x53, 0x50, 0x66, 0xdb, 0xf8, 0x61, 0x01, 0xe4, 0x00, 0x35,
	0xe8, 0xd9, 0xb5, 0x69, 0x68, 0x83, 0xec, 0x41, 0x73, 0x21, 0x54, 0x42, 0xe8, 0x7d, 0x37, 0x43,
	0xc9, 0x1d, 0xe4, 0x10, 0xfe, 0x93, 0x81, 0x1f, 0xb2, 0xe5, 0x97, 0x25, 0x5e, 0xe2, 0xd2, 0x2c,
	0x64, 0x47, 0xfb, 0xce, 0x94, 0x4b, 0x95, 0x5d, 0x2c, 0x99, 0x2f, 0xcd, 0x5e, 0xb4, 0xb1, 0x9e,
	0x7c, 0x2d, 0x9f, 0x3c, 0xfd, 0x65, 0x41, 0xbb, 0x88, 0x30, 0xe9, 0x41, 0x83, 0x71, 0x8e, 0x4c,
	0xa0, 0x26, 0xb7, 0xd3, 0x6b, 0x6f, 0x8c, 0xda, 0x64, 0xbc, 0x2b, 0x4d, 0xd6, 0x38, 0x72, 0x0c,
	0x75, 0xef, 0x2b, 0x0b, 0x7d, 0x9c, 0xa7, 0xac, 0xb7, 0xa7, 0x64, 0x30, 0xf2, 0x0c, 0xaa, 0x7e,
	0x14, 0x62, 0xfa, 0x86, 0xed, 0xf0, 0x14, 0xf3, 0xa6, 0x01, 0x76, 0x92, 0x72, 0xa3, 0x87, 0xb0,
	0xab, 0x59, 0x66, 0x37, 0xe5, 0x40, 0x25, 0x11, 0x99, 0x92, 0xd4, 0xa7, 0xd2, 0x60, 0x06, 0xd1,
	0x77, 0xde, 0xfb, 0x5d, 0x81, 0xda, 0x54, 0x95, 0x27, 0x27, 0x50, 0x37, 0x6a, 0x24, 0x6e, 0xde,
	0x71, 0x53, 0xb2, 0x9d, 0x07, 0x05, 0x11, 0x5d, 0x89, 0x96, 0x88, 0x0f, 0xed, 0x22, 0x4d, 0x91,
	0xa3, 0x8d, 0xa4, 0x6d, 0x92, 0xec, 0x3c, 0xfe, 0x1b, 0x6c, 0xdd, 0xe8, 0x13, 0xdc, 0xb9, 0x21,
	0x24, 0x42, 0xf3, 0xf4, 0x6d, 0x2a, 0xeb, 0xec, 0xdf, 0xaa, 0x0b, 0x5a, 0x3a, 0xb6, 0x48, 0x02,
	0x07, 0xb7, 0xdf, 0x2c, 0x79, 0x9e, 0x17, 0xf9, 0xa7, 0xeb, 0xee, 0x1c, 0x14, 0x2d, 0x30, 0xc7,
	0xa5, 0x6d, 0x5f, 0x81, 0x6d, 0x4e, 0xec, 0x7e, 0x8e, 0xde, 0x58, 0x67, 0xc7, 0xbd, 0x19, 0xc8,
	0x26, 0x32, 0xb3, 0xd3, 0xbf, 0xed, 0x8b, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8d, 0x89, 0xc3,
	0x41, 0x7a, 0x05, 0x00, 0x00,
}