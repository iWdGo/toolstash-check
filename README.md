[![Go Reference](https://pkg.go.dev/badge/iwdgo/toolstash-check.svg)](https://pkg.go.dev/iwdgo/toolstash-check)

# Using toolstash-check locally

`go get [-u] github.com/iwdgo/toolstash-check`

Toolstash-check automates running toolstash -cmp against a CL.

Usage:
`toolstash-check [options] [commit]`

Toolstash-check automates the following workflow for compiler
regression testing a CL:

1. Clone the $GOROOT Git repo into a new temporary directory.

2. Checks out the specified commit's parent.

3. Runs make.bash and toolstash save. (It sets GOROOT so that
   toolstash saves into the temporary directory.)

4. Checks out the specified commit itself.

5. Runs "go install std cmd".

6. Runs "go build -a -toolexec='toolstash -cmp' std cmd".

If no commit ID is specified, toolstash-check defaults to HEAD.

If -all is specified, toolstash-check instead runs
`golang.org/x/tools/cmd/toolstash/buildall` for the final step.

# Using Continuous Integration

## Flag `-repo <filepath to repo>`

Flag allows set file path to repository which defaults to `GOROOT`.

## Flag `-v`

Logs key steps for debugging purposes.

## Travis CI

The `compiler-travis` branch of [go-upon-ci](https://github.com/iWdGo/go-upon-ci) runs this fork of 
[toolstash-check](github.com/mdempsky/toolstash-check) on tip of
 [golang/go](https://github.com/golang/go) repository.

Patch files are produced using `git format-patch <your commit>` or `git format-patch -1`.

After `commit`, `git push` to your account to trigger CI.
Travis CI detects your file after account setup and granting rights to the repository.
