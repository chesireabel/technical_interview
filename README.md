# technical_interview
# ⚡ Technical Interview Project

![CI](https://img.shields.io/github/actions/workflow/status/chesireabel/technical_interview/ci.yml?label=CI&logo=github)
![Docker Pulls](https://img.shields.io/docker/pulls/abelchesire/technical-interview)
![Kubernetes](https://img.shields.io/badge/Kubernetes-Deployed-blue?logo=kubernetes)

---

## 📖 Overview
A simple Go service containerized with **Docker**, deployed to **Kubernetes** via **Helm**, and using **GitHub Actions** for CI/CD.

---

## 🚀 Features
- Dockerized app 
- Helm chart for easy deployments
- CI: build & push Docker images automatically
- CD: deploy to Kubernetes cluster
- Secrets managed via Kubernetes + GitHub Secrets

---

## 🛠 Tech Stack
- **Backend:** Golang
- **Database:** PostgreSQL
- **CI/CD:** GitHub Actions
- **Orchestration:** Kubernetes + Helm
- **Registry:** DockerHub / GHCR

---

## ⚡ Quickstart (local)
```bash
# 1. Clone your repo
git clone https://github.com/chesireabel/technical_interview.git
cd technical_interview

# 2. Build the Docker image and tag it
docker build -t abelchesire/technical-interview:dev .

# 3. Run the container locally, mapping port 8081
docker run -p 8081:8081 abelchesire/technical-interview:dev

