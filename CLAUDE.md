# CLAUDE.md

This is Andrey's personal web site and blog, with a custom single-file static site generator (`site.go`). I like to keep things simple and pragmatic, and my generator reflects that.

Content goes into `content` (served at site root), images, CSS and other stuff under `assets` (served under `/assets/`), stuff like layouts and partials goes into `theme/` subdirectories.

Content first uses a template from `theme/templates/`, and then an outer layout from `theme/layouts/`.

I use custom pseudo-tags starting with `x-` for shortcodes, e.g. `<x-textnav>` and `<x-cta>`.

The generator allows me the freedom to experiment with non-traditional navigation methods and layouts. My top-level navigation is a single paragraph of text built up from `data/nav/main.json`, with each sentence (or part of a sentence) linking to a separate article or section.

Site generator is invoked via `go run . -dev`, which is gonna serve the site on localhost (in development mode, meaning site is rebuilt on each request).

- The development server automatically reloads content from disk when `-dev` flag is used
- All paths are served from memory in production mode for performance
- Content files named `index.md` are served at their directory path
- Other content files are served with a trailing slash (e.g., `about.md` → `/about/`)
