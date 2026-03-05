




package semconv 

import "go.opentelemetry.io/otel/attribute"


const (
	
	
	
	
	
	
	
	
	
	
	
	AndroidOSAPILevelKey = attribute.Key("android.os.api_level")
)






func AndroidOSAPILevel(val string) attribute.KeyValue {
	return AndroidOSAPILevelKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	AspnetcoreRateLimitingResultKey = attribute.Key("aspnetcore.rate_limiting.result")

	
	
	
	
	
	
	
	
	
	
	
	AspnetcoreDiagnosticsHandlerTypeKey = attribute.Key("aspnetcore.diagnostics.handler.type")

	
	
	
	
	
	
	
	
	AspnetcoreDiagnosticsExceptionResultKey = attribute.Key("aspnetcore.diagnostics.exception.result")

	
	
	
	
	
	
	
	
	AspnetcoreRateLimitingPolicyKey = attribute.Key("aspnetcore.rate_limiting.policy")

	
	
	
	
	
	
	
	
	AspnetcoreRequestIsUnhandledKey = attribute.Key("aspnetcore.request.is_unhandled")

	
	
	
	
	
	
	
	
	AspnetcoreRoutingIsFallbackKey = attribute.Key("aspnetcore.routing.is_fallback")

	
	
	
	
	
	
	
	
	AspnetcoreRoutingMatchStatusKey = attribute.Key("aspnetcore.routing.match_status")
)

var (
	
	AspnetcoreRateLimitingResultAcquired = AspnetcoreRateLimitingResultKey.String("acquired")
	
	AspnetcoreRateLimitingResultEndpointLimiter = AspnetcoreRateLimitingResultKey.String("endpoint_limiter")
	
	AspnetcoreRateLimitingResultGlobalLimiter = AspnetcoreRateLimitingResultKey.String("global_limiter")
	
	AspnetcoreRateLimitingResultRequestCanceled = AspnetcoreRateLimitingResultKey.String("request_canceled")
)

var (
	
	AspnetcoreDiagnosticsExceptionResultHandled = AspnetcoreDiagnosticsExceptionResultKey.String("handled")
	
	AspnetcoreDiagnosticsExceptionResultUnhandled = AspnetcoreDiagnosticsExceptionResultKey.String("unhandled")
	
	AspnetcoreDiagnosticsExceptionResultSkipped = AspnetcoreDiagnosticsExceptionResultKey.String("skipped")
	
	AspnetcoreDiagnosticsExceptionResultAborted = AspnetcoreDiagnosticsExceptionResultKey.String("aborted")
)

var (
	
	AspnetcoreRoutingMatchStatusSuccess = AspnetcoreRoutingMatchStatusKey.String("success")
	
	AspnetcoreRoutingMatchStatusFailure = AspnetcoreRoutingMatchStatusKey.String("failure")
)






func AspnetcoreDiagnosticsHandlerType(val string) attribute.KeyValue {
	return AspnetcoreDiagnosticsHandlerTypeKey.String(val)
}




func AspnetcoreRateLimitingPolicy(val string) attribute.KeyValue {
	return AspnetcoreRateLimitingPolicyKey.String(val)
}




func AspnetcoreRequestIsUnhandled(val bool) attribute.KeyValue {
	return AspnetcoreRequestIsUnhandledKey.Bool(val)
}




func AspnetcoreRoutingIsFallback(val bool) attribute.KeyValue {
	return AspnetcoreRoutingIsFallbackKey.Bool(val)
}


const (
	
	
	
	
	
	
	
	
	AWSRequestIDKey = attribute.Key("aws.request_id")
)




func AWSRequestID(val string) attribute.KeyValue {
	return AWSRequestIDKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	AWSDynamoDBAttributeDefinitionsKey = attribute.Key("aws.dynamodb.attribute_definitions")

	
	
	
	
	
	
	
	
	AWSDynamoDBAttributesToGetKey = attribute.Key("aws.dynamodb.attributes_to_get")

	
	
	
	
	
	
	
	AWSDynamoDBConsistentReadKey = attribute.Key("aws.dynamodb.consistent_read")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBConsumedCapacityKey = attribute.Key("aws.dynamodb.consumed_capacity")

	
	
	
	
	
	
	
	
	AWSDynamoDBCountKey = attribute.Key("aws.dynamodb.count")

	
	
	
	
	
	
	
	
	AWSDynamoDBExclusiveStartTableKey = attribute.Key("aws.dynamodb.exclusive_start_table")

	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBGlobalSecondaryIndexUpdatesKey = attribute.Key("aws.dynamodb.global_secondary_index_updates")

	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBGlobalSecondaryIndexesKey = attribute.Key("aws.dynamodb.global_secondary_indexes")

	
	
	
	
	
	
	
	
	AWSDynamoDBIndexNameKey = attribute.Key("aws.dynamodb.index_name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBItemCollectionMetricsKey = attribute.Key("aws.dynamodb.item_collection_metrics")

	
	
	
	
	
	
	
	
	AWSDynamoDBLimitKey = attribute.Key("aws.dynamodb.limit")

	
	
	
	
	
	
	
	
	
	
	
	
	AWSDynamoDBLocalSecondaryIndexesKey = attribute.Key("aws.dynamodb.local_secondary_indexes")

	
	
	
	
	
	
	
	
	
	AWSDynamoDBProjectionKey = attribute.Key("aws.dynamodb.projection")

	
	
	
	
	
	
	
	
	
	AWSDynamoDBProvisionedReadCapacityKey = attribute.Key("aws.dynamodb.provisioned_read_capacity")

	
	
	
	
	
	
	
	
	
	AWSDynamoDBProvisionedWriteCapacityKey = attribute.Key("aws.dynamodb.provisioned_write_capacity")

	
	
	
	
	
	
	
	AWSDynamoDBScanForwardKey = attribute.Key("aws.dynamodb.scan_forward")

	
	
	
	
	
	
	
	
	AWSDynamoDBScannedCountKey = attribute.Key("aws.dynamodb.scanned_count")

	
	
	
	
	
	
	
	
	AWSDynamoDBSegmentKey = attribute.Key("aws.dynamodb.segment")

	
	
	
	
	
	
	
	
	AWSDynamoDBSelectKey = attribute.Key("aws.dynamodb.select")

	
	
	
	
	
	
	
	
	AWSDynamoDBTableCountKey = attribute.Key("aws.dynamodb.table_count")

	
	
	
	
	
	
	
	
	AWSDynamoDBTableNamesKey = attribute.Key("aws.dynamodb.table_names")

	
	
	
	
	
	
	
	
	AWSDynamoDBTotalSegmentsKey = attribute.Key("aws.dynamodb.total_segments")
)





func AWSDynamoDBAttributeDefinitions(val ...string) attribute.KeyValue {
	return AWSDynamoDBAttributeDefinitionsKey.StringSlice(val)
}




func AWSDynamoDBAttributesToGet(val ...string) attribute.KeyValue {
	return AWSDynamoDBAttributesToGetKey.StringSlice(val)
}




func AWSDynamoDBConsistentRead(val bool) attribute.KeyValue {
	return AWSDynamoDBConsistentReadKey.Bool(val)
}




func AWSDynamoDBConsumedCapacity(val ...string) attribute.KeyValue {
	return AWSDynamoDBConsumedCapacityKey.StringSlice(val)
}




func AWSDynamoDBCount(val int) attribute.KeyValue {
	return AWSDynamoDBCountKey.Int(val)
}




func AWSDynamoDBExclusiveStartTable(val string) attribute.KeyValue {
	return AWSDynamoDBExclusiveStartTableKey.String(val)
}





func AWSDynamoDBGlobalSecondaryIndexUpdates(val ...string) attribute.KeyValue {
	return AWSDynamoDBGlobalSecondaryIndexUpdatesKey.StringSlice(val)
}





func AWSDynamoDBGlobalSecondaryIndexes(val ...string) attribute.KeyValue {
	return AWSDynamoDBGlobalSecondaryIndexesKey.StringSlice(val)
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





func AWSDynamoDBLocalSecondaryIndexes(val ...string) attribute.KeyValue {
	return AWSDynamoDBLocalSecondaryIndexesKey.StringSlice(val)
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




func AWSDynamoDBScanForward(val bool) attribute.KeyValue {
	return AWSDynamoDBScanForwardKey.Bool(val)
}




func AWSDynamoDBScannedCount(val int) attribute.KeyValue {
	return AWSDynamoDBScannedCountKey.Int(val)
}




func AWSDynamoDBSegment(val int) attribute.KeyValue {
	return AWSDynamoDBSegmentKey.Int(val)
}




func AWSDynamoDBSelect(val string) attribute.KeyValue {
	return AWSDynamoDBSelectKey.String(val)
}




func AWSDynamoDBTableCount(val int) attribute.KeyValue {
	return AWSDynamoDBTableCountKey.Int(val)
}




func AWSDynamoDBTableNames(val ...string) attribute.KeyValue {
	return AWSDynamoDBTableNamesKey.StringSlice(val)
}




func AWSDynamoDBTotalSegments(val int) attribute.KeyValue {
	return AWSDynamoDBTotalSegmentsKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	AWSECSTaskIDKey = attribute.Key("aws.ecs.task.id")

	
	
	
	
	
	
	
	
	
	AWSECSClusterARNKey = attribute.Key("aws.ecs.cluster.arn")

	
	
	
	
	
	
	
	
	
	
	AWSECSContainerARNKey = attribute.Key("aws.ecs.container.arn")

	
	
	
	
	
	
	
	
	AWSECSLaunchtypeKey = attribute.Key("aws.ecs.launchtype")

	
	
	
	
	
	
	
	
	
	
	
	AWSECSTaskARNKey = attribute.Key("aws.ecs.task.arn")

	
	
	
	
	
	
	
	
	
	
	AWSECSTaskFamilyKey = attribute.Key("aws.ecs.task.family")

	
	
	
	
	
	
	
	
	AWSECSTaskRevisionKey = attribute.Key("aws.ecs.task.revision")
)

var (
	
	AWSECSLaunchtypeEC2 = AWSECSLaunchtypeKey.String("ec2")
	
	AWSECSLaunchtypeFargate = AWSECSLaunchtypeKey.String("fargate")
)




func AWSECSTaskID(val string) attribute.KeyValue {
	return AWSECSTaskIDKey.String(val)
}




func AWSECSClusterARN(val string) attribute.KeyValue {
	return AWSECSClusterARNKey.String(val)
}





func AWSECSContainerARN(val string) attribute.KeyValue {
	return AWSECSContainerARNKey.String(val)
}





func AWSECSTaskARN(val string) attribute.KeyValue {
	return AWSECSTaskARNKey.String(val)
}






func AWSECSTaskFamily(val string) attribute.KeyValue {
	return AWSECSTaskFamilyKey.String(val)
}




func AWSECSTaskRevision(val string) attribute.KeyValue {
	return AWSECSTaskRevisionKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	AWSEKSClusterARNKey = attribute.Key("aws.eks.cluster.arn")
)




func AWSEKSClusterARN(val string) attribute.KeyValue {
	return AWSEKSClusterARNKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	AWSLogGroupARNsKey = attribute.Key("aws.log.group.arns")

	
	
	
	
	
	
	
	
	
	
	
	AWSLogGroupNamesKey = attribute.Key("aws.log.group.names")

	
	
	
	
	
	
	
	
	
	
	
	
	
	AWSLogStreamARNsKey = attribute.Key("aws.log.stream.arns")

	
	
	
	
	
	
	
	
	AWSLogStreamNamesKey = attribute.Key("aws.log.stream.names")
)




func AWSLogGroupARNs(val ...string) attribute.KeyValue {
	return AWSLogGroupARNsKey.StringSlice(val)
}




func AWSLogGroupNames(val ...string) attribute.KeyValue {
	return AWSLogGroupNamesKey.StringSlice(val)
}




func AWSLogStreamARNs(val ...string) attribute.KeyValue {
	return AWSLogStreamARNsKey.StringSlice(val)
}




func AWSLogStreamNames(val ...string) attribute.KeyValue {
	return AWSLogStreamNamesKey.StringSlice(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	AWSLambdaInvokedARNKey = attribute.Key("aws.lambda.invoked_arn")
)






func AWSLambdaInvokedARN(val string) attribute.KeyValue {
	return AWSLambdaInvokedARNKey.String(val)
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
	
	
	
	
	
	
	
	
	
	
	
	BrowserBrandsKey = attribute.Key("browser.brands")

	
	
	
	
	
	
	
	
	
	
	BrowserLanguageKey = attribute.Key("browser.language")

	
	
	
	
	
	
	
	
	
	
	
	BrowserMobileKey = attribute.Key("browser.mobile")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	BrowserPlatformKey = attribute.Key("browser.platform")
)




func BrowserBrands(val ...string) attribute.KeyValue {
	return BrowserBrandsKey.StringSlice(val)
}




func BrowserLanguage(val string) attribute.KeyValue {
	return BrowserLanguageKey.String(val)
}




func BrowserMobile(val bool) attribute.KeyValue {
	return BrowserMobileKey.Bool(val)
}




func BrowserPlatform(val string) attribute.KeyValue {
	return BrowserPlatformKey.String(val)
}









const (
	
	
	
	
	
	
	
	
	
	
	
	
	ClientAddressKey = attribute.Key("client.address")

	
	
	
	
	
	
	
	
	
	
	ClientPortKey = attribute.Key("client.port")
)





func ClientAddress(val string) attribute.KeyValue {
	return ClientAddressKey.String(val)
}



func ClientPort(val int) attribute.KeyValue {
	return ClientPortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	CloudAccountIDKey = attribute.Key("cloud.account.id")

	
	
	
	
	
	
	
	
	
	
	
	
	CloudAvailabilityZoneKey = attribute.Key("cloud.availability_zone")

	
	
	
	
	
	
	
	
	CloudPlatformKey = attribute.Key("cloud.platform")

	
	
	
	
	
	
	CloudProviderKey = attribute.Key("cloud.provider")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	CloudRegionKey = attribute.Key("cloud.region")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	CloudResourceIDKey = attribute.Key("cloud.resource_id")
)

var (
	
	CloudPlatformAlibabaCloudECS = CloudPlatformKey.String("alibaba_cloud_ecs")
	
	CloudPlatformAlibabaCloudFc = CloudPlatformKey.String("alibaba_cloud_fc")
	
	CloudPlatformAlibabaCloudOpenshift = CloudPlatformKey.String("alibaba_cloud_openshift")
	
	CloudPlatformAWSEC2 = CloudPlatformKey.String("aws_ec2")
	
	CloudPlatformAWSECS = CloudPlatformKey.String("aws_ecs")
	
	CloudPlatformAWSEKS = CloudPlatformKey.String("aws_eks")
	
	CloudPlatformAWSLambda = CloudPlatformKey.String("aws_lambda")
	
	CloudPlatformAWSElasticBeanstalk = CloudPlatformKey.String("aws_elastic_beanstalk")
	
	CloudPlatformAWSAppRunner = CloudPlatformKey.String("aws_app_runner")
	
	CloudPlatformAWSOpenshift = CloudPlatformKey.String("aws_openshift")
	
	CloudPlatformAzureVM = CloudPlatformKey.String("azure_vm")
	
	CloudPlatformAzureContainerApps = CloudPlatformKey.String("azure_container_apps")
	
	CloudPlatformAzureContainerInstances = CloudPlatformKey.String("azure_container_instances")
	
	CloudPlatformAzureAKS = CloudPlatformKey.String("azure_aks")
	
	CloudPlatformAzureFunctions = CloudPlatformKey.String("azure_functions")
	
	CloudPlatformAzureAppService = CloudPlatformKey.String("azure_app_service")
	
	CloudPlatformAzureOpenshift = CloudPlatformKey.String("azure_openshift")
	
	CloudPlatformGCPBareMetalSolution = CloudPlatformKey.String("gcp_bare_metal_solution")
	
	CloudPlatformGCPComputeEngine = CloudPlatformKey.String("gcp_compute_engine")
	
	CloudPlatformGCPCloudRun = CloudPlatformKey.String("gcp_cloud_run")
	
	CloudPlatformGCPKubernetesEngine = CloudPlatformKey.String("gcp_kubernetes_engine")
	
	CloudPlatformGCPCloudFunctions = CloudPlatformKey.String("gcp_cloud_functions")
	
	CloudPlatformGCPAppEngine = CloudPlatformKey.String("gcp_app_engine")
	
	CloudPlatformGCPOpenshift = CloudPlatformKey.String("gcp_openshift")
	
	CloudPlatformIbmCloudOpenshift = CloudPlatformKey.String("ibm_cloud_openshift")
	
	CloudPlatformTencentCloudCvm = CloudPlatformKey.String("tencent_cloud_cvm")
	
	CloudPlatformTencentCloudEKS = CloudPlatformKey.String("tencent_cloud_eks")
	
	CloudPlatformTencentCloudScf = CloudPlatformKey.String("tencent_cloud_scf")
)

var (
	
	CloudProviderAlibabaCloud = CloudProviderKey.String("alibaba_cloud")
	
	CloudProviderAWS = CloudProviderKey.String("aws")
	
	CloudProviderAzure = CloudProviderKey.String("azure")
	
	CloudProviderGCP = CloudProviderKey.String("gcp")
	
	CloudProviderHeroku = CloudProviderKey.String("heroku")
	
	CloudProviderIbmCloud = CloudProviderKey.String("ibm_cloud")
	
	CloudProviderTencentCloud = CloudProviderKey.String("tencent_cloud")
)




func CloudAccountID(val string) attribute.KeyValue {
	return CloudAccountIDKey.String(val)
}






func CloudAvailabilityZone(val string) attribute.KeyValue {
	return CloudAvailabilityZoneKey.String(val)
}




func CloudRegion(val string) attribute.KeyValue {
	return CloudRegionKey.String(val)
}










func CloudResourceID(val string) attribute.KeyValue {
	return CloudResourceIDKey.String(val)
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
	
	
	
	
	
	
	
	
	
	
	ContainerCommandKey = attribute.Key("container.command")

	
	
	
	
	
	
	
	
	
	ContainerCommandArgsKey = attribute.Key("container.command_args")

	
	
	
	
	
	
	
	
	
	ContainerCommandLineKey = attribute.Key("container.command_line")

	
	
	
	
	
	
	
	
	ContainerCPUStateKey = attribute.Key("container.cpu.state")

	
	
	
	
	
	
	
	
	
	
	ContainerIDKey = attribute.Key("container.id")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ContainerImageIDKey = attribute.Key("container.image.id")

	
	
	
	
	
	
	
	
	ContainerImageNameKey = attribute.Key("container.image.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ContainerImageRepoDigestsKey = attribute.Key("container.image.repo_digests")

	
	
	
	
	
	
	
	
	
	
	
	ContainerImageTagsKey = attribute.Key("container.image.tags")

	
	
	
	
	
	
	
	
	ContainerNameKey = attribute.Key("container.name")

	
	
	
	
	
	
	
	
	ContainerRuntimeKey = attribute.Key("container.runtime")
)

var (
	
	ContainerCPUStateUser = ContainerCPUStateKey.String("user")
	
	ContainerCPUStateSystem = ContainerCPUStateKey.String("system")
	
	ContainerCPUStateKernel = ContainerCPUStateKey.String("kernel")
)




func ContainerCommand(val string) attribute.KeyValue {
	return ContainerCommandKey.String(val)
}





func ContainerCommandArgs(val ...string) attribute.KeyValue {
	return ContainerCommandArgsKey.StringSlice(val)
}





func ContainerCommandLine(val string) attribute.KeyValue {
	return ContainerCommandLineKey.String(val)
}






func ContainerID(val string) attribute.KeyValue {
	return ContainerIDKey.String(val)
}




func ContainerImageID(val string) attribute.KeyValue {
	return ContainerImageIDKey.String(val)
}




func ContainerImageName(val string) attribute.KeyValue {
	return ContainerImageNameKey.String(val)
}




func ContainerImageRepoDigests(val ...string) attribute.KeyValue {
	return ContainerImageRepoDigestsKey.StringSlice(val)
}







func ContainerImageTags(val ...string) attribute.KeyValue {
	return ContainerImageTagsKey.StringSlice(val)
}




func ContainerName(val string) attribute.KeyValue {
	return ContainerNameKey.String(val)
}




func ContainerRuntime(val string) attribute.KeyValue {
	return ContainerRuntimeKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	DBClientConnectionsPoolNameKey = attribute.Key("db.client.connections.pool.name")

	
	
	
	
	
	
	
	
	DBClientConnectionsStateKey = attribute.Key("db.client.connections.state")

	
	
	
	
	
	
	
	
	
	
	
	
	
	DBCollectionNameKey = attribute.Key("db.collection.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DBNamespaceKey = attribute.Key("db.namespace")

	
	
	
	
	
	
	
	
	
	
	DBOperationNameKey = attribute.Key("db.operation.name")

	
	
	
	
	
	
	
	
	DBQueryTextKey = attribute.Key("db.query.text")

	
	
	
	
	
	
	
	
	
	
	
	DBSystemKey = attribute.Key("db.system")
)

var (
	
	DBClientConnectionsStateIdle = DBClientConnectionsStateKey.String("idle")
	
	DBClientConnectionsStateUsed = DBClientConnectionsStateKey.String("used")
)

var (
	
	DBSystemOtherSQL = DBSystemKey.String("other_sql")
	
	DBSystemMSSQL = DBSystemKey.String("mssql")
	
	DBSystemMssqlcompact = DBSystemKey.String("mssqlcompact")
	
	DBSystemMySQL = DBSystemKey.String("mysql")
	
	DBSystemOracle = DBSystemKey.String("oracle")
	
	DBSystemDB2 = DBSystemKey.String("db2")
	
	DBSystemPostgreSQL = DBSystemKey.String("postgresql")
	
	DBSystemRedshift = DBSystemKey.String("redshift")
	
	DBSystemHive = DBSystemKey.String("hive")
	
	DBSystemCloudscape = DBSystemKey.String("cloudscape")
	
	DBSystemHSQLDB = DBSystemKey.String("hsqldb")
	
	DBSystemProgress = DBSystemKey.String("progress")
	
	DBSystemMaxDB = DBSystemKey.String("maxdb")
	
	DBSystemHanaDB = DBSystemKey.String("hanadb")
	
	DBSystemIngres = DBSystemKey.String("ingres")
	
	DBSystemFirstSQL = DBSystemKey.String("firstsql")
	
	DBSystemEDB = DBSystemKey.String("edb")
	
	DBSystemCache = DBSystemKey.String("cache")
	
	DBSystemAdabas = DBSystemKey.String("adabas")
	
	DBSystemFirebird = DBSystemKey.String("firebird")
	
	DBSystemDerby = DBSystemKey.String("derby")
	
	DBSystemFilemaker = DBSystemKey.String("filemaker")
	
	DBSystemInformix = DBSystemKey.String("informix")
	
	DBSystemInstantDB = DBSystemKey.String("instantdb")
	
	DBSystemInterbase = DBSystemKey.String("interbase")
	
	DBSystemMariaDB = DBSystemKey.String("mariadb")
	
	DBSystemNetezza = DBSystemKey.String("netezza")
	
	DBSystemPervasive = DBSystemKey.String("pervasive")
	
	DBSystemPointbase = DBSystemKey.String("pointbase")
	
	DBSystemSqlite = DBSystemKey.String("sqlite")
	
	DBSystemSybase = DBSystemKey.String("sybase")
	
	DBSystemTeradata = DBSystemKey.String("teradata")
	
	DBSystemVertica = DBSystemKey.String("vertica")
	
	DBSystemH2 = DBSystemKey.String("h2")
	
	DBSystemColdfusion = DBSystemKey.String("coldfusion")
	
	DBSystemCassandra = DBSystemKey.String("cassandra")
	
	DBSystemHBase = DBSystemKey.String("hbase")
	
	DBSystemMongoDB = DBSystemKey.String("mongodb")
	
	DBSystemRedis = DBSystemKey.String("redis")
	
	DBSystemCouchbase = DBSystemKey.String("couchbase")
	
	DBSystemCouchDB = DBSystemKey.String("couchdb")
	
	DBSystemCosmosDB = DBSystemKey.String("cosmosdb")
	
	DBSystemDynamoDB = DBSystemKey.String("dynamodb")
	
	DBSystemNeo4j = DBSystemKey.String("neo4j")
	
	DBSystemGeode = DBSystemKey.String("geode")
	
	DBSystemElasticsearch = DBSystemKey.String("elasticsearch")
	
	DBSystemMemcached = DBSystemKey.String("memcached")
	
	DBSystemCockroachdb = DBSystemKey.String("cockroachdb")
	
	DBSystemOpensearch = DBSystemKey.String("opensearch")
	
	DBSystemClickhouse = DBSystemKey.String("clickhouse")
	
	DBSystemSpanner = DBSystemKey.String("spanner")
	
	DBSystemTrino = DBSystemKey.String("trino")
)







func DBClientConnectionsPoolName(val string) attribute.KeyValue {
	return DBClientConnectionsPoolNameKey.String(val)
}




func DBCollectionName(val string) attribute.KeyValue {
	return DBCollectionNameKey.String(val)
}




func DBNamespace(val string) attribute.KeyValue {
	return DBNamespaceKey.String(val)
}




func DBOperationName(val string) attribute.KeyValue {
	return DBOperationNameKey.String(val)
}




func DBQueryText(val string) attribute.KeyValue {
	return DBQueryTextKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	DBCassandraConsistencyLevelKey = attribute.Key("db.cassandra.consistency_level")

	
	
	
	
	
	
	
	
	DBCassandraCoordinatorDCKey = attribute.Key("db.cassandra.coordinator.dc")

	
	
	
	
	
	
	
	
	DBCassandraCoordinatorIDKey = attribute.Key("db.cassandra.coordinator.id")

	
	
	
	
	
	
	
	DBCassandraIdempotenceKey = attribute.Key("db.cassandra.idempotence")

	
	
	
	
	
	
	
	
	DBCassandraPageSizeKey = attribute.Key("db.cassandra.page_size")

	
	
	
	
	
	
	
	
	
	DBCassandraSpeculativeExecutionCountKey = attribute.Key("db.cassandra.speculative_execution_count")
)

var (
	
	DBCassandraConsistencyLevelAll = DBCassandraConsistencyLevelKey.String("all")
	
	DBCassandraConsistencyLevelEachQuorum = DBCassandraConsistencyLevelKey.String("each_quorum")
	
	DBCassandraConsistencyLevelQuorum = DBCassandraConsistencyLevelKey.String("quorum")
	
	DBCassandraConsistencyLevelLocalQuorum = DBCassandraConsistencyLevelKey.String("local_quorum")
	
	DBCassandraConsistencyLevelOne = DBCassandraConsistencyLevelKey.String("one")
	
	DBCassandraConsistencyLevelTwo = DBCassandraConsistencyLevelKey.String("two")
	
	DBCassandraConsistencyLevelThree = DBCassandraConsistencyLevelKey.String("three")
	
	DBCassandraConsistencyLevelLocalOne = DBCassandraConsistencyLevelKey.String("local_one")
	
	DBCassandraConsistencyLevelAny = DBCassandraConsistencyLevelKey.String("any")
	
	DBCassandraConsistencyLevelSerial = DBCassandraConsistencyLevelKey.String("serial")
	
	DBCassandraConsistencyLevelLocalSerial = DBCassandraConsistencyLevelKey.String("local_serial")
)




func DBCassandraCoordinatorDC(val string) attribute.KeyValue {
	return DBCassandraCoordinatorDCKey.String(val)
}




func DBCassandraCoordinatorID(val string) attribute.KeyValue {
	return DBCassandraCoordinatorIDKey.String(val)
}




func DBCassandraIdempotence(val bool) attribute.KeyValue {
	return DBCassandraIdempotenceKey.Bool(val)
}




func DBCassandraPageSize(val int) attribute.KeyValue {
	return DBCassandraPageSizeKey.Int(val)
}





func DBCassandraSpeculativeExecutionCount(val int) attribute.KeyValue {
	return DBCassandraSpeculativeExecutionCountKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	DBCosmosDBClientIDKey = attribute.Key("db.cosmosdb.client_id")

	
	
	
	
	
	
	
	DBCosmosDBConnectionModeKey = attribute.Key("db.cosmosdb.connection_mode")

	
	
	
	
	
	
	
	DBCosmosDBOperationTypeKey = attribute.Key("db.cosmosdb.operation_type")

	
	
	
	
	
	
	
	
	DBCosmosDBRequestChargeKey = attribute.Key("db.cosmosdb.request_charge")

	
	
	
	
	
	
	
	DBCosmosDBRequestContentLengthKey = attribute.Key("db.cosmosdb.request_content_length")

	
	
	
	
	
	
	
	
	DBCosmosDBStatusCodeKey = attribute.Key("db.cosmosdb.status_code")

	
	
	
	
	
	
	
	
	DBCosmosDBSubStatusCodeKey = attribute.Key("db.cosmosdb.sub_status_code")
)

var (
	
	DBCosmosDBConnectionModeGateway = DBCosmosDBConnectionModeKey.String("gateway")
	
	DBCosmosDBConnectionModeDirect = DBCosmosDBConnectionModeKey.String("direct")
)

var (
	
	DBCosmosDBOperationTypeInvalid = DBCosmosDBOperationTypeKey.String("Invalid")
	
	DBCosmosDBOperationTypeCreate = DBCosmosDBOperationTypeKey.String("Create")
	
	DBCosmosDBOperationTypePatch = DBCosmosDBOperationTypeKey.String("Patch")
	
	DBCosmosDBOperationTypeRead = DBCosmosDBOperationTypeKey.String("Read")
	
	DBCosmosDBOperationTypeReadFeed = DBCosmosDBOperationTypeKey.String("ReadFeed")
	
	DBCosmosDBOperationTypeDelete = DBCosmosDBOperationTypeKey.String("Delete")
	
	DBCosmosDBOperationTypeReplace = DBCosmosDBOperationTypeKey.String("Replace")
	
	DBCosmosDBOperationTypeExecute = DBCosmosDBOperationTypeKey.String("Execute")
	
	DBCosmosDBOperationTypeQuery = DBCosmosDBOperationTypeKey.String("Query")
	
	DBCosmosDBOperationTypeHead = DBCosmosDBOperationTypeKey.String("Head")
	
	DBCosmosDBOperationTypeHeadFeed = DBCosmosDBOperationTypeKey.String("HeadFeed")
	
	DBCosmosDBOperationTypeUpsert = DBCosmosDBOperationTypeKey.String("Upsert")
	
	DBCosmosDBOperationTypeBatch = DBCosmosDBOperationTypeKey.String("Batch")
	
	DBCosmosDBOperationTypeQueryPlan = DBCosmosDBOperationTypeKey.String("QueryPlan")
	
	DBCosmosDBOperationTypeExecuteJavascript = DBCosmosDBOperationTypeKey.String("ExecuteJavaScript")
)




func DBCosmosDBClientID(val string) attribute.KeyValue {
	return DBCosmosDBClientIDKey.String(val)
}




func DBCosmosDBRequestCharge(val float64) attribute.KeyValue {
	return DBCosmosDBRequestChargeKey.Float64(val)
}




func DBCosmosDBRequestContentLength(val int) attribute.KeyValue {
	return DBCosmosDBRequestContentLengthKey.Int(val)
}




func DBCosmosDBStatusCode(val int) attribute.KeyValue {
	return DBCosmosDBStatusCodeKey.Int(val)
}




func DBCosmosDBSubStatusCode(val int) attribute.KeyValue {
	return DBCosmosDBSubStatusCodeKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	DBElasticsearchClusterNameKey = attribute.Key("db.elasticsearch.cluster.name")

	
	
	
	
	
	
	
	
	
	DBElasticsearchNodeNameKey = attribute.Key("db.elasticsearch.node.name")
)




func DBElasticsearchClusterName(val string) attribute.KeyValue {
	return DBElasticsearchClusterNameKey.String(val)
}





func DBElasticsearchNodeName(val string) attribute.KeyValue {
	return DBElasticsearchNodeNameKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DeploymentEnvironmentKey = attribute.Key("deployment.environment")
)





func DeploymentEnvironment(val string) attribute.KeyValue {
	return DeploymentEnvironmentKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	AndroidStateKey = attribute.Key("android.state")
)

var (
	
	AndroidStateCreated = AndroidStateKey.String("created")
	
	AndroidStateBackground = AndroidStateKey.String("background")
	
	AndroidStateForeground = AndroidStateKey.String("foreground")
)









const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	DestinationAddressKey = attribute.Key("destination.address")

	
	
	
	
	
	
	
	
	DestinationPortKey = attribute.Key("destination.port")
)





func DestinationAddress(val string) attribute.KeyValue {
	return DestinationAddressKey.String(val)
}




func DestinationPort(val int) attribute.KeyValue {
	return DestinationPortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DeviceIDKey = attribute.Key("device.id")

	
	
	
	
	
	
	
	
	
	
	
	DeviceManufacturerKey = attribute.Key("device.manufacturer")

	
	
	
	
	
	
	
	
	
	
	
	DeviceModelIdentifierKey = attribute.Key("device.model.identifier")

	
	
	
	
	
	
	
	
	
	
	DeviceModelNameKey = attribute.Key("device.model.name")
)




func DeviceID(val string) attribute.KeyValue {
	return DeviceIDKey.String(val)
}




func DeviceManufacturer(val string) attribute.KeyValue {
	return DeviceManufacturerKey.String(val)
}




func DeviceModelIdentifier(val string) attribute.KeyValue {
	return DeviceModelIdentifierKey.String(val)
}




func DeviceModelName(val string) attribute.KeyValue {
	return DeviceModelNameKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	DiskIoDirectionKey = attribute.Key("disk.io.direction")
)

var (
	
	DiskIoDirectionRead = DiskIoDirectionKey.String("read")
	
	DiskIoDirectionWrite = DiskIoDirectionKey.String("write")
)


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	DNSQuestionNameKey = attribute.Key("dns.question.name")
)




func DNSQuestionName(val string) attribute.KeyValue {
	return DNSQuestionNameKey.String(val)
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
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ErrorTypeKey = attribute.Key("error.type")
)

var (
	
	ErrorTypeOther = ErrorTypeKey.String("_OTHER")
)


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	EventNameKey = attribute.Key("event.name")
)




func EventName(val string) attribute.KeyValue {
	return EventNameKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ExceptionEscapedKey = attribute.Key("exception.escaped")

	
	
	
	
	
	
	
	
	
	ExceptionMessageKey = attribute.Key("exception.message")

	
	
	
	
	
	
	
	
	
	
	
	
	
	ExceptionStacktraceKey = attribute.Key("exception.stacktrace")

	
	
	
	
	
	
	
	
	
	
	ExceptionTypeKey = attribute.Key("exception.type")
)





func ExceptionEscaped(val bool) attribute.KeyValue {
	return ExceptionEscapedKey.Bool(val)
}




func ExceptionMessage(val string) attribute.KeyValue {
	return ExceptionMessageKey.String(val)
}





func ExceptionStacktrace(val string) attribute.KeyValue {
	return ExceptionStacktraceKey.String(val)
}






func ExceptionType(val string) attribute.KeyValue {
	return ExceptionTypeKey.String(val)
}


const (
	
	
	
	
	
	
	
	FaaSColdstartKey = attribute.Key("faas.coldstart")

	
	
	
	
	
	
	
	
	
	FaaSCronKey = attribute.Key("faas.cron")

	
	
	
	
	
	
	
	
	
	
	FaaSDocumentCollectionKey = attribute.Key("faas.document.collection")

	
	
	
	
	
	
	
	
	
	FaaSDocumentNameKey = attribute.Key("faas.document.name")

	
	
	
	
	
	
	
	FaaSDocumentOperationKey = attribute.Key("faas.document.operation")

	
	
	
	
	
	
	
	
	
	
	FaaSDocumentTimeKey = attribute.Key("faas.document.time")

	
	
	
	
	
	
	
	
	
	
	FaaSInstanceKey = attribute.Key("faas.instance")

	
	
	
	
	
	
	
	
	FaaSInvocationIDKey = attribute.Key("faas.invocation_id")

	
	
	
	
	
	
	
	
	
	
	FaaSInvokedNameKey = attribute.Key("faas.invoked_name")

	
	
	
	
	
	
	
	
	
	FaaSInvokedProviderKey = attribute.Key("faas.invoked_provider")

	
	
	
	
	
	
	
	
	
	
	FaaSInvokedRegionKey = attribute.Key("faas.invoked_region")

	
	
	
	
	
	
	
	
	
	
	
	
	
	FaaSMaxMemoryKey = attribute.Key("faas.max_memory")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	FaaSNameKey = attribute.Key("faas.name")

	
	
	
	
	
	
	
	
	
	
	FaaSTimeKey = attribute.Key("faas.time")

	
	
	
	
	
	
	
	FaaSTriggerKey = attribute.Key("faas.trigger")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	FaaSVersionKey = attribute.Key("faas.version")
)

var (
	
	FaaSDocumentOperationInsert = FaaSDocumentOperationKey.String("insert")
	
	FaaSDocumentOperationEdit = FaaSDocumentOperationKey.String("edit")
	
	FaaSDocumentOperationDelete = FaaSDocumentOperationKey.String("delete")
)

var (
	
	FaaSInvokedProviderAlibabaCloud = FaaSInvokedProviderKey.String("alibaba_cloud")
	
	FaaSInvokedProviderAWS = FaaSInvokedProviderKey.String("aws")
	
	FaaSInvokedProviderAzure = FaaSInvokedProviderKey.String("azure")
	
	FaaSInvokedProviderGCP = FaaSInvokedProviderKey.String("gcp")
	
	FaaSInvokedProviderTencentCloud = FaaSInvokedProviderKey.String("tencent_cloud")
)

var (
	
	FaaSTriggerDatasource = FaaSTriggerKey.String("datasource")
	
	FaaSTriggerHTTP = FaaSTriggerKey.String("http")
	
	FaaSTriggerPubsub = FaaSTriggerKey.String("pubsub")
	
	FaaSTriggerTimer = FaaSTriggerKey.String("timer")
	
	FaaSTriggerOther = FaaSTriggerKey.String("other")
)




func FaaSColdstart(val bool) attribute.KeyValue {
	return FaaSColdstartKey.Bool(val)
}





func FaaSCron(val string) attribute.KeyValue {
	return FaaSCronKey.String(val)
}






func FaaSDocumentCollection(val string) attribute.KeyValue {
	return FaaSDocumentCollectionKey.String(val)
}





func FaaSDocumentName(val string) attribute.KeyValue {
	return FaaSDocumentNameKey.String(val)
}






func FaaSDocumentTime(val string) attribute.KeyValue {
	return FaaSDocumentTimeKey.String(val)
}





func FaaSInstance(val string) attribute.KeyValue {
	return FaaSInstanceKey.String(val)
}




func FaaSInvocationID(val string) attribute.KeyValue {
	return FaaSInvocationIDKey.String(val)
}




func FaaSInvokedName(val string) attribute.KeyValue {
	return FaaSInvokedNameKey.String(val)
}




func FaaSInvokedRegion(val string) attribute.KeyValue {
	return FaaSInvokedRegionKey.String(val)
}




func FaaSMaxMemory(val int) attribute.KeyValue {
	return FaaSMaxMemoryKey.Int(val)
}




func FaaSName(val string) attribute.KeyValue {
	return FaaSNameKey.String(val)
}






func FaaSTime(val string) attribute.KeyValue {
	return FaaSTimeKey.String(val)
}




func FaaSVersion(val string) attribute.KeyValue {
	return FaaSVersionKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	FeatureFlagKeyKey = attribute.Key("feature_flag.key")

	
	
	
	
	
	
	
	
	FeatureFlagProviderNameKey = attribute.Key("feature_flag.provider_name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	FeatureFlagVariantKey = attribute.Key("feature_flag.variant")
)




func FeatureFlagKey(val string) attribute.KeyValue {
	return FeatureFlagKeyKey.String(val)
}




func FeatureFlagProviderName(val string) attribute.KeyValue {
	return FeatureFlagProviderNameKey.String(val)
}





func FeatureFlagVariant(val string) attribute.KeyValue {
	return FeatureFlagVariantKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	FileDirectoryKey = attribute.Key("file.directory")

	
	
	
	
	
	
	
	
	
	
	FileExtensionKey = attribute.Key("file.extension")

	
	
	
	
	
	
	
	
	FileNameKey = attribute.Key("file.name")

	
	
	
	
	
	
	
	
	
	FilePathKey = attribute.Key("file.path")

	
	
	
	
	
	
	FileSizeKey = attribute.Key("file.size")
)




func FileDirectory(val string) attribute.KeyValue {
	return FileDirectoryKey.String(val)
}




func FileExtension(val string) attribute.KeyValue {
	return FileExtensionKey.String(val)
}




func FileName(val string) attribute.KeyValue {
	return FileNameKey.String(val)
}




func FilePath(val string) attribute.KeyValue {
	return FilePathKey.String(val)
}



func FileSize(val int) attribute.KeyValue {
	return FileSizeKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	GCPCloudRunJobExecutionKey = attribute.Key("gcp.cloud_run.job.execution")

	
	
	
	
	
	
	
	
	
	
	GCPCloudRunJobTaskIndexKey = attribute.Key("gcp.cloud_run.job.task_index")
)








func GCPCloudRunJobExecution(val string) attribute.KeyValue {
	return GCPCloudRunJobExecutionKey.String(val)
}






func GCPCloudRunJobTaskIndex(val int) attribute.KeyValue {
	return GCPCloudRunJobTaskIndexKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	GCPGceInstanceHostnameKey = attribute.Key("gcp.gce.instance.hostname")

	
	
	
	
	
	
	
	
	
	
	
	
	GCPGceInstanceNameKey = attribute.Key("gcp.gce.instance.name")
)





func GCPGceInstanceHostname(val string) attribute.KeyValue {
	return GCPGceInstanceHostnameKey.String(val)
}







func GCPGceInstanceName(val string) attribute.KeyValue {
	return GCPGceInstanceNameKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	
	GenAiCompletionKey = attribute.Key("gen_ai.completion")

	
	
	
	
	
	
	
	
	
	
	
	GenAiPromptKey = attribute.Key("gen_ai.prompt")

	
	
	
	
	
	
	
	
	GenAiRequestMaxTokensKey = attribute.Key("gen_ai.request.max_tokens")

	
	
	
	
	
	
	
	
	GenAiRequestModelKey = attribute.Key("gen_ai.request.model")

	
	
	
	
	
	
	
	
	GenAiRequestTemperatureKey = attribute.Key("gen_ai.request.temperature")

	
	
	
	
	
	
	
	
	GenAiRequestTopPKey = attribute.Key("gen_ai.request.top_p")

	
	
	
	
	
	
	
	
	
	GenAiResponseFinishReasonsKey = attribute.Key("gen_ai.response.finish_reasons")

	
	
	
	
	
	
	
	
	GenAiResponseIDKey = attribute.Key("gen_ai.response.id")

	
	
	
	
	
	
	
	
	GenAiResponseModelKey = attribute.Key("gen_ai.response.model")

	
	
	
	
	
	
	
	
	
	
	
	
	GenAiSystemKey = attribute.Key("gen_ai.system")

	
	
	
	
	
	
	
	
	GenAiUsageCompletionTokensKey = attribute.Key("gen_ai.usage.completion_tokens")

	
	
	
	
	
	
	
	
	GenAiUsagePromptTokensKey = attribute.Key("gen_ai.usage.prompt_tokens")
)

var (
	
	GenAiSystemOpenai = GenAiSystemKey.String("openai")
)




func GenAiCompletion(val string) attribute.KeyValue {
	return GenAiCompletionKey.String(val)
}




func GenAiPrompt(val string) attribute.KeyValue {
	return GenAiPromptKey.String(val)
}




func GenAiRequestMaxTokens(val int) attribute.KeyValue {
	return GenAiRequestMaxTokensKey.Int(val)
}




func GenAiRequestModel(val string) attribute.KeyValue {
	return GenAiRequestModelKey.String(val)
}




func GenAiRequestTemperature(val float64) attribute.KeyValue {
	return GenAiRequestTemperatureKey.Float64(val)
}




func GenAiRequestTopP(val float64) attribute.KeyValue {
	return GenAiRequestTopPKey.Float64(val)
}





func GenAiResponseFinishReasons(val ...string) attribute.KeyValue {
	return GenAiResponseFinishReasonsKey.StringSlice(val)
}




func GenAiResponseID(val string) attribute.KeyValue {
	return GenAiResponseIDKey.String(val)
}




func GenAiResponseModel(val string) attribute.KeyValue {
	return GenAiResponseModelKey.String(val)
}




func GenAiUsageCompletionTokens(val int) attribute.KeyValue {
	return GenAiUsageCompletionTokensKey.Int(val)
}




func GenAiUsagePromptTokens(val int) attribute.KeyValue {
	return GenAiUsagePromptTokensKey.Int(val)
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



const (
	
	
	
	
	
	
	
	
	HerokuAppIDKey = attribute.Key("heroku.app.id")

	
	
	
	
	
	
	
	
	HerokuReleaseCommitKey = attribute.Key("heroku.release.commit")

	
	
	
	
	
	
	
	
	HerokuReleaseCreationTimestampKey = attribute.Key("heroku.release.creation_timestamp")
)




func HerokuAppID(val string) attribute.KeyValue {
	return HerokuAppIDKey.String(val)
}




func HerokuReleaseCommit(val string) attribute.KeyValue {
	return HerokuReleaseCommitKey.String(val)
}




func HerokuReleaseCreationTimestamp(val string) attribute.KeyValue {
	return HerokuReleaseCreationTimestampKey.String(val)
}



const (
	
	
	
	
	
	
	
	HostArchKey = attribute.Key("host.arch")

	
	
	
	
	
	
	
	
	HostCPUCacheL2SizeKey = attribute.Key("host.cpu.cache.l2.size")

	
	
	
	
	
	
	
	
	HostCPUFamilyKey = attribute.Key("host.cpu.family")

	
	
	
	
	
	
	
	
	
	HostCPUModelIDKey = attribute.Key("host.cpu.model.id")

	
	
	
	
	
	
	
	
	HostCPUModelNameKey = attribute.Key("host.cpu.model.name")

	
	
	
	
	
	
	
	
	HostCPUSteppingKey = attribute.Key("host.cpu.stepping")

	
	
	
	
	
	
	
	
	
	
	
	HostCPUVendorIDKey = attribute.Key("host.cpu.vendor.id")

	
	
	
	
	
	
	
	
	
	
	HostIDKey = attribute.Key("host.id")

	
	
	
	
	
	
	
	
	HostImageIDKey = attribute.Key("host.image.id")

	
	
	
	
	
	
	
	
	HostImageNameKey = attribute.Key("host.image.name")

	
	
	
	
	
	
	
	
	
	HostImageVersionKey = attribute.Key("host.image.version")

	
	
	
	
	
	
	
	
	
	
	
	HostIPKey = attribute.Key("host.ip")

	
	
	
	
	
	
	
	
	
	
	
	
	HostMacKey = attribute.Key("host.mac")

	
	
	
	
	
	
	
	
	
	HostNameKey = attribute.Key("host.name")

	
	
	
	
	
	
	
	
	HostTypeKey = attribute.Key("host.type")
)

var (
	
	HostArchAMD64 = HostArchKey.String("amd64")
	
	HostArchARM32 = HostArchKey.String("arm32")
	
	HostArchARM64 = HostArchKey.String("arm64")
	
	HostArchIA64 = HostArchKey.String("ia64")
	
	HostArchPPC32 = HostArchKey.String("ppc32")
	
	HostArchPPC64 = HostArchKey.String("ppc64")
	
	HostArchS390x = HostArchKey.String("s390x")
	
	HostArchX86 = HostArchKey.String("x86")
)




func HostCPUCacheL2Size(val int) attribute.KeyValue {
	return HostCPUCacheL2SizeKey.Int(val)
}




func HostCPUFamily(val string) attribute.KeyValue {
	return HostCPUFamilyKey.String(val)
}





func HostCPUModelID(val string) attribute.KeyValue {
	return HostCPUModelIDKey.String(val)
}




func HostCPUModelName(val string) attribute.KeyValue {
	return HostCPUModelNameKey.String(val)
}




func HostCPUStepping(val string) attribute.KeyValue {
	return HostCPUSteppingKey.String(val)
}




func HostCPUVendorID(val string) attribute.KeyValue {
	return HostCPUVendorIDKey.String(val)
}






func HostID(val string) attribute.KeyValue {
	return HostIDKey.String(val)
}




func HostImageID(val string) attribute.KeyValue {
	return HostImageIDKey.String(val)
}




func HostImageName(val string) attribute.KeyValue {
	return HostImageNameKey.String(val)
}





func HostImageVersion(val string) attribute.KeyValue {
	return HostImageVersionKey.String(val)
}




func HostIP(val ...string) attribute.KeyValue {
	return HostIPKey.StringSlice(val)
}




func HostMac(val ...string) attribute.KeyValue {
	return HostMacKey.StringSlice(val)
}





func HostName(val string) attribute.KeyValue {
	return HostNameKey.String(val)
}




func HostType(val string) attribute.KeyValue {
	return HostTypeKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	HTTPConnectionStateKey = attribute.Key("http.connection.state")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRequestBodySizeKey = attribute.Key("http.request.body.size")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRequestMethodKey = attribute.Key("http.request.method")

	
	
	
	
	
	
	
	
	HTTPRequestMethodOriginalKey = attribute.Key("http.request.method_original")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRequestResendCountKey = attribute.Key("http.request.resend_count")

	
	
	
	
	
	
	
	
	
	
	HTTPRequestSizeKey = attribute.Key("http.request.size")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPResponseBodySizeKey = attribute.Key("http.response.body.size")

	
	
	
	
	
	
	
	
	
	
	HTTPResponseSizeKey = attribute.Key("http.response.size")

	
	
	
	
	
	
	
	
	
	HTTPResponseStatusCodeKey = attribute.Key("http.response.status_code")

	
	
	
	
	
	
	
	
	
	
	
	
	
	HTTPRouteKey = attribute.Key("http.route")
)

var (
	
	HTTPConnectionStateActive = HTTPConnectionStateKey.String("active")
	
	HTTPConnectionStateIdle = HTTPConnectionStateKey.String("idle")
)

var (
	
	HTTPRequestMethodConnect = HTTPRequestMethodKey.String("CONNECT")
	
	HTTPRequestMethodDelete = HTTPRequestMethodKey.String("DELETE")
	
	HTTPRequestMethodGet = HTTPRequestMethodKey.String("GET")
	
	HTTPRequestMethodHead = HTTPRequestMethodKey.String("HEAD")
	
	HTTPRequestMethodOptions = HTTPRequestMethodKey.String("OPTIONS")
	
	HTTPRequestMethodPatch = HTTPRequestMethodKey.String("PATCH")
	
	HTTPRequestMethodPost = HTTPRequestMethodKey.String("POST")
	
	HTTPRequestMethodPut = HTTPRequestMethodKey.String("PUT")
	
	HTTPRequestMethodTrace = HTTPRequestMethodKey.String("TRACE")
	
	HTTPRequestMethodOther = HTTPRequestMethodKey.String("_OTHER")
)








func HTTPRequestBodySize(val int) attribute.KeyValue {
	return HTTPRequestBodySizeKey.Int(val)
}




func HTTPRequestMethodOriginal(val string) attribute.KeyValue {
	return HTTPRequestMethodOriginalKey.String(val)
}




func HTTPRequestResendCount(val int) attribute.KeyValue {
	return HTTPRequestResendCountKey.Int(val)
}






func HTTPRequestSize(val int) attribute.KeyValue {
	return HTTPRequestSizeKey.Int(val)
}








func HTTPResponseBodySize(val int) attribute.KeyValue {
	return HTTPResponseBodySizeKey.Int(val)
}






func HTTPResponseSize(val int) attribute.KeyValue {
	return HTTPResponseSizeKey.Int(val)
}




func HTTPResponseStatusCode(val int) attribute.KeyValue {
	return HTTPResponseStatusCodeKey.Int(val)
}




func HTTPRoute(val string) attribute.KeyValue {
	return HTTPRouteKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	JvmBufferPoolNameKey = attribute.Key("jvm.buffer.pool.name")

	
	
	
	
	
	
	
	
	
	
	JvmGcActionKey = attribute.Key("jvm.gc.action")

	
	
	
	
	
	
	
	
	
	JvmGcNameKey = attribute.Key("jvm.gc.name")

	
	
	
	
	
	
	
	
	
	
	JvmMemoryPoolNameKey = attribute.Key("jvm.memory.pool.name")

	
	
	
	
	
	
	
	
	JvmMemoryTypeKey = attribute.Key("jvm.memory.type")

	
	
	
	
	
	
	
	JvmThreadDaemonKey = attribute.Key("jvm.thread.daemon")

	
	
	
	
	
	
	
	
	JvmThreadStateKey = attribute.Key("jvm.thread.state")
)

var (
	
	JvmMemoryTypeHeap = JvmMemoryTypeKey.String("heap")
	
	JvmMemoryTypeNonHeap = JvmMemoryTypeKey.String("non_heap")
)

var (
	
	JvmThreadStateNew = JvmThreadStateKey.String("new")
	
	JvmThreadStateRunnable = JvmThreadStateKey.String("runnable")
	
	JvmThreadStateBlocked = JvmThreadStateKey.String("blocked")
	
	JvmThreadStateWaiting = JvmThreadStateKey.String("waiting")
	
	JvmThreadStateTimedWaiting = JvmThreadStateKey.String("timed_waiting")
	
	JvmThreadStateTerminated = JvmThreadStateKey.String("terminated")
)




func JvmBufferPoolName(val string) attribute.KeyValue {
	return JvmBufferPoolNameKey.String(val)
}




func JvmGcAction(val string) attribute.KeyValue {
	return JvmGcActionKey.String(val)
}



func JvmGcName(val string) attribute.KeyValue {
	return JvmGcNameKey.String(val)
}




func JvmMemoryPoolName(val string) attribute.KeyValue {
	return JvmMemoryPoolNameKey.String(val)
}




func JvmThreadDaemon(val bool) attribute.KeyValue {
	return JvmThreadDaemonKey.Bool(val)
}


const (
	
	
	
	
	
	
	
	
	K8SClusterNameKey = attribute.Key("k8s.cluster.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	K8SClusterUIDKey = attribute.Key("k8s.cluster.uid")

	
	
	
	
	
	
	
	
	
	K8SContainerNameKey = attribute.Key("k8s.container.name")

	
	
	
	
	
	
	
	
	
	K8SContainerRestartCountKey = attribute.Key("k8s.container.restart_count")

	
	
	
	
	
	
	
	
	K8SContainerStatusLastTerminatedReasonKey = attribute.Key("k8s.container.status.last_terminated_reason")

	
	
	
	
	
	
	
	
	K8SCronJobNameKey = attribute.Key("k8s.cronjob.name")

	
	
	
	
	
	
	
	
	K8SCronJobUIDKey = attribute.Key("k8s.cronjob.uid")

	
	
	
	
	
	
	
	
	K8SDaemonSetNameKey = attribute.Key("k8s.daemonset.name")

	
	
	
	
	
	
	
	
	K8SDaemonSetUIDKey = attribute.Key("k8s.daemonset.uid")

	
	
	
	
	
	
	
	
	K8SDeploymentNameKey = attribute.Key("k8s.deployment.name")

	
	
	
	
	
	
	
	
	K8SDeploymentUIDKey = attribute.Key("k8s.deployment.uid")

	
	
	
	
	
	
	
	K8SJobNameKey = attribute.Key("k8s.job.name")

	
	
	
	
	
	
	
	K8SJobUIDKey = attribute.Key("k8s.job.uid")

	
	
	
	
	
	
	
	
	K8SNamespaceNameKey = attribute.Key("k8s.namespace.name")

	
	
	
	
	
	
	
	K8SNodeNameKey = attribute.Key("k8s.node.name")

	
	
	
	
	
	
	
	K8SNodeUIDKey = attribute.Key("k8s.node.uid")

	
	
	
	
	
	
	
	K8SPodNameKey = attribute.Key("k8s.pod.name")

	
	
	
	
	
	
	
	K8SPodUIDKey = attribute.Key("k8s.pod.uid")

	
	
	
	
	
	
	
	
	K8SReplicaSetNameKey = attribute.Key("k8s.replicaset.name")

	
	
	
	
	
	
	
	
	K8SReplicaSetUIDKey = attribute.Key("k8s.replicaset.uid")

	
	
	
	
	
	
	
	
	K8SStatefulSetNameKey = attribute.Key("k8s.statefulset.name")

	
	
	
	
	
	
	
	
	K8SStatefulSetUIDKey = attribute.Key("k8s.statefulset.uid")
)




func K8SClusterName(val string) attribute.KeyValue {
	return K8SClusterNameKey.String(val)
}




func K8SClusterUID(val string) attribute.KeyValue {
	return K8SClusterUIDKey.String(val)
}





func K8SContainerName(val string) attribute.KeyValue {
	return K8SContainerNameKey.String(val)
}





func K8SContainerRestartCount(val int) attribute.KeyValue {
	return K8SContainerRestartCountKey.Int(val)
}




func K8SContainerStatusLastTerminatedReason(val string) attribute.KeyValue {
	return K8SContainerStatusLastTerminatedReasonKey.String(val)
}




func K8SCronJobName(val string) attribute.KeyValue {
	return K8SCronJobNameKey.String(val)
}




func K8SCronJobUID(val string) attribute.KeyValue {
	return K8SCronJobUIDKey.String(val)
}




func K8SDaemonSetName(val string) attribute.KeyValue {
	return K8SDaemonSetNameKey.String(val)
}




func K8SDaemonSetUID(val string) attribute.KeyValue {
	return K8SDaemonSetUIDKey.String(val)
}




func K8SDeploymentName(val string) attribute.KeyValue {
	return K8SDeploymentNameKey.String(val)
}




func K8SDeploymentUID(val string) attribute.KeyValue {
	return K8SDeploymentUIDKey.String(val)
}



func K8SJobName(val string) attribute.KeyValue {
	return K8SJobNameKey.String(val)
}



func K8SJobUID(val string) attribute.KeyValue {
	return K8SJobUIDKey.String(val)
}




func K8SNamespaceName(val string) attribute.KeyValue {
	return K8SNamespaceNameKey.String(val)
}



func K8SNodeName(val string) attribute.KeyValue {
	return K8SNodeNameKey.String(val)
}



func K8SNodeUID(val string) attribute.KeyValue {
	return K8SNodeUIDKey.String(val)
}



func K8SPodName(val string) attribute.KeyValue {
	return K8SPodNameKey.String(val)
}



func K8SPodUID(val string) attribute.KeyValue {
	return K8SPodUIDKey.String(val)
}




func K8SReplicaSetName(val string) attribute.KeyValue {
	return K8SReplicaSetNameKey.String(val)
}




func K8SReplicaSetUID(val string) attribute.KeyValue {
	return K8SReplicaSetUIDKey.String(val)
}




func K8SStatefulSetName(val string) attribute.KeyValue {
	return K8SStatefulSetNameKey.String(val)
}




func K8SStatefulSetUID(val string) attribute.KeyValue {
	return K8SStatefulSetUIDKey.String(val)
}


const (
	
	
	
	
	
	
	
	LogIostreamKey = attribute.Key("log.iostream")
)

var (
	
	LogIostreamStdout = LogIostreamKey.String("stdout")
	
	LogIostreamStderr = LogIostreamKey.String("stderr")
)


const (
	
	
	
	
	
	
	
	LogFileNameKey = attribute.Key("log.file.name")

	
	
	
	
	
	
	
	
	LogFileNameResolvedKey = attribute.Key("log.file.name_resolved")

	
	
	
	
	
	
	
	LogFilePathKey = attribute.Key("log.file.path")

	
	
	
	
	
	
	
	
	LogFilePathResolvedKey = attribute.Key("log.file.path_resolved")
)




func LogFileName(val string) attribute.KeyValue {
	return LogFileNameKey.String(val)
}




func LogFileNameResolved(val string) attribute.KeyValue {
	return LogFileNameResolvedKey.String(val)
}




func LogFilePath(val string) attribute.KeyValue {
	return LogFilePathKey.String(val)
}




func LogFilePathResolved(val string) attribute.KeyValue {
	return LogFilePathResolvedKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	LogRecordUIDKey = attribute.Key("log.record.uid")
)




func LogRecordUID(val string) attribute.KeyValue {
	return LogRecordUIDKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingBatchMessageCountKey = attribute.Key("messaging.batch.message_count")

	
	
	
	
	
	
	
	
	MessagingClientIDKey = attribute.Key("messaging.client.id")

	
	
	
	
	
	
	
	
	MessagingDestinationAnonymousKey = attribute.Key("messaging.destination.anonymous")

	
	
	
	
	
	
	
	
	
	
	
	
	MessagingDestinationNameKey = attribute.Key("messaging.destination.name")

	
	
	
	
	
	
	
	
	
	MessagingDestinationPartitionIDKey = attribute.Key("messaging.destination.partition.id")

	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingDestinationTemplateKey = attribute.Key("messaging.destination.template")

	
	
	
	
	
	
	
	
	MessagingDestinationTemporaryKey = attribute.Key("messaging.destination.temporary")

	
	
	
	
	
	
	
	
	MessagingDestinationPublishAnonymousKey = attribute.Key("messaging.destination_publish.anonymous")

	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingDestinationPublishNameKey = attribute.Key("messaging.destination_publish.name")

	
	
	
	
	
	
	
	
	
	
	
	MessagingMessageBodySizeKey = attribute.Key("messaging.message.body.size")

	
	
	
	
	
	
	
	
	
	MessagingMessageConversationIDKey = attribute.Key("messaging.message.conversation_id")

	
	
	
	
	
	
	
	
	
	
	
	MessagingMessageEnvelopeSizeKey = attribute.Key("messaging.message.envelope.size")

	
	
	
	
	
	
	
	
	
	MessagingMessageIDKey = attribute.Key("messaging.message.id")

	
	
	
	
	
	
	
	
	MessagingOperationNameKey = attribute.Key("messaging.operation.name")

	
	
	
	
	
	
	
	
	MessagingOperationTypeKey = attribute.Key("messaging.operation.type")

	
	
	
	
	
	
	
	
	
	
	
	MessagingSystemKey = attribute.Key("messaging.system")
)

var (
	
	MessagingOperationTypePublish = MessagingOperationTypeKey.String("publish")
	
	MessagingOperationTypeCreate = MessagingOperationTypeKey.String("create")
	
	MessagingOperationTypeReceive = MessagingOperationTypeKey.String("receive")
	
	MessagingOperationTypeDeliver = MessagingOperationTypeKey.String("process")
	
	MessagingOperationTypeSettle = MessagingOperationTypeKey.String("settle")
)

var (
	
	MessagingSystemActivemq = MessagingSystemKey.String("activemq")
	
	MessagingSystemAWSSqs = MessagingSystemKey.String("aws_sqs")
	
	MessagingSystemEventgrid = MessagingSystemKey.String("eventgrid")
	
	MessagingSystemEventhubs = MessagingSystemKey.String("eventhubs")
	
	MessagingSystemServicebus = MessagingSystemKey.String("servicebus")
	
	MessagingSystemGCPPubsub = MessagingSystemKey.String("gcp_pubsub")
	
	MessagingSystemJms = MessagingSystemKey.String("jms")
	
	MessagingSystemKafka = MessagingSystemKey.String("kafka")
	
	MessagingSystemRabbitmq = MessagingSystemKey.String("rabbitmq")
	
	MessagingSystemRocketmq = MessagingSystemKey.String("rocketmq")
)





func MessagingBatchMessageCount(val int) attribute.KeyValue {
	return MessagingBatchMessageCountKey.Int(val)
}




func MessagingClientID(val string) attribute.KeyValue {
	return MessagingClientIDKey.String(val)
}





func MessagingDestinationAnonymous(val bool) attribute.KeyValue {
	return MessagingDestinationAnonymousKey.Bool(val)
}




func MessagingDestinationName(val string) attribute.KeyValue {
	return MessagingDestinationNameKey.String(val)
}





func MessagingDestinationPartitionID(val string) attribute.KeyValue {
	return MessagingDestinationPartitionIDKey.String(val)
}




func MessagingDestinationTemplate(val string) attribute.KeyValue {
	return MessagingDestinationTemplateKey.String(val)
}





func MessagingDestinationTemporary(val bool) attribute.KeyValue {
	return MessagingDestinationTemporaryKey.Bool(val)
}





func MessagingDestinationPublishAnonymous(val bool) attribute.KeyValue {
	return MessagingDestinationPublishAnonymousKey.Bool(val)
}




func MessagingDestinationPublishName(val string) attribute.KeyValue {
	return MessagingDestinationPublishNameKey.String(val)
}




func MessagingMessageBodySize(val int) attribute.KeyValue {
	return MessagingMessageBodySizeKey.Int(val)
}





func MessagingMessageConversationID(val string) attribute.KeyValue {
	return MessagingMessageConversationIDKey.String(val)
}




func MessagingMessageEnvelopeSize(val int) attribute.KeyValue {
	return MessagingMessageEnvelopeSizeKey.Int(val)
}





func MessagingMessageID(val string) attribute.KeyValue {
	return MessagingMessageIDKey.String(val)
}




func MessagingOperationName(val string) attribute.KeyValue {
	return MessagingOperationNameKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	MessagingKafkaConsumerGroupKey = attribute.Key("messaging.kafka.consumer.group")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	MessagingKafkaMessageKeyKey = attribute.Key("messaging.kafka.message.key")

	
	
	
	
	
	
	
	
	MessagingKafkaMessageOffsetKey = attribute.Key("messaging.kafka.message.offset")

	
	
	
	
	
	
	
	MessagingKafkaMessageTombstoneKey = attribute.Key("messaging.kafka.message.tombstone")
)





func MessagingKafkaConsumerGroup(val string) attribute.KeyValue {
	return MessagingKafkaConsumerGroupKey.String(val)
}







func MessagingKafkaMessageKey(val string) attribute.KeyValue {
	return MessagingKafkaMessageKeyKey.String(val)
}




func MessagingKafkaMessageOffset(val int) attribute.KeyValue {
	return MessagingKafkaMessageOffsetKey.Int(val)
}




func MessagingKafkaMessageTombstone(val bool) attribute.KeyValue {
	return MessagingKafkaMessageTombstoneKey.Bool(val)
}


const (
	
	
	
	
	
	
	
	
	MessagingRabbitmqDestinationRoutingKeyKey = attribute.Key("messaging.rabbitmq.destination.routing_key")

	
	
	
	
	
	
	
	
	MessagingRabbitmqMessageDeliveryTagKey = attribute.Key("messaging.rabbitmq.message.delivery_tag")
)




func MessagingRabbitmqDestinationRoutingKey(val string) attribute.KeyValue {
	return MessagingRabbitmqDestinationRoutingKeyKey.String(val)
}




func MessagingRabbitmqMessageDeliveryTag(val int) attribute.KeyValue {
	return MessagingRabbitmqMessageDeliveryTagKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	MessagingRocketmqClientGroupKey = attribute.Key("messaging.rocketmq.client_group")

	
	
	
	
	
	
	
	
	MessagingRocketmqConsumptionModelKey = attribute.Key("messaging.rocketmq.consumption_model")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqMessageDelayTimeLevelKey = attribute.Key("messaging.rocketmq.message.delay_time_level")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqMessageDeliveryTimestampKey = attribute.Key("messaging.rocketmq.message.delivery_timestamp")

	
	
	
	
	
	
	
	
	
	
	MessagingRocketmqMessageGroupKey = attribute.Key("messaging.rocketmq.message.group")

	
	
	
	
	
	
	
	
	MessagingRocketmqMessageKeysKey = attribute.Key("messaging.rocketmq.message.keys")

	
	
	
	
	
	
	
	
	MessagingRocketmqMessageTagKey = attribute.Key("messaging.rocketmq.message.tag")

	
	
	
	
	
	
	
	MessagingRocketmqMessageTypeKey = attribute.Key("messaging.rocketmq.message.type")

	
	
	
	
	
	
	
	
	
	MessagingRocketmqNamespaceKey = attribute.Key("messaging.rocketmq.namespace")
)

var (
	
	MessagingRocketmqConsumptionModelClustering = MessagingRocketmqConsumptionModelKey.String("clustering")
	
	MessagingRocketmqConsumptionModelBroadcasting = MessagingRocketmqConsumptionModelKey.String("broadcasting")
)

var (
	
	MessagingRocketmqMessageTypeNormal = MessagingRocketmqMessageTypeKey.String("normal")
	
	MessagingRocketmqMessageTypeFifo = MessagingRocketmqMessageTypeKey.String("fifo")
	
	MessagingRocketmqMessageTypeDelay = MessagingRocketmqMessageTypeKey.String("delay")
	
	MessagingRocketmqMessageTypeTransaction = MessagingRocketmqMessageTypeKey.String("transaction")
)





func MessagingRocketmqClientGroup(val string) attribute.KeyValue {
	return MessagingRocketmqClientGroupKey.String(val)
}





func MessagingRocketmqMessageDelayTimeLevel(val int) attribute.KeyValue {
	return MessagingRocketmqMessageDelayTimeLevelKey.Int(val)
}





func MessagingRocketmqMessageDeliveryTimestamp(val int) attribute.KeyValue {
	return MessagingRocketmqMessageDeliveryTimestampKey.Int(val)
}






func MessagingRocketmqMessageGroup(val string) attribute.KeyValue {
	return MessagingRocketmqMessageGroupKey.String(val)
}




func MessagingRocketmqMessageKeys(val ...string) attribute.KeyValue {
	return MessagingRocketmqMessageKeysKey.StringSlice(val)
}




func MessagingRocketmqMessageTag(val string) attribute.KeyValue {
	return MessagingRocketmqMessageTagKey.String(val)
}





func MessagingRocketmqNamespace(val string) attribute.KeyValue {
	return MessagingRocketmqNamespaceKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	MessagingGCPPubsubMessageAckDeadlineKey = attribute.Key("messaging.gcp_pubsub.message.ack_deadline")

	
	
	
	
	
	
	
	
	MessagingGCPPubsubMessageAckIDKey = attribute.Key("messaging.gcp_pubsub.message.ack_id")

	
	
	
	
	
	
	
	
	
	MessagingGCPPubsubMessageDeliveryAttemptKey = attribute.Key("messaging.gcp_pubsub.message.delivery_attempt")

	
	
	
	
	
	
	
	
	
	MessagingGCPPubsubMessageOrderingKeyKey = attribute.Key("messaging.gcp_pubsub.message.ordering_key")
)





func MessagingGCPPubsubMessageAckDeadline(val int) attribute.KeyValue {
	return MessagingGCPPubsubMessageAckDeadlineKey.Int(val)
}




func MessagingGCPPubsubMessageAckID(val string) attribute.KeyValue {
	return MessagingGCPPubsubMessageAckIDKey.String(val)
}




func MessagingGCPPubsubMessageDeliveryAttempt(val int) attribute.KeyValue {
	return MessagingGCPPubsubMessageDeliveryAttemptKey.Int(val)
}





func MessagingGCPPubsubMessageOrderingKey(val string) attribute.KeyValue {
	return MessagingGCPPubsubMessageOrderingKeyKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	MessagingServicebusDestinationSubscriptionNameKey = attribute.Key("messaging.servicebus.destination.subscription_name")

	
	
	
	
	
	
	
	
	MessagingServicebusDispositionStatusKey = attribute.Key("messaging.servicebus.disposition_status")

	
	
	
	
	
	
	
	
	
	MessagingServicebusMessageDeliveryCountKey = attribute.Key("messaging.servicebus.message.delivery_count")

	
	
	
	
	
	
	
	
	
	MessagingServicebusMessageEnqueuedTimeKey = attribute.Key("messaging.servicebus.message.enqueued_time")
)

var (
	
	MessagingServicebusDispositionStatusComplete = MessagingServicebusDispositionStatusKey.String("complete")
	
	MessagingServicebusDispositionStatusAbandon = MessagingServicebusDispositionStatusKey.String("abandon")
	
	MessagingServicebusDispositionStatusDeadLetter = MessagingServicebusDispositionStatusKey.String("dead_letter")
	
	MessagingServicebusDispositionStatusDefer = MessagingServicebusDispositionStatusKey.String("defer")
)






func MessagingServicebusDestinationSubscriptionName(val string) attribute.KeyValue {
	return MessagingServicebusDestinationSubscriptionNameKey.String(val)
}





func MessagingServicebusMessageDeliveryCount(val int) attribute.KeyValue {
	return MessagingServicebusMessageDeliveryCountKey.Int(val)
}





func MessagingServicebusMessageEnqueuedTime(val int) attribute.KeyValue {
	return MessagingServicebusMessageEnqueuedTimeKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	MessagingEventhubsConsumerGroupKey = attribute.Key("messaging.eventhubs.consumer.group")

	
	
	
	
	
	
	
	
	
	MessagingEventhubsMessageEnqueuedTimeKey = attribute.Key("messaging.eventhubs.message.enqueued_time")
)





func MessagingEventhubsConsumerGroup(val string) attribute.KeyValue {
	return MessagingEventhubsConsumerGroupKey.String(val)
}





func MessagingEventhubsMessageEnqueuedTime(val int) attribute.KeyValue {
	return MessagingEventhubsMessageEnqueuedTimeKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	NetworkCarrierIccKey = attribute.Key("network.carrier.icc")

	
	
	
	
	
	
	
	
	NetworkCarrierMccKey = attribute.Key("network.carrier.mcc")

	
	
	
	
	
	
	
	
	NetworkCarrierMncKey = attribute.Key("network.carrier.mnc")

	
	
	
	
	
	
	
	
	NetworkCarrierNameKey = attribute.Key("network.carrier.name")

	
	
	
	
	
	
	
	
	
	
	NetworkConnectionSubtypeKey = attribute.Key("network.connection.subtype")

	
	
	
	
	
	
	
	
	NetworkConnectionTypeKey = attribute.Key("network.connection.type")

	
	
	
	
	
	
	
	
	NetworkIoDirectionKey = attribute.Key("network.io.direction")

	
	
	
	
	
	
	
	
	
	NetworkLocalAddressKey = attribute.Key("network.local.address")

	
	
	
	
	
	
	
	
	NetworkLocalPortKey = attribute.Key("network.local.port")

	
	
	
	
	
	
	
	
	
	NetworkPeerAddressKey = attribute.Key("network.peer.address")

	
	
	
	
	
	
	
	
	NetworkPeerPortKey = attribute.Key("network.peer.port")

	
	
	
	
	
	
	
	
	
	
	NetworkProtocolNameKey = attribute.Key("network.protocol.name")

	
	
	
	
	
	
	
	
	
	
	
	
	NetworkProtocolVersionKey = attribute.Key("network.protocol.version")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	NetworkTransportKey = attribute.Key("network.transport")

	
	
	
	
	
	
	
	
	
	NetworkTypeKey = attribute.Key("network.type")
)

var (
	
	NetworkConnectionSubtypeGprs = NetworkConnectionSubtypeKey.String("gprs")
	
	NetworkConnectionSubtypeEdge = NetworkConnectionSubtypeKey.String("edge")
	
	NetworkConnectionSubtypeUmts = NetworkConnectionSubtypeKey.String("umts")
	
	NetworkConnectionSubtypeCdma = NetworkConnectionSubtypeKey.String("cdma")
	
	NetworkConnectionSubtypeEvdo0 = NetworkConnectionSubtypeKey.String("evdo_0")
	
	NetworkConnectionSubtypeEvdoA = NetworkConnectionSubtypeKey.String("evdo_a")
	
	NetworkConnectionSubtypeCdma20001xrtt = NetworkConnectionSubtypeKey.String("cdma2000_1xrtt")
	
	NetworkConnectionSubtypeHsdpa = NetworkConnectionSubtypeKey.String("hsdpa")
	
	NetworkConnectionSubtypeHsupa = NetworkConnectionSubtypeKey.String("hsupa")
	
	NetworkConnectionSubtypeHspa = NetworkConnectionSubtypeKey.String("hspa")
	
	NetworkConnectionSubtypeIden = NetworkConnectionSubtypeKey.String("iden")
	
	NetworkConnectionSubtypeEvdoB = NetworkConnectionSubtypeKey.String("evdo_b")
	
	NetworkConnectionSubtypeLte = NetworkConnectionSubtypeKey.String("lte")
	
	NetworkConnectionSubtypeEhrpd = NetworkConnectionSubtypeKey.String("ehrpd")
	
	NetworkConnectionSubtypeHspap = NetworkConnectionSubtypeKey.String("hspap")
	
	NetworkConnectionSubtypeGsm = NetworkConnectionSubtypeKey.String("gsm")
	
	NetworkConnectionSubtypeTdScdma = NetworkConnectionSubtypeKey.String("td_scdma")
	
	NetworkConnectionSubtypeIwlan = NetworkConnectionSubtypeKey.String("iwlan")
	
	NetworkConnectionSubtypeNr = NetworkConnectionSubtypeKey.String("nr")
	
	NetworkConnectionSubtypeNrnsa = NetworkConnectionSubtypeKey.String("nrnsa")
	
	NetworkConnectionSubtypeLteCa = NetworkConnectionSubtypeKey.String("lte_ca")
)

var (
	
	NetworkConnectionTypeWifi = NetworkConnectionTypeKey.String("wifi")
	
	NetworkConnectionTypeWired = NetworkConnectionTypeKey.String("wired")
	
	NetworkConnectionTypeCell = NetworkConnectionTypeKey.String("cell")
	
	NetworkConnectionTypeUnavailable = NetworkConnectionTypeKey.String("unavailable")
	
	NetworkConnectionTypeUnknown = NetworkConnectionTypeKey.String("unknown")
)

var (
	
	NetworkIoDirectionTransmit = NetworkIoDirectionKey.String("transmit")
	
	NetworkIoDirectionReceive = NetworkIoDirectionKey.String("receive")
)

var (
	
	NetworkTransportTCP = NetworkTransportKey.String("tcp")
	
	NetworkTransportUDP = NetworkTransportKey.String("udp")
	
	NetworkTransportPipe = NetworkTransportKey.String("pipe")
	
	NetworkTransportUnix = NetworkTransportKey.String("unix")
)

var (
	
	NetworkTypeIpv4 = NetworkTypeKey.String("ipv4")
	
	NetworkTypeIpv6 = NetworkTypeKey.String("ipv6")
)




func NetworkCarrierIcc(val string) attribute.KeyValue {
	return NetworkCarrierIccKey.String(val)
}




func NetworkCarrierMcc(val string) attribute.KeyValue {
	return NetworkCarrierMccKey.String(val)
}




func NetworkCarrierMnc(val string) attribute.KeyValue {
	return NetworkCarrierMncKey.String(val)
}




func NetworkCarrierName(val string) attribute.KeyValue {
	return NetworkCarrierNameKey.String(val)
}




func NetworkLocalAddress(val string) attribute.KeyValue {
	return NetworkLocalAddressKey.String(val)
}




func NetworkLocalPort(val int) attribute.KeyValue {
	return NetworkLocalPortKey.Int(val)
}




func NetworkPeerAddress(val string) attribute.KeyValue {
	return NetworkPeerAddressKey.String(val)
}




func NetworkPeerPort(val int) attribute.KeyValue {
	return NetworkPeerPortKey.Int(val)
}





func NetworkProtocolName(val string) attribute.KeyValue {
	return NetworkProtocolNameKey.String(val)
}




func NetworkProtocolVersion(val string) attribute.KeyValue {
	return NetworkProtocolVersionKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	OciManifestDigestKey = attribute.Key("oci.manifest.digest")
)





func OciManifestDigest(val string) attribute.KeyValue {
	return OciManifestDigestKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	OpentracingRefTypeKey = attribute.Key("opentracing.ref_type")
)

var (
	
	OpentracingRefTypeChildOf = OpentracingRefTypeKey.String("child_of")
	
	OpentracingRefTypeFollowsFrom = OpentracingRefTypeKey.String("follows_from")
)



const (
	
	
	
	
	
	
	
	
	OSBuildIDKey = attribute.Key("os.build_id")

	
	
	
	
	
	
	
	
	
	
	OSDescriptionKey = attribute.Key("os.description")

	
	
	
	
	
	
	
	OSNameKey = attribute.Key("os.name")

	
	
	
	
	
	
	OSTypeKey = attribute.Key("os.type")

	
	
	
	
	
	
	
	
	
	OSVersionKey = attribute.Key("os.version")
)

var (
	
	OSTypeWindows = OSTypeKey.String("windows")
	
	OSTypeLinux = OSTypeKey.String("linux")
	
	OSTypeDarwin = OSTypeKey.String("darwin")
	
	OSTypeFreeBSD = OSTypeKey.String("freebsd")
	
	OSTypeNetBSD = OSTypeKey.String("netbsd")
	
	OSTypeOpenBSD = OSTypeKey.String("openbsd")
	
	OSTypeDragonflyBSD = OSTypeKey.String("dragonflybsd")
	
	OSTypeHPUX = OSTypeKey.String("hpux")
	
	OSTypeAIX = OSTypeKey.String("aix")
	
	OSTypeSolaris = OSTypeKey.String("solaris")
	
	OSTypeZOS = OSTypeKey.String("z_os")
)




func OSBuildID(val string) attribute.KeyValue {
	return OSBuildIDKey.String(val)
}





func OSDescription(val string) attribute.KeyValue {
	return OSDescriptionKey.String(val)
}



func OSName(val string) attribute.KeyValue {
	return OSNameKey.String(val)
}





func OSVersion(val string) attribute.KeyValue {
	return OSVersionKey.String(val)
}


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
	
	
	
	
	
	
	
	
	OTelScopeNameKey = attribute.Key("otel.scope.name")

	
	
	
	
	
	
	
	
	OTelScopeVersionKey = attribute.Key("otel.scope.version")
)




func OTelScopeName(val string) attribute.KeyValue {
	return OTelScopeNameKey.String(val)
}




func OTelScopeVersion(val string) attribute.KeyValue {
	return OTelScopeVersionKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	PeerServiceKey = attribute.Key("peer.service")
)






func PeerService(val string) attribute.KeyValue {
	return PeerServiceKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	ProcessCommandKey = attribute.Key("process.command")

	
	
	
	
	
	
	
	
	
	
	
	
	ProcessCommandArgsKey = attribute.Key("process.command_args")

	
	
	
	
	
	
	
	
	
	
	
	ProcessCommandLineKey = attribute.Key("process.command_line")

	
	
	
	
	
	
	
	
	ProcessContextSwitchTypeKey = attribute.Key("process.context_switch_type")

	
	
	
	
	
	
	
	
	ProcessCreationTimeKey = attribute.Key("process.creation.time")

	
	
	
	
	
	
	
	
	
	
	ProcessExecutableNameKey = attribute.Key("process.executable.name")

	
	
	
	
	
	
	
	
	
	
	ProcessExecutablePathKey = attribute.Key("process.executable.path")

	
	
	
	
	
	
	
	
	ProcessExitCodeKey = attribute.Key("process.exit.code")

	
	
	
	
	
	
	
	
	ProcessExitTimeKey = attribute.Key("process.exit.time")

	
	
	
	
	
	
	
	
	
	ProcessGroupLeaderPIDKey = attribute.Key("process.group_leader.pid")

	
	
	
	
	
	
	
	ProcessInteractiveKey = attribute.Key("process.interactive")

	
	
	
	
	
	
	
	
	ProcessOwnerKey = attribute.Key("process.owner")

	
	
	
	
	
	
	
	
	ProcessPagingFaultTypeKey = attribute.Key("process.paging.fault_type")

	
	
	
	
	
	
	
	
	ProcessParentPIDKey = attribute.Key("process.parent_pid")

	
	
	
	
	
	
	
	ProcessPIDKey = attribute.Key("process.pid")

	
	
	
	
	
	
	
	
	ProcessRealUserIDKey = attribute.Key("process.real_user.id")

	
	
	
	
	
	
	
	
	ProcessRealUserNameKey = attribute.Key("process.real_user.name")

	
	
	
	
	
	
	
	
	
	ProcessRuntimeDescriptionKey = attribute.Key("process.runtime.description")

	
	
	
	
	
	
	
	
	
	ProcessRuntimeNameKey = attribute.Key("process.runtime.name")

	
	
	
	
	
	
	
	
	
	ProcessRuntimeVersionKey = attribute.Key("process.runtime.version")

	
	
	
	
	
	
	
	
	ProcessSavedUserIDKey = attribute.Key("process.saved_user.id")

	
	
	
	
	
	
	
	
	ProcessSavedUserNameKey = attribute.Key("process.saved_user.name")

	
	
	
	
	
	
	
	
	
	ProcessSessionLeaderPIDKey = attribute.Key("process.session_leader.pid")

	
	
	
	
	
	
	
	
	ProcessUserIDKey = attribute.Key("process.user.id")

	
	
	
	
	
	
	
	
	ProcessUserNameKey = attribute.Key("process.user.name")

	
	
	
	
	
	
	
	
	
	
	ProcessVpidKey = attribute.Key("process.vpid")
)

var (
	
	ProcessContextSwitchTypeVoluntary = ProcessContextSwitchTypeKey.String("voluntary")
	
	ProcessContextSwitchTypeInvoluntary = ProcessContextSwitchTypeKey.String("involuntary")
)

var (
	
	ProcessPagingFaultTypeMajor = ProcessPagingFaultTypeKey.String("major")
	
	ProcessPagingFaultTypeMinor = ProcessPagingFaultTypeKey.String("minor")
)






func ProcessCommand(val string) attribute.KeyValue {
	return ProcessCommandKey.String(val)
}








func ProcessCommandArgs(val ...string) attribute.KeyValue {
	return ProcessCommandArgsKey.StringSlice(val)
}







func ProcessCommandLine(val string) attribute.KeyValue {
	return ProcessCommandLineKey.String(val)
}




func ProcessCreationTime(val string) attribute.KeyValue {
	return ProcessCreationTimeKey.String(val)
}






func ProcessExecutableName(val string) attribute.KeyValue {
	return ProcessExecutableNameKey.String(val)
}






func ProcessExecutablePath(val string) attribute.KeyValue {
	return ProcessExecutablePathKey.String(val)
}




func ProcessExitCode(val int) attribute.KeyValue {
	return ProcessExitCodeKey.Int(val)
}




func ProcessExitTime(val string) attribute.KeyValue {
	return ProcessExitTimeKey.String(val)
}





func ProcessGroupLeaderPID(val int) attribute.KeyValue {
	return ProcessGroupLeaderPIDKey.Int(val)
}




func ProcessInteractive(val bool) attribute.KeyValue {
	return ProcessInteractiveKey.Bool(val)
}




func ProcessOwner(val string) attribute.KeyValue {
	return ProcessOwnerKey.String(val)
}




func ProcessParentPID(val int) attribute.KeyValue {
	return ProcessParentPIDKey.Int(val)
}



func ProcessPID(val int) attribute.KeyValue {
	return ProcessPIDKey.Int(val)
}




func ProcessRealUserID(val int) attribute.KeyValue {
	return ProcessRealUserIDKey.Int(val)
}




func ProcessRealUserName(val string) attribute.KeyValue {
	return ProcessRealUserNameKey.String(val)
}





func ProcessRuntimeDescription(val string) attribute.KeyValue {
	return ProcessRuntimeDescriptionKey.String(val)
}





func ProcessRuntimeName(val string) attribute.KeyValue {
	return ProcessRuntimeNameKey.String(val)
}





func ProcessRuntimeVersion(val string) attribute.KeyValue {
	return ProcessRuntimeVersionKey.String(val)
}




func ProcessSavedUserID(val int) attribute.KeyValue {
	return ProcessSavedUserIDKey.Int(val)
}




func ProcessSavedUserName(val string) attribute.KeyValue {
	return ProcessSavedUserNameKey.String(val)
}





func ProcessSessionLeaderPID(val int) attribute.KeyValue {
	return ProcessSessionLeaderPIDKey.Int(val)
}




func ProcessUserID(val int) attribute.KeyValue {
	return ProcessUserIDKey.Int(val)
}




func ProcessUserName(val string) attribute.KeyValue {
	return ProcessUserNameKey.String(val)
}




func ProcessVpid(val int) attribute.KeyValue {
	return ProcessVpidKey.Int(val)
}


const (
	
	
	
	
	
	
	
	ProcessCPUStateKey = attribute.Key("process.cpu.state")
)

var (
	
	ProcessCPUStateSystem = ProcessCPUStateKey.String("system")
	
	ProcessCPUStateUser = ProcessCPUStateKey.String("user")
	
	ProcessCPUStateWait = ProcessCPUStateKey.String("wait")
)


const (
	
	
	
	
	
	
	
	
	RPCConnectRPCErrorCodeKey = attribute.Key("rpc.connect_rpc.error_code")

	
	
	
	
	
	
	
	
	
	RPCGRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")

	
	
	
	
	
	
	
	
	RPCJsonrpcErrorCodeKey = attribute.Key("rpc.jsonrpc.error_code")

	
	
	
	
	
	
	
	
	RPCJsonrpcErrorMessageKey = attribute.Key("rpc.jsonrpc.error_message")

	
	
	
	
	
	
	
	
	
	
	
	RPCJsonrpcRequestIDKey = attribute.Key("rpc.jsonrpc.request_id")

	
	
	
	
	
	
	
	
	
	RPCJsonrpcVersionKey = attribute.Key("rpc.jsonrpc.version")

	
	
	
	
	
	
	
	RPCMessageCompressedSizeKey = attribute.Key("rpc.message.compressed_size")

	
	
	
	
	
	
	
	
	
	
	RPCMessageIDKey = attribute.Key("rpc.message.id")

	
	
	
	
	
	
	
	RPCMessageTypeKey = attribute.Key("rpc.message.type")

	
	
	
	
	
	
	
	RPCMessageUncompressedSizeKey = attribute.Key("rpc.message.uncompressed_size")

	
	
	
	
	
	
	
	
	
	
	
	
	
	RPCMethodKey = attribute.Key("rpc.method")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	RPCServiceKey = attribute.Key("rpc.service")

	
	
	
	
	
	
	
	RPCSystemKey = attribute.Key("rpc.system")
)

var (
	
	RPCConnectRPCErrorCodeCancelled = RPCConnectRPCErrorCodeKey.String("cancelled")
	
	RPCConnectRPCErrorCodeUnknown = RPCConnectRPCErrorCodeKey.String("unknown")
	
	RPCConnectRPCErrorCodeInvalidArgument = RPCConnectRPCErrorCodeKey.String("invalid_argument")
	
	RPCConnectRPCErrorCodeDeadlineExceeded = RPCConnectRPCErrorCodeKey.String("deadline_exceeded")
	
	RPCConnectRPCErrorCodeNotFound = RPCConnectRPCErrorCodeKey.String("not_found")
	
	RPCConnectRPCErrorCodeAlreadyExists = RPCConnectRPCErrorCodeKey.String("already_exists")
	
	RPCConnectRPCErrorCodePermissionDenied = RPCConnectRPCErrorCodeKey.String("permission_denied")
	
	RPCConnectRPCErrorCodeResourceExhausted = RPCConnectRPCErrorCodeKey.String("resource_exhausted")
	
	RPCConnectRPCErrorCodeFailedPrecondition = RPCConnectRPCErrorCodeKey.String("failed_precondition")
	
	RPCConnectRPCErrorCodeAborted = RPCConnectRPCErrorCodeKey.String("aborted")
	
	RPCConnectRPCErrorCodeOutOfRange = RPCConnectRPCErrorCodeKey.String("out_of_range")
	
	RPCConnectRPCErrorCodeUnimplemented = RPCConnectRPCErrorCodeKey.String("unimplemented")
	
	RPCConnectRPCErrorCodeInternal = RPCConnectRPCErrorCodeKey.String("internal")
	
	RPCConnectRPCErrorCodeUnavailable = RPCConnectRPCErrorCodeKey.String("unavailable")
	
	RPCConnectRPCErrorCodeDataLoss = RPCConnectRPCErrorCodeKey.String("data_loss")
	
	RPCConnectRPCErrorCodeUnauthenticated = RPCConnectRPCErrorCodeKey.String("unauthenticated")
)

var (
	
	RPCGRPCStatusCodeOk = RPCGRPCStatusCodeKey.Int(0)
	
	RPCGRPCStatusCodeCancelled = RPCGRPCStatusCodeKey.Int(1)
	
	RPCGRPCStatusCodeUnknown = RPCGRPCStatusCodeKey.Int(2)
	
	RPCGRPCStatusCodeInvalidArgument = RPCGRPCStatusCodeKey.Int(3)
	
	RPCGRPCStatusCodeDeadlineExceeded = RPCGRPCStatusCodeKey.Int(4)
	
	RPCGRPCStatusCodeNotFound = RPCGRPCStatusCodeKey.Int(5)
	
	RPCGRPCStatusCodeAlreadyExists = RPCGRPCStatusCodeKey.Int(6)
	
	RPCGRPCStatusCodePermissionDenied = RPCGRPCStatusCodeKey.Int(7)
	
	RPCGRPCStatusCodeResourceExhausted = RPCGRPCStatusCodeKey.Int(8)
	
	RPCGRPCStatusCodeFailedPrecondition = RPCGRPCStatusCodeKey.Int(9)
	
	RPCGRPCStatusCodeAborted = RPCGRPCStatusCodeKey.Int(10)
	
	RPCGRPCStatusCodeOutOfRange = RPCGRPCStatusCodeKey.Int(11)
	
	RPCGRPCStatusCodeUnimplemented = RPCGRPCStatusCodeKey.Int(12)
	
	RPCGRPCStatusCodeInternal = RPCGRPCStatusCodeKey.Int(13)
	
	RPCGRPCStatusCodeUnavailable = RPCGRPCStatusCodeKey.Int(14)
	
	RPCGRPCStatusCodeDataLoss = RPCGRPCStatusCodeKey.Int(15)
	
	RPCGRPCStatusCodeUnauthenticated = RPCGRPCStatusCodeKey.Int(16)
)

var (
	
	RPCMessageTypeSent = RPCMessageTypeKey.String("SENT")
	
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

var (
	
	RPCSystemGRPC = RPCSystemKey.String("grpc")
	
	RPCSystemJavaRmi = RPCSystemKey.String("java_rmi")
	
	RPCSystemDotnetWcf = RPCSystemKey.String("dotnet_wcf")
	
	RPCSystemApacheDubbo = RPCSystemKey.String("apache_dubbo")
	
	RPCSystemConnectRPC = RPCSystemKey.String("connect_rpc")
)




func RPCJsonrpcErrorCode(val int) attribute.KeyValue {
	return RPCJsonrpcErrorCodeKey.Int(val)
}




func RPCJsonrpcErrorMessage(val string) attribute.KeyValue {
	return RPCJsonrpcErrorMessageKey.String(val)
}







func RPCJsonrpcRequestID(val string) attribute.KeyValue {
	return RPCJsonrpcRequestIDKey.String(val)
}





func RPCJsonrpcVersion(val string) attribute.KeyValue {
	return RPCJsonrpcVersionKey.String(val)
}




func RPCMessageCompressedSize(val int) attribute.KeyValue {
	return RPCMessageCompressedSizeKey.Int(val)
}





func RPCMessageID(val int) attribute.KeyValue {
	return RPCMessageIDKey.Int(val)
}




func RPCMessageUncompressedSize(val int) attribute.KeyValue {
	return RPCMessageUncompressedSizeKey.Int(val)
}




func RPCMethod(val string) attribute.KeyValue {
	return RPCMethodKey.String(val)
}




func RPCService(val string) attribute.KeyValue {
	return RPCServiceKey.String(val)
}









const (
	
	
	
	
	
	
	
	
	
	
	
	
	ServerAddressKey = attribute.Key("server.address")

	
	
	
	
	
	
	
	
	
	
	ServerPortKey = attribute.Key("server.port")
)





func ServerAddress(val string) attribute.KeyValue {
	return ServerAddressKey.String(val)
}



func ServerPort(val int) attribute.KeyValue {
	return ServerPortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ServiceInstanceIDKey = attribute.Key("service.instance.id")

	
	
	
	
	
	
	
	
	
	
	
	
	
	ServiceNameKey = attribute.Key("service.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ServiceNamespaceKey = attribute.Key("service.namespace")

	
	
	
	
	
	
	
	
	
	ServiceVersionKey = attribute.Key("service.version")
)




func ServiceInstanceID(val string) attribute.KeyValue {
	return ServiceInstanceIDKey.String(val)
}




func ServiceName(val string) attribute.KeyValue {
	return ServiceNameKey.String(val)
}




func ServiceNamespace(val string) attribute.KeyValue {
	return ServiceNamespaceKey.String(val)
}





func ServiceVersion(val string) attribute.KeyValue {
	return ServiceVersionKey.String(val)
}












const (
	
	
	
	
	
	
	
	SessionIDKey = attribute.Key("session.id")

	
	
	
	
	
	
	
	
	SessionPreviousIDKey = attribute.Key("session.previous_id")
)



func SessionID(val string) attribute.KeyValue {
	return SessionIDKey.String(val)
}




func SessionPreviousID(val string) attribute.KeyValue {
	return SessionPreviousIDKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	SignalrConnectionStatusKey = attribute.Key("signalr.connection.status")

	
	
	
	
	
	
	
	
	
	SignalrTransportKey = attribute.Key("signalr.transport")
)

var (
	
	SignalrConnectionStatusNormalClosure = SignalrConnectionStatusKey.String("normal_closure")
	
	SignalrConnectionStatusTimeout = SignalrConnectionStatusKey.String("timeout")
	
	SignalrConnectionStatusAppShutdown = SignalrConnectionStatusKey.String("app_shutdown")
)

var (
	
	SignalrTransportServerSentEvents = SignalrTransportKey.String("server_sent_events")
	
	SignalrTransportLongPolling = SignalrTransportKey.String("long_polling")
	
	SignalrTransportWebSockets = SignalrTransportKey.String("web_sockets")
)









const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	SourceAddressKey = attribute.Key("source.address")

	
	
	
	
	
	
	
	SourcePortKey = attribute.Key("source.port")
)





func SourceAddress(val string) attribute.KeyValue {
	return SourceAddressKey.String(val)
}



func SourcePort(val int) attribute.KeyValue {
	return SourcePortKey.Int(val)
}


const (
	
	
	
	
	
	
	
	SystemDeviceKey = attribute.Key("system.device")
)



func SystemDevice(val string) attribute.KeyValue {
	return SystemDeviceKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	SystemCPULogicalNumberKey = attribute.Key("system.cpu.logical_number")

	
	
	
	
	
	
	
	
	SystemCPUStateKey = attribute.Key("system.cpu.state")
)

var (
	
	SystemCPUStateUser = SystemCPUStateKey.String("user")
	
	SystemCPUStateSystem = SystemCPUStateKey.String("system")
	
	SystemCPUStateNice = SystemCPUStateKey.String("nice")
	
	SystemCPUStateIdle = SystemCPUStateKey.String("idle")
	
	SystemCPUStateIowait = SystemCPUStateKey.String("iowait")
	
	SystemCPUStateInterrupt = SystemCPUStateKey.String("interrupt")
	
	SystemCPUStateSteal = SystemCPUStateKey.String("steal")
)




func SystemCPULogicalNumber(val int) attribute.KeyValue {
	return SystemCPULogicalNumberKey.Int(val)
}


const (
	
	
	
	
	
	
	
	
	SystemMemoryStateKey = attribute.Key("system.memory.state")
)

var (
	
	SystemMemoryStateUsed = SystemMemoryStateKey.String("used")
	
	SystemMemoryStateFree = SystemMemoryStateKey.String("free")
	
	SystemMemoryStateShared = SystemMemoryStateKey.String("shared")
	
	SystemMemoryStateBuffers = SystemMemoryStateKey.String("buffers")
	
	SystemMemoryStateCached = SystemMemoryStateKey.String("cached")
)


const (
	
	
	
	
	
	
	
	
	SystemPagingDirectionKey = attribute.Key("system.paging.direction")

	
	
	
	
	
	
	
	
	SystemPagingStateKey = attribute.Key("system.paging.state")

	
	
	
	
	
	
	
	
	SystemPagingTypeKey = attribute.Key("system.paging.type")
)

var (
	
	SystemPagingDirectionIn = SystemPagingDirectionKey.String("in")
	
	SystemPagingDirectionOut = SystemPagingDirectionKey.String("out")
)

var (
	
	SystemPagingStateUsed = SystemPagingStateKey.String("used")
	
	SystemPagingStateFree = SystemPagingStateKey.String("free")
)

var (
	
	SystemPagingTypeMajor = SystemPagingTypeKey.String("major")
	
	SystemPagingTypeMinor = SystemPagingTypeKey.String("minor")
)


const (
	
	
	
	
	
	
	
	
	SystemFilesystemModeKey = attribute.Key("system.filesystem.mode")

	
	
	
	
	
	
	
	
	SystemFilesystemMountpointKey = attribute.Key("system.filesystem.mountpoint")

	
	
	
	
	
	
	
	
	SystemFilesystemStateKey = attribute.Key("system.filesystem.state")

	
	
	
	
	
	
	
	
	SystemFilesystemTypeKey = attribute.Key("system.filesystem.type")
)

var (
	
	SystemFilesystemStateUsed = SystemFilesystemStateKey.String("used")
	
	SystemFilesystemStateFree = SystemFilesystemStateKey.String("free")
	
	SystemFilesystemStateReserved = SystemFilesystemStateKey.String("reserved")
)

var (
	
	SystemFilesystemTypeFat32 = SystemFilesystemTypeKey.String("fat32")
	
	SystemFilesystemTypeExfat = SystemFilesystemTypeKey.String("exfat")
	
	SystemFilesystemTypeNtfs = SystemFilesystemTypeKey.String("ntfs")
	
	SystemFilesystemTypeRefs = SystemFilesystemTypeKey.String("refs")
	
	SystemFilesystemTypeHfsplus = SystemFilesystemTypeKey.String("hfsplus")
	
	SystemFilesystemTypeExt4 = SystemFilesystemTypeKey.String("ext4")
)




func SystemFilesystemMode(val string) attribute.KeyValue {
	return SystemFilesystemModeKey.String(val)
}




func SystemFilesystemMountpoint(val string) attribute.KeyValue {
	return SystemFilesystemMountpointKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	SystemNetworkStateKey = attribute.Key("system.network.state")
)

var (
	
	SystemNetworkStateClose = SystemNetworkStateKey.String("close")
	
	SystemNetworkStateCloseWait = SystemNetworkStateKey.String("close_wait")
	
	SystemNetworkStateClosing = SystemNetworkStateKey.String("closing")
	
	SystemNetworkStateDelete = SystemNetworkStateKey.String("delete")
	
	SystemNetworkStateEstablished = SystemNetworkStateKey.String("established")
	
	SystemNetworkStateFinWait1 = SystemNetworkStateKey.String("fin_wait_1")
	
	SystemNetworkStateFinWait2 = SystemNetworkStateKey.String("fin_wait_2")
	
	SystemNetworkStateLastAck = SystemNetworkStateKey.String("last_ack")
	
	SystemNetworkStateListen = SystemNetworkStateKey.String("listen")
	
	SystemNetworkStateSynRecv = SystemNetworkStateKey.String("syn_recv")
	
	SystemNetworkStateSynSent = SystemNetworkStateKey.String("syn_sent")
	
	SystemNetworkStateTimeWait = SystemNetworkStateKey.String("time_wait")
)


const (
	
	
	
	
	
	
	
	
	
	SystemProcessStatusKey = attribute.Key("system.process.status")
)

var (
	
	SystemProcessStatusRunning = SystemProcessStatusKey.String("running")
	
	SystemProcessStatusSleeping = SystemProcessStatusKey.String("sleeping")
	
	SystemProcessStatusStopped = SystemProcessStatusKey.String("stopped")
	
	SystemProcessStatusDefunct = SystemProcessStatusKey.String("defunct")
)


const (
	
	
	
	
	
	
	
	TelemetrySDKLanguageKey = attribute.Key("telemetry.sdk.language")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	TelemetrySDKNameKey = attribute.Key("telemetry.sdk.name")

	
	
	
	
	
	
	
	
	TelemetrySDKVersionKey = attribute.Key("telemetry.sdk.version")

	
	
	
	
	
	
	
	
	
	
	
	
	TelemetryDistroNameKey = attribute.Key("telemetry.distro.name")

	
	
	
	
	
	
	
	
	
	TelemetryDistroVersionKey = attribute.Key("telemetry.distro.version")
)

var (
	
	TelemetrySDKLanguageCPP = TelemetrySDKLanguageKey.String("cpp")
	
	TelemetrySDKLanguageDotnet = TelemetrySDKLanguageKey.String("dotnet")
	
	TelemetrySDKLanguageErlang = TelemetrySDKLanguageKey.String("erlang")
	
	TelemetrySDKLanguageGo = TelemetrySDKLanguageKey.String("go")
	
	TelemetrySDKLanguageJava = TelemetrySDKLanguageKey.String("java")
	
	TelemetrySDKLanguageNodejs = TelemetrySDKLanguageKey.String("nodejs")
	
	TelemetrySDKLanguagePHP = TelemetrySDKLanguageKey.String("php")
	
	TelemetrySDKLanguagePython = TelemetrySDKLanguageKey.String("python")
	
	TelemetrySDKLanguageRuby = TelemetrySDKLanguageKey.String("ruby")
	
	TelemetrySDKLanguageRust = TelemetrySDKLanguageKey.String("rust")
	
	TelemetrySDKLanguageSwift = TelemetrySDKLanguageKey.String("swift")
	
	TelemetrySDKLanguageWebjs = TelemetrySDKLanguageKey.String("webjs")
)




func TelemetrySDKName(val string) attribute.KeyValue {
	return TelemetrySDKNameKey.String(val)
}




func TelemetrySDKVersion(val string) attribute.KeyValue {
	return TelemetrySDKVersionKey.String(val)
}




func TelemetryDistroName(val string) attribute.KeyValue {
	return TelemetryDistroNameKey.String(val)
}




func TelemetryDistroVersion(val string) attribute.KeyValue {
	return TelemetryDistroVersionKey.String(val)
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
	
	
	
	
	
	
	
	
	
	
	
	
	
	TLSCipherKey = attribute.Key("tls.cipher")

	
	
	
	
	
	
	
	
	
	
	TLSClientCertificateKey = attribute.Key("tls.client.certificate")

	
	
	
	
	
	
	
	
	
	
	
	TLSClientCertificateChainKey = attribute.Key("tls.client.certificate_chain")

	
	
	
	
	
	
	
	
	
	
	TLSClientHashMd5Key = attribute.Key("tls.client.hash.md5")

	
	
	
	
	
	
	
	
	
	
	TLSClientHashSha1Key = attribute.Key("tls.client.hash.sha1")

	
	
	
	
	
	
	
	
	
	
	
	TLSClientHashSha256Key = attribute.Key("tls.client.hash.sha256")

	
	
	
	
	
	
	
	
	
	
	
	TLSClientIssuerKey = attribute.Key("tls.client.issuer")

	
	
	
	
	
	
	
	
	TLSClientJa3Key = attribute.Key("tls.client.ja3")

	
	
	
	
	
	
	
	
	TLSClientNotAfterKey = attribute.Key("tls.client.not_after")

	
	
	
	
	
	
	
	
	TLSClientNotBeforeKey = attribute.Key("tls.client.not_before")

	
	
	
	
	
	
	
	
	
	TLSClientServerNameKey = attribute.Key("tls.client.server_name")

	
	
	
	
	
	
	
	
	
	TLSClientSubjectKey = attribute.Key("tls.client.subject")

	
	
	
	
	
	
	
	
	
	TLSClientSupportedCiphersKey = attribute.Key("tls.client.supported_ciphers")

	
	
	
	
	
	
	
	
	TLSCurveKey = attribute.Key("tls.curve")

	
	
	
	
	
	
	
	
	
	TLSEstablishedKey = attribute.Key("tls.established")

	
	
	
	
	
	
	
	
	
	
	TLSNextProtocolKey = attribute.Key("tls.next_protocol")

	
	
	
	
	
	
	
	
	
	TLSProtocolNameKey = attribute.Key("tls.protocol.name")

	
	
	
	
	
	
	
	
	
	
	TLSProtocolVersionKey = attribute.Key("tls.protocol.version")

	
	
	
	
	
	
	
	
	TLSResumedKey = attribute.Key("tls.resumed")

	
	
	
	
	
	
	
	
	
	
	TLSServerCertificateKey = attribute.Key("tls.server.certificate")

	
	
	
	
	
	
	
	
	
	
	
	TLSServerCertificateChainKey = attribute.Key("tls.server.certificate_chain")

	
	
	
	
	
	
	
	
	
	
	TLSServerHashMd5Key = attribute.Key("tls.server.hash.md5")

	
	
	
	
	
	
	
	
	
	
	TLSServerHashSha1Key = attribute.Key("tls.server.hash.sha1")

	
	
	
	
	
	
	
	
	
	
	
	TLSServerHashSha256Key = attribute.Key("tls.server.hash.sha256")

	
	
	
	
	
	
	
	
	
	
	
	TLSServerIssuerKey = attribute.Key("tls.server.issuer")

	
	
	
	
	
	
	
	
	TLSServerJa3sKey = attribute.Key("tls.server.ja3s")

	
	
	
	
	
	
	
	
	TLSServerNotAfterKey = attribute.Key("tls.server.not_after")

	
	
	
	
	
	
	
	
	TLSServerNotBeforeKey = attribute.Key("tls.server.not_before")

	
	
	
	
	
	
	
	
	
	TLSServerSubjectKey = attribute.Key("tls.server.subject")
)

var (
	
	TLSProtocolNameSsl = TLSProtocolNameKey.String("ssl")
	
	TLSProtocolNameTLS = TLSProtocolNameKey.String("tls")
)





func TLSCipher(val string) attribute.KeyValue {
	return TLSCipherKey.String(val)
}






func TLSClientCertificate(val string) attribute.KeyValue {
	return TLSClientCertificateKey.String(val)
}






func TLSClientCertificateChain(val ...string) attribute.KeyValue {
	return TLSClientCertificateChainKey.StringSlice(val)
}






func TLSClientHashMd5(val string) attribute.KeyValue {
	return TLSClientHashMd5Key.String(val)
}






func TLSClientHashSha1(val string) attribute.KeyValue {
	return TLSClientHashSha1Key.String(val)
}






func TLSClientHashSha256(val string) attribute.KeyValue {
	return TLSClientHashSha256Key.String(val)
}






func TLSClientIssuer(val string) attribute.KeyValue {
	return TLSClientIssuerKey.String(val)
}




func TLSClientJa3(val string) attribute.KeyValue {
	return TLSClientJa3Key.String(val)
}




func TLSClientNotAfter(val string) attribute.KeyValue {
	return TLSClientNotAfterKey.String(val)
}




func TLSClientNotBefore(val string) attribute.KeyValue {
	return TLSClientNotBeforeKey.String(val)
}





func TLSClientServerName(val string) attribute.KeyValue {
	return TLSClientServerNameKey.String(val)
}




func TLSClientSubject(val string) attribute.KeyValue {
	return TLSClientSubjectKey.String(val)
}




func TLSClientSupportedCiphers(val ...string) attribute.KeyValue {
	return TLSClientSupportedCiphersKey.StringSlice(val)
}




func TLSCurve(val string) attribute.KeyValue {
	return TLSCurveKey.String(val)
}





func TLSEstablished(val bool) attribute.KeyValue {
	return TLSEstablishedKey.Bool(val)
}






func TLSNextProtocol(val string) attribute.KeyValue {
	return TLSNextProtocolKey.String(val)
}






func TLSProtocolVersion(val string) attribute.KeyValue {
	return TLSProtocolVersionKey.String(val)
}




func TLSResumed(val bool) attribute.KeyValue {
	return TLSResumedKey.Bool(val)
}






func TLSServerCertificate(val string) attribute.KeyValue {
	return TLSServerCertificateKey.String(val)
}






func TLSServerCertificateChain(val ...string) attribute.KeyValue {
	return TLSServerCertificateChainKey.StringSlice(val)
}






func TLSServerHashMd5(val string) attribute.KeyValue {
	return TLSServerHashMd5Key.String(val)
}






func TLSServerHashSha1(val string) attribute.KeyValue {
	return TLSServerHashSha1Key.String(val)
}






func TLSServerHashSha256(val string) attribute.KeyValue {
	return TLSServerHashSha256Key.String(val)
}






func TLSServerIssuer(val string) attribute.KeyValue {
	return TLSServerIssuerKey.String(val)
}




func TLSServerJa3s(val string) attribute.KeyValue {
	return TLSServerJa3sKey.String(val)
}




func TLSServerNotAfter(val string) attribute.KeyValue {
	return TLSServerNotAfterKey.String(val)
}




func TLSServerNotBefore(val string) attribute.KeyValue {
	return TLSServerNotBeforeKey.String(val)
}




func TLSServerSubject(val string) attribute.KeyValue {
	return TLSServerSubjectKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	URLDomainKey = attribute.Key("url.domain")

	
	
	
	
	
	
	
	
	
	
	
	
	URLExtensionKey = attribute.Key("url.extension")

	
	
	
	
	
	
	
	
	URLFragmentKey = attribute.Key("url.fragment")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	URLFullKey = attribute.Key("url.full")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	URLOriginalKey = attribute.Key("url.original")

	
	
	
	
	
	
	
	
	
	
	URLPathKey = attribute.Key("url.path")

	
	
	
	
	
	
	
	URLPortKey = attribute.Key("url.port")

	
	
	
	
	
	
	
	
	
	
	URLQueryKey = attribute.Key("url.query")

	
	
	
	
	
	
	
	
	
	
	
	
	URLRegisteredDomainKey = attribute.Key("url.registered_domain")

	
	
	
	
	
	
	
	
	
	URLSchemeKey = attribute.Key("url.scheme")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	URLSubdomainKey = attribute.Key("url.subdomain")

	
	
	
	
	
	
	
	
	
	URLTemplateKey = attribute.Key("url.template")

	
	
	
	
	
	
	
	
	
	
	
	
	URLTopLevelDomainKey = attribute.Key("url.top_level_domain")
)




func URLDomain(val string) attribute.KeyValue {
	return URLDomainKey.String(val)
}




func URLExtension(val string) attribute.KeyValue {
	return URLExtensionKey.String(val)
}




func URLFragment(val string) attribute.KeyValue {
	return URLFragmentKey.String(val)
}




func URLFull(val string) attribute.KeyValue {
	return URLFullKey.String(val)
}




func URLOriginal(val string) attribute.KeyValue {
	return URLOriginalKey.String(val)
}




func URLPath(val string) attribute.KeyValue {
	return URLPathKey.String(val)
}



func URLPort(val int) attribute.KeyValue {
	return URLPortKey.Int(val)
}




func URLQuery(val string) attribute.KeyValue {
	return URLQueryKey.String(val)
}




func URLRegisteredDomain(val string) attribute.KeyValue {
	return URLRegisteredDomainKey.String(val)
}





func URLScheme(val string) attribute.KeyValue {
	return URLSchemeKey.String(val)
}







func URLSubdomain(val string) attribute.KeyValue {
	return URLSubdomainKey.String(val)
}





func URLTemplate(val string) attribute.KeyValue {
	return URLTemplateKey.String(val)
}





func URLTopLevelDomain(val string) attribute.KeyValue {
	return URLTopLevelDomainKey.String(val)
}


const (
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	UserAgentNameKey = attribute.Key("user_agent.name")

	
	
	
	
	
	
	
	
	
	
	
	
	
	UserAgentOriginalKey = attribute.Key("user_agent.original")

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	UserAgentVersionKey = attribute.Key("user_agent.version")
)




func UserAgentName(val string) attribute.KeyValue {
	return UserAgentNameKey.String(val)
}






func UserAgentOriginal(val string) attribute.KeyValue {
	return UserAgentOriginalKey.String(val)
}




func UserAgentVersion(val string) attribute.KeyValue {
	return UserAgentVersionKey.String(val)
}



const (
	
	
	
	
	
	
	
	
	
	
	WebEngineDescriptionKey = attribute.Key("webengine.description")

	
	
	
	
	
	
	
	WebEngineNameKey = attribute.Key("webengine.name")

	
	
	
	
	
	
	
	
	WebEngineVersionKey = attribute.Key("webengine.version")
)





func WebEngineDescription(val string) attribute.KeyValue {
	return WebEngineDescriptionKey.String(val)
}




func WebEngineName(val string) attribute.KeyValue {
	return WebEngineNameKey.String(val)
}




func WebEngineVersion(val string) attribute.KeyValue {
	return WebEngineVersionKey.String(val)
}
