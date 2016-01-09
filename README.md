# Here's what happened.

The code was being compiled during the staging process just fine, as we suspected
because we didn't see any errors. The name of the executable was what was messed
up, we were on the right track but didn't look at the right part of the build pack.

The build pack uses the ImportPath in Godeps.json to name the executable

https://github.com/cloudfoundry/go-buildpack/blob/master/bin/compile#L121 is where
it picks up the name

https://github.com/cloudfoundry/go-buildpack/blob/master/bin/compile#L176 is where
it makes a directory with the name, and then by convention go names the executable
(because it's got a main) to the name of the directory.

So, I had to change the imports.json. To reference the path of the project.
https://github.com/krujos/goservice/commit/995619bea518b40a817b2d8e4f76139b1b68bf3c (you might have to click "show diff")

I'm not sure how you ended up in this situation I suspect it's because an unusual  pattern / workflow by creating code without repo or something that caused godeps to use "." as the import path.

# Other stuff I tried
I tried a couple of things and failed as an expierement. At first I thought it
might be that you had app.go instead of main.go. But that made no difference
https://github.com/krujos/goservice/commit/a217dffeb82704d07560349ca5023d25ee974331
