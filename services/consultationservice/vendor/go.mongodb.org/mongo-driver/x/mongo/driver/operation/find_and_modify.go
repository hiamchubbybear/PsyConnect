





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type FindAndModify struct {
	authenticator            driver.Authenticator
	arrayFilters             bsoncore.Array
	bypassDocumentValidation *bool
	collation                bsoncore.Document
	comment                  bsoncore.Value
	fields                   bsoncore.Document
	maxTime                  *time.Duration
	newDocument              *bool
	query                    bsoncore.Document
	remove                   *bool
	sort                     bsoncore.Document
	update                   bsoncore.Value
	upsert                   *bool
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	crypt                    driver.Crypt
	hint                     bsoncore.Value
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	timeout                  *time.Duration
	bypassEmptyTsReplacement *bool

	result FindAndModifyResult
}


type LastErrorObject struct {
	
	UpdatedExisting bool
	
	Upserted interface{}
}


type FindAndModifyResult struct {
	
	Value bsoncore.Document
	
	LastErrorObject LastErrorObject
}

func buildFindAndModifyResult(response bsoncore.Document) (FindAndModifyResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return FindAndModifyResult{}, err
	}
	famr := FindAndModifyResult{}
	for _, element := range elements {
		switch element.Key() {
		case "value":
			var ok bool
			famr.Value, ok = element.Value().DocumentOK()

			
			if element.Value().Type != bsontype.Null && !ok {
				return famr, fmt.Errorf("response field 'value' is type document or null, but received BSON type %s", element.Value().Type)
			}
		case "lastErrorObject":
			valDoc, ok := element.Value().DocumentOK()
			if !ok {
				return famr, fmt.Errorf("response field 'lastErrorObject' is type document, but received BSON type %s", element.Value().Type)
			}

			var leo LastErrorObject
			if err = bson.Unmarshal(valDoc, &leo); err != nil {
				return famr, err
			}
			famr.LastErrorObject = leo
		}
	}
	return famr, nil
}


func NewFindAndModify(query bsoncore.Document) *FindAndModify {
	return &FindAndModify{
		query: query,
	}
}


func (fam *FindAndModify) Result() FindAndModifyResult { return fam.result }

func (fam *FindAndModify) processResponse(info driver.ResponseInfo) error {
	var err error

	fam.result, err = buildFindAndModifyResult(info.ServerResponse)
	return err

}


func (fam *FindAndModify) Execute(ctx context.Context) error {
	if fam.deployment == nil {
		return errors.New("the FindAndModify operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         fam.command,
		ProcessResponseFn: fam.processResponse,

		RetryMode:      fam.retry,
		Type:           driver.Write,
		Client:         fam.session,
		Clock:          fam.clock,
		CommandMonitor: fam.monitor,
		Database:       fam.database,
		Deployment:     fam.deployment,
		MaxTime:        fam.maxTime,
		Selector:       fam.selector,
		WriteConcern:   fam.writeConcern,
		Crypt:          fam.crypt,
		ServerAPI:      fam.serverAPI,
		Timeout:        fam.timeout,
		Name:           driverutil.FindAndModifyOp,
		Authenticator:  fam.authenticator,
	}.Execute(ctx)

}

func (fam *FindAndModify) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "findAndModify", fam.collection)
	if fam.arrayFilters != nil {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(6) {
			return nil, errors.New("the 'arrayFilters' command parameter requires a minimum server wire version of 6")
		}
		dst = bsoncore.AppendArrayElement(dst, "arrayFilters", fam.arrayFilters)
	}
	if fam.bypassDocumentValidation != nil {

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *fam.bypassDocumentValidation)
	}
	if fam.collation != nil {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", fam.collation)
	}
	if fam.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", fam.comment)
	}
	if fam.fields != nil {

		dst = bsoncore.AppendDocumentElement(dst, "fields", fam.fields)
	}
	if fam.newDocument != nil {

		dst = bsoncore.AppendBooleanElement(dst, "new", *fam.newDocument)
	}
	if fam.query != nil {

		dst = bsoncore.AppendDocumentElement(dst, "query", fam.query)
	}
	if fam.remove != nil {

		dst = bsoncore.AppendBooleanElement(dst, "remove", *fam.remove)
	}
	if fam.sort != nil {

		dst = bsoncore.AppendDocumentElement(dst, "sort", fam.sort)
	}
	if fam.update.Data != nil {
		dst = bsoncore.AppendValueElement(dst, "update", fam.update)
	}
	if fam.upsert != nil {

		dst = bsoncore.AppendBooleanElement(dst, "upsert", *fam.upsert)
	}
	if fam.hint.Type != bsontype.Type(0) {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(8) {
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 8")
		}
		if !fam.writeConcern.Acknowledged() {
			return nil, errUnacknowledgedHint
		}
		dst = bsoncore.AppendValueElement(dst, "hint", fam.hint)
	}
	if fam.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", fam.let)
	}
	if fam.bypassEmptyTsReplacement != nil {
		dst = bsoncore.AppendBooleanElement(dst, "bypassEmptyTsReplacement", *fam.bypassEmptyTsReplacement)
	}

	return dst, nil
}


func (fam *FindAndModify) ArrayFilters(arrayFilters bsoncore.Array) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.arrayFilters = arrayFilters
	return fam
}


func (fam *FindAndModify) BypassDocumentValidation(bypassDocumentValidation bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.bypassDocumentValidation = &bypassDocumentValidation
	return fam
}


func (fam *FindAndModify) Collation(collation bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.collation = collation
	return fam
}


func (fam *FindAndModify) Comment(comment bsoncore.Value) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.comment = comment
	return fam
}


func (fam *FindAndModify) Fields(fields bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.fields = fields
	return fam
}


func (fam *FindAndModify) MaxTime(maxTime *time.Duration) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.maxTime = maxTime
	return fam
}


func (fam *FindAndModify) NewDocument(newDocument bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.newDocument = &newDocument
	return fam
}


func (fam *FindAndModify) Query(query bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.query = query
	return fam
}


func (fam *FindAndModify) Remove(remove bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.remove = &remove
	return fam
}


func (fam *FindAndModify) Sort(sort bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.sort = sort
	return fam
}


func (fam *FindAndModify) Update(update bsoncore.Value) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.update = update
	return fam
}


func (fam *FindAndModify) Upsert(upsert bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.upsert = &upsert
	return fam
}


func (fam *FindAndModify) Session(session *session.Client) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.session = session
	return fam
}


func (fam *FindAndModify) ClusterClock(clock *session.ClusterClock) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.clock = clock
	return fam
}


func (fam *FindAndModify) Collection(collection string) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.collection = collection
	return fam
}


func (fam *FindAndModify) CommandMonitor(monitor *event.CommandMonitor) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.monitor = monitor
	return fam
}


func (fam *FindAndModify) Database(database string) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.database = database
	return fam
}


func (fam *FindAndModify) Deployment(deployment driver.Deployment) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.deployment = deployment
	return fam
}


func (fam *FindAndModify) ServerSelector(selector description.ServerSelector) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.selector = selector
	return fam
}


func (fam *FindAndModify) WriteConcern(writeConcern *writeconcern.WriteConcern) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.writeConcern = writeConcern
	return fam
}




func (fam *FindAndModify) Retry(retry driver.RetryMode) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.retry = &retry
	return fam
}


func (fam *FindAndModify) Crypt(crypt driver.Crypt) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.crypt = crypt
	return fam
}


func (fam *FindAndModify) Hint(hint bsoncore.Value) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.hint = hint
	return fam
}


func (fam *FindAndModify) ServerAPI(serverAPI *driver.ServerAPIOptions) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.serverAPI = serverAPI
	return fam
}


func (fam *FindAndModify) Let(let bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.let = let
	return fam
}


func (fam *FindAndModify) Timeout(timeout *time.Duration) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.timeout = timeout
	return fam
}


func (fam *FindAndModify) Authenticator(authenticator driver.Authenticator) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.authenticator = authenticator
	return fam
}


func (fam *FindAndModify) BypassEmptyTsReplacement(bypassEmptyTsReplacement bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.bypassEmptyTsReplacement = &bypassEmptyTsReplacement
	return fam
}
