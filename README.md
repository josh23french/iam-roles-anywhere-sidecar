# iam-roles-anywhere-sidecar

A sidecar container that provides credentials (AWS_CONTAINER_CREDENTIALS_FULL_URI) using the credential helper via a simple HTTP server.

## Usage

Set this container up as a sidecar to your application container. In Kubernetes, this would be another container in the same pod, but this will work for any container solution.

For this sidecar container, specify the following environment variables (non-specific ones like DEBUG and WITH_PROXY have sane defaults as shown below and do not need to be specified unless you're changing them):

```yaml
env:
# These three have to do with the location of your certs on disk
- name: PRIVATE_KEY_ID
  value: /var/run/wherever.crt # path to TLS certificate
- name: CERTIFICATE_ID
  value: /var/run/wherever.key # path to matching key
- name: CERTIFICATE_BUNDLE_ID
  value: # not sure what this does right now

# These mirror specifics of your AWS account/region
- name: AWS_REGION
  value: us-east-1
- name: ENDPOINT
  value: rolesanywhere.us-east-1.amazonaws.com # change if you are using VPC endpoints without private DNS
- name: NO_VERIFY_SSL
  value: false
- name: WITH_PROXY
  value: ""
- name: ROLE_ARN
  value: arn:aws:iam::123456789012:role/MyRole
- name: PROFILE_ARN
  value: arn:aws:rolesanywhere:us-east-1:123456789012:profile/701edfc3-c651-40dc-bf73-87c6428dffb4
- name: TRUST_ANCHOR_ID
  value: arn:aws:rolesanywhere:us-east-1:123456789012:trust-anchor/5934afd0-0a02-40c0-8abc-bd6dce01496c
- name: DEBUG
  value: false
- name: CREDENTIAL_VERSION
  value: 1
- name: SESSION_DURATION
  value: 900 # Number of seconds requested credentials should be valid for

# This controls the HTTP server itself
- name: LISTEN
  value: [::1]:8080
```

In your main container, add the following environment variable, substituting what you set for LISTEN if changed:

```yaml
env:
- name: AWS_CONTAINER_CREDENTIALS_FULL_URI
  value: http://localhost:8080/creds
```

That's it! Just use the AWS SDK as normal in the application container without specifying credentials.
