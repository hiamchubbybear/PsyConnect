





//go:build cse
// +build cse

package mongocrypt






import "C"
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/httputil"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/creds"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

type kmsProvider interface {
	GetCredentialsDoc(context.Context) (bsoncore.Document, error)
}

type MongoCrypt struct {
	wrapped      *C.mongocrypt_t
	kmsProviders map[string]kmsProvider
	httpClient   *http.Client
}



func Version() string {
	str := C.GoString(C.mongocrypt_version(nil))
	return str
}


func NewMongoCrypt(opts *options.MongoCryptOptions) (*MongoCrypt, error) {
	
	wrapped := C.mongocrypt_new()
	if wrapped == nil {
		return nil, errors.New("could not create new mongocrypt object")
	}
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = httputil.DefaultHTTPClient
	}
	kmsProviders := make(map[string]kmsProvider)
	if needsKmsProvider(opts.KmsProviders, "gcp") {
		kmsProviders["gcp"] = creds.NewGCPCredentialProvider(httpClient)
	}
	if needsKmsProvider(opts.KmsProviders, "aws") {
		kmsProviders["aws"] = creds.NewAWSCredentialProvider(httpClient)
	}
	if needsKmsProvider(opts.KmsProviders, "azure") {
		kmsProviders["azure"] = creds.NewAzureCredentialProvider(httpClient)
	}
	crypt := &MongoCrypt{
		wrapped:      wrapped,
		kmsProviders: kmsProviders,
		httpClient:   httpClient,
	}

	
	if err := crypt.setProviderOptions(opts.KmsProviders); err != nil {
		return nil, err
	}
	if err := crypt.setLocalSchemaMap(opts.LocalSchemaMap); err != nil {
		return nil, err
	}
	if err := crypt.setEncryptedFieldsMap(opts.EncryptedFieldsMap); err != nil {
		return nil, err
	}

	if opts.BypassQueryAnalysis {
		C.mongocrypt_setopt_bypass_query_analysis(wrapped)
	}

	
	
	if !opts.CryptSharedLibDisabled {
		systemStr := C.CString("$SYSTEM")
		defer C.free(unsafe.Pointer(systemStr))
		C.mongocrypt_setopt_append_crypt_shared_lib_search_path(crypt.wrapped, systemStr)

		if opts.CryptSharedLibOverridePath != "" {
			cryptSharedLibOverridePathStr := C.CString(opts.CryptSharedLibOverridePath)
			defer C.free(unsafe.Pointer(cryptSharedLibOverridePathStr))
			C.mongocrypt_setopt_set_crypt_shared_lib_path_override(crypt.wrapped, cryptSharedLibOverridePathStr)
		}
	}

	C.mongocrypt_setopt_use_need_kms_credentials_state(crypt.wrapped)

	
	if !C.mongocrypt_init(crypt.wrapped) {
		return nil, crypt.createErrorFromStatus()
	}

	return crypt, nil
}


func (m *MongoCrypt) CreateEncryptionContext(db string, cmd bsoncore.Document) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	cmdBinary := newBinaryFromBytes(cmd)
	defer cmdBinary.close()
	dbStr := C.CString(db)
	defer C.free(unsafe.Pointer(dbStr))

	if ok := C.mongocrypt_ctx_encrypt_init(ctx.wrapped, dbStr, C.int32_t(-1), cmdBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}


func (m *MongoCrypt) CreateDecryptionContext(cmd bsoncore.Document) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	cmdBinary := newBinaryFromBytes(cmd)
	defer cmdBinary.close()

	if ok := C.mongocrypt_ctx_decrypt_init(ctx.wrapped, cmdBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}



func lookupString(doc bsoncore.Document, key string) string {
	strVal, _ := doc.Lookup(key).StringValueOK()
	return strVal
}

func setAltName(ctx *Context, altName string) error {
	
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendStringElement(doc, "keyAltName", altName)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	keyAltBinary := newBinaryFromBytes(doc)
	defer keyAltBinary.close()

	if ok := C.mongocrypt_ctx_setopt_key_alt_name(ctx.wrapped, keyAltBinary.wrapped); !ok {
		return ctx.createErrorFromStatus()
	}
	return nil
}

func setKeyMaterial(ctx *Context, keyMaterial []byte) error {
	
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendBinaryElement(doc, "keyMaterial", 0x00, keyMaterial)
	doc, err := bsoncore.AppendDocumentEnd(doc, idx)
	if err != nil {
		return err
	}

	keyMaterialBinary := newBinaryFromBytes(doc)
	defer keyMaterialBinary.close()

	if ok := C.mongocrypt_ctx_setopt_key_material(ctx.wrapped, keyMaterialBinary.wrapped); !ok {
		return ctx.createErrorFromStatus()
	}
	return nil
}

func rewrapDataKey(ctx *Context, filter []byte) error {
	filterBinary := newBinaryFromBytes(filter)
	defer filterBinary.close()

	if ok := C.mongocrypt_ctx_rewrap_many_datakey_init(ctx.wrapped, filterBinary.wrapped); !ok {
		return ctx.createErrorFromStatus()
	}
	return nil
}


func (m *MongoCrypt) CreateDataKeyContext(kmsProvider string, opts *options.DataKeyOptions) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	
	var masterKey bsoncore.Document
	switch {
	case opts.MasterKey != nil:
		
		
		masterKey = opts.MasterKey[:len(opts.MasterKey)-1]
		masterKey = bsoncore.AppendStringElement(masterKey, "provider", kmsProvider)
		masterKey, _ = bsoncore.AppendDocumentEnd(masterKey, 0)
	default:
		masterKey = bsoncore.NewDocumentBuilder().AppendString("provider", kmsProvider).Build()
	}

	masterKeyBinary := newBinaryFromBytes(masterKey)
	defer masterKeyBinary.close()

	if ok := C.mongocrypt_ctx_setopt_key_encryption_key(ctx.wrapped, masterKeyBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}

	for _, altName := range opts.KeyAltNames {
		if err := setAltName(ctx, altName); err != nil {
			return nil, err
		}
	}

	if opts.KeyMaterial != nil {
		if err := setKeyMaterial(ctx, opts.KeyMaterial); err != nil {
			return nil, err
		}
	}

	if ok := C.mongocrypt_ctx_datakey_init(ctx.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}

const (
	IndexTypeUnindexed = 1
	IndexTypeIndexed   = 2
)


func (m *MongoCrypt) createExplicitEncryptionContext(opts *options.ExplicitEncryptionOptions) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	if opts.KeyID != nil {
		keyIDBinary := newBinaryFromBytes(opts.KeyID.Data)
		defer keyIDBinary.close()

		if ok := C.mongocrypt_ctx_setopt_key_id(ctx.wrapped, keyIDBinary.wrapped); !ok {
			return nil, ctx.createErrorFromStatus()
		}
	}
	if opts.KeyAltName != nil {
		if err := setAltName(ctx, *opts.KeyAltName); err != nil {
			return nil, err
		}
	}

	if opts.RangeOptions != nil {
		idx, mongocryptDoc := bsoncore.AppendDocumentStart(nil)
		if opts.RangeOptions.Min != nil {
			mongocryptDoc = bsoncore.AppendValueElement(mongocryptDoc, "min", *opts.RangeOptions.Min)
		}
		if opts.RangeOptions.Max != nil {
			mongocryptDoc = bsoncore.AppendValueElement(mongocryptDoc, "max", *opts.RangeOptions.Max)
		}
		if opts.RangeOptions.Precision != nil {
			mongocryptDoc = bsoncore.AppendInt32Element(mongocryptDoc, "precision", *opts.RangeOptions.Precision)
		}
		if opts.RangeOptions.Sparsity != nil {
			mongocryptDoc = bsoncore.AppendInt64Element(mongocryptDoc, "sparsity", *opts.RangeOptions.Sparsity)
		}
		if opts.RangeOptions.TrimFactor != nil {
			mongocryptDoc = bsoncore.AppendInt32Element(mongocryptDoc, "trimFactor", *opts.RangeOptions.TrimFactor)
		}

		mongocryptDoc, err := bsoncore.AppendDocumentEnd(mongocryptDoc, idx)
		if err != nil {
			return nil, err
		}

		mongocryptBinary := newBinaryFromBytes(mongocryptDoc)
		defer mongocryptBinary.close()

		if ok := C.mongocrypt_ctx_setopt_algorithm_range(ctx.wrapped, mongocryptBinary.wrapped); !ok {
			return nil, ctx.createErrorFromStatus()
		}
	}

	algoStr := C.CString(opts.Algorithm)
	defer C.free(unsafe.Pointer(algoStr))

	if ok := C.mongocrypt_ctx_setopt_algorithm(ctx.wrapped, algoStr, -1); !ok {
		return nil, ctx.createErrorFromStatus()
	}

	if opts.QueryType != "" {
		queryStr := C.CString(opts.QueryType)
		defer C.free(unsafe.Pointer(queryStr))
		if ok := C.mongocrypt_ctx_setopt_query_type(ctx.wrapped, queryStr, -1); !ok {
			return nil, ctx.createErrorFromStatus()
		}
	}

	if opts.ContentionFactor != nil {
		if ok := C.mongocrypt_ctx_setopt_contention_factor(ctx.wrapped, C.int64_t(*opts.ContentionFactor)); !ok {
			return nil, ctx.createErrorFromStatus()
		}
	}
	return ctx, nil
}


func (m *MongoCrypt) CreateExplicitEncryptionContext(doc bsoncore.Document, opts *options.ExplicitEncryptionOptions) (*Context, error) {
	ctx, err := m.createExplicitEncryptionContext(opts)
	if err != nil {
		return ctx, err
	}
	docBinary := newBinaryFromBytes(doc)
	defer docBinary.close()
	if ok := C.mongocrypt_ctx_explicit_encrypt_init(ctx.wrapped, docBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}

	return ctx, nil
}


func (m *MongoCrypt) CreateExplicitEncryptionExpressionContext(doc bsoncore.Document, opts *options.ExplicitEncryptionOptions) (*Context, error) {
	ctx, err := m.createExplicitEncryptionContext(opts)
	if err != nil {
		return ctx, err
	}
	docBinary := newBinaryFromBytes(doc)
	defer docBinary.close()
	if ok := C.mongocrypt_ctx_explicit_encrypt_expression_init(ctx.wrapped, docBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}

	return ctx, nil
}


func (m *MongoCrypt) CreateExplicitDecryptionContext(doc bsoncore.Document) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	docBinary := newBinaryFromBytes(doc)
	defer docBinary.close()

	if ok := C.mongocrypt_ctx_explicit_decrypt_init(ctx.wrapped, docBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}



func (m *MongoCrypt) CryptSharedLibVersion() uint64 {
	return uint64(C.mongocrypt_crypt_shared_lib_version(m.wrapped))
}



func (m *MongoCrypt) CryptSharedLibVersionString() string {
	
	
	len := C.uint(0)
	str := C.GoString(C.mongocrypt_crypt_shared_lib_version_string(m.wrapped, &len))
	return str
}


func (m *MongoCrypt) Close() {
	C.mongocrypt_destroy(m.wrapped)
	if m.httpClient == httputil.DefaultHTTPClient {
		httputil.CloseIdleHTTPConnections(m.httpClient)
	}
}


func (m *MongoCrypt) RewrapDataKeyContext(filter []byte, opts *options.RewrapManyDataKeyOptions) (*Context, error) {
	const masterKey = "masterKey"
	const providerKey = "provider"

	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	if opts.MasterKey != nil && opts.Provider == nil {
		
		return nil, fmt.Errorf("expected 'Provider' to be set to identify type of 'MasterKey'")
	}

	if opts.Provider != nil {
		
		
		
		idx, mongocryptDoc := bsoncore.AppendDocumentStart(nil)
		mongocryptDoc = bsoncore.AppendStringElement(mongocryptDoc, providerKey, *opts.Provider)

		if opts.MasterKey != nil {
			mongocryptDoc = opts.MasterKey[:len(opts.MasterKey)-1]
			mongocryptDoc = bsoncore.AppendStringElement(mongocryptDoc, providerKey, *opts.Provider)
		}

		mongocryptDoc, err := bsoncore.AppendDocumentEnd(mongocryptDoc, idx)
		if err != nil {
			return nil, err
		}

		mongocryptBinary := newBinaryFromBytes(mongocryptDoc)
		defer mongocryptBinary.close()

		
		if ok := C.mongocrypt_ctx_setopt_key_encryption_key(ctx.wrapped, mongocryptBinary.wrapped); !ok {
			return nil, ctx.createErrorFromStatus()
		}
	}

	return ctx, rewrapDataKey(ctx, filter)
}

func (m *MongoCrypt) setProviderOptions(kmsProviders bsoncore.Document) error {
	providersBinary := newBinaryFromBytes(kmsProviders)
	defer providersBinary.close()

	if ok := C.mongocrypt_setopt_kms_providers(m.wrapped, providersBinary.wrapped); !ok {
		return m.createErrorFromStatus()
	}
	return nil
}


func (m *MongoCrypt) setLocalSchemaMap(schemaMap map[string]bsoncore.Document) error {
	if len(schemaMap) == 0 {
		return nil
	}

	
	schemaMapBSON, err := bson.Marshal(schemaMap)
	if err != nil {
		return fmt.Errorf("error marshalling SchemaMap: %v", err)
	}

	schemaMapBinary := newBinaryFromBytes(schemaMapBSON)
	defer schemaMapBinary.close()

	if ok := C.mongocrypt_setopt_schema_map(m.wrapped, schemaMapBinary.wrapped); !ok {
		return m.createErrorFromStatus()
	}
	return nil
}


func (m *MongoCrypt) setEncryptedFieldsMap(encryptedfieldsMap map[string]bsoncore.Document) error {
	if len(encryptedfieldsMap) == 0 {
		return nil
	}

	
	encryptedfieldsMapBSON, err := bson.Marshal(encryptedfieldsMap)
	if err != nil {
		return fmt.Errorf("error marshalling EncryptedFieldsMap: %v", err)
	}

	encryptedfieldsMapBinary := newBinaryFromBytes(encryptedfieldsMapBSON)
	defer encryptedfieldsMapBinary.close()

	if ok := C.mongocrypt_setopt_encrypted_field_config_map(m.wrapped, encryptedfieldsMapBinary.wrapped); !ok {
		return m.createErrorFromStatus()
	}
	return nil
}


func (m *MongoCrypt) createErrorFromStatus() error {
	status := C.mongocrypt_status_new()
	defer C.mongocrypt_status_destroy(status)
	C.mongocrypt_status(m.wrapped, status)
	return errorFromStatus(status)
}



func needsKmsProvider(kmsProviders bsoncore.Document, provider string) bool {
	val, err := kmsProviders.LookupErr(provider)
	if err != nil {
		
		return false
	}
	doc, ok := val.DocumentOK()
	
	
	return ok && len(doc) == 5
}



func (m *MongoCrypt) GetKmsProviders(ctx context.Context) (bsoncore.Document, error) {
	builder := bsoncore.NewDocumentBuilder()
	for k, p := range m.kmsProviders {
		doc, err := p.GetCredentialsDoc(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve %s credentials: %w", k, err)
		}
		builder.AppendDocument(k, doc)
	}
	return builder.Build(), nil
}
