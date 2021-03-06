#!/usr/bin/env groovy

//Author: Carlos Alvarado
//Jenkins Pipeline to handle the Continuous Integration and Continuous Deployment with Golang and ArgoCD


node('worker0') {
    
    env.HOME = "${WORKSPACE}"
    env.CONTAINER_IMAGE = 'docker.io/calvarado2004/bookings'
    env.CONTAINER_IMAGE_INIT = 'docker.io/calvarado2004/bookings-init'
    env.GIT_PROJECT_CI = 'github.com/calvarado2004/bookings.git'
    env.GIT_PROJECT_CD = 'github.com/calvarado2004/bookings-cd.git'

    stage ('Download the source code from GitHub'){
            git branch: 'main', url: "https://${GIT_PROJECT_CI}"
    }
    
    stage ('Build Container'){
    
           sh "sudo buildah bud -f Dockerfile -t ${CONTAINER_IMAGE}:${BUILD_TAG}"
           sh "sudo buildah bud -f Dockerfile -t ${CONTAINER_IMAGE}:latest"

    }
    
    stage ('Build Init Container') {
           sh "sudo buildah bud -f Dockerfile-init -t ${CONTAINER_IMAGE_INIT}:${BUILD_TAG}"
           sh "sudo buildah bud -f Dockerfile-init -t ${CONTAINER_IMAGE_INIT}:latest"
    }

    stage ('Push Container Images'){
        withCredentials([usernamePassword(credentialsId: 'docker-hub', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
            
            sh "sudo buildah login --username ${USERNAME} --password ${PASSWORD} docker.io"
            sh "sudo buildah push  ${CONTAINER_IMAGE}:${BUILD_TAG}"
            sh "sudo buildah push  ${CONTAINER_IMAGE}:latest"
            sh "sudo buildah push  ${CONTAINER_IMAGE_INIT}:${BUILD_TAG}"
            sh "sudo buildah push  ${CONTAINER_IMAGE_INIT}:latest"
        }
    }
    
    stage ('Modify Bookings image on Deployment') {
      withCredentials([string(credentialsId: 'git-token', variable: 'TOKEN')]) {
        git branch: 'main', url: "https://${GIT_PROJECT_CD}"
        sh 'cat manifest/app-deployment.j2 | sed "s#{{ CONTAINER_IMAGE }}:{{ USED_TAG }}#${CONTAINER_IMAGE}:${BUILD_TAG}#g" > manifest/app-deployment-temp.yaml'
        sh 'cat manifest/app-deployment-temp.yaml | sed "s#{{ CONTAINER_IMAGE_INIT }}#${CONTAINER_IMAGE_INIT}#g" > manifest/app-deployment.yaml && rm -f manifest/app-deployment-temp.yaml'
        sh "git add . && git commit -m 'Change of image version of Golang App to ${BUILD_TAG}'"
        sh "git push https://${TOKEN}@${GIT_PROJECT_CD}"
      }
    }
    
} 