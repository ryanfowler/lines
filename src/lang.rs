// The MIT License (MIT)
//
// Copyright (c) 2022 Ryan Fowler
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// copies of the Software, and to permit persons to whom the Software is
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

use phf::phf_map;
use serde::Serializer;
use std::fmt;

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub enum Language {
    Ada,
    AppleScript,
    Assembly,
    C,
    Clojure,
    ClojureScript,
    Cobol,
    CoffeeScript,
    Cpp,
    CSharp,
    Css,
    //D,
    Dart,
    Elixir,
    Elm,
    Erlang,
    Fortran,
    Go,
    Groovy,
    Handlebars,
    Haskell,
    Html,
    Java,
    JavaScript,
    //JavaScriptReact,
    Json,
    Julia,
    Kotlin,
    Less,
    Lua,
    Markdown,
    Mustache,
    ObjectiveC,
    OCaml,
    Pascal,
    Perl,
    Php,
    Prolog,
    ProtocolBuffer,
    Python,
    R,
    Racket,
    ReasonMl,
    Ruby,
    Rust,
    Sass,
    Scala,
    Shell,
    Sql,
    Stylus,
    Swift,
    Toml,
    TypeScript,
    //TypeScriptReact,
    Vue,
    WebAssembly,
    Xml,
    Yaml,
}

impl Language {
    pub fn as_str(&self) -> &str {
        match self {
            Language::Ada => "Ada",
            Language::AppleScript => "AppleScript",
            Language::Assembly => "Assembly",
            Language::C => "C",
            Language::Clojure => "Clojure",
            Language::ClojureScript => "ClojureScript",
            Language::Cobol => "COBOL",
            Language::CoffeeScript => "CoffeScript",
            Language::Cpp => "C++",
            Language::CSharp => "C#",
            Language::Css => "CSS",
            //Language::D => "D",
            Language::Dart => "Dart",
            Language::Elixir => "Elixir",
            Language::Elm => "Elm",
            Language::Erlang => "Erlang",
            Language::Fortran => "Fortran",
            Language::Go => "Go",
            Language::Groovy => "Groovy",
            Language::Handlebars => "Handlebars",
            Language::Haskell => "Haskell",
            Language::Html => "HTML",
            Language::Java => "Java",
            Language::JavaScript => "JavaScript",
            //Language::JavaScriptReact => "JavaScriptReact",
            Language::Json => "JSON",
            Language::Julia => "Julia",
            Language::Kotlin => "Kotlin",
            Language::Less => "LESS",
            Language::Lua => "Lua",
            Language::Markdown => "Markdown",
            Language::Mustache => "Mustache",
            Language::ObjectiveC => "Objective-C",
            Language::OCaml => "OCaml",
            Language::Pascal => "Pascal",
            Language::Perl => "Perl",
            Language::Php => "PHP",
            Language::Prolog => "Prolog",
            Language::ProtocolBuffer => "ProtocolBuffer",
            Language::Python => "Python",
            Language::R => "R",
            Language::Racket => "Racket",
            Language::ReasonMl => "ReasonML",
            Language::Ruby => "Ruby",
            Language::Rust => "Rust",
            Language::Sass => "SASS",
            Language::Scala => "Scala",
            Language::Shell => "Shell",
            Language::Sql => "SQL",
            Language::Stylus => "Stylus",
            Language::Swift => "Swift",
            Language::Toml => "TOML",
            Language::TypeScript => "TypeScript",
            //Language::TypeScriptReact => "TypeScriptReact",
            Language::Vue => "Vue",
            Language::WebAssembly => "WebAssembly",
            Language::Xml => "XML",
            Language::Yaml => "YAML",
        }
    }
}

impl serde::Serialize for Language {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(self.as_str())
    }
}

impl fmt::Display for Language {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        fmt::Debug::fmt(self, f)
    }
}

pub fn get_language(s: &str) -> Option<Language> {
    EXT_LANGS.get(s).copied()
}

static EXT_LANGS: phf::Map<&str, Language> = phf_map! {
    "adb" => Language::Ada,
    "ads" => Language::Ada,
    "applescript" => Language::AppleScript,
    "asm" => Language::Assembly,
    "c" => Language::C,
    "cbl" => Language::Cobol,
    "clj" => Language::Clojure,
    "cljs" => Language::ClojureScript,
    "cob" => Language::Cobol,
    "coffee" => Language::CoffeeScript,
    "cpp" => Language::Cpp,
    "cpy" => Language::Cobol,
    "cs" => Language::CSharp,
    "css" => Language::Css,
    //"d" => Language::D,
    "dart" => Language::Dart,
    "elm" => Language::Elm,
    "erl" => Language::Erlang,
    "ex" => Language::Elixir,
    "exs" => Language::Elixir,
    "f" => Language::Fortran,
    "f90" => Language::Fortran,
    "for" => Language::Fortran,
    "go" => Language::Go,
    "groovy" => Language::Groovy,
    "gsh" => Language::Groovy,
    "gvy" => Language::Groovy,
    "gy" => Language::Groovy,
    "handlebars" => Language::Handlebars,
    "hbs" => Language::Handlebars,
    "hrl" => Language::Erlang,
    "hs" => Language::Haskell,
    "htm" => Language::Html,
    "html" => Language::Html,
    "inc" => Language::Pascal,
    "java" => Language::Java,
    "jl" => Language::Julia,
    "js" => Language::JavaScript,
    "json" => Language::Json,
    "jsx" => Language::JavaScript,
    "kt" => Language::Kotlin,
    "kts" => Language::Kotlin,
    "less" => Language::Less,
    "lua" => Language::Lua,
    "m" => Language::ObjectiveC,
    "md" => Language::Markdown,
    "ml" => Language::OCaml,
    "mli" => Language::OCaml,
    "mustache" => Language::Mustache,
    "p" => Language::Prolog,
    "pas" => Language::Pascal,
    "php" => Language::Php,
    "pl" => Language::Perl,
    "pm" => Language::Perl,
    "pp" => Language::Pascal,
    "pro" => Language::Prolog,
    "proto" => Language::ProtocolBuffer,
    "py" => Language::Python,
    "r" => Language::R,
    "rb" => Language::Ruby,
    "re" => Language::ReasonMl,
    "rkt" => Language::Racket,
    "rs" => Language::Rust,
    "s" => Language::Assembly,
    "sass" => Language::Sass,
    "scala" => Language::Scala,
    "scpt" => Language::AppleScript,
    "scptd" => Language::AppleScript,
    "scss" => Language::Sass,
    "sh" => Language::Shell,
    "sql" => Language::Sql,
    "styl" => Language::Stylus,
    "swift" => Language::Swift,
    "toml" => Language::Toml,
    "ts" => Language::TypeScript,
    "tsx" => Language::TypeScript,
    "vue" => Language::Vue,
    "wat" => Language::WebAssembly,
    "xml" => Language::Xml,
    "yaml" => Language::Yaml,
    "yml" => Language::Yaml,
};
