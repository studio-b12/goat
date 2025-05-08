// Populate the sidebar
//
// This is a script, and not included directly in the page, to control the total size of the book.
// The TOC contains an entry for each page, so if each page includes a copy of the TOC,
// the total size of the page becomes O(n**2).
class MDBookSidebarScrollbox extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.innerHTML = '<ol class="chapter"><li class="chapter-item expanded "><a href="introduction/index.html"><strong aria-hidden="true">1.</strong> Introduction</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="introduction/glossary.html"><strong aria-hidden="true">1.1.</strong> Glossary</a></li></ol></li><li class="chapter-item expanded "><a href="getting-started/index.html"><strong aria-hidden="true">2.</strong> Getting Started</a></li><li class="chapter-item expanded "><a href="command-line-tool/index.html"><strong aria-hidden="true">3.</strong> Command Line Tool</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="command-line-tool/profiles.html"><strong aria-hidden="true">3.1.</strong> Profiles</a></li></ol></li><li class="chapter-item expanded "><a href="explanations/index.html"><strong aria-hidden="true">4.</strong> How does it work?</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="explanations/state.html"><strong aria-hidden="true">4.1.</strong> State Management</a></li><li class="chapter-item expanded "><a href="explanations/lifecycle.html"><strong aria-hidden="true">4.2.</strong> Lifecycle</a></li></ol></li><li class="chapter-item expanded "><a href="goatfile/index.html"><strong aria-hidden="true">5.</strong> Goatfile</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="goatfile/comments.html"><strong aria-hidden="true">5.1.</strong> Comment</a></li><li class="chapter-item expanded "><a href="goatfile/import-statement.html"><strong aria-hidden="true">5.2.</strong> Import Statement</a></li><li class="chapter-item expanded "><a href="goatfile/execute-statement.html"><strong aria-hidden="true">5.3.</strong> Execute Statement</a></li><li class="chapter-item expanded "><a href="goatfile/sections.html"><strong aria-hidden="true">5.4.</strong> Section</a></li><li class="chapter-item expanded "><a href="goatfile/defaults-section.html"><strong aria-hidden="true">5.5.</strong> Defaults Section</a></li><li class="chapter-item expanded "><a href="goatfile/logsections.html"><strong aria-hidden="true">5.6.</strong> Log Section</a></li><li class="chapter-item expanded "><a href="goatfile/requests/index.html"><strong aria-hidden="true">5.7.</strong> Request</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="goatfile/requests/method-and-url.html"><strong aria-hidden="true">5.7.1.</strong> Method and URL</a></li><li class="chapter-item expanded "><a href="goatfile/requests/options.html"><strong aria-hidden="true">5.7.2.</strong> Options</a></li><li class="chapter-item expanded "><a href="goatfile/requests/header.html"><strong aria-hidden="true">5.7.3.</strong> Headers</a></li><li class="chapter-item expanded "><a href="goatfile/requests/query-params.html"><strong aria-hidden="true">5.7.4.</strong> Query Parameters</a></li><li class="chapter-item expanded "><a href="goatfile/requests/auth.html"><strong aria-hidden="true">5.7.5.</strong> Auth</a></li><li class="chapter-item expanded "><a href="goatfile/requests/body.html"><strong aria-hidden="true">5.7.6.</strong> Body</a></li><li class="chapter-item expanded "><a href="goatfile/requests/formdata.html"><strong aria-hidden="true">5.7.7.</strong> FormData</a></li><li class="chapter-item expanded "><a href="goatfile/requests/formurlencoded.html"><strong aria-hidden="true">5.7.8.</strong> FormUrlEncoded</a></li><li class="chapter-item expanded "><a href="goatfile/requests/prescript.html"><strong aria-hidden="true">5.7.9.</strong> PreScript</a></li><li class="chapter-item expanded "><a href="goatfile/requests/script.html"><strong aria-hidden="true">5.7.10.</strong> Script</a></li></ol></li></ol></li><li class="chapter-item expanded "><a href="templating/index.html"><strong aria-hidden="true">6.</strong> Templating</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="templating/builtins.html"><strong aria-hidden="true">6.1.</strong> Built-ins</a></li></ol></li><li class="chapter-item expanded "><a href="scripting/index.html"><strong aria-hidden="true">7.</strong> Scripting</a></li><li><ol class="section"><li class="chapter-item expanded "><a href="scripting/builtins.html"><strong aria-hidden="true">7.1.</strong> Built-ins</a></li></ol></li><li class="chapter-item expanded "><a href="project-structure/index.html"><strong aria-hidden="true">8.</strong> Project Structure</a></li></ol>';
        // Set the current, active page, and reveal it if it's hidden
        let current_page = document.location.href.toString().split("#")[0].split("?")[0];
        if (current_page.endsWith("/")) {
            current_page += "index.html";
        }
        var links = Array.prototype.slice.call(this.querySelectorAll("a"));
        var l = links.length;
        for (var i = 0; i < l; ++i) {
            var link = links[i];
            var href = link.getAttribute("href");
            if (href && !href.startsWith("#") && !/^(?:[a-z+]+:)?\/\//.test(href)) {
                link.href = path_to_root + href;
            }
            // The "index" page is supposed to alias the first chapter in the book.
            if (link.href === current_page || (i === 0 && path_to_root === "" && current_page.endsWith("/index.html"))) {
                link.classList.add("active");
                var parent = link.parentElement;
                if (parent && parent.classList.contains("chapter-item")) {
                    parent.classList.add("expanded");
                }
                while (parent) {
                    if (parent.tagName === "LI" && parent.previousElementSibling) {
                        if (parent.previousElementSibling.classList.contains("chapter-item")) {
                            parent.previousElementSibling.classList.add("expanded");
                        }
                    }
                    parent = parent.parentElement;
                }
            }
        }
        // Track and set sidebar scroll position
        this.addEventListener('click', function(e) {
            if (e.target.tagName === 'A') {
                sessionStorage.setItem('sidebar-scroll', this.scrollTop);
            }
        }, { passive: true });
        var sidebarScrollTop = sessionStorage.getItem('sidebar-scroll');
        sessionStorage.removeItem('sidebar-scroll');
        if (sidebarScrollTop) {
            // preserve sidebar scroll position when navigating via links within sidebar
            this.scrollTop = sidebarScrollTop;
        } else {
            // scroll sidebar to current active section when navigating via "next/previous chapter" buttons
            var activeSection = document.querySelector('#sidebar .active');
            if (activeSection) {
                activeSection.scrollIntoView({ block: 'center' });
            }
        }
        // Toggle buttons
        var sidebarAnchorToggles = document.querySelectorAll('#sidebar a.toggle');
        function toggleSection(ev) {
            ev.currentTarget.parentElement.classList.toggle('expanded');
        }
        Array.from(sidebarAnchorToggles).forEach(function (el) {
            el.addEventListener('click', toggleSection);
        });
    }
}
window.customElements.define("mdbook-sidebar-scrollbox", MDBookSidebarScrollbox);
