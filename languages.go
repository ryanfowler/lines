//
// lines
//
// Copyright (c) 2015 Ryan Fowler
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

type Language struct {
	Name  string
	Ext   string
	lCom  string
	bComS string
	bComE string
}

var (
	// language extensions and info
	LANGS = map[string]*Language{
		".c": &Language{
			Name:  "C",
			Ext:   ".c",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".clj": &Language{
			Name:  "Clojure",
			Ext:   ".clj",
			lCom:  `;`,
			bComS: "",
			bComE: "",
		},
		".cljs": &Language{
			Name:  "ClojureScript",
			Ext:   ".cljs",
			lCom:  `;`,
			bComS: "",
			bComE: "",
		},
		".coffee": &Language{
			Name:  "CoffeeScript",
			Ext:   ".coffee",
			lCom:  `#`,
			bComS: `#{3}`,
			bComE: `#{3}`,
		},
		".cpp": &Language{
			Name:  "C++",
			Ext:   ".cpp",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".cs": &Language{
			Name:  "C#",
			Ext:   ".cs",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".css": &Language{
			Name:  "CSS",
			Ext:   ".css",
			lCom:  "",
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".go": &Language{
			Name:  "Go",
			Ext:   ".go",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".h": &Language{
			Name:  "C",
			Ext:   ".h",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".handlebars": &Language{
			Name:  "Handlebars",
			Ext:   ".handlebars",
			lCom:  "",
			bComS: `\{\{\!`,
			bComE: `\}\}`,
		},
		".hbs": &Language{
			Name:  "Handlebars",
			Ext:   ".hbs",
			lCom:  "",
			bComS: `\{\{\!`,
			bComE: `\}\}`,
		},
		".hs": &Language{
			Name:  "Haskell",
			Ext:   ".hs",
			lCom:  `--`,
			bComS: `\{-`,
			bComE: `-\}`,
		},
		".htm": &Language{
			Name:  "HTML",
			Ext:   ".htm",
			lCom:  "",
			bComS: `<\!--`,
			bComE: `-->`,
		},
		".html": &Language{
			Name:  "HTML",
			Ext:   ".html",
			lCom:  "",
			bComS: `<\!--`,
			bComE: `-->`,
		},
		".java": &Language{
			Name:  "Java",
			Ext:   ".java",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".js": &Language{
			Name:  "Javascript",
			Ext:   ".js",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".less": &Language{
			Name:  "LESS",
			Ext:   ".less",
			lCom:  "",
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".m": &Language{
			Name:  "Objective-C",
			Ext:   ".m",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".mustache": &Language{
			Name:  "Mustache",
			Ext:   ".mustache",
			lCom:  "",
			bComS: `\{\{\!`,
			bComE: `\}\}`,
		},
		".php": &Language{
			Name:  "PHP",
			Ext:   ".php",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".py": &Language{
			Name:  "Python",
			Ext:   ".py",
			lCom:  `#`,
			bComS: `\"{3}|\'{3}`,
			bComE: `\"{3}|\'{3}`,
		},
		".rb": &Language{
			Name:  "Ruby",
			Ext:   ".rb",
			lCom:  `#`,
			bComS: `\=begin`,
			bComE: `\=end`,
		},
		".rkt": &Language{
			Name:  "Racket",
			Ext:   ".rkt",
			lCom:  `;`,
			bComS: `#\|`,
			bComE: `\|#`,
		},
		".rs": &Language{
			Name:  "Rust",
			Ext:   ".rs",
			lCom:  `\/{2,3}`,
			bComS: "",
			bComE: "",
		},
		".sass": &Language{
			Name:  "SASS",
			Ext:   ".sass",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".scala": &Language{
			Name:  "Scala",
			Ext:   ".scala",
			lCom:  "",
			bComS: `\/\*\*`,
			bComE: `\*\/`,
		},
		".styl": &Language{
			Name:  "Stylus",
			Ext:   ".styl",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".swift": &Language{
			Name:  "Swift",
			Ext:   ".swift",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".vue": &Language{
			Name:  "Vue",
			Ext:   ".vue",
			lCom:  `\/\/`,
			bComS: `\/\*`,
			bComE: `\*\/`,
		},
		".xml": &Language{
			Name:  "XML",
			Ext:   ".xml",
			lCom:  "",
			bComS: `<\!--`,
			bComE: `-->`,
		},
	}
)
