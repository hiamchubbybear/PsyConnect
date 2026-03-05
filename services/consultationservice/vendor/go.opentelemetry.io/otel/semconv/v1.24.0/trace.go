




package semconv 

import "go.opentelemetry.io/otel/attribute"


const (
	
	
	
	
	
	
	
	
	
	
	PeerServiceKey = attribute.Key("peer.service")
)






func PeerService(val string) attribute.KeyValue {
	return PeerServiceKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	EnduserIDKey = attribute.Key("enduser.id")

	
	
	
	
	
	
	
	
	
	EnduserRoleKey = attribute.Key("enduser.role")

	
	
	
	
	
	
	
	
	
	
	
	
	
	EnduserScopeKey = attribute.Key("enduser.scope")
)






func EnduserID(val string) attribute.KeyValue {
	return EnduserIDKey.String(val)
}





func EnduserRole(val string) attribute.KeyValue {
	return EnduserRoleKey.String(val)
}









func EnduserScope(val string) attribute.KeyValue {
	return EnduserScopeKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	CodeColumnKey = attribute.Key("code.column")

	
	
	
	
	
	
	
	
	
	CodeFilepathKey = attribute.Key("code.filepath")

	
	
	
	
	
	
	
	
	CodeFunctionKey = attribute.Key("code.function")

	
	
	
	
	
	
	
	
	
	CodeLineNumberKey = attribute.Key("code.lineno")

	
	
	
	
	
	
	
	
	
	
	CodeNamespaceKey = attribute.Key("code.namespace")

	
	
	
	
	
	
	
	
	
	
	
	
	CodeStacktraceKey = attribute.Key("code.stacktrace")
)





func CodeColumn(val int) attribute.KeyValue {
	return CodeColumnKey.Int(val)
}





func CodeFilepath(val string) attribute.KeyValue {
	return CodeFilepathKey.String(val)
}




func CodeFunction(val string) attribute.KeyValue {
	return CodeFunctionKey.String(val)
}





func CodeLineNumber(val int) attribute.KeyValue {
	return CodeLineNumberKey.Int(val)
}






func CodeNamespace(val string) attribute.KeyValue {
	return CodeNamespaceKey.String(val)
}





func CodeStacktrace(val string) attribute.KeyValue {
	return CodeStacktraceKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	ThreadIDKey = attribute.Key("thread.id")

	
	
	
	
	
	
	
	ThreadNameKey = attribute.Key("thread.name")
)




func ThreadID(val int) attribute.KeyValue {
	return ThreadIDKey.Int(val)
}



func ThreadName(val string) attribute.KeyValue {
	return ThreadNameKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	
	AWSLambdaInvokedARNKey = attribute.Key("aws.lambda.invoked_arn")
)






func AWSLambdaInvokedARN(val string) attribute.KeyValue {
	return AWSLambdaInvokedARNKey.String(val)
}





const (
	
	
	
	
	
	
	
	
	
	CloudeventsEventIDKey = attribute.Key("cloudevents.event_id")

	
	
	
	
	
	
	
	
	
	
	CloudeventsEventSourceKey = attribute.Key("cloudevents.event_source")

	
	
	
	
	
	
	
	
	
	
	CloudeventsEventSpecVersionKey = attribute.Key("cloudevents.event_spec_version")

	
	
	
	
	
	
	
	
	
	
	CloudeventsEventSubjectKey = attribute.Key("cloudevents.event_subject")

	
	
	
	
	
	
	
	
	
	
	
	CloudeventsEventTypeKey = attribute.Key("cloudevents.event_type")
)





func CloudeventsEventID(val string) attribute.KeyValue {
	return CloudeventsEventIDKey.String(val)
}





func CloudeventsEventSource(val string) attribute.KeyValue {
	return CloudeventsEventSourceKey.String(val)
}






func CloudeventsEventSpecVersion(val string) attribute.KeyValue {
	return CloudeventsEventSpecVersionKey.String(val)
}





func CloudeventsEventSubject(val string) attribute.KeyValue {
	return CloudeventsEventSubjectKey.String(val)
}






func CloudeventsEventType(val string) attribute.KeyValue {
	return CloudeventsEventTypeKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	OpentracingRefTypeKey = attribute.Key("opentracing.ref_type")
)

var (
	
	OpentracingRefTypeChildOf = OpentracingRefTypeKey.String("child_of")
	
	OpentracingRefTypeFollowsFrom = OpentracingRefTypeKey.String("follows_from")
)



const (
	
	
	
	
	
	
	
	
	OTelStatusCodeKey = attribute.Key("otel.status_code")

	
	
	
	
	
	
	
	
	OTelStatusDescriptionKey = attribute.Key("otel.status_description")
)

var (
	
	OTelStatusCodeOk = OTelStatusCodeKey.String("OK")
	
	OTelStatusCodeError = OTelStatusCodeKey.String("ERROR")
)




func OTelStatusDescription(val string) attribute.KeyValue {
	return OTelStatusDescriptionKey.String(val)
}




const (
	
	
	
	
	
	
	
	
	FaaSInvocationIDKey = attribute.Key("faas.invocation_id")
)




func FaaSInvocationID(val string) attribute.KeyValue {
	return FaaSInvocationIDKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	FaaSDocumentCollectionKey = attribute.Key("faas.document.collection")

	
	
	
	
	
	
	
	
	
	FaaSDocumentNameKey = attribute.Key("faas.document.name")

	
	
	
	
	
	
	
	FaaSDocumentOperationKey = attribute.Key("faas.document.operation")

	
	
	
	
	
	
	
	
	
	
	FaaSDocumentTimeKey = attribute.Key("faas.document.time")
)

var (
	
	FaaSDocumentOperationInsert = FaaSDocumentOperationKey.String("insert")
	
	FaaSDocumentOperationEdit = FaaSDocumentOperationKey.String("edit")
	
	FaaSDocumentOperationDelete = FaaSDocumentOperationKey.String("delete")
)






func FaaSDocumentCollection(val string) attribute.KeyValue {
	return FaaSDocumentCollectionKey.String(val)
}





func FaaSDocumentName(val string) attribute.KeyValue {
	return FaaSDocumentNameKey.String(val)
}






func FaaSDocumentTime(val string) attribute.KeyValue {
	return FaaSDocumentTimeKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	FaaSCronKey = attribute.Key("faas.cron")

	
	
	
	
	
	
	
	
	
	
	FaaSTimeKey = attribute.Key("faas.time")
)





func FaaSCron(val string) attribute.KeyValue {
	return FaaSCronKey.String(val)
}






func FaaSTime(val string) attribute.KeyValue {
	return FaaSTimeKey.String(val)
}


const (
	
	
	
	
	
	
	
	FaaSColdstartKey = attribute.Key("faas.coldstart")
)




func FaaSColdstart(val bool) attribute.KeyValue {
	return FaaSColdstartKey.Bool(val)
}








const (
	
	
	
	
	
	
	
	
	AWSRequestIDKey = attribute.Key("aws.request_id")
)




func AWSRequestID(val string) attribute.KeyValue {
	return AWSRequestIDKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	AWSDynamoDBAttributesToGetKey = attribute.Key("aws.dynamodb.attributes_to_get")

	
	
	
	
	
	
	
	AWSDynamoDBConsistentReadKey = attribute.Key("aws.dynamodb.consistent_read")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBConsumedCapacityKey = attribute.Key("aws.dynamodb.consumed_capacity")

	
	
	
	
	
	
	
	
	AWSDynamoDBIndexNameKey = attribute.Key("aws.dynamodb.index_name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBItemCollectionMetricsKey = attribute.Key("aws.dynamodb.item_collection_metrics")

	
	
	
	
	
	
	
	
	AWSDynamoDBLimitKey = attribute.Key("aws.dynamodb.limit")

	
	
	
	
	
	
	
	
	
	AWSDynamoDBProjectionKey = attribute.Key("aws.dynamodb.projection")

	
	
	
	
	
	
	
	
	
	AWSDynamoDBProvisionedReadCapacityKey = attribute.Key("aws.dynamodb.provisioned_read_capacity")

	
	
	
	
	
	
	
	
	
	AWSDynamoDBProvisionedWriteCapacityKey = attribute.Key("aws.dynamodb.provisioned_write_capacity")

	
	
	
	
	
	
	
	
	AWSDynamoDBSelectKey = attribute.Key("aws.dynamodb.select")

	
	
	
	
	
	
	
	
	AWSDynamoDBTableNamesKey = attribute.Key("aws.dynamodb.table_names")
)




func AWSDynamoDBAttributesToGet(val ...string) attribute.KeyValue {
	return AWSDynamoDBAttributesToGetKey.StringSlice(val)
}




func AWSDynamoDBConsistentRead(val bool) attribute.KeyValue {
	return AWSDynamoDBConsistentReadKey.Bool(val)
}




func AWSDynamoDBConsumedCapacity(val ...string) attribute.KeyValue {
	return AWSDynamoDBConsumedCapacityKey.StringSlice(val)
}




func AWSDynamoDBIndexName(val string) attribute.KeyValue {
	return AWSDynamoDBIndexNameKey.String(val)
}





func AWSDynamoDBItemCollectionMetrics(val string) attribute.KeyValue {
	return AWSDynamoDBItemCollectionMetricsKey.String(val)
}




func AWSDynamoDBLimit(val int) attribute.KeyValue {
	return AWSDynamoDBLimitKey.Int(val)
}




func AWSDynamoDBProjection(val string) attribute.KeyValue {
	return AWSDynamoDBProjectionKey.String(val)
}





func AWSDynamoDBProvisionedReadCapacity(val float64) attribute.KeyValue {
	return AWSDynamoDBProvisionedReadCapacityKey.Float64(val)
}





func AWSDynamoDBProvisionedWriteCapacity(val float64) attribute.KeyValue {
	return AWSDynamoDBProvisionedWriteCapacityKey.Float64(val)
}




func AWSDynamoDBSelect(val string) attribute.KeyValue {
	return AWSDynamoDBSelectKey.String(val)
}




func AWSDynamoDBTableNames(val ...string) attribute.KeyValue {
	return AWSDynamoDBTableNamesKey.StringSlice(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBGlobalSecondaryIndexesKey = attribute.Key("aws.dynamodb.global_secondary_indexes")

	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBLocalSecondaryIndexesKey = attribute.Key("aws.dynamodb.local_secondary_indexes")
)





func AWSDynamoDBGlobalSecondaryIndexes(val ...string) attribute.KeyValue {
	return AWSDynamoDBGlobalSecondaryIndexesKey.StringSlice(val)
}





func AWSDynamoDBLocalSecondaryIndexes(val ...string) attribute.KeyValue {
	return AWSDynamoDBLocalSecondaryIndexesKey.StringSlice(val)
}


const (
	
	
	
	
	
	
	
	
	AWSDynamoDBExclusiveStartTableKey = attribute.Key("aws.dynamodb.exclusive_start_table")

	
	
	
	
	
	
	
	
	AWSDynamoDBTableCountKey = attribute.Key("aws.dynamodb.table_count")
)




func AWSDynamoDBExclusiveStartTable(val string) attribute.KeyValue {
	return AWSDynamoDBExclusiveStartTableKey.String(val)
}




func AWSDynamoDBTableCount(val int) attribute.KeyValue {
	return AWSDynamoDBTableCountKey.Int(val)
}


const (
	
	
	
	
	
	
	
	AWSDynamoDBScanForwardKey = attribute.Key("aws.dynamodb.scan_forward")
)




func AWSDynamoDBScanForward(val bool) attribute.KeyValue {
	return AWSDynamoDBScanForwardKey.Bool(val)
}


const (
	
	
	
	
	
	
	
	
	AWSDynamoDBCountKey = attribute.Key("aws.dynamodb.count")

	
	
	
	
	
	
	
	
	AWSDynamoDBScannedCountKey = attribute.Key("aws.dynamodb.scanned_count")

	
	
	
	
	
	
	
	
	AWSDynamoDBSegmentKey = attribute.Key("aws.dynamodb.segment")

	
	
	
	
	
	
	
	
	AWSDynamoDBTotalSegmentsKey = attribute.Key("aws.dynamodb.total_segments")
)




func AWSDynamoDBCount(val int) attribute.KeyValue {
	return AWSDynamoDBCountKey.Int(val)
}




func AWSDynamoDBScannedCount(val int) attribute.KeyValue {
	return AWSDynamoDBScannedCountKey.Int(val)
}




func AWSDynamoDBSegment(val int) attribute.KeyValue {
	return AWSDynamoDBSegmentKey.Int(val)
}




func AWSDynamoDBTotalSegments(val int) attribute.KeyValue {
	return AWSDynamoDBTotalSegmentsKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	AWSDynamoDBAttributeDefinitionsKey = attribute.Key("aws.dynamodb.attribute_definitions")

	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBGlobalSecondaryIndexUpdatesKey = attribute.Key("aws.dynamodb.global_secondary_index_updates")
)





func AWSDynamoDBAttributeDefinitions(val ...string) attribute.KeyValue {
	return AWSDynamoDBAttributeDefinitionsKey.StringSlice(val)
}





func AWSDynamoDBGlobalSecondaryIndexUpdates(val ...string) attribute.KeyValue {
	return AWSDynamoDBGlobalSecondaryIndexUpdatesKey.StringSlice(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSS3BucketKey = attribute.Key("aws.s3.bucket")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSS3CopySourceKey = attribute.Key("aws.s3.copy_source")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSS3DeleteKey = attribute.Key("aws.s3.delete")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSS3KeyKey = attribute.Key("aws.s3.key")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSS3PartNumberKey = attribute.Key("aws.s3.part_number")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSS3UploadIDKey = attribute.Key("aws.s3.upload_id")
)






func AWSS3Bucket(val string) attribute.KeyValue {
	return AWSS3BucketKey.String(val)
}




func AWSS3CopySource(val string) attribute.KeyValue {
	return AWSS3CopySourceKey.String(val)
}




func AWSS3Delete(val string) attribute.KeyValue {
	return AWSS3DeleteKey.String(val)
}






func AWSS3Key(val string) attribute.KeyValue {
	return AWSS3KeyKey.String(val)
}





func AWSS3PartNumber(val int) attribute.KeyValue {
	return AWSS3PartNumberKey.Int(val)
}




func AWSS3UploadID(val string) attribute.KeyValue {
	return AWSS3UploadIDKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	GraphqlDocumentKey = attribute.Key("graphql.document")

	
	
	
	
	
	
	
	
	GraphqlOperationNameKey = attribute.Key("graphql.operation.name")

	
	
	
	
	
	
	
	
	GraphqlOperationTypeKey = attribute.Key("graphql.operation.type")
)

var (
	
	GraphqlOperationTypeQuery = GraphqlOperationTypeKey.String("query")
	
	GraphqlOperationTypeMutation = GraphqlOperationTypeKey.String("mutation")
	
	GraphqlOperationTypeSubscription = GraphqlOperationTypeKey.String("subscription")
)




func GraphqlDocument(val string) attribute.KeyValue {
	return GraphqlDocumentKey.String(val)
}




func GraphqlOperationName(val string) attribute.KeyValue {
	return GraphqlOperationNameKey.String(val)
}
