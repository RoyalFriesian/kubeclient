# kubeclient

Sample commands

## To find all the pods in all the namespaces
go run . --kubeconfig "C:/Users/tarun/.kube/config" 

go run . --kubeconfig "C:/Users/tarun/.kube/config" --name "kube-proxy-5l4wz"
# To find a perticular pod

go run . --kubeconfig "C:/Users/tarun/.kube/config" --namespace "foo"
# To find all the pods in a given namespace

go run . --kubeconfig "C:/Users/tarun/.kube/config" --namespace "foo" --name "kube-proxy-514wz"
# To find the pod within the given namespace

go run . --kubeconfig "C:/Users/tarun/.kube/config"  --namespace "foo" --age 10  
# To find all the pods running for minimum 10 seconds and within the given namespace

go run . --kubeconfig "C:/Users/tarun/.kube/config"  --age 10  
# To find all the pods running for minimum 10 seconds across all the namespaces
