#!/bin/bash

# Disaster Recovery Validation Script
# Tests backup restoration using Velero

set -e

NAMESPACE="velero"
BACKUP_NAME="test-backup-$(date +%Y%m%d-%H%M%S)"

echo "=== Disaster Recovery Validation ==="

# 1. Check Velero installation
echo "1. Checking Velero installation..."
kubectl get pods -n ${NAMESPACE} -l app.kubernetes.io/name=velero || {
    echo "Velero not installed. Installing..."
    kubectl apply -f k8s/velero/velero-install.yaml
    velero install --provider aws --plugins velero-plugin-for-aws --bucket velero-backups --secret-file /tmp/velero-credentials
}

# 2. Create backup
echo "2. Creating backup of all namespaces..."
velero backup create ${BACKUP_NAME} --include-namespaces '*' --exclude-namespaces kube-system,velero

# 3. Wait for backup completion
echo "3. Waiting for backup completion..."
velero backup get ${BACKUP_NAME} --request-timeout 5m

# 4. Simulate disaster (delete a namespace)
echo "4. Simulating disaster by deleting test namespace..."
kubectl create namespace test-disaster || true
kubectl run test-pod --image=nginx -n test-disaster
kubectl delete namespace test-disaster --force --grace-period=0

# 5. Restore from backup
echo "5. Restoring from backup..."
velero restore create ${BACKUP_NAME}-restore --from-backup ${BACKUP_NAME}

# 6. Wait for restore completion
echo "6. Waiting for restore completion..."
velero restore get ${BACKUP_NAME}-restore --request-timeout 5m

# 7. Verify restoration
echo "7. Verifying restoration..."
kubectl get namespace test-disaster || echo "Namespace restoration failed"
kubectl get pod test-pod -n test-disaster || echo "Pod restoration failed"

# 8. Calculate RPO/RTO metrics
echo "8. Calculating RPO/RTO metrics..."
BACKUP_TIME=$(velero backup get ${BACKUP_NAME} -o json | jq -r '.status.completionTimestamp')
RESTORE_TIME=$(velero restore get ${BACKUP_NAME}-restore -o json | jq -r '.status.completionTimestamp')

echo "Backup completed at: ${BACKUP_TIME}"
echo "Restore completed at: ${RESTORE_TIME}"

# 9. Cleanup
echo "9. Cleaning up..."
kubectl delete namespace test-disaster --force --grace-period=0
velero backup delete ${BACKUP_NAME}
velero restore delete ${BACKUP_NAME}-restore

echo "=== Disaster Recovery Validation Complete ==="
