# Kubernetes Integration Guide

Learn how to integrate OpenComply with your Kubernetes environment.

This guide shows you how to connect OpenComply to a Kubernetes cluster using a `kubeconfig` file. Each integration corresponds to a single Kubernetes cluster. Once integrated, OpenComply will discover and assess key Kubernetes resources—such as pods, deployments, services, namespaces, RBAC policies, and configurations—enabling compliance and visibility across your Kubernetes infrastructure.

## Prerequisites

- OpenComply installed and running
- Access to a Kubernetes cluster
- A valid `kubeconfig` file with appropriate read-only permissions

## Configure Integration in OpenComply

1. In the OpenComply dashboard, go to **Integrations > Kubernetes**.
2. Select **kubeconfig integration**.
3. Upload the `kubeconfig` file.
4. Click **Save** to establish the connection.

Once connected, OpenComply will scan your Kubernetes environment for security and compliance insights.