someutils
=========

Some CLI utilities written in Go.

 * Mainly intended as Unix-like commands for Windows, but cross-platform anyway. 
 * Covers similar ground to coreutils, but not intended as a replacement. (Won't ever support all commands & options).
 * Just because.

Installation.
---------

### Method 1: download precompiled binaries (recommended for when you're not currently running go):

 * Grab some recent (v0.3.0) Windows binaries zipped up for [32-bit Windows](http://dl.bintray.com/laher/utils/someutils_0.3.0_windows_386.zip) and [64-bit Windows](http://dl.bintray.com/laher/utils/someutils_0.3.0_windows_amd64.zip). 
 * Just unzip somewhere on your %PATH% (environment variable). 
 * (These binaries were built and uploaded with `goxc`, ofcourse)

### Method 2: `go get` all binaries into your GOPATH
   
NOTE: On Unix systems in particular, **be careful that your system PATH elements come BEFORE `GOPATH` within your PATH environment variable**

   1. `go get github.com/laher/someutils/./...`
   2. `ls`, `pwd` etc

You could also use `go get` to pick out a subset of those commands.

### Method 3: just install the `some` command, (similar to `busybox`):

   1. `go get github.com/laher/someutils/cmd/some`
   2. `some ls`, `some pwd` etc
   3. Optionally, use `alias` or `doskey` to make `some` behave more like `busybox`.


Scope etc
---------
My initial target is to get my Windows CLI a bit closer to being as productive as my Linux CLI, by creating many small utilities under one umbrella. Then, who knows. I do like the idea of a Go/Linux, as opposed to Gnu/Linux :)

Some commands are not included, because they're either ubiquitous anyway (such as echo,cd,whoami), just too big, or hard to acheive with pure Go.

I'll just keep adding stuff as I need it. Contributions welcome!

### Progress

So far, limited versions of the following commands are available:
You can also use 'some [cmd] [args...]' for any of these.
 
 Command | Options supported | STDIN support  | Notes
 --------|-------------------|----------------|------------------------
 basename|                   |                | TODO: -a, -z
 cat     | -Ens              | Yes            | 
 cp      | -r                | n/a            | TODO: check symlink behaviour. Test large file support
 dirname |                   | n/a            |
 grep    | -nvHi -E -P       | Yes            | TODO: binary files support. !!No support for BRE - uses -E by default.
 gunzip  | -k                | TODO           | TODO: -f, prompt when file exists
 gzip    | -k                | TODO           | TODO: -f, prompt when file exists
 head    | -n                | Yes            | TODO: -c
 ls      | -lahr -1          | Yes            | TODO: -p -t
 mv      |                   | n/a            | TODO: check symlink behaviour
 pwd     |                   | n/a            | 
 rm      | -r                | n/a            | TODO: check symlink behaviour
 scp     | -r -P             | ?              | INCOMPLETE - see [scp-go](https://github.com/laher/scp-go) .
 sleep   |                   | n/a            |
 tail    | -n -F             | Yes            | TODO: -c, -f (by descriptor rather than by name). Bug: won't currently print last line unless terminated by a CR.
 tar     | -cvf -x -t -r     | Yes (IN+OUT)   | Just the core functionality so far.
 tee     | -a                | Yes            | TODO: -i
 touch   |                   | n/a            | 
 unzip   | -t                | TODO(STDOUT)   | Password support would not be straightforward (not supported by standard lib)
 wc      | -c -l -w          |                | 
 which   | -a                | n/a            | 
 wget    | -c -o             | n/a            | TODO: multi-threading? (not part of real wget). See [wget-go](https://github.com/laher/wget-go)
 zip     |                   | TODO           | Password support would not be straightforward (not supported by standard lib)
 

### ToMaybeDo
 * stat,size,file,type
 * split,join
 * chmod/chown (relevant? Yes I think so)
 * diff (too big? Maybe a minimal version would be good here)
 * more (how easy is it?)
 * du/dh (need OS-specifics: syscall would probably cover it for Unix and Windows)
 * find/locate (find is a bit of a monster. locate is probably a stretch)
 * ln (would it need some non-Go stuff for Windows? Yes - maybe an 'exec' at this stage)
 * ps,kill,pgrep,pkill (need to explore mileage of os.FindProcess, syscall.Kill)
 * id,w (is it doable cross-platform?)
 * sshd (minimal version, for hosting file transfers etc), ssh (maybe just for running remote commands. Terminal handling might be too challenging for now)
 * traceroute (see wtn. Requires setuid & therefore chowning + chmodding on Unix - on Windows I think you'd just need to run as administrator)
 * ping (see above)
 * dig (I think. Raw DNS requests & collect responses. Hmm, investigate go.net packages)
 * chroot (chroot possible for unix via syscall - see gobox)
 
### TooBig?
 * less
 * a text editor
 * top
 * dd (I guess. Maybe not)
 * awk, sed
 * xargs, find -exec
 * rsync (algorithms might be a bit hard)
 * cron (I guess service handling is another chapter aswell)
 
### Not possible/easy with pure Go
 * bg,fg 
 * fsck

See Also
--------

 * I have separated out the 'flag' functionality, and some other relevant behaviour, into a separate package, [uggo](https://github.com/laher/uggo)
 * I drew inspiration early on from a couple of similar projects - many thanks to [gobox](https://github.com/surma/gobox) and [go-coreutils](https://github.com/sepeth/go-coreutils). In both cases I considered a fork but my focus is just a little too different to make it feasible. Cheers guys
