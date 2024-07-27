# ubuntu-k8s-playbook

kubernetes cluster build on ubuntu 22.04

## prerequisite

Programs below should be installed on Ubuntu.

- sshd (OpenSSH): ansible uses ssh for logging into node
- python3: ansible uses python when running playbook on a node

## initialize baremetal nodes

After checking prerequisite above, edit .envrc file from sample.envrc and generate host.yml according to your environment.

```bash
mage downloadSecretEnvFile
# if you use direnv then run below (make sure you enable .env load option)
direnv allow .
mage generateHostsYaml
# initialize all node and create ansible user
ansible-playbook -i hosts.yml playbook.yml -t initialize 
# or you can specify a node you want to setup with --limit option
ansible-playbook -i hosts.yml --limit node1 playbook.yml -t initialize 
```

**After running initialize playbook successfully, ssh port of target nodes will be changed for security. If you want to run initialize playbook again, make sure to change ssh port (ansible_port) in hosts.yml**

## Install k8s via k3s

### setup control-plane nodes

```bash
ansible-playbook -i hosts.yml playbook.yml -t k3s-control-plane

# if you want to try for specific node
ansible-playbook -i hosts.yml --limit control_plane1 playbook.yml -t k3s-control-plane
```

### setup worker nodes

```bash
# run sudo k3s token create --print-join-command in control plane node and get latest token & ca-cert
ansible-playbook -i hosts.yml k3s-worker.yml

# if you want to try for specific node
ansible-playbook -i hosts.yml --limit worker1 playbook.yml -t k3s-worker
```

## deploy Argo CD

```bash
# apply manifest on one node is sufficient
ansible-playbook -i hosts.yml --limit control_plane1 playbook.yml -t k8s-argocd
# you can access to UI via port-forwarding. https://argo-cd.readthedocs.io/en/stable/getting_started/
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

## upgrate kubernetes

Upgrading k8s can be done automatically using [system-upgrade-controller](https://docs.k3s.io/upgrades/automated)

```bash
ansible-playbook -i hosts.yml --limit control_plane1 playbook.yml -t system-upgrade-controller
```

## what's next?

You can start progressive delivery using your own k8s manifest repository via [Argo CD](https://argo-cd.readthedocs.io/en/stable/).

Please check [6. Create An Application From A Git Repository](https://argo-cd.readthedocs.io/en/stable/getting_started/#6-create-an-application-from-a-git-repository).
