# Exhume - Dig up and export your dead LiveJournal

## Importing your LiveJournal into the Hugo site generator

If you're like me, you can no longer use LiveJournal given their
latest terms of use.

I wrote this tool to import my LiveJournal into a new blog.  I chose
the [Hugo](http://gohugo.io) blogging engine, as it is super fast,
runs on any PC, and you can host the output on AWS or a bunch of free hosting sites.

You should dump your blog
with [LJDump](http://hewgill.com/ljdump/)
which will produce a file for each entry named (eg.) L-99 and a file
for each Comment.

Comments (in files C-99 etc.) will be read and processed if present.
Threading is not implemented.   Spam, Banned and deleted comments will
be ignored by default (use `exhume post --help` to see how to change this).

If you run this script on the entry files, it will output markdown
files.

## Example of use

First install and configure [Hugo](http://gohugo.io)
and [LJDump](http://hewgill.com/ljdump/) from their authors.

You should read about how to install themes for hugo.

Next, install exhume

```
   go get -v github.com/unixbigot/exhume
```

The main event, dump and convert your LiveJournal:

```
   hugo new site exhumed-journal
   mkdir exhumed-journal/import
   cd exhumed-journal/import
   ljdump.py
   exhume post L-*
   mv L-*.md ../content/post
   cd ..
   hugo server
   # Now open http://localhost:1313/ to see if it worked
```

Once you've done the above, you can install a theme for 
Hugo, and upload your blog.

```
   # Generate a static website
   hugo --theme my-chosen-theme
   
   # Upload the files in public/
   rsync public/ my_hosting_server:/my/hosting/path
   # or
   aws s3 sync --acl public-read public/ s3://name.of.my.s3.bucket
```
   
For an example, here is [My old Livejournal](http://rant.unixbigot.id.au.s3-website-ap-southeast-2.amazonaws.com) now preserved under my own control at Amazon.
