package model

// InteractiveApps contains the settings needed for interactive apps across all
// steps in a Job.
type InteractiveApps struct {
	ProxyImage  string //The docker image for the reverse proxy that runs on the cluster with the job steps.
	ProxyName   string //The name of the container for the reverse proxy.
	FrontendURL string //The URL for the frontend of the application. Will get prefixed with the job id.
	CASURL      string //The base URL for the CAS server.
	CASValidate string //The path to the validate endpoint on the CAS server.
	SSLCertPath string //The path to the SSL cert file on the Condor nodes.
	SSLKeyPath  string //The path to the SSL key file on the Condor nodes.
}
