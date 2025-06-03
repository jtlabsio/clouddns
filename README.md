# clouddns

Cloud DNS app for Dynamic DNS with Google Cloud

## Prerequisites

- Go 1.16 or later
- Google Cloud Project with billing enabled
- Google Cloud CLI
- Google Cloud DNS API enabled
- A domain name and a subdomain to use for dynamic DNS updates

### Environment Preparation

Once you have setup a Google Cloud project, the following instructions will guide you through setting up the gcloud CLI: <https://cloud.google.com/sdk/docs/install>

Once the above step is complete, enable the Cloud Resource Manager and IAM APIs for your Google Cloud Project:

```bash
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable iam.googleapis.com
```

After completing the above, make note of your PROJECT_ID (this value can be found in the "select a project" menu in the Google Cloud Console). To generate an access token for use with this tool in your project, you'll need to create a service account by following these steps (replace PROJECT_ID with your actual project ID):

```bash
PROJECT_ID="your-project-id"
gcloud iam service-accounts create clouddns-service-account --display-name "CloudDNS Service Account"
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:clouddns-service-account@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/dns.admin"
gcloud iam service-accounts keys create credentials.json --iam-account "clouddns-service-account@${PROJECT_ID}.iam.gserviceaccount.com"
```

The above commands will generate a `credentials.json` file for a service account with the ability to modify the Cloud DNS service. This file should be kept secure and not shared publicly.

## Installation

### Docker

To run this project using Docker, follow these steps:

Step 1. Create your configuration:

  Create a YAML file in `./settings` with the name of your zone... e.g., `example-zone.yaml`. The contents should look like this:

  ```yaml
  google:
    projectID: your-project-id # Replace with your project ID from Google Cloud Console
    managedZone: example-zone  # Replace with your managed zone name as defined in Google Cloud DNS
    record: myip.example.com.  # Replace with the record you want to update (be sure there is a . at the end of the record)
  ```

Step 2. Build the Docker image:

  ```bash
  make docker
  ```

  __note__: If you do not have`make` installed, you can use `docker build -t clouddns-server .`

Step 3. Run the Docker container:

  Replace the `example-zone` with the zone name used in the step above (do not include the `.yaml` extension).

  ```bash
  docker run -d --name clouddns-server -p 8000:8000 -e GO_ENV=example-zone --restart unless-stopped clouddns-server
  ```

  __note__: If you are running this on a local machine, you may need to adjust the port mapping and environment variable accordingly.

### Running Locally with Go

Step 1. Create your configuration:

  Create a YAML file in `./settings` with the name of your zone... e.g., `example-zone.yaml`. The contents should look like this:

  ```yaml
  google:
    projectID: your-project-id # Replace with your project ID from Google Cloud Console
    managedZone: example-zone  # Replace with your managed zone name as defined in Google Cloud DNS
    record: myip.example.com.  # Replace with the record you want to update (be sure there is a . at the end of the record)
  ```

Step 2. Run the server:

  ```bash
  GO_ENV=example-zone go run cmd/main.go
  ```

## Status

The application exposes a status endpoint on port `8000` at `/status` that returns the current IP address and the last update time:

```bash
curl http://localhost:8000/status
```
