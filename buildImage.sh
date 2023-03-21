IMAGE_TAG=1
podman build --layers=false -t ca.icr.io/mverhiel/ace-webhook-validation:$IMAGE_TAG .
podman push ca.icr.io/mverhiel/ace-webhook-validation:$IMAGE_TAG
