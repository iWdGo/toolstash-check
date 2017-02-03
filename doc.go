// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Toolstash-check automates running toolstash -cmp against a CL.

Usage:
	toolstash-check [-all] [commit]

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
golang.org/x/tools/cmd/toolstash/buildall for the final step.
*/
package main
