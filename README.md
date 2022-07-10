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
ansible-playbook -i hosts.yml --limit node1 initialize.yml
# or you can specify a node you want to setup with --limit option
ansible-playbook -i hosts.yml --limit node1 initialize.yml
```

**After running initialize playbook successfully, ssh port of target nodes will be changed for security. If you want to run initialize playbook again, make sure to change ssh port (ansible_port) in hosts.yml**

## setup control plane node

```bash
ansible-playbook -i hosts.yml k8s-contol-plane.yml
```

## setup worker node

```bash
ansible-playbook -i hosts.yml k8s-worker.yml
```

## update kubernetes

//TBD

## what's next?

[Argo CD](https://argo-cd.readthedocs.io/en/stable/) will be also installed in this playbook. You can start progressive delivery using your own k8s manifest repository.
