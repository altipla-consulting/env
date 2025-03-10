package env

import (
	"os"
)

// IsLocal returns true if we are running inside a local debug environment instead
// of a production container. It dependes on Version() working correctly.
func IsLocal() bool {
	return Version() == ""
}

// IsProduction returns true if we are running inside a production environment.
// It depends on Version() working correctly.
func IsProduction() bool {
	return Version() != ""
}

// IsJenkins detects if we are running as a step of a Jenkins build.
func IsJenkins() bool {
	return os.Getenv("BUILD_ID") != ""
}

// IsCI detects CI environments like Jenkins or GitHub Actions.
func IsCI() bool {
	return IsJenkins() || os.Getenv("CI") != ""
}

// IsCloudRun detects if we are running inside a Cloud Run app.
func IsCloudRun() bool {
	return os.Getenv("K_CONFIGURATION") != "" || os.Getenv("CLOUD_RUN_JOB") != ""
}

// IsAzureFunction detects if we are running inside an Azure Function app.
func IsAzureFunction() bool {
	return os.Getenv("APPSETTING_WEBSITE_SITE_NAME") != ""
}

// IsKubernetes detects if we are running inside a Kubernetes pod.
func IsKubernetes() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

// Version returns the application version. In supported environment it may extract
// the info from files or environment variables. Otherwise it will use the env variable
// VERSION that should be set manually to the desired value.
func Version() string {
	// Explictly set from the outside in an environment variable.
	if v := os.Getenv("VERSION"); v != "" {
		return v
	}

	// Cloud Run default revision name.
	if e := os.Getenv("K_REVISION"); e != "" {
		return e
	}
	if e := os.Getenv("CLOUD_RUN_EXECUTION"); e != "" {
		return e
	}

	// Azure Function has a Kubu file one level up of the working directory.
	if IsAzureFunction() {
		v, err := os.ReadFile("../deployments/active")
		if err == nil {
			return string(v)
		} else if !os.IsNotExist(err) {
			panic(err)
		}
	}

	return ""
}

// ServiceName returns the version of the application from the environment-dependent
// variables or the hostname if none is available.
func ServiceName() string {
	// Cloud Run.
	if v := os.Getenv("K_SERVICE"); v != "" {
		return v
	}
	if v := os.Getenv("CLOUD_RUN_JOB"); v != "" {
		return v
	}

	// Azure Functions.
	if v := os.Getenv("APPSETTING_WEBSITE_SITE_NAME"); v != "" {
		return v
	}

	// Fly.io.
	if v := os.Getenv("FLY_APP_NAME"); v != "" {
		return v
	}

	// Last option is a hostname which should always exist, though is not very useful usually.
	v, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return v
}
