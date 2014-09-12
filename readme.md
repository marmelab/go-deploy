# Deployed PR

Deployed PR is a **webservice** used to create a **comment** on all pull requests contained on a **deployment** to a given **target**.    
The idea is that when deploying project, the deployment tool used should launch a request to Deployed PR (DPR for the next), indicating the **branch** or **tag** deployed, the **URL of the github repository** and the deployment target (prod, preprod, dev …).    
Then DPR will identify all the PR merged on the deployed code (since the last deploy), and automatically write a new comment on **github** PR, indicating that this PR has been deployed and on which target.    
The goal is to help developers to more easily identify which code is deployed, and where.    

## Install

Project depends on some packages (mainly about github api)

    make install

## Configure

You must configure the project on which you want to test DPR. For that, you have to copy config.dist to config.json. When it's done, you must add informations about :

* **Owner**: the github username of the project owner
* **Repository**:   the repository name of the project for which you wish to comment deployed PR 
* **AccessToken**:  a personnal access token (generate on github Settings>Application). This access token is required to access private projects, but mostly to add comments

You can add as many projects as you want to test

## Use DPR

First, you have to launch server :

    go run main.go

Then, you should make a POST request at http://localhost:8080 with Content-Type set to "application/json" with a json file formated as :

    {
        "Owner":"the github username of the project (which must therefore be in config.yml)", 
        "Repo": "the repository name of the project that is deployed (which must also be in config.yml)", 
        "BaseType": "the type of "marker" deployment : branch or tag (Tag doesn't work for the moment!)", 
        "BaseName": "the name of the branch or tag. Ex: master, preprod, v1 ...", 
        "Target": "name of the target (server) on which the code is deployed"
    }

## CAUTION !!! This is a POC
The project is just started, and still in really beta. For example, there is no history for comments, so the same PR will be commented as many times you will call the webservice with identical parameters.     
Also, there is no historical for deployments, so all PR are concerned, not only these which have been merged between two deployments (but it will soon be the case !).
