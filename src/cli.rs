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

use clap::Parser;
use num_format::{Locale, ToFormattedString};
use serde::Serialize;
use std::path::PathBuf;
use std::str::FromStr;
use std::string::ToString;
use tabled::{
    settings::{
        object::{Columns, Rows},
        Alignment, Border, Modify, Style,
    },
    Table, Tabled,
};

use crate::lang;

#[derive(Clone, Debug)]
pub enum Format {
    Json,
    Table,
}

impl FromStr for Format {
    type Err = String;
    fn from_str(format: &str) -> Result<Self, Self::Err> {
        match format {
            "json" => Ok(Format::Json),
            "table" => Ok(Format::Table),
            _ => Err(format.to_string()),
        }
    }
}

/// Count lines of code.
#[derive(Debug, Parser)]
#[clap(version, about)]
pub struct Args {
    /// Output format ("table" or "json").
    #[clap(short = 'o', long = "output", default_value = "table")]
    pub format: Format,

    /// Show timing information.
    #[clap(short, long)]
    pub timing: bool,

    /// Directory or file to scan.
    #[clap(default_value = ".")]
    pub path: PathBuf,
}

#[derive(Debug, Serialize)]
pub struct Output {
    pub languages: Vec<LangOut>,
    pub total_num_files: u64,
    pub total_num_lines: u64,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub elapsed_ms: Option<u64>,
}

#[derive(Debug, Serialize)]
pub struct LangOut {
    pub language: lang::Language,
    pub num_files: u64,
    pub num_lines: u64,
}

pub fn get_options() -> Args {
    Args::parse()
}

pub fn write_output(out: &Output, format: Format) {
    match format {
        Format::Json => write_json_pretty(out),
        Format::Table => write_table(out),
    }
}

fn write_json_pretty(out: &Output) {
    println!("{}", serde_json::to_string_pretty(&out).unwrap());
}

#[derive(Tabled)]
struct Row {
    #[tabled(rename = "Language")]
    language: String,
    #[tabled(rename = "Files")]
    files: String,
    #[tabled(rename = "Lines")]
    lines: String,
}

fn write_table(out: &Output) {
    let mut data = Vec::new();
    for lang in &out.languages {
        data.push(Row {
            language: lang.language.as_str().to_string(),
            files: lang.num_files.to_formatted_string(&Locale::en),
            lines: lang.num_lines.to_formatted_string(&Locale::en),
        });
    }

    if out.languages.len() != 1 {
        data.push(Row {
            language: "Total".to_string(),
            files: out.total_num_files.to_formatted_string(&Locale::en),
            lines: out.total_num_lines.to_formatted_string(&Locale::en),
        });
    }

    let mut table = Table::new(&data);
    table
        .with(Style::psql())
        .with(Modify::new(Columns::first()).with(Alignment::left()))
        .with(Modify::new(Columns::new(1..=2)).with(Alignment::right()))
        .with(Modify::new(Rows::first()).with(Alignment::left()));

    if out.languages.len() != 1 {
        table.with(Modify::new(Rows::last()).with(Border::default().top('-')));
    }

    println!("{}", table);

    if let Some(elapsed_ms) = out.elapsed_ms {
        println!("\nTook: {}ms", elapsed_ms);
    }
}
