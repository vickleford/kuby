kuby
===

Ever had to manage multiple Kubernetes clusters with differing versions?

Meet kuby.

Installation
---

go get github.com/vickleford/kuby

Usage
---

Use kuby the same you would use kubectl, except quit worrying about what version it is. Kuby finds out which version of kubectl it needs by making a call to the version endpoint of a Kubernetes server then passes through the arguments to the matching kubectl version.

```
kuby --context context1 deploy -f whatever.yaml
```

WIP
---

Obviously, this is a work in progress, but hopefully this illustrates the vision and direction!

Contributing
---

kuby uses Go Modules. Please use Go 1.11 and clone this repository outside `GOPATH`.