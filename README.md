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
 
 Command | Options supported | STDIN support  | Notes
 --------|-------------------|----------------|------------------------
 cat     | -Ens              | Yes            | 
 cp      | -r                | n/a            | TODO: check symlink behaviour
 grep    | -nvHi -E -P       | Yes            | TODO: binary files support. !!No support for BRE - uses -E by default.
 ls      | -lahr -1          | Yes            | TODO: -p -t
 mv      |                   | n/a            | TODO: check symlink behaviour
 pwd     |                   | n/a            | 
 rm      | -r                | n/a            | TODO: check symlink behaviour
 [scp](https://github.com/laher/scp-go)     | -r -P             | ?              | INCOMPLETE AND NOT WORKING WITH ALL SSH SERVERS.
 touch   |                   | n/a            | 
 unzip   | -t                | TODO(STDOUT)   | 
 which   | -a                | n/a            | TODO: Windows treats current dir above PATH variables.
 [wget](https://github.com/laher/wget-go)    | -c -o             | n/a            | TODO: multi-threading? (not part of real wget)
 zip     |                   | TODO           | No password support. 
 
You can also use 'some [cmd] [args...]' for any of the above.

### ToMaybeDo
 * tar,gzip,gunzip
 * stat,size,file,split,type
 * tee,split,join,head,tail (tail -f??)
 * chmod/chown (relevant?)
 * diff (too big? Maybe a minimal version would be good here)
 * more (how easy is it?)
 * du/dh (need OS-specifics?)
 * find/locate
 * ln (would it need some non-Go stuff for Windows? YES at this stage)
 * ps,kill,pgrep,pkill
 * id,w,which
 * tailf (tail -f). How easy would this be?
 * ssh (terminal handling might be too challenging)

### TooBig?
 * less
 * a text editor
 * top
 * dd
 * awk, sed
 * xargs, find -exec
 * rsync (algorithms might be a bit hard)
 
### Not possible/easy with pure Go
 * cron,chroot,bg,fg
 * ps
 * fsck
 * dig (unless I can construct raw DNS requests & collect responses)

See Also
--------

 * I have separated out the 'flag' functionality, and some other relevant behaviour, into a separate package, [uggo](https://github.com/laher/uggo)
 * I drew inspiration early on from a couple of similar projects - many thanks to [gobox](https://github.com/surma/gobox) and [go-coreutils](https://github.com/sepeth/go-coreutils). In both cases I considered a fork but my focus is just a little too different to make it feasible. Cheers guys
