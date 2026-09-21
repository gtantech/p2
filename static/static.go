package static

import _ "embed"

//go:embed home_style.css
var StaticHomeCss []byte

//go:embed home_script.js
var StaticHomeScript []byte
