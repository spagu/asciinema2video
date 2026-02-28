# Changelog

All notable changes to this project will be documented in this file.

## [1.0.1] - 2026-02-28

### Fixed
- README demo now displays correctly on GitHub (changed from video to GIF format)

## [1.0.0] - 2026-02-28

### Added
- Initial release
- Parse asciinema v2 format (.cast files from [asciinema.org](https://asciinema.org/))
- VT100 terminal emulation with ANSI color support
- Output formats: MP4, GIF, WebP, WebM, MOV (ProRes 4444)
- Video codec selection: H.264 (default) or H.265 for MP4
- Transparent background support (WebM, MOV)
- Color themes: default, monokai, dracula, solarized-dark, nord, gruvbox
- Custom theme support via JSON file
- Custom TTF font support
- Terminal size override (--cols, --rows)
- Border with customizable width and color
- Rounded corners with configurable radius
- Outer background color/transparency
- Configurable padding
- Quiet mode (--quiet) with timing and file size output
- CLI with cobra framework
- curl install script
- Release packages: deb, rpm, brew, snap, freebsd
- Man pages
