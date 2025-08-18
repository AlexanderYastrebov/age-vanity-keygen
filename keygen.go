// Copyright (c) 2025, Alexander Yastrebov
// Copyright (c) 2019, The age Authors
//
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"golang.org/x/term"
)

const usage = `Usage:
    age-vanity-keygen [-o OUTPUT] PREFIX

Options:
    -o, --output OUTPUT       Write the result to the file at path OUTPUT.
	

age-vanity-keygen generates a new vanity X25519 key pair with recipient PREFIX,
and outputs it to standard output or to the OUTPUT file.

PREFIX can not contain character '1'. It is transformed to lower case,
characters 'b', 'i' and 'o' are replaced for '6', '7' and '0' respectively.

If an OUTPUT file is specified, the public key is printed to standard error.
If OUTPUT already exists, it is not overwritten.

Examples:

    $ age-vanity-keygen 23456
    Found age123456... in 0s after 15855390 attempts (43267922 attempts/s)
    # created: 2025-08-18T18:18:18+02:00
    # public key: age123456gpgacec4alqvqnfdacx6djhx98wzwn4l3eh5q5n5ec2evdsfzn7tn
    AGE-SECRET-KEY-1XRTF5T02CR2HEC29RAH29Y46DPHQ7EAPK5EEPYKTFE3682LPWSCS4CXJSX

    $ age-vanity-keygen -o key.txt 23456
    Found age123456... in 2s after 74446587 attempts (43327977 attempts/s)
    Public key: age123456l7nurcmk5xfp5009lu65drh2ull7hghpkd0xlp3f4l7vv4sg8fdga
`

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module. See
// golang.org/issue/29814 and golang.org/issue/29228.
var Version string

func main() {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	var (
		versionFlag bool
		outFlag     string
	)

	flag.BoolVar(&versionFlag, "version", false, "print the version")
	flag.BoolVar(&versionFlag, "v", false, "print the version")
	flag.StringVar(&outFlag, "o", "", "output to `FILE` (default stdout)")
	flag.StringVar(&outFlag, "output", "", "output to `FILE` (default stdout)")
	flag.Parse()

	if versionFlag {
		if Version != "" {
			fmt.Println(Version)
			return
		}
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(buildInfo.Main.Version)
			return
		}
		fmt.Println("(unknown)")
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "PREFIX required")
		flag.Usage()
		os.Exit(1)
	}

	prefix := flag.Arg(0)

	out := os.Stdout
	if outFlag != "" {
		f, err := os.OpenFile(outFlag, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			errorf("failed to open output file %q: %v", outFlag, err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				errorf("failed to close output file %q: %v", outFlag, err)
			}
		}()
		out = f
	}

	if fi, err := out.Stat(); err == nil && fi.Mode().IsRegular() && fi.Mode().Perm()&0o004 != 0 {
		warning("writing secret key to a world-readable file")
	}
	generate(out, prefix)
}

func generate(out *os.File, prefix string) {
	k, err := generateX25519Identity(prefix)
	if err != nil {
		errorf("internal error: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Found %s... in %s after %d attempts (%.0f attempts/s)\n",
		k.prefix, k.elapsed.Round(time.Second), k.attempts, float64(k.attempts)/k.elapsed.Seconds())

	if !term.IsTerminal(int(out.Fd())) {
		fmt.Fprintf(os.Stderr, "Public key: %s\n", k.Recipient())
	}

	fmt.Fprintf(out, "# created: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(out, "# public key: %s\n", k.Recipient())
	fmt.Fprintf(out, "%s\n", k)
}

func errorf(format string, v ...interface{}) {
	log.Printf("age-vanity-keygen: error: "+format, v...)
	log.Fatalf("age-vanity-keygen: report unexpected or unhelpful errors at https://github.com/AlexanderYastrebov/age-vanity-keygen")
}

func warning(msg string) {
	log.Printf("age-vanity-keygen: warning: %s", msg)
}
