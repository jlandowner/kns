# kns - Easy switch Kubernetes namespaces in kubeconfig
This tool can change namespaces easily.
You can change namespace seemlessly in kubectl commands operation.

## How to install
```
$ git clone xxx
$ cp -p ./bin/kns /usr/local/bin/
$ export PATH=$PATH:/usr/local/bin/
$ kns help
```
## How to use
### Interactive action & List available namespaces

```
$ kns
** List of Namespaces in the Current-context Cluster.
0 :  default
1 :  kube-system
2 :  registry
** Which namespace do you want to switch? (exit: q)
Select[n] => 1
** Completed: Switch namespace  kube-system
```

### Select the Specific namespace by a param
```
$ kns exist-namepsace
** Completed: Switch namespace  exist-namepsace

$ kns not-exist
Namespace not-exist does NOT Exist in the Cluster.
```

### Templete actions
#### Switch to default namespace
```
$ kns default
$ kns reset
'** Completed: Switch namespace  default
```

#### Switch to kube-system namespace
```
$ kns kube-system
$ kns kube
$ kns system
$ kns sys
** Completed: Switch namespace  kube-system
```

### Others
$ kns help
$ kns --help
$ kns version
