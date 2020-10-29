# Notes on contributing to Kubernetes project

## Resources

- [Development Guide](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md)
- [git workflow](https://github.com/kubernetes/community/blob/master/contributors/guide/github-workflow.md)
- [Contributor Summit NA 2019: Keeping the Bar High - How to be a bad ass Code Reviewer](https://www.youtube.com/watch?v=OZVv7-o8i40)

## Understand vanity redirects 

- go imports (k8s.io/foo --> github.com/kubernetes/foo and sigs.k8s.io/bar --> github.com/kubernetes-sigs/bar)
- gitHub website redirects (git.k8s.io --> github.com/kubernetes) 

## Prereqs

- git
- docker
- direnv
- gimme
- go

# Setup develop environment

```
# Install direnv by following instructions here. 
https://direnv.net/

# Install gimmie

mkdir ~/bin
curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme

gimmie -k # list of go version it can install
gimmie -l # current version

mkdir k8s-development
cd k8s-development
mkdir -p go/src

gimme stable

# gimme stable >> .envrc if you are using direnv

vi .envrc
export GOPATH="${PWD}/go"
export PATH=$PATH:${GOPATH}/bin

mkdir $GOPATH/src/k8s.io
cd $GOPATH/src/k8s.io
git clone git@github.com:$GITHUB_USER/kubernetes.git
cd $GOPATH/src/k8s.io/kubernetes
git remote add upstream https://github.com/kubernetes/kubernetes.git
git remote set-url --push upstream no_push
git fetch upstream
git checkout master
git rebase upstream/master

# install kubetest
cd $GOPATH/src/k8s.io
git clone https://github.com/kubernetes/test-infra.git
cd test-infra/
GO111MODULE=on go install ./kubetest

```

## Confirm setup

Compile a few of the binaries

```
# Compile kubectl
make WHAT=cmd/kubectl
_output/bin/kubectl version

# Compile other components
make WHAT=cmd/kubelet
_output/bin/kubelet -h

make WHAT=cmd/kube-scheduler
_output/bin/kube-scheduler -h

```

# Code

- Create branch for work `git checkout -b feature-branch-name`
- Make changes
- Test changes using kind (see steps below)
- Add changes to repo `git add <file>`
- Commit changes `git commit -m "ISSUE: Updated ..."
- Push changes to origin `git push origin feature-branch-name`
- Create Pull Request
- Additional modifications `git commit --amend`
- Push up modifications `git push -f feature-branch-name`

## Test code changes

Create a container images from the code changes made.

```
# Create a node-image.
kind build node-image --image=feature-branch-name

# Create cluster from feature branch.
kind create cluster --image=feature-branch-name

Run e2e tests
# --provider string    Kubernetes provider such as gce, gke, aws, etc
# --test               Run Ginkgo tests.
# --extract exactStrategies       Extract k8s binaries from the specified release location
# --test_args string        Space-separated list of arguments to pass to Ginkgo test runner.

K8S_VERSION=$(kubectl version -o json | jq -r '.serverVersion.gitVersion')
export KUBERNETES_CONFORMANCE_TEST=y
export KUBECONFIG="$HOME/.kube/config

# took a while to complete.
kubetest --provider=skeleton \
	--test \
	--test_args=”--ginkgo.focus=\[Conformance\]” \
	--extract ${K8S_VERSION} | tee test.out

```


## Non-Code contribution

- https://github.com/kubernetes/community/blob/master/contributors/guide/non-code-contributions.md
- [Release team roles](https://github.com/kubernetes/sig-release/tree/master/release-team)
- [Release Special Interest Group](https://github.com/kubernetes/community/tree/master/sig-release)
- [Shadow Roles Throughout the Kubernetes Ecosystem](https://github.com/kubernetes/community/blob/master/mentoring/programs/shadow-roles.md)

