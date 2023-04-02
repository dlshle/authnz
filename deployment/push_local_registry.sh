#!/bin/sh

cd ..
APP=${1:-authz}
CTX=${1:-local-k3s}

sudo docker build -t $APP .

LOCAL_REGISTRY=192.168.0.164:15000/tucom/$APP

sudo docker tag $APP $LOCAL_REGISTRY
if [ $? -ne 0 ];then
    echo 'tagging failed, exiting'
    exit 0
fi

sudo docker push $LOCAL_REGISTRY

kubectl config use-context $CTX

cd deployment