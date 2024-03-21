#!/bin/bash

git clone https://github.com/prometheus-operator/kube-prometheus.git tmp/kube-prometheus
k apply --server-side -f tmp/kube-prometheus/manifests/setup  
sed -i 's/replicas: [0-9]*/replicas: 1/' tmp/kube-prometheus/manifests/prometheus-prometheus.yaml
kubectl apply -f 'tmp/kube-prometheus/manifests/prometheusOperator-*'

# Add permissions to prometheus in its cluster role to scrape endpoints

kubectl apply -f 'tmp/kube-prometheus/manifests/prometheus-*'
kubectl apply -f 'resources/prometheus/*'