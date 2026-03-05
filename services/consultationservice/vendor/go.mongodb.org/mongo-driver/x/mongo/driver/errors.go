





package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)




const LegacyNotPrimaryErrMsg = "not master"

var (
	retryableCodes          = []int32{11600, 11602, 10107, 13435, 13436, 189, 91, 7, 6, 89, 9001, 262}
	nodeIsRecoveringCodes   = []int32{11600, 11602, 13436, 189, 91}
	notPrimaryCodes         = []int32{10107, 13435, 10058}
	nodeIsShuttingDownCodes = []int32{11600, 91}

	unknownReplWriteConcernCode   = int32(79)
	unsatisfiableWriteConcernCode = int32(100)
)

var (
	
	UnknownTransactionCommitResult = "UnknownTransactionCommitResult"
	
	TransientTransactionError = "TransientTransactionError"
	
	NetworkError = "NetworkError"
	
	RetryableWriteError = "RetryableWriteError"
	
	NoWritesPerformed = "NoWritesPerformed"
	
	ErrCursorNotFound = errors.New("cursor not found")
	
	
	ErrUnacknowledgedWrite = errors.New("unacknowledged write")
	
	
	ErrUnsupportedStorageEngine = errors.New("this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string")
	
	
	
	ErrDeadlineWouldBeExceeded = fmt.Errorf(
		"operation not sent to server, as Timeout would be exceeded: %w",
		context.DeadlineExceeded)
	
	ErrNegativeMaxTime = errors.New("a negative value was provided for MaxTime on an operation")
)


type QueryFailureError struct {
	Message  string
	Response bsoncore.Document
	Wrapped  error
}


func (e QueryFailureError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Response)
}


func (e QueryFailureError) Unwrap() error {
	return e.Wrapped
}


type ResponseError struct {
	Message string
	Wrapped error
}


func NewCommandResponseError(msg string, err error) ResponseError {
	return ResponseError{Message: msg, Wrapped: err}
}


func (e ResponseError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Wrapped)
	}
	return e.Message
}


type WriteCommandError struct {
	WriteConcernError *WriteConcernError
	WriteErrors       WriteErrors
	Labels            []string
	Raw               bsoncore.Document
}



func (wce WriteCommandError) UnsupportedStorageEngine() bool {
	for _, writeError := range wce.WriteErrors {
		if writeError.Code == 20 && strings.HasPrefix(strings.ToLower(writeError.Message), "transaction numbers") {
			return true
		}
	}
	return false
}

func (wce WriteCommandError) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write command error: [")
	fmt.Fprintf(&buf, "{%s}, ", wce.WriteErrors)
	fmt.Fprintf(&buf, "{%s}]", wce.WriteConcernError)
	return buf.String()
}


func (wce WriteCommandError) Retryable(wireVersion *description.VersionRange) bool {
	for _, label := range wce.Labels {
		if label == RetryableWriteError {
			return true
		}
	}
	if wireVersion != nil && wireVersion.Max >= 9 {
		return false
	}

	if wce.WriteConcernError == nil {
		return false
	}
	return wce.WriteConcernError.Retryable()
}


func (wce WriteCommandError) HasErrorLabel(label string) bool {
	if wce.Labels != nil {
		for _, l := range wce.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}



type WriteConcernError struct {
	Name            string
	Code            int64
	Message         string
	Details         bsoncore.Document
	Labels          []string
	TopologyVersion *description.TopologyVersion
	Raw             bsoncore.Document
}

func (wce WriteConcernError) Error() string {
	if wce.Name != "" {
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	}
	return wce.Message
}


func (wce WriteConcernError) Retryable() bool {
	for _, code := range retryableCodes {
		if wce.Code == int64(code) {
			return true
		}
	}

	return false
}


func (wce WriteConcernError) NodeIsRecovering() bool {
	for _, code := range nodeIsRecoveringCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, "node is recovering")
}


func (wce WriteConcernError) NodeIsShuttingDown() bool {
	for _, code := range nodeIsShuttingDownCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, "node is shutting down")
}


func (wce WriteConcernError) NotPrimary() bool {
	for _, code := range notPrimaryCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, LegacyNotPrimaryErrMsg)
}



type WriteError struct {
	Index   int64
	Code    int64
	Message string
	Details bsoncore.Document
	Raw     bsoncore.Document
}

func (we WriteError) Error() string { return we.Message }



type WriteErrors []WriteError

func (we WriteErrors) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write errors: [")
	for idx, err := range we {
		if idx != 0 {
			fmt.Fprintf(&buf, ", ")
		}
		fmt.Fprintf(&buf, "{%s}", err)
	}
	fmt.Fprint(&buf, "]")
	return buf.String()
}


type Error struct {
	Code            int32
	Message         string
	Labels          []string
	Name            string
	Wrapped         error
	TopologyVersion *description.TopologyVersion
	Raw             bsoncore.Document
}


func (e Error) UnsupportedStorageEngine() bool {
	return e.Code == 20 && strings.HasPrefix(strings.ToLower(e.Message), "transaction numbers")
}


func (e Error) Error() string {
	var msg string
	if e.Name != "" {
		msg = fmt.Sprintf("(%v)", e.Name)
	}
	msg += " " + e.Message
	if e.Wrapped != nil {
		msg += ": " + e.Wrapped.Error()
	}
	return msg
}


func (e Error) Unwrap() error {
	return e.Wrapped
}


func (e Error) HasErrorLabel(label string) bool {
	if e.Labels != nil {
		for _, l := range e.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}


func (e Error) RetryableRead() bool {
	for _, label := range e.Labels {
		if label == NetworkError {
			return true
		}
	}
	for _, code := range retryableCodes {
		if e.Code == code {
			return true
		}
	}

	return false
}


func (e Error) RetryableWrite(wireVersion *description.VersionRange) bool {
	for _, label := range e.Labels {
		if label == NetworkError || label == RetryableWriteError {
			return true
		}
	}
	if wireVersion != nil && wireVersion.Max >= 9 {
		return false
	}
	for _, code := range retryableCodes {
		if e.Code == code {
			return true
		}
	}

	return false
}


func (e Error) NetworkError() bool {
	for _, label := range e.Labels {
		if label == NetworkError {
			return true
		}
	}
	return false
}


func (e Error) NodeIsRecovering() bool {
	for _, code := range nodeIsRecoveringCodes {
		if e.Code == code {
			return true
		}
	}
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, "node is recovering")
}


func (e Error) NodeIsShuttingDown() bool {
	for _, code := range nodeIsShuttingDownCodes {
		if e.Code == code {
			return true
		}
	}
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, "node is shutting down")
}


func (e Error) NotPrimary() bool {
	for _, code := range notPrimaryCodes {
		if e.Code == code {
			return true
		}
	}
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, LegacyNotPrimaryErrMsg)
}


func (e Error) NamespaceNotFound() bool {
	return e.Code == 26 || e.Message == "ns not found"
}



func ExtractErrorFromServerResponse(ctx context.Context, doc bsoncore.Document) error {
	var errmsg, codeName string
	var code int32
	var labels []string
	var ok bool
	var tv *description.TopologyVersion
	var wcError WriteCommandError
	elems, err := doc.Elements()
	if err != nil {
		return err
	}

	for _, elem := range elems {
		switch elem.Key() {
		case "ok":
			switch elem.Value().Type {
			case bson.TypeInt32:
				if elem.Value().Int32() == 1 {
					ok = true
				}
			case bson.TypeInt64:
				if elem.Value().Int64() == 1 {
					ok = true
				}
			case bson.TypeDouble:
				if elem.Value().Double() == 1 {
					ok = true
				}
			case bson.TypeBoolean:
				if elem.Value().Boolean() {
					ok = true
				}
			}
		case "errmsg":
			if str, okay := elem.Value().StringValueOK(); okay {
				errmsg = str
			}
		case "codeName":
			if str, okay := elem.Value().StringValueOK(); okay {
				codeName = str
			}
		case "code":
			if c, okay := elem.Value().Int32OK(); okay {
				code = c
			}
		case "errorLabels":
			if arr, okay := elem.Value().ArrayOK(); okay {
				vals, err := arr.Values()
				if err != nil {
					continue
				}
				for _, val := range vals {
					if str, ok := val.StringValueOK(); ok {
						labels = append(labels, str)
					}
				}

			}
		case "writeErrors":
			arr, exists := elem.Value().ArrayOK()
			if !exists {
				break
			}
			vals, err := arr.Values()
			if err != nil {
				continue
			}
			for _, val := range vals {
				var we WriteError
				doc, exists := val.DocumentOK()
				if !exists {
					continue
				}
				if index, exists := doc.Lookup("index").AsInt64OK(); exists {
					we.Index = index
				}
				if code, exists := doc.Lookup("code").AsInt64OK(); exists {
					we.Code = code
				}
				if msg, exists := doc.Lookup("errmsg").StringValueOK(); exists {
					we.Message = msg
				}
				if info, exists := doc.Lookup("errInfo").DocumentOK(); exists {
					we.Details = make([]byte, len(info))
					copy(we.Details, info)
				}
				we.Raw = doc
				wcError.WriteErrors = append(wcError.WriteErrors, we)
			}
		case "writeConcernError":
			doc, exists := elem.Value().DocumentOK()
			if !exists {
				break
			}
			wcError.WriteConcernError = new(WriteConcernError)
			wcError.WriteConcernError.Raw = doc
			if code, exists := doc.Lookup("code").AsInt64OK(); exists {
				wcError.WriteConcernError.Code = code
			}
			if name, exists := doc.Lookup("codeName").StringValueOK(); exists {
				wcError.WriteConcernError.Name = name
			}
			if msg, exists := doc.Lookup("errmsg").StringValueOK(); exists {
				wcError.WriteConcernError.Message = msg
			}
			if info, exists := doc.Lookup("errInfo").DocumentOK(); exists {
				wcError.WriteConcernError.Details = make([]byte, len(info))
				copy(wcError.WriteConcernError.Details, info)
			}
			if errLabels, exists := doc.Lookup("errorLabels").ArrayOK(); exists {
				vals, err := errLabels.Values()
				if err != nil {
					continue
				}
				for _, val := range vals {
					if str, ok := val.StringValueOK(); ok {
						labels = append(labels, str)
					}
				}
			}
		case "topologyVersion":
			doc, ok := elem.Value().DocumentOK()
			if !ok {
				break
			}
			version, err := description.NewTopologyVersion(bson.Raw(doc))
			if err == nil {
				tv = version
			}
		}
	}

	if !ok {
		if errmsg == "" {
			errmsg = "command failed"
		}

		err := Error{
			Code:            code,
			Message:         errmsg,
			Name:            codeName,
			Labels:          labels,
			TopologyVersion: tv,
			Raw:             doc,
		}

		
		
		
		
		
		
		
		
		if csot.IsTimeoutContext(ctx) && err.Code == 50 {
			err.Wrapped = context.DeadlineExceeded
		}

		return err
	}

	if len(wcError.WriteErrors) > 0 || wcError.WriteConcernError != nil {
		wcError.Labels = labels
		if wcError.WriteConcernError != nil {
			wcError.WriteConcernError.TopologyVersion = tv
		}
		wcError.Raw = doc
		return wcError
	}

	return nil
}
