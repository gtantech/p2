package static

import "embed"

//go:embed home_style.css
var StaticHomeCss embed.FS

//go:embed home_script.js
var StaticHomeScript embed.FS
