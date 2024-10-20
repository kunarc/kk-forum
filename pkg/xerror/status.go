package xerror

import (
	"context"
	"strconv"

	"pkg/xerror/pb"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// StatusFromError
//
//	生成 status by err
//	@return pb.Status
//	@author kunarc
//	@update 2024-10-20 06:49:02

// type Status struct {
// 	st *pb.Status
// }

// func (s *Status) Code() int32 {
// 	return s.st.Code
// }

// func (s *Status) Message() string {
// 	return s.st.Message
// }

// StatusFromError
//
//	生成带有自定义 status 的 status
//	@param err error
//	@return *status.Status
//	@author kunarc
//	@update 2024-10-20 06:57:07
func GrpcStatusFromError(err error) *status.Status {
	err = errors.Cause(err)
	if ex, ok := err.(XError); ok {
		grpcStatus, e := XError2GrpcStatus(ex)
		if e != nil {
			return status.New(codes.Internal, err.Error())
		}
		return grpcStatus
	}
	var grpcStatus *status.Status
	switch err {
	case context.Canceled:
		grpcStatus, _ = XError2GrpcStatus(Canceled)
	case context.DeadlineExceeded:
		grpcStatus, _ = XError2GrpcStatus(Deadline)
	default: // 不为业务错误，例如序列化出错
		grpcStatus = status.New(codes.Internal, err.Error())
	}
	return grpcStatus
}

// XError2GrpcStatus
//
//	@param err XError
//	@return *status.Status
//	@return error
//	@author kunarc
//	@update 2024-10-20 07:23:16
func XError2GrpcStatus(err XError) (*status.Status, error) {
	// 暂时不考虑错误为status， 这里默认err为xcode
	st := &pb.Status{
		Code:    int32(err.Code()),
		Message: err.Message(),
	}
	stat := status.New(codes.Unknown, strconv.Itoa(int(st.Code)))
	return stat.WithDetails(st)
}

// GrpcStatus2XError
//
//	@param status *status.Status
//	@return XError
//	@author kunarc
//	@update 2024-10-20 10:17:55
func GrpcStatus2XError(status *status.Status) XError {
	details := status.Details()
	for i := len(details) - 1; i >= 0; i-- {
		if pb, ok := details[i].(proto.Message); ok {
			return fromProto(pb)
		}
	}
	// 兜底处理，防止服务端传来未带有业务错误码的status
	return toXError(status)
}

// fromProto 转换成业务错误码
//
//	@param msg proto.Message
//	@return XError
//	@author kunarc
//	@update 2024-10-20 10:01:55
func fromProto(msg proto.Message) XError {
	if st, ok := msg.(*pb.Status); ok {
		return fromStatus(st)
	}
	return Errorf(ServerErr, "invalid proto message get %v", msg)
}

// toXError
//
//	@param status *status.Status
//	@return xe XError
//	@author kunarc
//	@update 2024-10-20 09:58:18
func toXError(status *status.Status) (xe XError) {
	grpcCode := status.Code()
	switch grpcCode {
	case codes.OK:
		xe = OK
	case codes.Canceled:
		xe = Canceled
	case codes.DeadlineExceeded:
		xe = Deadline
	default:
		xe = Errorf(ServerErr, "rpc error code = %s desc = %s", status.Code(), status.Message())
	}
	return
}
