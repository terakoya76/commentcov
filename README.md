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
  - "**/**_test.go"
mode: scope
```

Then, run commentcov. You could get the comment coverage of the `.go` files in csv format.
```bash
$ commentcov coverage
{"@level":"info","@message":"Install Plugin","@module":"commentcov","@timestamp":"2022-06-27T14:38:53.552977+09:00","plugin":"commentcov-plugin-for-go"}
{"@level":"debug","@message":"starting plugin","@module":"commentcov","@timestamp":"2022-06-27T14:38:56.009873+09:00","args":["commentcov-plugin-go"],"path":"/home/terakoya76/go/bin/commentcov-plugin-go"}
{"@level":"debug","@message":"plugin started","@module":"commentcov","@timestamp":"2022-06-27T14:38:56.010089+09:00","path":"/home/terakoya76/go/bin/commentcov-plugin-go","pid":359483}
{"@level":"debug","@message":"waiting for RPC address","@module":"commentcov","@timestamp":"2022-06-27T14:38:56.010123+09:00","path":"/home/terakoya76/go/bin/commentcov-plugin-go"}
{"@level":"debug","@message":"plugin address","@module":"commentcov.commentcov-plugin-go","@timestamp":"2022-06-27T14:38:56.013775+09:00","address":"/tmp/plugin341157163","network":"unix","timestamp":"2022-06-27T14:38:56.013+0900"}
{"@level":"debug","@message":"using plugin","@module":"commentcov","@timestamp":"2022-06-27T14:38:56.013824+09:00","version":1}
{"@level":"trace","@message":"waiting for stdio data","@module":"commentcov.stdio","@timestamp":"2022-06-27T14:38:56.014527+09:00"}
,FILE,10.67036890122967
,PRIVATE_CLASS,45.61276287164612
,PRIVATE_FUNCTION,30.317776735459663
,PRIVATE_TYPE,42.21311475409836
,PRIVATE_VARIABLE,16.287672723316287
,PUBLIC_CLASS,87.85892224990167
,PUBLIC_FUNCTION,53.6091994076139
,PUBLIC_TYPE,88.03245436105476
,PUBLIC_VARIABLE,69.56668923493568
{"@level":"debug","@message":"received EOF, stopping recv loop","@module":"commentcov.stdio","@timestamp":"2022-06-27T14:38:57.536686+09:00","err":"rpc error: code = Unavailable desc = error reading from server: EOF"}
{"@level":"info","@message":"plugin process exited","@module":"commentcov","@timestamp":"2022-06-27T14:38:57.542710+09:00","path":"/home/terakoya76/go/bin/commentcov-plugin-go","pid":359483}
{"@level":"debug","@message":"plugin exited","@module":"commentcov","@timestamp":"2022-06-27T14:38:57.542746+09:00"}
```
