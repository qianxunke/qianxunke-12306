// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

package task

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Task struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	UserId               string   `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	SeatTypes            string   `protobuf:"bytes,3,opt,name=seat_types,json=seatTypes,proto3" json:"seat_types,omitempty"`
	TrainDates           string   `protobuf:"bytes,4,opt,name=train_dates,json=trainDates,proto3" json:"train_dates,omitempty"`
	FindFrom             string   `protobuf:"bytes,5,opt,name=find_from,json=findFrom,proto3" json:"find_from,omitempty"`
	FindTo               string   `protobuf:"bytes,6,opt,name=find_to,json=findTo,proto3" json:"find_to,omitempty"`
	OkNo                 string   `protobuf:"bytes,10,opt,name=ok_no,json=okNo,proto3" json:"ok_no,omitempty"`
	Status               int64    `protobuf:"varint,11,opt,name=status,proto3" json:"status,omitempty"`
	CreatedTime          string   `protobuf:"bytes,12,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	Type                 string   `protobuf:"bytes,13,opt,name=type,proto3" json:"type,omitempty"`
	Trips                string   `protobuf:"bytes,14,opt,name=trips,proto3" json:"trips,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *Task) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Task) GetSeatTypes() string {
	if m != nil {
		return m.SeatTypes
	}
	return ""
}

func (m *Task) GetTrainDates() string {
	if m != nil {
		return m.TrainDates
	}
	return ""
}

func (m *Task) GetFindFrom() string {
	if m != nil {
		return m.FindFrom
	}
	return ""
}

func (m *Task) GetFindTo() string {
	if m != nil {
		return m.FindTo
	}
	return ""
}

func (m *Task) GetOkNo() string {
	if m != nil {
		return m.OkNo
	}
	return ""
}

func (m *Task) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Task) GetCreatedTime() string {
	if m != nil {
		return m.CreatedTime
	}
	return ""
}

func (m *Task) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Task) GetTrips() string {
	if m != nil {
		return m.Trips
	}
	return ""
}

type TaskPassenger struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TaskId               string   `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	IdNum                string   `protobuf:"bytes,4,opt,name=id_num,json=idNum,proto3" json:"id_num,omitempty"`
	TelNum               string   `protobuf:"bytes,5,opt,name=tel_num,json=telNum,proto3" json:"tel_num,omitempty"`
	Type                 string   `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	SeatNum              string   `protobuf:"bytes,7,opt,name=seat_num,json=seatNum,proto3" json:"seat_num,omitempty"`
	AllEncStr            string   `protobuf:"bytes,8,opt,name=allEncStr,proto3" json:"allEncStr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskPassenger) Reset()         { *m = TaskPassenger{} }
func (m *TaskPassenger) String() string { return proto.CompactTextString(m) }
func (*TaskPassenger) ProtoMessage()    {}
func (*TaskPassenger) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}

func (m *TaskPassenger) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskPassenger.Unmarshal(m, b)
}
func (m *TaskPassenger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskPassenger.Marshal(b, m, deterministic)
}
func (m *TaskPassenger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskPassenger.Merge(m, src)
}
func (m *TaskPassenger) XXX_Size() int {
	return xxx_messageInfo_TaskPassenger.Size(m)
}
func (m *TaskPassenger) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskPassenger.DiscardUnknown(m)
}

var xxx_messageInfo_TaskPassenger proto.InternalMessageInfo

func (m *TaskPassenger) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TaskPassenger) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *TaskPassenger) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TaskPassenger) GetIdNum() string {
	if m != nil {
		return m.IdNum
	}
	return ""
}

func (m *TaskPassenger) GetTelNum() string {
	if m != nil {
		return m.TelNum
	}
	return ""
}

func (m *TaskPassenger) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *TaskPassenger) GetSeatNum() string {
	if m != nil {
		return m.SeatNum
	}
	return ""
}

func (m *TaskPassenger) GetAllEncStr() string {
	if m != nil {
		return m.AllEncStr
	}
	return ""
}

type TaskDetails struct {
	Task                 *Task            `protobuf:"bytes,2,opt,name=task,proto3" json:"task,omitempty"`
	TaskPassenger        []*TaskPassenger `protobuf:"bytes,3,rep,name=task_passenger,json=taskPassenger,proto3" json:"task_passenger,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *TaskDetails) Reset()         { *m = TaskDetails{} }
func (m *TaskDetails) String() string { return proto.CompactTextString(m) }
func (*TaskDetails) ProtoMessage()    {}
func (*TaskDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{2}
}

func (m *TaskDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskDetails.Unmarshal(m, b)
}
func (m *TaskDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskDetails.Marshal(b, m, deterministic)
}
func (m *TaskDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskDetails.Merge(m, src)
}
func (m *TaskDetails) XXX_Size() int {
	return xxx_messageInfo_TaskDetails.Size(m)
}
func (m *TaskDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TaskDetails proto.InternalMessageInfo

func (m *TaskDetails) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

func (m *TaskDetails) GetTaskPassenger() []*TaskPassenger {
	if m != nil {
		return m.TaskPassenger
	}
	return nil
}

type Error struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{3}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type In_GetTaskInfo struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *In_GetTaskInfo) Reset()         { *m = In_GetTaskInfo{} }
func (m *In_GetTaskInfo) String() string { return proto.CompactTextString(m) }
func (*In_GetTaskInfo) ProtoMessage()    {}
func (*In_GetTaskInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{4}
}

func (m *In_GetTaskInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_In_GetTaskInfo.Unmarshal(m, b)
}
func (m *In_GetTaskInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_In_GetTaskInfo.Marshal(b, m, deterministic)
}
func (m *In_GetTaskInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_In_GetTaskInfo.Merge(m, src)
}
func (m *In_GetTaskInfo) XXX_Size() int {
	return xxx_messageInfo_In_GetTaskInfo.Size(m)
}
func (m *In_GetTaskInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_In_GetTaskInfo.DiscardUnknown(m)
}

var xxx_messageInfo_In_GetTaskInfo proto.InternalMessageInfo

func (m *In_GetTaskInfo) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

type Out_GetTaskInfo struct {
	Error                *Error       `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	TaskDetails          *TaskDetails `protobuf:"bytes,2,opt,name=taskDetails,proto3" json:"taskDetails,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Out_GetTaskInfo) Reset()         { *m = Out_GetTaskInfo{} }
func (m *Out_GetTaskInfo) String() string { return proto.CompactTextString(m) }
func (*Out_GetTaskInfo) ProtoMessage()    {}
func (*Out_GetTaskInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{5}
}

func (m *Out_GetTaskInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Out_GetTaskInfo.Unmarshal(m, b)
}
func (m *Out_GetTaskInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Out_GetTaskInfo.Marshal(b, m, deterministic)
}
func (m *Out_GetTaskInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Out_GetTaskInfo.Merge(m, src)
}
func (m *Out_GetTaskInfo) XXX_Size() int {
	return xxx_messageInfo_Out_GetTaskInfo.Size(m)
}
func (m *Out_GetTaskInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_Out_GetTaskInfo.DiscardUnknown(m)
}

var xxx_messageInfo_Out_GetTaskInfo proto.InternalMessageInfo

func (m *Out_GetTaskInfo) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *Out_GetTaskInfo) GetTaskDetails() *TaskDetails {
	if m != nil {
		return m.TaskDetails
	}
	return nil
}

type In_GetTaskInfoList struct {
	Limit                int64    `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Pages                int64    `protobuf:"varint,2,opt,name=pages,proto3" json:"pages,omitempty"`
	SearchKey            string   `protobuf:"bytes,3,opt,name=search_key,json=searchKey,proto3" json:"search_key,omitempty"`
	Status               string   `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	StartTime            string   `protobuf:"bytes,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime              string   `protobuf:"bytes,6,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *In_GetTaskInfoList) Reset()         { *m = In_GetTaskInfoList{} }
func (m *In_GetTaskInfoList) String() string { return proto.CompactTextString(m) }
func (*In_GetTaskInfoList) ProtoMessage()    {}
func (*In_GetTaskInfoList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{6}
}

func (m *In_GetTaskInfoList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_In_GetTaskInfoList.Unmarshal(m, b)
}
func (m *In_GetTaskInfoList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_In_GetTaskInfoList.Marshal(b, m, deterministic)
}
func (m *In_GetTaskInfoList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_In_GetTaskInfoList.Merge(m, src)
}
func (m *In_GetTaskInfoList) XXX_Size() int {
	return xxx_messageInfo_In_GetTaskInfoList.Size(m)
}
func (m *In_GetTaskInfoList) XXX_DiscardUnknown() {
	xxx_messageInfo_In_GetTaskInfoList.DiscardUnknown(m)
}

var xxx_messageInfo_In_GetTaskInfoList proto.InternalMessageInfo

func (m *In_GetTaskInfoList) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *In_GetTaskInfoList) GetPages() int64 {
	if m != nil {
		return m.Pages
	}
	return 0
}

func (m *In_GetTaskInfoList) GetSearchKey() string {
	if m != nil {
		return m.SearchKey
	}
	return ""
}

func (m *In_GetTaskInfoList) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *In_GetTaskInfoList) GetStartTime() string {
	if m != nil {
		return m.StartTime
	}
	return ""
}

func (m *In_GetTaskInfoList) GetEndTime() string {
	if m != nil {
		return m.EndTime
	}
	return ""
}

type Out_GetTaskInfoList struct {
	Error                *Error         `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Limit                int64          `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Pages                int64          `protobuf:"varint,3,opt,name=pages,proto3" json:"pages,omitempty"`
	Total                int64          `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
	TaskDetailsList      []*TaskDetails `protobuf:"bytes,5,rep,name=task_details_list,json=taskDetailsList,proto3" json:"task_details_list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Out_GetTaskInfoList) Reset()         { *m = Out_GetTaskInfoList{} }
func (m *Out_GetTaskInfoList) String() string { return proto.CompactTextString(m) }
func (*Out_GetTaskInfoList) ProtoMessage()    {}
func (*Out_GetTaskInfoList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{7}
}

func (m *Out_GetTaskInfoList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Out_GetTaskInfoList.Unmarshal(m, b)
}
func (m *Out_GetTaskInfoList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Out_GetTaskInfoList.Marshal(b, m, deterministic)
}
func (m *Out_GetTaskInfoList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Out_GetTaskInfoList.Merge(m, src)
}
func (m *Out_GetTaskInfoList) XXX_Size() int {
	return xxx_messageInfo_Out_GetTaskInfoList.Size(m)
}
func (m *Out_GetTaskInfoList) XXX_DiscardUnknown() {
	xxx_messageInfo_Out_GetTaskInfoList.DiscardUnknown(m)
}

var xxx_messageInfo_Out_GetTaskInfoList proto.InternalMessageInfo

func (m *Out_GetTaskInfoList) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *Out_GetTaskInfoList) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *Out_GetTaskInfoList) GetPages() int64 {
	if m != nil {
		return m.Pages
	}
	return 0
}

func (m *Out_GetTaskInfoList) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *Out_GetTaskInfoList) GetTaskDetailsList() []*TaskDetails {
	if m != nil {
		return m.TaskDetailsList
	}
	return nil
}

type In_UpdateTaskInfo struct {
	TaskDetails          *TaskDetails `protobuf:"bytes,1,opt,name=taskDetails,proto3" json:"taskDetails,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *In_UpdateTaskInfo) Reset()         { *m = In_UpdateTaskInfo{} }
func (m *In_UpdateTaskInfo) String() string { return proto.CompactTextString(m) }
func (*In_UpdateTaskInfo) ProtoMessage()    {}
func (*In_UpdateTaskInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{8}
}

func (m *In_UpdateTaskInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_In_UpdateTaskInfo.Unmarshal(m, b)
}
func (m *In_UpdateTaskInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_In_UpdateTaskInfo.Marshal(b, m, deterministic)
}
func (m *In_UpdateTaskInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_In_UpdateTaskInfo.Merge(m, src)
}
func (m *In_UpdateTaskInfo) XXX_Size() int {
	return xxx_messageInfo_In_UpdateTaskInfo.Size(m)
}
func (m *In_UpdateTaskInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_In_UpdateTaskInfo.DiscardUnknown(m)
}

var xxx_messageInfo_In_UpdateTaskInfo proto.InternalMessageInfo

func (m *In_UpdateTaskInfo) GetTaskDetails() *TaskDetails {
	if m != nil {
		return m.TaskDetails
	}
	return nil
}

type Out_UpdateTaskInfo struct {
	Error                *Error   `protobuf:"bytes,1,opt,name=Error,proto3" json:"Error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Out_UpdateTaskInfo) Reset()         { *m = Out_UpdateTaskInfo{} }
func (m *Out_UpdateTaskInfo) String() string { return proto.CompactTextString(m) }
func (*Out_UpdateTaskInfo) ProtoMessage()    {}
func (*Out_UpdateTaskInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{9}
}

func (m *Out_UpdateTaskInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Out_UpdateTaskInfo.Unmarshal(m, b)
}
func (m *Out_UpdateTaskInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Out_UpdateTaskInfo.Marshal(b, m, deterministic)
}
func (m *Out_UpdateTaskInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Out_UpdateTaskInfo.Merge(m, src)
}
func (m *Out_UpdateTaskInfo) XXX_Size() int {
	return xxx_messageInfo_Out_UpdateTaskInfo.Size(m)
}
func (m *Out_UpdateTaskInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_Out_UpdateTaskInfo.DiscardUnknown(m)
}

var xxx_messageInfo_Out_UpdateTaskInfo proto.InternalMessageInfo

func (m *Out_UpdateTaskInfo) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type In_AddTask struct {
	TaskDetails          *TaskDetails `protobuf:"bytes,1,opt,name=taskDetails,proto3" json:"taskDetails,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *In_AddTask) Reset()         { *m = In_AddTask{} }
func (m *In_AddTask) String() string { return proto.CompactTextString(m) }
func (*In_AddTask) ProtoMessage()    {}
func (*In_AddTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{10}
}

func (m *In_AddTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_In_AddTask.Unmarshal(m, b)
}
func (m *In_AddTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_In_AddTask.Marshal(b, m, deterministic)
}
func (m *In_AddTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_In_AddTask.Merge(m, src)
}
func (m *In_AddTask) XXX_Size() int {
	return xxx_messageInfo_In_AddTask.Size(m)
}
func (m *In_AddTask) XXX_DiscardUnknown() {
	xxx_messageInfo_In_AddTask.DiscardUnknown(m)
}

var xxx_messageInfo_In_AddTask proto.InternalMessageInfo

func (m *In_AddTask) GetTaskDetails() *TaskDetails {
	if m != nil {
		return m.TaskDetails
	}
	return nil
}

type Out_AddTask struct {
	Error                *Error       `protobuf:"bytes,1,opt,name=Error,proto3" json:"Error,omitempty"`
	TaskDetails          *TaskDetails `protobuf:"bytes,2,opt,name=taskDetails,proto3" json:"taskDetails,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Out_AddTask) Reset()         { *m = Out_AddTask{} }
func (m *Out_AddTask) String() string { return proto.CompactTextString(m) }
func (*Out_AddTask) ProtoMessage()    {}
func (*Out_AddTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{11}
}

func (m *Out_AddTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Out_AddTask.Unmarshal(m, b)
}
func (m *Out_AddTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Out_AddTask.Marshal(b, m, deterministic)
}
func (m *Out_AddTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Out_AddTask.Merge(m, src)
}
func (m *Out_AddTask) XXX_Size() int {
	return xxx_messageInfo_Out_AddTask.Size(m)
}
func (m *Out_AddTask) XXX_DiscardUnknown() {
	xxx_messageInfo_Out_AddTask.DiscardUnknown(m)
}

var xxx_messageInfo_Out_AddTask proto.InternalMessageInfo

func (m *Out_AddTask) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *Out_AddTask) GetTaskDetails() *TaskDetails {
	if m != nil {
		return m.TaskDetails
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "task.Task")
	proto.RegisterType((*TaskPassenger)(nil), "task.Task_passenger")
	proto.RegisterType((*TaskDetails)(nil), "task.TaskDetails")
	proto.RegisterType((*Error)(nil), "task.Error")
	proto.RegisterType((*In_GetTaskInfo)(nil), "task.In_GetTaskInfo")
	proto.RegisterType((*Out_GetTaskInfo)(nil), "task.Out_GetTaskInfo")
	proto.RegisterType((*In_GetTaskInfoList)(nil), "task.In_GetTaskInfoList")
	proto.RegisterType((*Out_GetTaskInfoList)(nil), "task.Out_GetTaskInfoList")
	proto.RegisterType((*In_UpdateTaskInfo)(nil), "task.In_UpdateTaskInfo")
	proto.RegisterType((*Out_UpdateTaskInfo)(nil), "task.Out_UpdateTaskInfo")
	proto.RegisterType((*In_AddTask)(nil), "task.In_AddTask")
	proto.RegisterType((*Out_AddTask)(nil), "task.Out_AddTask")
}

func init() { proto.RegisterFile("task.proto", fileDescriptor_ce5d8dd45b4a91ff) }

var fileDescriptor_ce5d8dd45b4a91ff = []byte{
	// 754 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xcb, 0x6e, 0x13, 0x4b,
	0x10, 0xb5, 0x3d, 0x1e, 0x3f, 0x6a, 0x12, 0xe7, 0xba, 0x93, 0xdc, 0x74, 0x7c, 0x2f, 0x90, 0xcc,
	0x2a, 0x6c, 0x22, 0xe4, 0x08, 0xb1, 0x00, 0x16, 0x91, 0x12, 0x88, 0x05, 0x0a, 0x68, 0x62, 0xd6,
	0xa3, 0xc1, 0xd3, 0x49, 0x1a, 0xcf, 0xc3, 0xea, 0x6e, 0x23, 0xe5, 0x27, 0xf8, 0x0e, 0xb6, 0x2c,
	0xf9, 0x05, 0xbe, 0x0a, 0x55, 0xf5, 0xd8, 0x1e, 0x5b, 0x46, 0x44, 0x59, 0xb9, 0xeb, 0x9c, 0xee,
	0xea, 0xaa, 0x53, 0xa7, 0xc7, 0x00, 0x26, 0xd2, 0xe3, 0xe3, 0x89, 0xca, 0x4d, 0xce, 0xea, 0xb8,
	0xf6, 0xbf, 0xd7, 0xa0, 0x3e, 0x8c, 0xf4, 0x98, 0xed, 0x41, 0x13, 0x81, 0x50, 0xc6, 0xbc, 0x7a,
	0x50, 0x3d, 0x6a, 0x07, 0x0d, 0x0c, 0x07, 0x31, 0x12, 0x53, 0x2d, 0x14, 0x12, 0x35, 0x4b, 0x60,
	0x38, 0x88, 0xd9, 0x23, 0x00, 0x2d, 0x22, 0x13, 0x9a, 0xbb, 0x89, 0xd0, 0xdc, 0x21, 0xae, 0x8d,
	0xc8, 0x10, 0x01, 0xf6, 0x04, 0x3c, 0xa3, 0x22, 0x99, 0x85, 0x71, 0x64, 0x84, 0xe6, 0x75, 0xe2,
	0x81, 0xa0, 0x33, 0x44, 0xd8, 0x7f, 0xd0, 0xbe, 0x96, 0x59, 0x1c, 0x5e, 0xab, 0x3c, 0xe5, 0x2e,
	0xd1, 0x2d, 0x04, 0xde, 0xa8, 0x3c, 0xc5, 0x5b, 0x89, 0x34, 0x39, 0x6f, 0xd8, 0x5b, 0x31, 0x1c,
	0xe6, 0x6c, 0x1b, 0xdc, 0x7c, 0x1c, 0x66, 0x39, 0x07, 0x82, 0xeb, 0xf9, 0xf8, 0x32, 0x67, 0xff,
	0x42, 0x43, 0x9b, 0xc8, 0x4c, 0x35, 0xf7, 0x0e, 0xaa, 0x47, 0x4e, 0x50, 0x44, 0xec, 0x10, 0x36,
	0x46, 0x4a, 0x44, 0x46, 0xc4, 0xa1, 0x91, 0xa9, 0xe0, 0x1b, 0x74, 0xc6, 0x2b, 0xb0, 0xa1, 0x4c,
	0x05, 0x63, 0x50, 0xc7, 0x06, 0xf8, 0xa6, 0x4d, 0x87, 0x6b, 0xb6, 0x03, 0xae, 0x51, 0x72, 0xa2,
	0x79, 0x87, 0x40, 0x1b, 0xf8, 0xbf, 0xaa, 0xd0, 0x41, 0xa9, 0xc2, 0x49, 0xa4, 0xb5, 0xc8, 0x6e,
	0x84, 0x62, 0x1d, 0xa8, 0xcd, 0xf5, 0xaa, 0xc9, 0xb8, 0x2c, 0x62, 0x6d, 0x49, 0x44, 0x06, 0xf5,
	0x2c, 0x4a, 0x45, 0xa1, 0x12, 0xad, 0xd9, 0x2e, 0x34, 0x64, 0x1c, 0x66, 0xd3, 0xb4, 0xd0, 0xc6,
	0x95, 0xf1, 0xe5, 0x94, 0x3a, 0x37, 0x22, 0x21, 0xdc, 0x2d, 0x72, 0x88, 0x04, 0x89, 0x59, 0xa5,
	0x8d, 0x52, 0xa5, 0xfb, 0xd0, 0xa2, 0x19, 0xe0, 0xee, 0x26, 0xe1, 0x4d, 0x8c, 0x71, 0xfb, 0xff,
	0xd0, 0x8e, 0x92, 0xe4, 0x3c, 0x1b, 0x5d, 0x19, 0xc5, 0x5b, 0x76, 0x3a, 0x73, 0xc0, 0xff, 0x02,
	0x1e, 0xf6, 0x72, 0x26, 0x4c, 0x24, 0x13, 0xcd, 0x1e, 0x03, 0xd9, 0x81, 0xaa, 0xf6, 0xfa, 0x70,
	0x4c, 0x3e, 0xc1, 0x0d, 0x01, 0xe1, 0xec, 0x25, 0x74, 0xcc, 0x52, 0xeb, 0xdc, 0x39, 0x70, 0x8e,
	0xbc, 0xfe, 0xce, 0x62, 0xe7, 0x82, 0x0b, 0x36, 0x11, 0xfc, 0x38, 0x0b, 0xfd, 0xe7, 0xe0, 0x9e,
	0x2b, 0x95, 0x2b, 0xec, 0x60, 0x94, 0xc7, 0x82, 0x04, 0x73, 0x03, 0x5a, 0x33, 0x0e, 0xcd, 0x54,
	0x68, 0x1d, 0xdd, 0x88, 0x42, 0xb2, 0x59, 0xe8, 0x3f, 0x85, 0xce, 0x20, 0x0b, 0xdf, 0x0a, 0x83,
	0xd9, 0x07, 0xd9, 0x75, 0xfe, 0x47, 0x8f, 0xfa, 0x12, 0xb6, 0x3e, 0x4c, 0xcd, 0xd2, 0xde, 0x43,
	0x70, 0x05, 0x5e, 0x4a, 0x3b, 0xbd, 0xbe, 0x67, 0x0b, 0xa5, 0x3a, 0x02, 0xcb, 0xb0, 0x13, 0xf0,
	0xcc, 0x42, 0x83, 0xa2, 0xf7, 0xee, 0xa2, 0xa3, 0x82, 0x08, 0xca, 0xbb, 0xfc, 0x1f, 0x55, 0x60,
	0xcb, 0x65, 0xbd, 0x97, 0xda, 0xa0, 0x65, 0x12, 0x99, 0x4a, 0x43, 0xd7, 0x39, 0x81, 0x0d, 0x10,
	0x9d, 0x44, 0x37, 0xc2, 0xe6, 0x76, 0x02, 0x1b, 0x14, 0x0f, 0x47, 0x8d, 0x6e, 0xc3, 0xb1, 0xb8,
	0x2b, 0x3d, 0x1c, 0x35, 0xba, 0x7d, 0x27, 0xee, 0x4a, 0x66, 0xb6, 0xbe, 0x98, 0x99, 0x19, 0x8f,
	0x99, 0x48, 0x19, 0x6b, 0x65, 0xb7, 0x38, 0x86, 0x08, 0x19, 0x79, 0x1f, 0x5a, 0x22, 0x2b, 0x7c,
	0x6e, 0x2d, 0xd2, 0x14, 0x19, 0x79, 0xdc, 0xff, 0x59, 0x85, 0xed, 0x15, 0x7d, 0xa8, 0xe8, 0x7b,
	0x68, 0x34, 0xef, 0xab, 0xb6, 0xb6, 0x2f, 0xa7, 0xdc, 0x17, 0x3e, 0x9b, 0xdc, 0x44, 0x09, 0xd5,
	0xed, 0x04, 0x36, 0x60, 0xaf, 0xa1, 0x4b, 0x43, 0x8b, 0xad, 0x80, 0x61, 0x22, 0xb5, 0xe1, 0x2e,
	0xb9, 0x67, 0x8d, 0xd6, 0x5b, 0x25, 0xad, 0xb1, 0x46, 0xff, 0x02, 0xba, 0x83, 0x2c, 0xfc, 0x34,
	0xc1, 0xaf, 0xc8, 0x7c, 0xb8, 0x2b, 0x93, 0xab, 0xde, 0x6b, 0x72, 0x2f, 0x80, 0xa1, 0x08, 0x2b,
	0xa9, 0x0e, 0x0b, 0x73, 0xae, 0xd5, 0x80, 0x7e, 0xfc, 0x53, 0x80, 0x41, 0x16, 0x9e, 0xc6, 0x31,
	0x7d, 0x28, 0x1f, 0x74, 0xb7, 0x00, 0x0f, 0xef, 0x9e, 0xe5, 0xf8, 0xfb, 0xa5, 0x0f, 0x32, 0x67,
	0xff, 0x5b, 0xcd, 0x3e, 0xeb, 0x2b, 0xa1, 0xbe, 0xca, 0x91, 0x60, 0xaf, 0xc0, 0x2b, 0xbf, 0x89,
	0xe2, 0xb5, 0x2e, 0xdb, 0xb7, 0xb7, 0x6b, 0xd1, 0x15, 0x83, 0xf8, 0x15, 0x76, 0x01, 0x5b, 0xab,
	0x8e, 0xe1, 0xeb, 0x32, 0x20, 0xd3, 0xdb, 0x5f, 0x9b, 0x85, 0x46, 0x58, 0x61, 0xe7, 0xd0, 0x59,
	0x91, 0x7d, 0x6f, 0x9e, 0x68, 0x99, 0xe8, 0xf1, 0x45, 0x9e, 0x65, 0xc6, 0xaf, 0xb0, 0x67, 0xd0,
	0x9c, 0x29, 0xf8, 0xcf, 0xfc, 0x7c, 0x81, 0xf4, 0xba, 0x8b, 0x83, 0x05, 0xe4, 0x57, 0x3e, 0x37,
	0xe8, 0xbf, 0xee, 0xe4, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb4, 0x66, 0x44, 0xce, 0xf9, 0x06,
	0x00, 0x00,
}
