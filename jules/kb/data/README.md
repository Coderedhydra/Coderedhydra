# Jules Knowledge Base Data

This directory is intended to store the offline data sources for the Jules security toolkit's knowledge base.
This includes, but is not limited to:

-   **CVE Archives**: JSON dumps or structured text files of CVE details (e.g., from NVD).
-   **OWASP Documentation**: Markdown or text versions of OWASP Top 10, Testing Guide, Cheat Sheets, etc.
-   **Exploit-DB Dumps**: Information about public exploits (e.g., selected entries from Exploit-DB).
-   **Vulnerability Write-ups**: Collections of public blog posts or articles detailing specific vulnerabilities or bypass techniques.
-   **Custom Knowledge**: User-added notes, patterns, or internal findings.

## Data Format

The `LocalKBService` (placeholder) currently expects:
-   Individual articles as `.json` files, where each file's content can be unmarshalled into the `kb.KBArticle` struct.
-   (Or) A more structured approach would involve a local database or an indexed document store.

## Population

Scripts or manual processes will be needed to populate this directory with relevant, up-to-date (periodically) information from public sources, ensuring it remains useful for the LLM's "self-evolving intelligence" feature.

**Note:** Ensure any data included here respects the licensing terms of its original source.

## Data Licensing

When populating this directory with data from external sources (e.g., CVE details from NVD, OWASP materials, Exploit-DB entries), you **must** respect the original licensing terms of that data. 

It is recommended to include a copy of the relevant licenses alongside the data files you add here, or clearly document the source and its license if the data is transformed. For example, if you include a dump of CVEs, ensure you comply with the NVD's data usage guidelines.

The Jules framework provides the tools to *use* this data; you are responsible for sourcing and licensing it appropriately.
