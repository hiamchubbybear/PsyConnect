





package driverutil

import (
	"os"
	"strings"
)

const AwsLambdaPrefix = "AWS_Lambda_"

const (
	

	
	EnvVarAWSExecutionEnv = "AWS_EXECUTION_ENV"
	
	EnvVarAWSLambdaRuntimeAPI = "AWS_LAMBDA_RUNTIME_API"
	
	EnvVarFunctionsWorkerRuntime = "FUNCTIONS_WORKER_RUNTIME"
	
	EnvVarKService = "K_SERVICE"
	
	EnvVarFunctionName = "FUNCTION_NAME"
	
	EnvVarVercel = "VERCEL"
	
	EnvVarK8s = "KUBERNETES_SERVICE_HOST"
)

const (
	

	
	EnvVarAWSRegion = "AWS_REGION"
	
	EnvVarAWSLambdaFunctionMemorySize = "AWS_LAMBDA_FUNCTION_MEMORY_SIZE"
	
	EnvVarFunctionMemoryMB = "FUNCTION_MEMORY_MB"
	
	EnvVarFunctionTimeoutSec = "FUNCTION_TIMEOUT_SEC"
	
	EnvVarFunctionRegion = "FUNCTION_REGION"
	
	EnvVarVercelRegion = "VERCEL_REGION"
)

const (
	

	
	EnvNameAWSLambda = "aws.lambda"
	
	EnvNameAzureFunc = "azure.func"
	
	EnvNameGCPFunc = "gcp.func"
	
	EnvNameVercel = "vercel"
)







func GetFaasEnvName() string {
	envVars := []string{
		EnvVarAWSExecutionEnv,
		EnvVarAWSLambdaRuntimeAPI,
		EnvVarFunctionsWorkerRuntime,
		EnvVarKService,
		EnvVarFunctionName,
		EnvVarVercel,
	}

	
	
	names := make(map[string]struct{})

	for _, envVar := range envVars {
		val := os.Getenv(envVar)
		if val == "" {
			continue
		}

		var name string

		switch envVar {
		case EnvVarAWSExecutionEnv:
			if !strings.HasPrefix(val, AwsLambdaPrefix) {
				continue
			}

			name = EnvNameAWSLambda
		case EnvVarAWSLambdaRuntimeAPI:
			name = EnvNameAWSLambda
		case EnvVarFunctionsWorkerRuntime:
			name = EnvNameAzureFunc
		case EnvVarKService, EnvVarFunctionName:
			name = EnvNameGCPFunc
		case EnvVarVercel:
			
			delete(names, EnvNameAWSLambda)

			name = EnvNameVercel
		}

		names[name] = struct{}{}
		if len(names) > 1 {
			
			
			names = nil

			break
		}
	}

	for name := range names {
		return name
	}

	return ""
}
