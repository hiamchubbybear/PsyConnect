





package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/codecutil"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)


var ErrUnacknowledgedWrite = errors.New("unacknowledged write")


var ErrClientDisconnected = errors.New("client is disconnected")


var ErrNilDocument = errors.New("document is nil")


var ErrNilValue = errors.New("value is nil")


var ErrEmptySlice = errors.New("must provide at least one element in input slice")


type ErrMapForOrderedArgument struct {
	ParamName string
}


func (e ErrMapForOrderedArgument) Error() string {
	return fmt.Sprintf("multi-key map passed in for ordered parameter %v", e.ParamName)
}

func replaceErrors(err error) error {
	
	if err == nil {
		return nil
	}

	if errors.Is(err, topology.ErrTopologyClosed) {
		return ErrClientDisconnected
	}
	if de, ok := err.(driver.Error); ok {
		return CommandError{
			Code:    de.Code,
			Message: de.Message,
			Labels:  de.Labels,
			Name:    de.Name,
			Wrapped: de.Wrapped,
			Raw:     bson.Raw(de.Raw),
		}
	}
	if qe, ok := err.(driver.QueryFailureError); ok {
		
		ce := CommandError{
			Name:    qe.Message,
			Wrapped: qe.Wrapped,
			Raw:     bson.Raw(qe.Response),
		}

		dollarErr, err := qe.Response.LookupErr("$err")
		if err == nil {
			ce.Message, _ = dollarErr.StringValueOK()
		}
		code, err := qe.Response.LookupErr("code")
		if err == nil {
			ce.Code, _ = code.Int32OK()
		}

		return ce
	}
	if me, ok := err.(mongocrypt.Error); ok {
		return MongocryptError{Code: me.Code, Message: me.Message}
	}

	if errors.Is(err, codecutil.ErrNilValue) {
		return ErrNilValue
	}

	if marshalErr, ok := err.(codecutil.MarshalError); ok {
		return MarshalError{
			Value: marshalErr.Value,
			Err:   marshalErr.Err,
		}
	}

	return err
}


func IsDuplicateKeyError(err error) bool {
	if se := ServerError(nil); errors.As(err, &se) {
		return se.HasErrorCode(11000) || 
			se.HasErrorCode(11001) || 
			
			se.HasErrorCode(12582) ||
			
			
			se.HasErrorCodeWithMessage(16460, " E11000 ")
	}
	return false
}


var timeoutErrs = [...]error{
	context.DeadlineExceeded,
	driver.ErrDeadlineWouldBeExceeded,
	topology.ErrServerSelectionTimeout,
}



func IsTimeout(err error) bool {
	
	for _, target := range timeoutErrs {
		if errors.Is(err, target) {
			return true
		}
	}

	
	
	if errors.As(err, &topology.WaitQueueTimeoutError{}) {
		return true
	}
	if ce := (CommandError{}); errors.As(err, &ce) && ce.IsMaxTimeMSExpiredError() {
		return true
	}
	if we := (WriteException{}); errors.As(err, &we) && we.WriteConcernError != nil && we.WriteConcernError.IsMaxTimeMSExpiredError() {
		return true
	}
	if ne := net.Error(nil); errors.As(err, &ne) {
		return ne.Timeout()
	}
	
	if le := LabeledError(nil); errors.As(err, &le) {
		if le.HasErrorLabel("NetworkTimeoutError") || le.HasErrorLabel("ExceededTimeLimitError") {
			return true
		}
	}

	return false
}


func unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}


func errorHasLabel(err error, label string) bool {
	for ; err != nil; err = unwrap(err) {
		if le, ok := err.(LabeledError); ok && le.HasErrorLabel(label) {
			return true
		}
	}
	return false
}


func IsNetworkError(err error) bool {
	return errorHasLabel(err, "NetworkError")
}


type MongocryptError struct {
	Code    int32
	Message string
}


func (m MongocryptError) Error() string {
	return fmt.Sprintf("mongocrypt error %d: %v", m.Code, m.Message)
}



type EncryptionKeyVaultError struct {
	Wrapped error
}


func (ekve EncryptionKeyVaultError) Error() string {
	return fmt.Sprintf("key vault communication error: %v", ekve.Wrapped)
}


func (ekve EncryptionKeyVaultError) Unwrap() error {
	return ekve.Wrapped
}


type MongocryptdError struct {
	Wrapped error
}


func (e MongocryptdError) Error() string {
	return fmt.Sprintf("mongocryptd communication error: %v", e.Wrapped)
}


func (e MongocryptdError) Unwrap() error {
	return e.Wrapped
}


type LabeledError interface {
	error
	
	HasErrorLabel(string) bool
}



type ServerError interface {
	LabeledError
	
	HasErrorCode(int) bool
	
	HasErrorMessage(string) bool
	
	HasErrorCodeWithMessage(int, string) bool

	serverError()
}

var _ ServerError = CommandError{}
var _ ServerError = WriteError{}
var _ ServerError = WriteException{}
var _ ServerError = BulkWriteException{}


type CommandError struct {
	Code    int32
	Message string
	Labels  []string 
	Name    string   
	Wrapped error    
	Raw     bson.Raw 
}


func (e CommandError) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("(%v) %v", e.Name, e.Message)
	}
	return e.Message
}


func (e CommandError) Unwrap() error {
	return e.Wrapped
}


func (e CommandError) HasErrorCode(code int) bool {
	return int(e.Code) == code
}


func (e CommandError) HasErrorLabel(label string) bool {
	if e.Labels != nil {
		for _, l := range e.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}


func (e CommandError) HasErrorMessage(message string) bool {
	return strings.Contains(e.Message, message)
}


func (e CommandError) HasErrorCodeWithMessage(code int, message string) bool {
	return int(e.Code) == code && strings.Contains(e.Message, message)
}


func (e CommandError) IsMaxTimeMSExpiredError() bool {
	return e.Code == 50 || e.Name == "MaxTimeMSExpired"
}


func (e CommandError) serverError() {}



type WriteError struct {
	
	Index int

	Code    int
	Message string
	Details bson.Raw

	
	Raw bson.Raw
}

func (we WriteError) Error() string {
	msg := we.Message
	if len(we.Details) > 0 {
		msg = fmt.Sprintf("%s: %s", msg, we.Details.String())
	}
	return msg
}


func (we WriteError) HasErrorCode(code int) bool {
	return we.Code == code
}



func (we WriteError) HasErrorLabel(string) bool {
	return false
}


func (we WriteError) HasErrorMessage(message string) bool {
	return strings.Contains(we.Message, message)
}


func (we WriteError) HasErrorCodeWithMessage(code int, message string) bool {
	return we.Code == code && strings.Contains(we.Message, message)
}


func (we WriteError) serverError() {}


type WriteErrors []WriteError


func (we WriteErrors) Error() string {
	errs := make([]error, len(we))
	for i := 0; i < len(we); i++ {
		errs[i] = we[i]
	}
	
	return "write errors: " + joinBatchErrors(errs)
}

func writeErrorsFromDriverWriteErrors(errs driver.WriteErrors) WriteErrors {
	wes := make(WriteErrors, 0, len(errs))
	for _, err := range errs {
		wes = append(wes, WriteError{
			Index:   int(err.Index),
			Code:    int(err.Code),
			Message: err.Message,
			Details: bson.Raw(err.Details),
			Raw:     bson.Raw(err.Raw),
		})
	}
	return wes
}



type WriteConcernError struct {
	Name    string
	Code    int
	Message string
	Details bson.Raw
	Raw     bson.Raw 
}


func (wce WriteConcernError) Error() string {
	if wce.Name != "" {
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	}
	return wce.Message
}


func (wce WriteConcernError) IsMaxTimeMSExpiredError() bool {
	return wce.Code == 50
}



type WriteException struct {
	
	WriteConcernError *WriteConcernError

	
	WriteErrors WriteErrors

	
	Labels []string

	
	Raw bson.Raw
}


func (mwe WriteException) Error() string {
	causes := make([]string, 0, 2)
	if mwe.WriteConcernError != nil {
		causes = append(causes, "write concern error: "+mwe.WriteConcernError.Error())
	}
	if len(mwe.WriteErrors) > 0 {
		
		
		causes = append(causes, mwe.WriteErrors.Error())
	}

	message := "write exception: "
	if len(causes) == 0 {
		return message + "no causes"
	}
	return message + strings.Join(causes, ", ")
}


func (mwe WriteException) HasErrorCode(code int) bool {
	if mwe.WriteConcernError != nil && mwe.WriteConcernError.Code == code {
		return true
	}
	for _, we := range mwe.WriteErrors {
		if we.Code == code {
			return true
		}
	}
	return false
}


func (mwe WriteException) HasErrorLabel(label string) bool {
	if mwe.Labels != nil {
		for _, l := range mwe.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}


func (mwe WriteException) HasErrorMessage(message string) bool {
	if mwe.WriteConcernError != nil && strings.Contains(mwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range mwe.WriteErrors {
		if strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}


func (mwe WriteException) HasErrorCodeWithMessage(code int, message string) bool {
	if mwe.WriteConcernError != nil &&
		mwe.WriteConcernError.Code == code && strings.Contains(mwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range mwe.WriteErrors {
		if we.Code == code && strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}


func (mwe WriteException) serverError() {}

func convertDriverWriteConcernError(wce *driver.WriteConcernError) *WriteConcernError {
	if wce == nil {
		return nil
	}

	return &WriteConcernError{
		Name:    wce.Name,
		Code:    int(wce.Code),
		Message: wce.Message,
		Details: bson.Raw(wce.Details),
		Raw:     bson.Raw(wce.Raw),
	}
}



type BulkWriteError struct {
	WriteError            
	Request    WriteModel 
}


func (bwe BulkWriteError) Error() string {
	return bwe.WriteError.Error()
}


type BulkWriteException struct {
	
	WriteConcernError *WriteConcernError

	
	WriteErrors []BulkWriteError

	
	Labels []string
}


func (bwe BulkWriteException) Error() string {
	causes := make([]string, 0, 2)
	if bwe.WriteConcernError != nil {
		causes = append(causes, "write concern error: "+bwe.WriteConcernError.Error())
	}
	if len(bwe.WriteErrors) > 0 {
		errs := make([]error, len(bwe.WriteErrors))
		for i := 0; i < len(bwe.WriteErrors); i++ {
			errs[i] = &bwe.WriteErrors[i]
		}
		causes = append(causes, "write errors: "+joinBatchErrors(errs))
	}

	message := "bulk write exception: "
	if len(causes) == 0 {
		return message + "no causes"
	}
	return "bulk write exception: " + strings.Join(causes, ", ")
}


func (bwe BulkWriteException) HasErrorCode(code int) bool {
	if bwe.WriteConcernError != nil && bwe.WriteConcernError.Code == code {
		return true
	}
	for _, we := range bwe.WriteErrors {
		if we.Code == code {
			return true
		}
	}
	return false
}


func (bwe BulkWriteException) HasErrorLabel(label string) bool {
	if bwe.Labels != nil {
		for _, l := range bwe.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}


func (bwe BulkWriteException) HasErrorMessage(message string) bool {
	if bwe.WriteConcernError != nil && strings.Contains(bwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range bwe.WriteErrors {
		if strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}


func (bwe BulkWriteException) HasErrorCodeWithMessage(code int, message string) bool {
	if bwe.WriteConcernError != nil &&
		bwe.WriteConcernError.Code == code && strings.Contains(bwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range bwe.WriteErrors {
		if we.Code == code && strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}


func (bwe BulkWriteException) serverError() {}





type returnResult int

const (
	rrNone returnResult = 1 << iota 
	rrOne                           
	rrMany                          

	rrAll returnResult = rrOne | rrMany 
)






func processWriteError(err error) (returnResult, error) {
	switch {
	case errors.Is(err, driver.ErrUnacknowledgedWrite):
		return rrAll, ErrUnacknowledgedWrite
	case err != nil:
		switch tt := err.(type) {
		case driver.WriteCommandError:
			return rrMany, WriteException{
				WriteConcernError: convertDriverWriteConcernError(tt.WriteConcernError),
				WriteErrors:       writeErrorsFromDriverWriteErrors(tt.WriteErrors),
				Labels:            tt.Labels,
				Raw:               bson.Raw(tt.Raw),
			}
		default:
			return rrNone, replaceErrors(err)
		}
	default:
		return rrAll, nil
	}
}




const batchErrorsTargetLength = 2000








func joinBatchErrors(errs []error) string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "[")
	for idx, err := range errs {
		if idx != 0 {
			fmt.Fprint(&buf, ", ")
		}
		
		
		if buf.Len() > batchErrorsTargetLength {
			fmt.Fprintf(&buf, "+%d more errors...", len(errs)-idx)
			break
		}
		fmt.Fprint(&buf, err.Error())
	}
	fmt.Fprint(&buf, "]")

	return buf.String()
}
