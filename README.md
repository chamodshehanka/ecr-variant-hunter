# ECR Variant Hunter

ECR Variant Hunter is a Go script that uses the AWS SDK to find outdated ECR images and prune them. It is designed to be run as a Kubernetes CronJob.

## Configuration

The script uses the following environment variables:

- `AWS_REGION`: The AWS region where the ECR repositories are located.
- `AWS_ACCESS_KEY_ID`: The AWS access key ID.
- `AWS_SECRET_ACCESS_KEY`: The AWS secret access key.
- `ECR_REPOSITORIES`: A comma-separated list of ECR repositories to scan.

## Usage

```bash
go run main.go
```

## Deployment

The script is designed to be run as a Kubernetes CronJob. The following is an example of a CronJob manifest:

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: ecr-variant-hunter
spec:
  schedule: "0 0 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: ecr-variant-hunter
              image: ghcr.io/chamodshehanka/ecr-variant-hunter:latest
              env:
                - name: AWS_REGION
                  value: <region>
                - name: AWS_ACCESS_KEY_ID
                  value: <access-key-id>
                - name: AWS_SECRET_ACCESS_KEY
                  value: <secret-access-key>
                - name: ECR_REPOSITORIES
                  value: <repositories>
          restartPolicy: OnFailure
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

