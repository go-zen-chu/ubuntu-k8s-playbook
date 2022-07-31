# ubuntu-k8s-playbook

kubernetes cluster build on ubuntu 22.04

## prerequisite

Programs below should be installed on nodes.

- sshd (OpenSSH): ansible uses ssh for logging into node
- python3: ansible uses python when running playbook on a node

## initialize

After checking prerequisite above, edit .envrc file from sample.envrc and generate host.yml according to your environment.

```bash
cp sample.envrc .envrc
# edit according to your env
vim .envrc
# if you use direnv then run below
direnv allow .
# or if not run below
source .envrc
# envsubst read env vars and substitute them
envsubst < hosts-template.yml > hosts.yml
# initialize all node and create ansible user
ansible-playbook -i hosts.yml initialize.yml
# or you can specify a node you want to setup with --limit option
ansible-playbook -i hosts.yml --limit node1 initialize.yml
```

**After running initialize playbook successfully, ssh port of target nodes will be changed for security. If you want to run initialize playbook again, make sure to change ssh port (ansible_port) in hosts.yml**

## setup control plane node

```bash
ansible-playbook -i hosts.yml k8s-contol-plane.yml
# if you want to try for specific node
ansible-playbook -i hosts.yml --limit control_plane1 k8s-contol-plane.yml
```

## setup worker node

```bash
# run sudo kubeadm token create --print-join-command in control plane node and get latest token & ca-cert
vim .envrc
envsubst < hosts-template.yml > hosts.yml
ansible-playbook -i hosts.yml k8s-worker.yml
# if you want to try for specific node
ansible-playbook -i hosts.yml --limit worker1 k8s-contol-plane.yml
```

## deploy ops components run on k8s

```bash
# apply manifest on one node is enough
$ ansible-playbook -i hosts.yml --limit control_plane1 k8s-ops.yml

# you can access to UI via port-forwarding. https://argo-cd.readthedocs.io/en/stable/getting_started/
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

## update kubernetes

//TBD

## what's next?

You can start progressive delivery using your own k8s manifest repository via [Argo CD](https://argo-cd.readthedocs.io/en/stable/).

Please check [6. Create An Application From A Git Repository](https://argo-cd.readthedocs.io/en/stable/getting_started/#6-create-an-application-from-a-git-repository).
