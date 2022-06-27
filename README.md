# commentcov
pluggable comment coverage generator

## How to use

Specify plugin info and filepaths info on your commentcov configuraiton (default is `<path_to_project>/.commentcov.yaml`).

```bash
$ pwd
/home/commentcov/workspace/github.com/kubernetes/kubernetes

$ cat .commentcov.yaml
plugins:
  - extension: .go
    install_command: go install github.com/commentcov/commentcov-plugin-go@latest
    execute_command: commentcov-plugin-go
target_path: .
exclude_paths:
  - "vendor/**/**"
mode: scope
```

Then, run commentcov. You could get the comment coverage of the `.go` files in csv format.
```bash
$ commentcov coverage
{"@level":"info","@message":"Install Plugin","@module":"commentcov","@timestamp":"2022-06-04T09:49:57.276670+09:00","plugin":"commentcov-plugin-for-go"}
{"@level":"debug","@message":"starting plugin","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.637761+09:00","args":["commentcov-plugin-go"],"path":"/home/commentcov/go/bin/commentcov-plugin-go"}
{"@level":"debug","@message":"plugin started","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.637963+09:00","path":"/home/commentcov/go/bin/commentcov-plugin-go","pid":302317}
{"@level":"debug","@message":"waiting for RPC address","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.638000+09:00","path":"/home/commentcov/go/bin/commentcov-plugin-go"}
{"@level":"debug","@message":"plugin address","@module":"commentcov.commentcov-plugin-go","@timestamp":"2022-06-04T09:49:59.641455+09:00","address":"/tmp/plugin3894134805","network":"unix","timestamp":"2022-06-04T09:49:59.641+0900"}
{"@level":"debug","@message":"using plugin","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.641511+09:00","version":1}
{"@level":"trace","@message":"waiting for stdio data","@module":"commentcov.stdio","@timestamp":"2022-06-04T09:49:59.642178+09:00"}
,PUBLIC_TYPE,86.67992047713717
,PRIVATE_FUNCTION,26.929215170859933
,PUBLIC_FUNCTION,46.19191641462744
,PUBLIC_CLASS,86.1316662413881
,PUBLIC_VARIABLE,69.15447974449488
,PRIVATE_TYPE,36.721311475409834
,FILE,8.397812854637367
,PRIVATE_CLASS,38.396509408235616
,PRIVATE_VARIABLE,14.998234463276836
```
