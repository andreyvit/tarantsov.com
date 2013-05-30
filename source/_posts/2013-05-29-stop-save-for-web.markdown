---
layout: post
title: "STOP SAVE FOR WEB"
date: 2013-05-30 11:13
comments: true
categories: WorkflowThu
permalink: /WorkflowThu/stop-save-for-web/
---

Today we black out Photoshop to protest the time people waste doing Save for Web.

<video id="WorkflowThu03-StopSAVEFORWEB" class="sublime" poster="http://assets.livereload.com/WorkflowThu03-StopSAVEFORWEB-poster.gif" width="800" height="450" title="#WorkflowThu 03 Stop Save For Web" data-uid="WorkflowThu03-StopSAVEFORWEB" preload="none">
  <source src="http://assets.livereload.com/WorkflowThu03-StopSAVEFORWEB.mp4" />
</video>

Here's how *you* can help:

1. Support the anti-SAVEFORWEB tools:

    Mac:

    * [Enigma64](http://getenigma64.com) by Piffle ($20)
    * [Keyboard Maestro](http://www.keyboardmaestro.com/) by Stairways Software ($40)
    * [LiveReload for Mac](http://livereload.com/) by yours truly ($10) — [alternatives](http://go.livereload.com/alternatives)
    * [Hazel](http://www.noodlesoft.com/hazel.php) by Noodlesoft ($25)
    * [ImageOptim](http://imageoptim.com) by porneL (free)

    Windows:

    * [Enigma64](http://getenigma64.com) by Piffle ($20)
    * [Sikuli](http://www.sikuli.org) by MIT and CU-Boulder (free) — an automation tool
    * [LiveReload for Windows](http://livereload.com/) by yours truly (free while in alpha)
    * [Belvedere](http://lifehacker.com/341950/belvedere-automates-your-self+cleaning-pc) by Adam Pash (free) — a port of Hazel

2. Educate others! [Click to tweet with #StopSAVEFORWEB](http://clicktotweet.com/ygq9j).

3. [Subscribe](/workflow-thursday-subscribe/) to be prepared when Save for Web returns.

Which of the following topics do you want covered in Ep 05? Let me know in the comments!

* Sharing code between projects (a spotlight on git-submodule, git-subdir and Git client apps)
* Static site generators
* Getting started with SASS and Compass


## Transcript

Contents:

* Imagine a world without free kitten images
* The shenanigans of old hardware
* Today we fight back
* Steamy, yes?
* Bonus: Hazel


## 1. Imagine a world without free kitten images

If you're like me, you spend a lot of time working on your images in Photoshop.

Let's say a new questionable law has appeared and you want to send a message to your visitors:

<img src="http://assets.livereload.com/WorkflowThu03-images/Photoshop-kitten-censored.png" width="800" height="189">

So you make the change, and then go:

...Save for Web...

...Save...

...Save...

...Replace...

...and then maybe you have to refresh your browser on top of that.

And the way it _should_ work is: make a change, push a button, boom, make another change, push a button, boom.

(You really need to see the video if you want to appreciate how nice that can be.)

This is possible with a piece of old German equipment called Enigma. (It's also an app that adds one-click image exporting options to Photoshop.)

Here it is:

<img src="http://assets.livereload.com/WorkflowThu03-images/Enigma-front.png" width="796" height="319">

Let me walk you through my settings:

<img src="http://assets.livereload.com/WorkflowThu03-images/Enigma-settings.png" width="800" height="359">

1. “Prompt File Name” is off, because otherwise Enigma would be asking me for a file name every time.

2. PNG and JPEG settings can be adjusted to your liking.

3. On the main screen, I choose the image format (JPEG in this case) and the output folder.

4. There's also an option to trim the transparent pixels around your image and to optimize it. I don't care about the optimization one because I see this as a development process. Before a release, you would run your images through something like ImageOptim anyway.

When you click the “Visible Canvas” button, the visible layers are saved into a file in the output folder with the same name as the PSD file.

Alternatively, if you have separate layers for separate images, you can click the “Selected Layers” button to export a file with the same name as the _layer_.

Now, this is WorkflowThu, so we cannot simply leave it at that, can we?


## 2. The shenanigans of old hardware

There are actually several problems with this workflow that require additional tools to solve.

The first problem is that Enigma does not support hotkeys. When I'm working on something, I really don't want to move my mouse to click the button, because I might be adjusting a setting somewhere and I want to continue doing that while previewing the results.

The second one is that Enigma never overwrites files. So if you configure it to output into your actual images folder, it won't overwrite them, it will simply create additional files with numeric suffixes.

And the third issue is pushing the changes into your browser window, but we all know how this is going to be solved (hint hint).


## 3. Today we fight back

Let's handle the hot keys first. We'll use a tool called Keyboard Maestro.

(This part is going to be Mac-only. On Windows, there's [Sikuli](http://www.sikuli.org) that can do all these tricks, but it's not nearly as user-friendly, and there are some crazily expensive tools for automated testing, but I haven't been able to find a good solution for a regular user. Let me know in the comments if you know any.)

Create a new macro called e.g. “Enigma Export”, assign a hot key like F6, and add a new “Move or Click Mouse” action from the “Interface Control” section:

<img src="http://assets.livereload.com/WorkflowThu03-images/KeyboardMaestro-macro-step1.png" width="800" height="356">

Here's the magic part: set “relative to the found image's center”, then make a screenshot of the “Visible Canvas” button of Enigma, and drag the screenshot into the macro action's image well:

<img src="http://assets.livereload.com/WorkflowThu03-images/KeyboardMaestro-macro-step1b.png" width="800" height="356">

(Btw I have the Desktop folder as a stack on my Dock, really helps to get to the screenshots faster.)

Enable “Restore Mouse Location”, and we're almost done — the action now works!

There is a problem, though: after you click the Enigma button (or run our macro), the keyboard focus moves into the Enigma window, and many Photoshop shortcuts (the ones without modifier keys: 1..9, Shift-1..Shift-9, tool selection keys) stop working until you click one of the native Photoshop windows.

We can extend the macro to perform that second click for us too. A safe place to click on is, say, somewhere on the toolbar of the layers view.

<img src="http://assets.livereload.com/WorkflowThu03-images/KeyboardMaestro-macro-step2.png" width="800" height="357">

Done! Pressing F6 now runs the export, and the Photoshop hotkeys still work after that. Nice.

You can enable Display checkboxes in the macro editor to see the winning match (green) and other potential matches (orange):

<img src="http://assets.livereload.com/WorkflowThu03-images/KeyboardMaestro-displaying-matches.png" width="800" height="268">

You may want to drag the fuzziness sliders all the way to the left so that there is no risk of unwanted matches:

<img src="http://assets.livereload.com/WorkflowThu03-images/KeyboardMaestro-macro-fuzziness.png" width="800" height="269">

Are you already thinking about all sorts of ways you can automate stuff using this app?


## 4. Steamy, yes?

So how about the file names issue?

Well, Enigma is set up to put the files into a separate ‘exported’ folder, and we'll use another tool to move them into the images folder, overwriting the existing files. In the example, we need to move images from ‘exported’ into ‘img/examples’:

<img src="http://assets.livereload.com/WorkflowThu03-images/folder-layout.png" width="800" height="145">

I'm going to show two different ways to move the files: the first one uses LiveReload, and the second one involves an app called Hazel.

We'll now set up LiveReload to push image changes to the browser, and also to move the files:

First, drag the entire folder into LiveReload. (This is actually all the configuration we need to do to push the changes into the browser.)

Second, let's add a post-processing command to move all files from the “exported” folder into the images folder: `mv exported/* img/examples/`:

<img src="http://assets.livereload.com/WorkflowThu03-images/LiveReload-mv.png" width="800" height="450">

Third, open our site in the web browser and enable the extension:

<img src="http://assets.livereload.com/WorkflowThu03-images/LiveReload-extension-button.png" width="800" height="167">

Note that because I'm working with local files, I need to have “Allow access to file URLs” enabled on the Chrome Extensions page; Safari does not provide a similar option at all.

And everything should work now!

Switch to the Photoshop, change something, hit F6, boom.

Isn't it just awesome? I think it is.


## Bonus: Hazel

If you do not want to use LiveReload, or find that figuring out a right post-processing command is too complicated for your needs, there's an app called Hazel.

Hazel is actually my first choice for any automated file manipulations. It's a Mac app that monitors a given set of folders and does stuff based on a set of rules. There's also a similar free Windows app called Belvedere. Hazel has a _huge_ following, and is literally the difference between a clutterred and a clean Mac. I use it to clean up my Downloads, Desktop and Trash folders and to highlight the apps I haven't used for a while.

Let's add the “exported” folder to Hazel and create a new rule to “move all images to destinations”. Set it to match all files, and move them into the images folder with overwriting enabled:

<img src="http://assets.livereload.com/WorkflowThu03-images/Hazel-rule.png" width="800" height="450">

That's it. As soon as a file appears in the “exported” folder, Hazel will notice it and move into “img/examples”.

There is a couple of seconds of delay, though, because Hazel does not process the changes immediately. This is the downside of the approach; on the other hand, you can easily set up very complicated rules in Hazel which could get tricky to do in LiveReload.


## Moar

The next time we'll dig into what makes Sublime Text great.

<img src="http://assets.livereload.com/WorkflowThu03-images/WorkflowThu04-preview.png" width="800" height="450">

Other episodes of WorkflowThu:

* Ep. 01: [Sublime Text Workflow That Beats Coda and Espresso](/blog/2012/02/sublime-text-workflow-that-beats-coda-and-espresso/)
* Ep. 02: [7-minute Guide to Source Maps With CoffeeScript and Uglify.js](/WorkflowThu/source-maps-with-coffeescript-and-uglify-js/)

[Click here to subscribe!](/workflow-thursday-subscribe/)

Which of the following topics do you want covered in Ep 05? Let me know in the comments!

* Sharing code between projects (a spotlight on git-submodule, git-subdir and Git client apps)
* Static site generators
* Getting started with SASS and Compass

Acknowledgements:

* “Steamy, yes?” is a reference to [The Oatmeal](http://theoatmeal.com/sopa) (mildly NSFW).

* Kitten image by <a href="http://flickr.com/photos/bigtallguy/">BigTallGuy</a>, courtesy of <a href="http://placekitten.com">{placekitten}</a> — a leading source of free kitten images that power the web.
