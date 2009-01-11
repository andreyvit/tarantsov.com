#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import logging
import wsgiref.handlers
import urllib
import re
import markdown
from datetime import datetime, timedelta

from google.appengine.ext import webapp
from google.appengine.ext.webapp import template
from google.appengine.api import users

template_path = os.path.join(os.path.dirname(__file__), 'templates')
pages_path = os.path.join(os.path.dirname(__file__), 'pages')

class FinishRequest(Exception):
  pass

class prolog(object):
  def __init__(decor, path_components = [], fetch = [], config_needed = True):
    decor.config_needed = config_needed
    decor.path_components = path_components
    decor.fetch = fetch
    pass

  def __call__(decor, original_func):
    def decoration(self, *args):
      try:
        return original_func(self, *args)
      except FinishRequest:
        pass
    decoration.__name__ = original_func.__name__
    decoration.__dict__ = original_func.__dict__
    decoration.__doc__  = original_func.__doc__
    return decoration

class BaseHandler(webapp.RequestHandler):
  def __init__(self):
    self.now = datetime.now()
    self.data = dict(now = self.now)
    
  def redirect_and_finish(self, url, flash = None):
    self.redirect(url)
    raise FinishRequest
    
  def render_and_finish(self, *path_components):
    self.response.out.write(template.render(os.path.join(template_path, *path_components), self.data))
    raise FinishRequest
    
  def access_denied(self, message = None, attemp_login = True):
    if attemp_login and self.user == None and self.request.method == 'GET':
      self.redirect_and_finish(users.create_login_url(self.request.uri))
    self.die(403, 'access_denied.html', message=message)

  def not_found(self, message = None):
    self.die(404, 'not_found.html', message=message)

  def invalid_request(self, message = None):
    self.die(400, 'invalid_request.html', message=message)
    
  def die(self, code, template, message = None):
    if message:
      logging.warning("%d: %s" % (code, message))
    self.error(code)
    self.data.update(message = message)
    self.render_and_finish('errors', template)


########################################################################################################
# handlers

META_RE = re.compile(r'^[ ]{0,3}(?P<key>[A-Za-z0-9_-]+):\s*(?P<value>.*)')
TITLE_RE = re.compile(r'#\s*(.*)')
H1_RE = re.compile(r'<h1[^>]*>(?P<data>.+?)</h1>')
PAGE_LINK_RE = re.compile(r'\[\[(.+?)\|(.+?)]]')

class component(object):
  
  def __init__(self, name, href):
    self.name = name
    self.href = href

def parse_meta(lines):
  meta = {}
  new_lines = []
  for line in lines:
    m1 = META_RE.match(line)
    if m1:
      key = m1.group('key').lower().strip()
      meta[key] = m1.group('value').strip()
    else:
      new_lines.append(line)
  return meta, new_lines

def parse_title(lines):
  new_lines = []
  title = None
  for line in lines:
    m = TITLE_RE.match(line)
    if m:
      title = m.group(1).strip()
    else:
      new_lines.append(line)
  return title, new_lines
  
def find_page(path, page):
  folder = os.path.join(pages_path, path)
  file = os.path.join(folder, page)
  if os.path.isdir(file) or os.path.isfile(file):
    return os.path.join(path, page)
  if len(path) > 0:
    pos = path.rfind('/')
    parent = path[0:pos] if pos >= 0 else ''
    return find_page(parent, page)
  else:
    return None
    
class IndexHandler(BaseHandler):
  @prolog()
  def get(self, path):
    if len(path) > 0:
      if path[-1] == '/':
        path = path[0:-1]
      else:
        self.redirect_and_finish('/%s/' % path)
    
    file = os.path.join(pages_path, path)
    if os.path.isdir(file):
      file = os.path.join(file, 'index')
    if not os.path.isfile(file):
      self.data.update(path = path)
      self.render_and_finish('page-not-found.html')
    
    f = open(file)
    try:
      content = unicode(f.read(), 'utf-8')
    finally:
      f.close()
    
    lines = content.split("\n")
    meta, lines = parse_meta(lines)
    content = "\n".join(lines)
    
    html = markdown.markdown(content)
    title = None
    if 'title' in meta:
      title = meta['title']
    if not title:
      t = [title]
      def title_recorder(m):
        t[0] = m.group(1)
        return ''
      html = re.sub(H1_RE, title_recorder, html)
      title = t[0]
    else:
      title = 'Untitled'
        
    def page_linker(m):
      caption = m.group(1)
      page = m.group(2)
      url = find_page(path, page)
      if url:
        return "<a href=\"/%s/\">%s</a>" % (url, caption)
      else:
        return "%s <span style=\"color: red;\">(%s)</span>" % (caption, page)
      
    html = re.sub(PAGE_LINK_RE, page_linker, html)
    
    components = [component('home', '/')]
    # components = []
    if len(path) > 0:
      # components = [component('home', '/')]
      cur_components = []
      for c in path.split('/'):
        cur_components = cur_components + [c]
        components.append(component(c, "/%s/" % '/'.join(cur_components)))
        
    html = re.sub('img src="', lambda m: '%s/static/images/' % m.group(), html)
        
    self.data.update(title = title, content = html, components = components)
    self.render_and_finish('page.html')

url_mapping = [
  ('/([a-zA-Z0-9/-]*)', IndexHandler),
]

def main():
  application = webapp.WSGIApplication(url_mapping, debug=True)
  wsgiref.handlers.CGIHandler().run(application)

if __name__ == '__main__':
  main()
