# Image clone controller

## Goal
We’d like to be safe against the risk of public container images disappearing from the registry while we use them, breaking our deployments. 

## Problem
We have a Kubernetes cluster on which we can run applications. These applications will often use publicly available container images, like official images of popular programs, e.g. Jenkins, PostgreSQL, and so on. Since the images reside in repositories over which we have no control, it is possible that the owner of the repo deletes the image while our pods are configured to use it. In the case of a subsequent node rotation, the locally cached copies of the images would be deleted and Kubernetes would be unable to re-download them in order to re-provision the applications. 

# Demo  
[![asciicast](https://asciinema.org/a/c951Ym7zJfuZVXiV0EV27Sx28.svg)](https://asciinema.org/a/c951Ym7zJfuZVXiV0EV27Sx28)

## Idea
Have a controller which watches the applications and “caches” the images by re-uploading to our own registry repository and reconfiguring the applications to use these copies.

## Installing the controller

```bash
   kubectl apply -f config/manager.yaml 
   kubectl apply -f config/rbac/leader_election_role.yaml 
   kubectl apply -f config/rbac/leader_election_role.yaml
```
