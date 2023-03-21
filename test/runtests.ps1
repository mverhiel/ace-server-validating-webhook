###############################################################################
# Script that runs ACE validating web hook tests.
#
# If running locally:
# 1. Make tests directory the current directory
# 2. Login to the OCP cluster as a user that can create and delete projects.
#
# If running from an AzDo pipeline:
# 1. Use the redhat.openshift-vsts.oc-setup-task.oc-setup@2 taskto set up the
#    ocp CLI. Use a service connection configured with a service account
#    that can create and delete projects.
# 2. Make the tests directory the working directory. That can be set on the
#    PowerShell task in the pipeline.
#
# Parameters:
# isProductionParam - set to true if a production OCP cluster else set to false
#                     default: false
###############################################################################
param ([string] $isProductionParam = "false", [string] $kubeConfigPath)

Write-Host "[debug] Setting KUBECONFIG environment variable BACK to $kubeConfigPath before running any oc or helm commands"
$env:KUBECONFIG = "$kubeConfigPath"

[bool] $isProduction = [System.Convert]::ToBoolean($isProductionParam)

$project = "ace-webhook-test"
Write-Host "[debug] Creating project $project"

# Need to create the project with the metlife.com/kind: cp4i label so that the
# validating webhook is invoked!
oc process -f ./project.yaml -p EAI_CODE="13574" -p PROJECT_NAME=$project -p PROJECT_DESCRIPTION="ACE Webhook Test Project" -p PROJECT_DISPLAYNAME="ACE Webhook Test Project" | oc create -f -

[int]$numberTestCases = 0
[int]$numberFailures = 0
[int]$numberPassed = 0

try {
  # We are testing a production environment.
  if ( $isProduction ) {
    Write-Host "[debug] We are testing a production environment."
    $failPath = "./production/fail"
    $passPath = "./production/pass"
  }
  # We are testing a non-production enviornment.
  else {
    Write-Host "[debug] We are testing a non-production environment."
    $failPath = "./non-production/fail"
    $passPath = "./non-production/pass"
  }

  $files = Get-ChildItem $failPath

  Write-Host "[debug] Running error cases"

  foreach ($f in $files) {
    $numberTestCases++
    $fullName = $f.FullName
    Write-Host "[debug] Test case $fullName"
    oc process -f $fullname -p PROJECT_NAME=$project | oc apply -f -

    if ( $? ) {
      Write-Host "[error] Test incorrectly passed validation."
      $numberFailures++
    } else {
      Write-Host "[debug] Test case passed."
      $numberPassed++
    }
  }

  Write-Host "[debug] Running valid cases"
  $files = Get-ChildItem $passPath

  Write-Host "[debug] Running valid cases"

  # We are testing the CREATE validation so we need to be sure that the integration server
  # does not exist and is removed after each test case.
  foreach ($f in $files) {
    $numberTestCases++
    $fullName = $f.FullName
    Write-Host "[debug]  Test case $fullName"

    try {
      oc process -f $fullname -p PROJECT_NAME=$project | oc apply -f -

      if ( $? ) {
        Write-Host "[debug] Test case passed."
        $numberPassed++
      } else {
        Write-Host "[error] Test incorrectly failed validation"
        $numberFailures++
      }
    }
    finally {
      Write-Host "[debug] Deleting IntegrationServer"
      oc delete IntegrationServer ace -n $project
    }
  }
}
finally {
  Write-Host "[debug] Removing integration servers in project $project"
  oc delete IntegrationServer -n $project --all

  Write-Host "[debug] Removing integration server configurations in project $project"
  oc delete Configuration -n $project --all

  Write-Host "[debug] Removing project $project"
  oc delete project $project
}

Write-Host "[debug] Total number of test cases: $numberTestCases"
Write-Host "[debug] Total number of failed test cases: $numberFailures"
Write-Host "[debug] Total number of passed test cases: $numberPassed"

if ($numberFailures -gt 0) {
    exit 8
} else {
    exit 0
}