# Sphinx configuration for nwpc-data-client documentation.
#
# Build from the repository root:
#
#   pip install -r docs/requirements.txt
#   sphinx-build -b html docs docs/_build/html
#
# Or simply `make docs` (which also regenerates the CLI reference pages).

from pathlib import Path

repo_root = Path(__file__).resolve().parent.parent

# -- Project information -----------------------------------------------------

project = "nwpc-data-client"
author = "CEMC"
copyright = "CEMC"

# Keep in sync with scripts/version.sh: VERSION holds the next release
# version; append git describe info when building from a non-tagged commit.
version = (repo_root / "VERSION").read_text(encoding="utf-8").strip()
release = version

# -- General configuration ---------------------------------------------------

extensions = [
    "myst_parser",
    "sphinx_copybutton",
    "sphinx_design",
]

myst_enable_extensions = [
    "colon_fence",
    "deflist",
    "fieldlist",
    "tasklist",
]

myst_heading_anchors = 3

exclude_patterns = [
    "_build",
    "Thumbs.db",
    ".DS_Store",
    "requirements.txt",
]

language = "zh_CN"

# -- HTML output -------------------------------------------------------------

html_theme = "pydata_sphinx_theme"
html_title = "nwpc-data-client"

html_theme_options = {
    "github_url": "https://github.com/cemc-oper/nwpc-data-client",
    "navbar_end": ["theme-switcher", "navbar-icon-links"],
    "show_toc_level": 2,
    "show_version_warning_banner": False,
}

# -- MyST / copybutton tweaks ------------------------------------------------

# Strip the shell prompt when copying command examples.
copybutton_prompt_text = "$ "
