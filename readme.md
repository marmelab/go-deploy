# Deployed PR

Deployed PR is a **webservice** used to create a **comment** on all pull requests contained on a **deployment** to a given **target**.    
The idea is that when deploying project, the deployment tool used should launch a request to Deployed PR (DPR for the next), indicating the **branch** or **tag** deployed, the **URL of the github repository** and the deployment target (prod, preprod, dev …).    
Then DPR will identify all the PR merged on the deployed code (since the last deploy), and automatically write a new comment on **github** PR, indicating that this PR has been deployed and on which target.    
The goal is to help developers to more easily identify which code is deployed, and where.    

## This is a POC
The project is just started, and do nothing for the moment.    
