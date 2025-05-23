# jules/cli/crawler.py
import requests # Assuming requests is used for fetching pages
from bs4 import BeautifulSoup # For HTML parsing
from urllib.parse import urljoin, urlparse

class DiscoveredInput:
    def __init__(self, name, input_type, current_value, html_context_details=None):
        self.name = name
        self.type = input_type # e.g., text, hidden, submit, textarea, select
        self.current_value = current_value
        # html_context_details could be a dict: e.g. {'tag': 'input', 'parent_form_id': 'loginForm'}
        self.html_context_details = html_context_details if html_context_details else {}

    def __repr__(self):
        return f"Input(name='{self.name}', type='{self.type}', value='{self.current_value}')"

class DiscoveredForm:
    def __init__(self, action, method, form_id=None, form_class=None):
        self.action = action
        self.method = method.upper() if method else "GET" # Default to GET if not specified
        self.form_id = form_id
        self.form_class = form_class
        self.inputs = [] # List of DiscoveredInput objects

    def add_input(self, discovered_input):
        self.inputs.append(discovered_input)

    def __repr__(self):
        return f"Form(action='{self.action}', method='{self.method}', inputs={len(self.inputs)})"

class Crawler:
    def __init__(self, base_url):
        self.base_url = base_url
        self.parsed_base_url = urlparse(base_url)
        self.visited_urls = set()
        self.sitemap = []  # List of discovered URLs
        self.discovered_forms = [] # List of DiscoveredForm objects
        self.session = requests.Session()
        # TODO: Add User-Agent from config

    def _get_page_content(self, url):
        try:
            response = self.session.get(url, timeout=5) # Add timeout
            response.raise_for_status() # Raise HTTPError for bad responses (4xx or 5xx)
            if 'text/html' in response.headers.get('Content-Type', ''):
                return response.text
            return None # Not HTML content
        except requests.RequestException as e:
            print(f"Error fetching {url}: {e}")
            return None

    def _parse_html_for_links_and_forms(self, html_content, current_url):
        if not html_content:
            return [], []

        soup = BeautifulSoup(html_content, 'html.parser')
        links = []
        forms = []

        # Extract links
        for anchor in soup.find_all('a', href=True):
            href = anchor['href']
            absolute_url = urljoin(current_url, href)
            # Basic filter: stay on the same domain, avoid mailto/tel etc.
            if urlparse(absolute_url).netloc == self.parsed_base_url.netloc:
                links.append(absolute_url)

        # Extract forms
        for form_tag in soup.find_all('form'):
            action = form_tag.get('action', '')
            method = form_tag.get('method', 'GET')
            form_id = form_tag.get('id')
            form_class = form_tag.get('class') # Returns a list of classes

            # Resolve form action URL relative to the current page's URL
            absolute_form_action = urljoin(current_url, action)
            
            discovered_form = DiscoveredForm(absolute_form_action, method, form_id, form_class)

            # Find all input types within the form
            for input_tag in form_tag.find_all(['input', 'textarea', 'select', 'button']):
                name = input_tag.get('name')
                input_type = input_tag.name # 'input', 'textarea', 'select', 'button'
                if input_tag.name == 'input':
                    input_type = input_tag.get('type', 'text') # More specific type for <input>
                
                value = input_tag.get('value', '')
                if input_tag.name == 'textarea':
                    value = input_tag.string or ''
                
                # TODO: For select, capture options. For buttons, capture type (submit, button).
                
                if name: # Only consider inputs with a 'name' attribute for now
                    # Basic HTML context (can be expanded)
                    html_context = {'tag': input_tag.name}
                    if form_id:
                        html_context['form_id'] = form_id

                    inp = DiscoveredInput(name, input_type, value, html_context_details=html_context)
                    discovered_form.add_input(inp)
            
            if discovered_form.inputs: # Only add forms that have named inputs
                forms.append(discovered_form)
                print(f"Discovered form on {current_url}: Action='{absolute_form_action}', Method='{method}', Inputs={len(discovered_form.inputs)}")


        return list(set(links)), forms


    def crawl(self, max_pages=10):
        if not self.base_url:
            print("Base URL is not set. Cannot start crawl.")
            return

        queue = [self.base_url]
        self.visited_urls.add(self.base_url)
        self.sitemap.append(self.base_url)

        pages_crawled = 0
        while queue and pages_crawled < max_pages:
            current_url = queue.pop(0)
            print(f"Crawling: {current_url}")
            pages_crawled += 1

            html_content = self._get_page_content(current_url)
            if html_content:
                new_links, new_forms = self._parse_html_for_links_and_forms(html_content, current_url)
                
                self.sitemap.extend(l for l in new_links if l not in self.sitemap) # Add to sitemap
                self.discovered_forms.extend(f for f in new_forms if f not in self.discovered_forms) # Avoid duplicates if crawling same page via different routes

                for link in new_links:
                    if link not in self.visited_urls:
                        self.visited_urls.add(link)
                        queue.append(link)
            
            # TODO: Add delay between requests if configured

        print(f"Crawl finished. Visited {len(self.visited_urls)} URLs. Found {len(self.sitemap)} unique URLs in sitemap.")
        print(f"Discovered {len(self.discovered_forms)} forms with inputs.")
        # For debugging:
        # for form_obj in self.discovered_forms:
        #     print(f"  Form Action: {form_obj.action}, Method: {form_obj.method}")
        #     for input_obj in form_obj.inputs:
        #         print(f"    Input: {input_obj.name}, Type: {input_obj.type}, Value: {input_obj.current_value}")


def start_crawl(url, max_pages=10):
    if not url.startswith(('http://', 'https://')):
        print(f"Invalid URL scheme for: {url}. Please use http:// or https://")
        return [], []
        
    crawler = Crawler(url)
    crawler.crawl(max_pages=max_pages)
    # The crawler now stores discovered_forms internally.
    # This function could return them, or they could be accessed from the crawler instance.
    return crawler.sitemap, crawler.discovered_forms

if __name__ == '__main__':
    # Example Usage (for direct testing of this script)
    # test_url = "http://testphp.vulnweb.com" # Example site with forms
    test_url = "http://localhost:8000" # if you have a simple local server with forms
    print(f"Starting crawl on {test_url}...")
    sitemap_results, forms_results = start_crawl(test_url, max_pages=5)
    
    print("\nSitemap:")
    for item in sitemap_results:
        print(item)
    
    print("\nDiscovered Forms:")
    for form in forms_results:
        print(f"  Action: {form.action}, Method: {form.method}, ID: {form.form_id}, Class: {form.form_class}")
        for inp in form.inputs:
            print(f"    Input: {inp.name}, Type: {inp.type}, Value: '{inp.current_value}', Context: {inp.html_context_details}")
