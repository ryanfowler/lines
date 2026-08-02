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

    /// Exclude regex patterns (can be used multiple times).
    #[clap(short = 'e', long = "exclude", action = clap::ArgAction::Append)]
    pub exclude: Vec<String>,

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

struct Row {
    language: &'static str,
    files: String,
    lines: String,
}

fn write_table(out: &Output) {
    println!("{}", format_table(out));

    if let Some(elapsed_ms) = out.elapsed_ms {
        println!("\nTook: {elapsed_ms}ms");
    }
}

fn format_table(out: &Output) -> String {
    let mut rows = Vec::new();
    for lang in &out.languages {
        rows.push(Row {
            language: lang.language.as_str(),
            files: lang.num_files.to_formatted_string(&Locale::en),
            lines: lang.num_lines.to_formatted_string(&Locale::en),
        });
    }

    let has_total = out.languages.len() != 1;
    if has_total {
        rows.push(Row {
            language: "Total",
            files: out.total_num_files.to_formatted_string(&Locale::en),
            lines: out.total_num_lines.to_formatted_string(&Locale::en),
        });
    }

    let widths = [
        column_width("Language", rows.iter().map(|row| row.language)),
        column_width("Files", rows.iter().map(|row| row.files.as_str())),
        column_width("Lines", rows.iter().map(|row| row.lines.as_str())),
    ];

    let mut output = vec![format_header(&widths), format_separator(&widths, '+')];
    for (index, row) in rows.iter().enumerate() {
        if has_total && index == rows.len() - 1 {
            output.push(format_separator(&widths, ' '));
        }
        output.push(format_row(row, &widths));
    }

    output.join("\n")
}

fn column_width<'a>(header: &'a str, values: impl Iterator<Item = &'a str>) -> usize {
    values
        .map(display_width)
        .chain(std::iter::once(display_width(header)))
        .max()
        .unwrap_or(0)
        + 2
}

fn display_width(value: &str) -> usize {
    value.chars().count()
}

fn format_header(widths: &[usize; 3]) -> String {
    [
        format_left_cell("Language", widths[0]),
        format_left_cell("Files", widths[1]),
        format_left_cell("Lines", widths[2]),
    ]
    .join("|")
}

fn format_separator(widths: &[usize; 3], separator: char) -> String {
    let separator = separator.to_string();
    [
        "-".repeat(widths[0]),
        "-".repeat(widths[1]),
        "-".repeat(widths[2]),
    ]
    .join(&separator)
}

fn format_row(row: &Row, widths: &[usize; 3]) -> String {
    [
        format_left_cell(row.language, widths[0]),
        format_right_cell(&row.files, widths[1]),
        format_right_cell(&row.lines, widths[2]),
    ]
    .join("|")
}

fn format_left_cell(value: &str, width: usize) -> String {
    format!(" {value}{} ", " ".repeat(width - display_width(value) - 2),)
}

fn format_right_cell(value: &str, width: usize) -> String {
    format!(" {}{value} ", " ".repeat(width - display_width(value) - 2),)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn format_table_aligns_rows_and_total() {
        let output = Output {
            languages: vec![
                LangOut {
                    language: lang::Language::Rust,
                    num_files: 12,
                    num_lines: 123,
                },
                LangOut {
                    language: lang::Language::Markdown,
                    num_files: 3,
                    num_lines: 45,
                },
            ],
            total_num_files: 15,
            total_num_lines: 168,
            elapsed_ms: None,
        };

        assert_eq!(
            format_table(&output),
            concat!(
                " Language | Files | Lines \n",
                "----------+-------+-------\n",
                " Rust     |    12 |   123 \n",
                " Markdown |     3 |    45 \n",
                "---------- ------- -------\n",
                " Total    |    15 |   168 ",
            )
        );
    }

    #[test]
    fn format_table_omits_total_for_one_language() {
        let output = Output {
            languages: vec![LangOut {
                language: lang::Language::Rust,
                num_files: 1,
                num_lines: 9,
            }],
            total_num_files: 1,
            total_num_lines: 9,
            elapsed_ms: Some(7),
        };

        assert_eq!(
            format_table(&output),
            concat!(
                " Language | Files | Lines \n",
                "----------+-------+-------\n",
                " Rust     |     1 |     9 ",
            )
        );
    }

    #[test]
    fn format_table_handles_no_languages() {
        let output = Output {
            languages: Vec::new(),
            total_num_files: 0,
            total_num_lines: 0,
            elapsed_ms: None,
        };

        assert_eq!(
            format_table(&output),
            concat!(
                " Language | Files | Lines \n",
                "----------+-------+-------\n",
                "---------- ------- -------\n",
                " Total    |     0 |     0 ",
            )
        );
    }

    #[test]
    fn format_table_handles_large_formatted_numbers() {
        let output = Output {
            languages: vec![LangOut {
                language: lang::Language::Rust,
                num_files: 1_234_567,
                num_lines: 987_654_321,
            }],
            total_num_files: 1_234_567,
            total_num_lines: 987_654_321,
            elapsed_ms: None,
        };

        assert_eq!(
            format_table(&output),
            concat!(
                " Language | Files     | Lines       \n",
                "----------+-----------+-------------\n",
                " Rust     | 1,234,567 | 987,654,321 ",
            )
        );
    }
}
