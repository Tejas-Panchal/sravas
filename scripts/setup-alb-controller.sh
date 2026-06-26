#!/usr/bin/env bash
# Install AWS Load Balancer Controller on EKS via Terraform IRSA
set -euo pipefail

ENV="${1:-dev}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TF_DIR="${SCRIPT_DIR}/../infrastructure/terraform"
CLUSTER_NAME="sravas-${ENV}"
CONTROLLER_VERSION="v2.7.1"
NAMESPACE="kube-system"

cd "$TF_DIR"

# Get IRSA role ARN from Terraform output
ROLE_ARN=$(terraform output -raw alb_controller_role_arn)
echo "IRSA role ARN: ${ROLE_ARN}"

# Get VPC ID from cluster
VPC_ID=$(aws eks describe-cluster --name "${CLUSTER_NAME}" \
  --query "cluster.resourcesVpcConfig.vpcId" --output text)
echo "VPC ID: ${VPC_ID}"

# Install CRDs
echo "Installing ALB Controller CRDs..."
kubectl apply -f "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/${CONTROLLER_VERSION}/config/crd/bases/ingress.k8s.aws_ingressclassparams.yaml" \
  -f "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/${CONTROLLER_VERSION}/config/crd/bases/ingress.k8s.aws_targetgroupbindings.yaml" \
  -f "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/${CONTROLLER_VERSION}/config/crd/bases/ingress.k8saws_awsv2targetgroups.yaml"

# Download and patch manifest
echo "Installing ALB Controller ${CONTROLLER_VERSION}..."
TMP_MANIFEST=$(mktemp)
curl -sL "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/${CONTROLLER_VERSION}/docs/install/v2_7_1.yaml" \
  | sed "s|eks.amazonaws.com/role-arn:.*|eks.amazonaws.com/role-arn: ${ROLE_ARN}|" \
  | sed "s|clusterName: .*|clusterName: ${CLUSTER_NAME}|" \
  > "${TMP_MANIFEST}"

# Apply controller manifest
kubectl apply -f "${TMP_MANIFEST}"
rm -f "${TMP_MANIFEST}"

# Wait for controller readiness
echo "Waiting for ALB Controller to become ready..."
kubectl -n "${NAMESPACE}" rollout status deployment/aws-load-balancer-controller --timeout=120s

# Verify
echo "ALB Controller installation complete:"
kubectl -n "${NAMESPACE}" get deployment aws-load-balancer-controller
kubectl -n "${NAMESPACE}" get pods -l app.kubernetes.io/name=aws-load-balancer-controller
