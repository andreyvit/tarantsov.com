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
# utilities

META_RE = re.compile(r'^[ ]{0,3}(?P<key>[A-Za-z0-9_-]+):\s*(?P<value>.*)')
TITLE_RE = re.compile(r'#\s*(.*)')
H1_RE = re.compile(r'<h1[^>]*>(?P<data>.+?)</h1>')
PAGE_LINK_RE = re.compile(r'\[\[(.+?)\|(.+?)]]')
ALT_LINK_RE = re.compile(r'\[(.+?)]\(\((.+?)\)\)')

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
  
def find_pages(path, page):
  folder = os.path.join(pages_path, path)
  file = os.path.join(folder, page)
  result = [os.path.join(path, page)] if os.path.isdir(file) or os.path.isfile(file) else []
  if len(path) > 0:
    pos = path.rfind('/')
    parent = path[0:pos] if pos >= 0 else ''
    return find_pages(parent, page) + result
  else:
    return result

def find_page(path, page):
  data = find_pages(path, page)
  return data[-1] if len(data) > 0 else None

def read_file(file):    
  f = open(file)
  try:
    content = unicode(f.read(), 'utf-8')
  finally:
    f.close()
  return content
    
def determine_title(html, meta):
  if 'title' in meta:
    return meta['title']
  t = [None]
  def title_recorder(m):
    t[0] = m.group(1)
    return ''
  html = re.sub(H1_RE, title_recorder, html)
  return t[0], html
    
def place_links_to_pages(path, html):
  def page_linker(m):
    caption = m.group(1)
    page = m.group(2)
    url = find_page(path, page)
    if url:
      return "<a href=\"/%s/\">%s</a>" % (url, caption)
    else:
      return "%s <span style=\"color: red;\">(%s)</span>" % (caption, page)
  html = re.sub(PAGE_LINK_RE, page_linker, html)
  html = re.sub(ALT_LINK_RE, page_linker, html)
  return html
  
def determine_path_components(path):
  components = [component('home', '/')]
  # components = []
  if len(path) > 0:
    # components = [component('home', '/')]
    cur_components = []
    for c in path.split('/'):
      cur_components = cur_components + [c]
      components.append(component(c, "/%s/" % '/'.join(cur_components)))
  return components
  
def relink_images(html):
  return re.sub('img src="', lambda m: '%s/static/images/' % m.group(), html)
  
def textualize(path):
  file = os.path.join(pages_path, path)
  if os.path.isdir(file):
    file = os.path.join(file, 'index')
  if not os.path.isfile(file):
    return None, {}

  content = read_file(file)
  lines = content.split("\n")
  meta, lines = parse_meta(lines)
  return "\n".join(lines), meta
  
def fix_typography(html):
  html = re.sub(u'&laquo;', '«', html)
  html = re.sub(u'<p>«', '<p><font class="hlaquo">&laquo;</font>', html)
  html = re.sub(u' «', '<font class="slaquo"> </font><font class="hlaquo">&laquo;</font>', html)
  html = re.sub(u' \\(', '<font class="sbrace"> </font><font class="hbrace">(</font>', html)
  html = re.sub(u'~', '&nbsp;', html)
  return html
  
COL_WIDTH = 40
COL_MARGIN = 10
  
def calculate_width_and_margin(colspan):
  return COL_WIDTH * colspan - COL_MARGIN, COL_MARGIN

def split_width_and_margin(num_cols, divider_size = 1, total_cols = 24):
  colspan = (total_cols - divider_size * (num_cols - 1)) / num_cols
  col_width         = COL_WIDTH * colspan - COL_MARGIN
  margin_right      = COL_MARGIN + divider_size * COL_WIDTH
  last_margin_right = (total_cols - colspan * num_cols - divider_size * (num_cols - 1)) * COL_WIDTH
  return col_width, margin_right, last_margin_right
  
def format_sidebyside_cols(text, col_count):
  cols = re.split(r'(?m)-+\[col\]-+$', text)
  col_width, margin_right, last_margin_right = split_width_and_margin(len(cols), 1)
  
  html = ''
  for col in cols[0:-1]:
    html += """<div style="float: left; width: %dpx; margin-right: %dpx;">\n\n%s\n\n</div>\n\n""" % (col_width, margin_right, col.strip())
  html += """<div style="float: left; width: %dpx; margin-right: %dpx;">\n\n%s\n\n</div>\n\n""" % (col_width, last_margin_right, cols[-1].strip())
  return html
  
def format_sidebyside(m):
  text = m.group(1)
  rows = re.split(r'(?m)^=+\[row\]=+$', text.strip())
  col_count = len(rows[0])
  html = "\n".join([format_sidebyside_cols(row, col_count) for row in rows])
  html += """\n<div class="clear"></div>\n\n"""
  return html
  
def htmlize_file(path):
  file = os.path.join(pages_path, path)
  if os.path.isdir(file):
    file = os.path.join(file, 'index')
  if not os.path.isfile(file):
    return None, {}

  content = read_file(file)
  lines = content.split("\n")
  meta, lines = parse_meta(lines)
  content = "\n".join(lines)
  
  content = place_links_to_pages(path, content)
  content = re.sub(r'(?s)(?m)^=+\[sidebyside\]=+(.*?)=+\[/sidebyside\]=+$', format_sidebyside, content)
  html = markdown.markdown(content)
  html = relink_images(html)
  html = fix_typography(html)
  html = re.sub(r'<clear>', '<div class="clear"></div>', html)
  return html, meta

def find_and_htmlize(context_path, page):
  path = find_page(context_path, page)
  if path:
    return htmlize_file(path)
  else:
    return None, {}
    
def read_options(context_path, page):
  meta = {}
  for path in find_pages(context_path, page):
    text, local_meta = textualize(path)
    meta.update(**local_meta)
  return meta

########################################################################################################
# handlers
  
class IndexHandler(BaseHandler):
  @prolog()
  def get(self, path):
    if len(path) > 0:
      if path[-1] == '/':
        path = path[0:-1]
      else:
        self.redirect_and_finish('/%s/' % path)
        
    html, meta = htmlize_file(path)
    if not html:
      self.data.update(path = path)
      self.render_and_finish('page-not-found.html')
      
    title, html = determine_title(html, meta)
    components = determine_path_components(path)
    options = read_options(path, 'index')
    site_title = options['site-title'] if 'site-title' in options else 'Site-Title missing in .options'
    copyright_year = options['copyright-year']
    copyright_email = options['copyright-email']
    copyright_name = options['copyright-name']
    
    self.data.update(title = title, content = html, components = components, site_title = site_title,
      copyright_year=copyright_year, copyright_email=copyright_email, copyright_name=copyright_name)
    self.render_and_finish('page.html')

url_mapping = [
  ('/([a-zA-Z0-9/-]*)', IndexHandler),
]

def main():
  application = webapp.WSGIApplication(url_mapping, debug=True)
  wsgiref.handlers.CGIHandler().run(application)

if __name__ == '__main__':
  main()
