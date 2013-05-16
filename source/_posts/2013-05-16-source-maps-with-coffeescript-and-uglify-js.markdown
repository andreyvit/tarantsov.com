---
layout: post
title: "7-minute Guide to Source Maps With CoffeeScript and Uglify.js"
date: 2013-05-16 17:33
comments: true
draft: true
categories: WorkflowThu
permalink: /WorkflowThu/source-maps-with-coffeescript-and-uglify-js/
---

Hello, and welcome back to the rebooted #WorkflowThu series.

Today we're going to see how ridiculously easy it is to debug CoffeeScript files with source maps. Watch the video (07:14) or read the article below.

<video id="WorkflowThu02-CoffeeScriptSourceMaps" class="sublime" poster="http://assets.livereload.com/WorkflowThu02-CoffeeScriptSourceMaps-poster.gif" width="800" height="500" title="#WorkflowThu 02 CoffeeScript Source Maps" data-uid="WorkflowThu02-CoffeeScriptSourceMaps" preload="none">
  <source src="http://assets.livereload.com/WorkflowThu02-CoffeeScriptSourceMaps.mp4" />
</video>

Tools mentioned in this video:

* [Uglify.js v2](https://github.com/mishoo/UglifyJS2) by Mihai Bazon (free)
* [LiveReload](http://livereload.com/) by yours truly ($10) — alternatives include [CodeKit](http://incident57.com/codekit/index.php) by Brian Jones ($25), [Fire.app](http://fireapp.handlino.com) by Handlino ($14), and [others](http://go.livereload.com/alternatives)
* [coffee-rails-source-maps gem](https://github.com/markbates/coffee-rails-source-maps) be Mark Bates for source maps support in Rails
* [chardin.js](http://heelhook.github.io/chardin.js/) by Pablo Fernandez — download the exact version used in this video [here](http://assets.livereload.com/WorkflowThu02-chardinjs.zip)

[Subscribe to #WorkflowThu](/workflow-thursday-subscribe/) if you don't want to miss future screencasts.


## Transcript

Contents:

* Intro
* Producing Source Maps with CoffeeScript
* Important Notes
* Producing Source Maps with Uglify.js
* Workflow Discussion
* Preview of WorkflowThu 03


### Intro

Let's start with a copy of chardin.js, a neat library for adding simple overlay instructions to your apps:

![Screenshot of chardin.js overlay in action](http://assets.livereload.com/WorkflowThu02-images/ChardinOverlay.png)

It's written in CoffeeScript. If we open the developer tools, we will only see the compiled JavaScript, of course. While it's readable, it's still a mess: the line numbers don't match, some constructs are less than obvious, and so on.

![](http://assets.livereload.com/WorkflowThu02-images/ChromeDevTools-raw-js.png)

Thankfully, now we have a tool to deal with that, at least in Google Chrome. Let me show you how.


### Producing Source Maps with CoffeeScript

We need to change CoffeeScript compiler options to generate a source map. I'll be using LiveReload here because, you know, I've made it, and also because it's a very simple app. But pretty much every similar tool will support source maps soon.

Let's go to Compiler Options in LiveReload and enable “Generate source map”. (If you were using CoffeeScript on the command line, you would simply pass `--map`.)

![](http://assets.livereload.com/WorkflowThu02-images/LiveReload-CoffeeScript-options.png)

The next time we save a `.coffee` file, there is an additional `.map` file generated:

![](http://assets.livereload.com/WorkflowThu02-images/SublimeText-source-map.png)

Let's open Chrome Developer Tools again. We can now see the source CoffeeScript file, set breakpoints and step through the original code:

![](http://assets.livereload.com/WorkflowThu02-images/ChromeDevTools-debugging-CoffeeScript.png)

That's it.

If you're running a web app server like Rails, Django or Node.js, there is an additional step for you because you need to make sure that the map files and CoffeeScript source files are exposed to the web browser. Rails apps can use [coffee-rails-source-maps gem](https://github.com/markbates/coffee-rails-source-maps),  otherwise Google for your framework name plus “source maps”.

For most PHP apps and static designs files, however, the steps described above should be enough.


### Important Notes

First: Source maps must be enabled in Chrome Dev Tools options. Click the gears button in the bottom-right corner and then enable the checkbox:

![](http://assets.livereload.com/WorkflowThu02-images/ChromeDevTools-enable-source-maps.png)

Second: You need a recent version of LiveReload; I have v2.3.26 here, which is available both on the Mac App Store and [on the support web site](http://go.livereload.com/trial).


### Producing Source Maps with Uglify.js

Let's now minify the library, but still keep the mapping to the original source lines.

We'll be using Uglify version 2 for that, which I already have installed (it's as simple as `sudo npm install -g uglify-js` as long as you have Node.js).

I'm going to set up LiveReload to run minification on every change. I don't normally recommend that, but it's appropriate for the demo. The syntax is `uglifyjs source.js -o minified.js`:

![](http://assets.livereload.com/WorkflowThu02-images/LiveReload-uglify-1.png)

If we re-save the coffee file, we now get `chardinjs.min.js`. Let's also fix index.html to include the minified file, and then reload Chrome:

![](http://assets.livereload.com/WorkflowThu02-images/ChromeDevTools-minified-js.png)

Not very promising for debugging, is it? We can fix it, though, in two easy steps.

Let's ask Uglify.js to generate a source map for us by adding `--source-map chardinjs.min.js.map`:

![](http://assets.livereload.com/WorkflowThu02-images/LiveReload-uglify-2.png)

Re-save the coffee file, reload Chrome and things are looking better now:

![](http://assets.livereload.com/WorkflowThu02-images/ChromeDevTools-minified-js-with-map.png)

We see the unminified version, can set breakpoints and single-step. But we wanted to see CoffeeScript here, not JavaScript, right?

That's the next step. Instead of giving Uglify.js an input JavaScript file, let's give it an input source map (produced by LiveReload) instead: `--in-source-map chardinjs.map`:

![](http://assets.livereload.com/WorkflowThu02-images/LiveReload-uglify-3.png)

And here we get the holy grail of web development — a minified JavaScript file that can be debugged as if it was the original CoffeeScript:

![](http://assets.livereload.com/WorkflowThu02-images/ChromeDevTools-minified-js-with-proper-map.png)


### Workflow Discussion

So that's source maps in CoffeeScript. Yes, you should be using them; it's ridiculously easy, and you basically have no excuse.

If you've been waiting to try CoffeeScript, this is probably a good time to see if it clicks with you and your team. There are no semantic differences between JavaScript and CoffeeScript, so the choice is purely a matter of personal taste and aesthetic.

Now, the minification workflow we saw is a nice way to debug an occasional production issue. I still don't think it's a good idea to use minified files during development, at least not until all your target browsers support source maps. Right now only Google Chrome and WebKit Nightlies have the support; hopefully Firefox will join the bunch soon.


## Preview of WorkflowThu 03

As I've mentioned, this series is resumed and will continue for a while. In the next installment, I'm going to show a ming-bogging images workflow involving Photoshop, Enigma64 and LiveReload.

Note: I promise there will be episodes that don't feature LiveReload. :-) I will obviously mention my products often, but the series is focused on using hot / new / frequently-asked-about technologies in everyday work, and not on any specific products.

![](http://assets.livereload.com/WorkflowThu02-images/WorkflowThu03-sneak-peek.png)


P.S. If you're excited, [you should subscribe here!](/workflow-thursday-subscribe/)

P.P.S. Sharing this article with your friends would be very nice, too.
