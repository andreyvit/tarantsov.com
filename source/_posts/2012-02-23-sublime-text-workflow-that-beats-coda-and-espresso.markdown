---
layout: post
title: "Sublime Text Workflow That Beats Coda and Espresso"
date: 2012-02-23 19:45
comments: true
categories: WorkflowThu
permalink: /blog/2012/02/sublime-text-workflow-that-beats-coda-and-espresso
translations:
  - language: Spanish
    url: http://www.webhostinghub.com/support/es/misc/un-workflow-de-texto
---

Welcome to the new #WorkflowThu series which helps web developers try new great things.

In this episode, we're talking about jumping into Sublime Text 2 and setting up a workflow that beats traditional tools like Coda and Espresso. Watch the video (10:30) or read the article below.

<video class="sublime" width="800" height="500" poster="http://assets.livereload.com/WorkflowThu01-SublimeTextWorkflow.png" preload="none">
  <source src="http://assets.livereload.com/WorkflowThu01-SublimeTextWorkflow.mov" />
  <source src="http://assets.livereload.com/WorkflowThu01-SublimeTextWorkflow.mp4" />
</video>

App and plugins mentioned in this video:

* [Sublime Text 2](http://www.sublimetext.com/) ($59 after an unlimited free trial)
* [Package Control](http://wbond.net/sublime_packages/package_control) (free)
* [SideBarEnhancements](https://github.com/titoBouzout/SideBarEnhancements) (free)
* [Divvy](http://mizage.com/divvy/) ($14) — other alternatives are: [Moom](http://manytricks.com/moom/) ($5), [Optimal Layout](http://most-advantageous.com/optimal-layout/) ($14), [Window Tidy](http://www.lightpillar.com/macos/windowtidy/index.html) ($8)
* [LiveReload](http://livereload.com/) ($10)
* [Sublime SFTP](http://wbond.net/sublime_packages/sftp) ($16)

## Transcript

Contents:

* Installing Plugins
* Managing Windows
* Live Preview
* FTP/SFTP
* Switching Projects
* Preferences and Key Bindings


### Part 1. Installing Plugins

We're starting with a freshly installed copy of Sublime Text 2:

![A pointless screenshot](http://assets.livereload.com/WorkflowThu01/SublimeDefault.png)

Unfortunately, it does not ship packages for languages like LESS and CoffeeScript out of the box, so one of the first things you need to learn is how to install packages.

For that, you need a package manager called [Package Control](http://wbond.net/sublime_packages/package_control). Copy the strange-looking piece of code from [the installation instructions](http://wbond.net/sublime_packages/package_control/installation), switch to Sublime Text, open the Python console (View » Show Console), paste the code into the input field and press Enter:

![](http://assets.livereload.com/WorkflowThu01/PackageControlSnippetInConsole.png)

After you restart Sublime Text, you are ready to use this package manager to install the plugins that you want. For that, you need to learn a central concept of Sublime called a “Command Palette” (Tools » Command Palette, or Command-Shift-P). It lists all the commands that Sublime can do:

![](http://assets.livereload.com/WorkflowThu01/SublimeCommandPalette.png)

In particular, you can see that Package Control has added a bunch of commands of its own. What you're interested in is the Install command, so you just type Install to search for that command and press Enter. You get a list of packages; there are lots of them available for Sublime, which is one of the reasons it is awesome:

![](http://assets.livereload.com/WorkflowThu01/PackageControlListOfPackages.png)

Choose the LESS package, and in a few moments it is installed and ready to use. Syntax highlighting works for .less files now:

![](http://assets.livereload.com/WorkflowThu01/LessSyntaxHighlighting.png)

Do the same with CoffeeScript (you really want to learn that Command-Shift-P shortcut). Another plugin I absolutely recommend to install is [SideBarEnhancements](https://github.com/titoBouzout/SideBarEnhancements); it adds many useful commands to the context menu of the sidebar and changes how New File command works.


### Part 2. Managing Windows

This will be a key part of our workflow because we really want to run the editor and the browser side by side. There are many apps that can help us with that, see the beginning of the post. My favorite one is [Divvy](http://mizage.com/divvy/); the way I normally use it is by setting up a set of keyboard shortcuts to quickly make the current window full screen or position it on the left or right half of the screen:

![](http://assets.livereload.com/WorkflowThu01/DivvyShortcuts.png)

So let's put Sublime and the browser side by side. Additionaly, to save some space, open the files you want to edit in tabs (by double-clicking them in the sidebar) and then hide the sidebar (View » Sidebar » Hide Sidebar or Command-K Command-B).

The second trick you need to learn is to split the window into groups, which allows you to see several files at the same time. There are multiple split layouts to choose from (View » Layout submenu). In our case, we want to see the PHP file at the top and the LESS file at the bottom:

![](http://assets.livereload.com/WorkflowThu01/SplitGroups.png)

You can drag'n'drop tabs between groups, although I recommend you to learn the shortcuts under View » Focus Group and View » Move File To Group submenus.


### Part 3. Live Preview

To really take advantage of this layout, you want the browser to be refreshed automatically when you change a file. For that, we'll be using an app called [LiveReload](http://livereload.com/) (made by yours truly). Run it and drag your project folder onto LiveReload's menu bar icon. Additionally, enable compilation of LESS and CoffeeScript:

![](http://assets.livereload.com/WorkflowThu01/LiveReload.png)

Follow that link you see in the window and install browser extensions. Switch to your browser and enable LiveReload by clicking on the toolbar button:

![](http://assets.livereload.com/WorkflowThu01/LiveReloadEnable.png)

Now, when you make a change on the page, the browser is refreshed automatically. When you change a LESS file, it is automatically compiled and applied to the page without reloading it.


### Part 4. FTP/SFTP

Chances are you are using SFTP to publish your site or even to edit it directly on the server. You may be using Coda or Espresso or Transmit for that, but Sublime has a good solution too: the [SFTP plugin](http://wbond.net/sublime_packages/sftp). Install it via Package Control; it is not free, but, like Sublime Text, it has an unlimited free trial.

After installing the plugin, right-click your web site's root folder in the sidebar and choose SFTP/FTP » Map to Remote. A JSON file appears, which you can enter all the connection settings in. Don't forget to set remote_path, and be sure to enable upload_on_save if you need it:

![](http://assets.livereload.com/WorkflowThu01/SublimeSFTP.png)

Save the config, make a change to your php file and watch it being uploaded:

![](http://assets.livereload.com/WorkflowThu01/SublimeSFTPUpload.png)

There are sync/upload/download commands available in the sidebar context menu and in the Command Palette for those who only publish their changes occasionally.

If you work with a remote site using the upload_on_save option, you need to configure LiveReload to handle that — otherwise, it will try to reload stuff while the upload is still running. Click the monitoring Options button, set up a small delay for full page reloads and enable CSS overriding:

![](http://assets.livereload.com/WorkflowThu01/LiveReloadRemote.png)

Now when you change a PHP file, it is uploaded to your server and then the browser is refreshed. When you change a LESS file, those changes are applied immediately.

You still need to upload the compiled CSS file later; for that, SFTP plugin has a monitoring option: open the CSS file (be sure to double-click it — single click does not create a tab), invoke Start Monitoring command from the Command Palette, and don't close the tab. Or you can simply use the Sync command when you are satisfied with your stylesheet changes, relying on LiveReload's override behavior to preview them.


### Part 5. Switching Projects

This is a small tip: you really want to save your projects using Project » Save Project As, because you can later switch between recently saved projects using Project » Switch Project In Window command (use the Command-Ctrl-P shortcut).

It's a real killer when working on several projects: hit Command-Ctrl-P, type a few letters of the project name, Enter, and you have that project open, with all of the open files preserved from the last time:

![](http://assets.livereload.com/WorkflowThu01/SublimeSwitchProjects.png)


### Part 6. Preferences and Key Bindings

Customizing Sublime Text may seem a bit nerdy, so you need to get used to it. Let's say you want to make the font size a little bit larger. When you choose Preferences, you get this empty JSON file:

![](http://assets.livereload.com/WorkflowThu01/PreferencesEmpty.png)

Turns out, there is another file with exactly the same format which contains the default values of all available settings. Open it using Preferences » Settings - Default. The settings are pretty well commented. Find the one you want, then copy it into your personal settings file.

![](http://assets.livereload.com/WorkflowThu01/PreferencesDefault.png)

Your new preferences are applied the moment you save that file:

![](http://assets.livereload.com/WorkflowThu01/PreferencesUser.png)

One setting that I recommend you to enable is auto_complete_commit_on_tab, explained in the comments:

![](http://assets.livereload.com/WorkflowThu01/PreferencesAutoCompleteCommitOnTab.png)

Sublime Text's completion feature is really awesome: for example, you can type “accot” and hit TAB to complete auto_complete_commit_on_tab — and before you finish typing half of that, Sublime usually has a proposal right there on your screen:

![](http://assets.livereload.com/WorkflowThu01/Completion.png)

You can customize the key bindings in a similar way. For example, I show and hide the sidebar very often so I prefer Control-S to the default Command-K Command-B shortcut. To set that up, you need to search for side bar in the default key bindings file (Preferences » Key Bindings - Default):

![](http://assets.livereload.com/WorkflowThu01/KeyBindingsDefault.png)

then paste that line into your personal key bindings file (Preferences » Key Bindings - User) and change the shortcut to Ctrl-S:

![](http://assets.livereload.com/WorkflowThu01/KeyBindingsUser.png)

As soon as you save the file, Ctrl-S starts working. (The default shortcut works too, so this is adding a shortcut rather than changing it.)


## The End

You can Google many more plugins and tricks for Sublime Text. Thanks for watching (err, reading), and see you next week!

Be sure to subscribe to this blog if you don't want to miss future screencasts. Or <a href="/workflow-thursday-subscribe/">subscribe via email</a>.
