* unit tests:

```
[shajmakh@shajmakh numalign-rewritten (main)]$ go test -v ./...
?       github.com/shajmakh/numaalign-rewritten/cmd     [no test files]
?       github.com/shajmakh/numaalign-rewritten/internal        [no test files]
?       github.com/shajmakh/numaalign-rewritten/pkg/numa        [no test files]
=== RUN   TestCheckNumaCpuMapping
--- PASS: TestCheckNumaCpuMapping (0.00s)
PASS
ok      github.com/shajmakh/numaalign-rewritten/internal/cpu    (cached)
=== RUN   TestCheckPciDevicesAlignment
--- PASS: TestCheckPciDevicesAlignment (0.00s)
PASS
ok      github.com/shajmakh/numaalign-rewritten/internal/device (cached)
=== RUN   TestCheckAlignmentWith
--- PASS: TestCheckAlignmentWith (0.00s)
PASS
ok      github.com/shajmakh/numaalign-rewritten/internal/memory (cached)
=== RUN   TestResourcesNumaAlign
--- PASS: TestResourcesNumaAlign (0.00s)
PASS
ok      github.com/shajmakh/numaalign-rewritten/internal/tests  (cached)
[shajmakh@shajmakh numalign-rewritten (main)]$
```


* manual tests on k8s cluster:

topology manager policy set to single-numa-node.
using this image: quay.io/rhn_support_shajmakh/numalign-rewritten:latest.

```
[shajmakh@shajmakh p1-k8s-testing]$ oc get po 
NAME    READY   STATUS    RESTARTS   AGE
be      1/1     Running   0          7s
bu      1/1     Running   0          41s
pod-1   1/1     Running   0          3m37s
```

gu-pod:

```
[shajmakh@shajmakh p1-k8s-testing]$ oc logs pod-1
Is Aligned: true
NUMA: 0
```

best-effort and burstable pods:

```
[shajmakh@shajmakh p1-k8s-testing]$ oc logs bu
Is Aligned: false
NUMA: -1

[shajmakh@shajmakh p1-k8s-testing]$ oc logs be
Is Aligned: false
NUMA: -1
```
