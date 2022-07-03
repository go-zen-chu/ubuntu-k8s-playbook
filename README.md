# ubuntu-k8s-playbook

kubernetes cluster build on ubuntu 22.04

## prerequisite

Programs below should be installed manually on nodes manually before running this playbook.

- sshd: ansible uses ssh for logging into node
- python3: ansible uses python when running playbook on a node
- dhcpcd: optional. network connection from ansible origin is required

```bash
# TBD
```

## initialize

After the manual installation above, edit .envrc file from sample.envrc and generate host.yml according to your environment.

```bash
# edit according to your env
vim .envrc
# if you use direnv then run below
direnv allow .
# if not run below
source .envrc
# envsubst read env vars and substitute them
envsubst < hosts-template.yml > hosts.yml
# specify a node you want to setup 
ansible-playbook -i hosts.yml --limit node1 initialize.yml
```

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
