package treesitter

// preferredSyntaxOverrides mirrors the preferred syntax emphasis rules used for
// Zed so tree-sitter-based targets stay visually aligned.
var preferredSyntaxOverrides = map[string]syntaxOverride{
	"keyword": {
		FontWeight: 700,
		FontStyle:  "italic",
	},
	"keyword.control": {
		FontWeight: 700,
	},
	"keyword.storage": {
		FontWeight: 400,
		FontStyle:  "normal",
	},
	"lexical_declaration kind": {
		FontStyle: "italic",
	},
	"operator": {
		FontWeight: 400,
		FontStyle:  "normal",
	},
	"variable": {
		FontWeight: 400,
		FontStyle:  "normal",
	},
	"variable.builtin": {
		FontWeight: 700,
		FontStyle:  "italic",
	},
	"variable.parameter": {
		FontStyle: "italic",
	},
	"variable.parameter.builtin": {
		FontWeight: 700,
		FontStyle:  "italic",
	},
	"variable.member": {
		FontStyle: "normal",
	},
	"property": {
		FontWeight: 400,
		FontStyle:  "normal",
	},
	"function": {
		FontWeight: 700,
	},
	"method": {
		FontWeight: 700,
	},
	"function.builtin": {
		FontWeight: 700,
	},
	"function.method": {
		FontWeight: 700,
	},
	"type": {
		FontWeight: 700,
		FontStyle:  "italic",
	},
	"class": {
		FontStyle: "italic",
	},
	"constructor": {
		FontStyle: "italic",
	},
	"body.kind": {
		FontWeight: 700,
	},
	"string": {
		FontWeight: 400,
		FontStyle:  "normal",
	},
	"comment": {
		FontStyle: "italic",
	},
	"comment.doc": {
		FontStyle: "italic",
	},
	"constant": {
		FontWeight: 700,
	},
	"boolean": {
		FontWeight: 700,
	},
	"number": {
		FontWeight: 400,
		FontStyle:  "normal",
	},
	"punctuation": {
		FontWeight: 700,
	},
	"punctuation.bracket": {
		FontWeight: 700,
	},
	"punctuation.delimiter": {
		FontWeight: 700,
	},
	"tag": {
		FontWeight: 700,
	},
	"attribute": {
		FontStyle: "italic",
	},
}
