
# go-kafka-on-helm

Goal of this repository is to practice helm charts. There will be three helm charts
* Go-Kafka
* Go-Producer
* Go-Consumer

# go-producer

more details in `ci-script.sh` - test, build binary and docker image, push to registry  
`ci-script.sh` should be executed by CI infrastructure  

# go-consumer

more details in `ci-script.sh` - test, build binary and docker image, push to registry  
`ci-script.sh` should be executed by CI infrastructure  

# helm

deployment with helm chart see
* `go-apps-install.sh`
* `kafka-install.sh`

# Useful Links

[deploy go app with helm on kubernetes tutorial](https://docs.bitnami.com/kubernetes/how-to/deploy-go-application-kubernetes-helm/)  
[golang project structure](https://github.com/golang-standards/project-layout)