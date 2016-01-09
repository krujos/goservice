# Here's what happened.

You can walk the commit history to see what I tried. Below is an annotated summary.

##What's the actual problem, our logs are telling us nothing?
The code was being compiled during the staging process just fine, as we suspected because we didn't see any errors. But we could not find an executable named `goservice` as we expected to (note, I changed the name from `go_service` in the zip you sent me to comply with golang convention). I couldn't think of anything more I could do to get CF to tell me what it was doing, we used all the diagnostic commands I could think of `cf events` and `cf logs`. So I need to see what it's working with. Lets download and untar the dropplet and stare at it a bit.

```
➜ goservice git:(master) ✗ cf curl /v2/apps/`cf app gocf --guid`/droplet/download > /tmp/my_droplet
➜ goservice git:(master) ✗ pushd /tmp
➜  /tmp  tar -xvf my_droplet
x ./: Can't update time for .
x ./staging_info.yml
x ./logs/
x ./app/
x ./app/Procfile
x ./app/widget.go
x ./app/main.go
x ./app/Godeps/
x ./app/Godeps/_workspace/
x ./app/Godeps/_workspace/pkg/
x ./app/Godeps/_workspace/pkg/darwin_amd64/
x ./app/Godeps/_workspace/pkg/darwin_amd64/github.com/
x ./app/Godeps/_workspace/pkg/darwin_amd64/github.com/julienschmidt/
x ./app/Godeps/_workspace/pkg/darwin_amd64/github.com/julienschmidt/httprouter.a
x ./app/Godeps/_workspace/src/
x ./app/Godeps/_workspace/src/github.com/
x ./app/Godeps/_workspace/src/github.com/julienschmidt/
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/tree.go
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/LICENSE
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/README.md
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/.travis.yml
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/router.go
x ./app/Godeps/_workspace/src/github.com/julienschmidt/httprouter/path.go
x ./app/Godeps/Godeps.json
x ./app/Godeps/Readme
x ./app/bin/
x ./app/bin/src
x ./app/.profile.d/
x ./app/.profile.d/concurrency.sh
x ./app/.profile.d/go.sh
x ./tmp/
tar: Error exit delayed from previous errors.

➜ /tmp # Scratch head for an inordinate amount of time, write off the tar error
➜ /tmp # as red herring based on gut, CF would have blown up differently for a
➜ /tmp # bad droplet. the app directory. "goservice" should be there, but I see
➜ /tmp # Spend a lot of time looking into the app directory, getting confused by
➜ /tmp # a "src" directory with nothing in it... wait, it's not a directory.

➜ /tmp  file app/bin/src
app/bin/src: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked (uses shared libs), not stripped
```

So I changed the `Procfile` to call an executable named `src` and pushed again. [Viola! The app started](https://github.com/krujos/goservice/commit/a217dffeb82704d07560349ca5023d25ee974331)! So suspicion confirmed, we've done something to cause an unexpected executable name.

## How did we get here?

I needed to figure out how the buildpack names a go app. Usually the name of the directory is used if `go build|install` does not receive a `-o` flag. I could not figure out how CF would know what to feed for `-o`, so I started looking at how the [go buildpack](https://github.com/cloudfoundry/go-buildpack) lays out the project directories. [I discovered the build pack uses the `ImportPath` in Godeps.json to name the executable](https://github.com/cloudfoundry/go-buildpack/blob/master/bin/compile#L121).

So the buildpack must do something with name right? [Line 176 & 177](https://github.com/cloudfoundry/go-buildpack/blob/master/bin/compile#L176) is where it makes a directory with the name, and then by convention go names the executable to the name of the directory. So standard go conventions so far.

So, I know something goofy must be happening in Godeps.json, lets look at or ImportPath. [We reference the current directory. None of my other go projects seem to do it that way, they all point to the parent repo. Lets change that to something more conventional.](https://github.com/krujos/goservice/commit/995619bea518b40a817b2d8e4f76139b1b68bf3c)(you might have to click "show diff")!

After that I pushed the app again and it worked as expected, cf started it, health checks passed and I could `curl` widgets to my hearts content.  

I'm not sure how you ended up in this situation I suspect it's because an unusual workflow by creating code without repo or something that caused Godep to use "." as the import path. Was it maybe that the project was initially in the root directory of $GOPATH? Godep got really confused and did the wrong thing, at least as far as CF is concerned, and I don't think it makes sense to be self referential given go's import semantics... curious for others thoughts here. I'm noodling on what I want to do to handle this case in the future. I think we should catch it in the buildpack and error out (as three hours and having to know a lot about go conventions is a little much for troubleshooting), but I'd like to figure out how we got here before I submit that to the team.

# Other stuff I tried
I tried a couple of things and failed as an experiment. At first I thought it might be that you had [app.go instead of main.go. But that made no difference](https://github.com/krujos/goservice/commit/a217dffeb82704d07560349ca5023d25ee974331)
