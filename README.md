someutils
=========

Some CLI utilities written in Go.

 * Mainly intended as Unix-like commands for Windows, but cross-platform anyway. 
 * Covers similar ground to coreutils, but not intended as a replacement. (Won't ever support all commands & options).
 * Just because.

Scope etc
---------
My target is to get my Windows CLI a bit closer to being as productive as my Linux CLI, by creating many small utilities under one umbrella.
Some commands are ubiquitous anyway (such as echo,cd,whoami), some are just too big, and some are hard to acheive with pure Go:
 
### Progress

So far, limited versions of the following commands are available:
 
 Command | Options supported | STDIN support  | TODO (or !!WONTDO)
 --------|-------------------|----------------|------------------------
 cat     | -Ens              | Yes            | 
 cp      | -r                | n/a            | check symlink behaviour
 ls      | -lahr -1          | Yes            | -p -t
 grep    | -nvHi -E -P       | Yes            | TODO: binary files support. !!No support for BRE - uses -E by default.
 mv      |                   | n/a            | check symlink behaviour
 pwd     |                   | n/a            | 
 rm      | -r                | n/a            | check symlink behaviour
 touch   |                   | n/a            | 
 unzip   | -t                | TODO(STDOUT)   | 
 which   | -a                | n/a            | 
 zip     |                   | TODO           |  
 
You can also use 'some [cmd] [args...]' for any of the above.

### ToMaybeDo
 * tar,gzip,gunzip
 * stat,size,file,split,type
 * tee,split,join,head,tail (tail -f??)
 * chmod/chown (relevant?)
 * diff (too big? Maybe a minimal version would be good here)
 * more (how easy is it?)
 * du/dh
 * find/locate
 * ln (would it need some non-Go stuff for Windows?)
 * ps,kill,pgrep,pkill
 * id,w,which
 * tailf (tail -f). How easy would this be?
 * scp
 * ssh (terminal handling might be too challenging)

### TooBig
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