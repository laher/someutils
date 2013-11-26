someutils
=========

Some CLI utilities written in Go.

 * Mainly intended as Unix-like commands for Windows, but cross-platform anyway. 
 * Covers similar ground to coreutils, but not intended as a replacement. (Won't ever support all commands & options).
 * Just because.

Installation.
---------

### Method 1: download precompiled binaries

 * Grab some recent (v0.3.0) Windows binaries zipped up for [32-bit Windows](http://dl.bintray.com/laher/utils/someutils_0.3.0_windows_386.zip) and [64-bit Windows](http://dl.bintray.com/laher/utils/someutils_0.3.0_windows_amd64.zip). 
 * Just unzip somewhere on your %PATH% (environment variable). 
 * (These binaries were built and uploaded with `goxc`, ofcourse)

### Method 2: `some`, recommended for Unix systems

   1. `go get github.com/laher/someutils/cmd/some`
   2. `some ls`, `some pwd` etc

### Method 3: `go install ./...` installs all binaries into your GOPATH
   
NOTE: On Unix systems in particular, **be careful that your system PATH elements come BEFORE `GOPATH` within your PATH environment variable**

   1. `go get github.com/laher/someutils`
   2. `cd <....>/someutils`
   3. `go install ./...`
   4. `ls`, `pwd` etc


Scope etc
---------
My main target is to get my Windows CLI a bit closer to being as productive as my Linux CLI, by creating many small utilities under one umbrella.
Some commands are excluded because they're either ubiquitous anyway (such as echo,cd,whoami), just too big, or hard to acheive with pure Go.

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

 * I have separated out the 'flag' functionality into a separate package, [uggo](https://github.com/laher/uggo)
