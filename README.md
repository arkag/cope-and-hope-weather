# 🌤️ Cope & Hope Weather

*Is your weather miserable? We'll find somewhere better (or worse).*

Cope & Hope is a serverless, event-driven web application that lets users check their local weather and instantly compare it to alternative cities around the world. Depending on your mood, you can choose to "Cope" by finding a city where the weather is even worse, or "Hope" by finding a city where the weather is significantly better.

## 🏗️ Architecture & Tech Stack

The application is built on a modern, fully serverless AWS architecture, utilizing a microservices pattern to separate the frontend UI from the backend processing logic.

### 1. Web UI (Frontend)
- **Framework:** Python / Flask
- **Styling:** Vanilla CSS featuring premium glassmorphism aesthetics, dark mode, and dynamic animations.
- **Deployment:** AWS Lambda via the `Mangum` ASGI adapter and `a2wsgi`.
- **Functionality:** Serves the HTML frontend and acts as an internal proxy to securely invoke the Go backend. 

### 2. Core API (Backend)
- **Framework:** Go (Golang)
- **Deployment:** AWS Lambda using `aws-lambda-go-api-proxy`.
- **Functionality:** Handles high-performance concurrent fetching of weather data from OpenWeatherMap, caching (via DynamoDB), and the "Cope/Hope" comparison logic.
- **Security:** The Go API is completely **private**. It is disconnected from the public internet (no API Gateway routes) and is invoked securely over the AWS internal network using `boto3` and IAM role authorization.

### 3. Infrastructure (Terraform)
- **Provider:** AWS
- **State Management:** S3 Backend with DynamoDB state locking.
- **Networking:** API Gateway HTTP API for public web traffic routing.
- **Custom Domain:** Integrated with AWS Certificate Manager (ACM) for fully automated TLS certificates and DNS mapping (e.g., `demo.kagno.com`).

### 4. CI/CD (GitHub Actions)
Fully automated pipelines handle the entire lifecycle of the application:
- **`terraform.yml`**: Provisions shared infrastructure, validates domains, and configures IAM roles.
- **`api.yml`**: Lints (`golangci-lint`), tests (with race detectors), compiles, and zips the Go backend for Lambda deployment.
- **`web.yml`**: Runs `flake8` strict linting, packages the Flask app natively, and deploys to the Web UI Lambda.

---

## 🚀 Roadmap: Application Security Testing

As the infrastructure matures, integrating continuous security testing directly into our GitHub Actions pipelines is the next priority.

### Phase 1: SAST (Static Application Security Testing)
SAST tools will analyze our source code and Terraform configuration for vulnerabilities before deployment. These will be added as jobs in our existing GitHub Actions:

- **Terraform Security (`checkov` or `tfsec`)**
  - Add to `terraform.yml` to automatically scan for misconfigurations (e.g., overly permissive IAM roles, unencrypted data).
- **Go Security (`gosec`)**
  - Add to `api.yml` to inspect the Go AST for issues like hardcoded credentials or unsafe memory pointers.
- **Python Security (`Bandit`)**
  - Add to `web.yml` to analyze the Flask application for common Python flaws (e.g., shell injection, weak cryptography).
- **Global Scanning (GitHub CodeQL)**
  - Implement a new `codeql.yml` workflow to deeply analyze all code branches for semantic security vulnerabilities on every Pull Request.

### Phase 2: DAST (Dynamic Application Security Testing)
DAST tools will probe the live, running application from the outside to identify runtime vulnerabilities. 

- **OWASP ZAP (Zed Attack Proxy) Integration**
  - Create a new `dast.yml` workflow that triggers *after* a successful deployment to our API Gateway endpoint.
  - ZAP will run automated baseline scans against the live Flask UI and attempt common web attacks (XSS, CSRF, header injection).
  - Configure the workflow to fail and alert the team if high-severity vulnerabilities are detected in the live environment.

### Phase 3: Dependency Scanning (SCA)
- Enable **Dependabot** and **GitHub Advanced Security** to automatically monitor `go.mod` and `requirements.txt` for known vulnerable packages and automatically generate PRs for patches.
