#!/bin/bash

NEXT_RELEASE_VERSION="${1#"v"}"

function update_version_helm_chart_files() {
  echo "It's going to update appVersion in helm chart to version ${NEXT_RELEASE_VERSION}"
  yq w -i helm/check-christmas-lottery-numbers/Chart.yaml "appVersion" "${NEXT_RELEASE_VERSION}"

  echo "It's going to update version in helm chart to version ${NEXT_RELEASE_VERSION}"
  yq w -i helm/check-christmas-lottery-numbers/Chart.yaml "version" "${NEXT_RELEASE_VERSION}"

  echo "It's going to update image tag in helm chart values to version ${NEXT_RELEASE_VERSION}"
  yq w -i helm/check-christmas-lottery-numbers/values.yaml "image.tag" "${NEXT_RELEASE_VERSION}"
}

function update_version_package_json() {
  echo "It's going to update version in package.json to version ${NEXT_RELEASE_VERSION}"
  jq --arg NEXT_RELEASE_VERSION "$NEXT_RELEASE_VERSION" '.version = $NEXT_RELEASE_VERSION' package.json > /tmp/package.json
  mv /tmp/package.json package.json
}

function update_versions() {
  update_version_package_json
  update_version_helm_chart_files
}

update_versions