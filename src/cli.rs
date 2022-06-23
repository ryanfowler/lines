use num_format::{Locale, ToFormattedString};
use serde;
use serde::Serialize;
use serde_json;
use std::path::PathBuf;
use std::str::FromStr;
use std::string::ToString;
use structopt::StructOpt;
use tabled::{
    object::{Columns, Rows},
    style::Border,
    Alignment, Modify, Style, Table, Tabled,
};

use crate::lang;

#[derive(Debug)]
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

#[derive(Debug, StructOpt)]
#[structopt(name = "lines", about = "Count lines of code.")]
pub struct Opt {
    /// Output format.
    #[structopt(short = "o", long = "output", default_value = "table")]
    pub format: Format,

    /// Show timing
    #[structopt(short, long)]
    pub timing: bool,

    /// Directory or file to scan.
    #[structopt(parse(from_os_str))]
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

pub fn get_options() -> Opt {
    Opt::from_args()
}

pub fn write_output(out: &Output, format: Format) {
    match format {
        Format::Json => write_json_pretty(&out),
        Format::Table => write_table(&out),
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

    let mut table = Table::new(&data)
        .with(Style::psql())
        .with(Modify::new(Columns::first()).with(Alignment::left()))
        .with(Modify::new(Columns::new(1..=2)).with(Alignment::right()))
        .with(Modify::new(Rows::first()).with(Alignment::left()));

    if out.languages.len() != 1 {
        table = table.with(Modify::new(Rows::last()).with(Border::default().top('-')));
    }

    println!("{}", table);

    if let Some(elapsed_ms) = out.elapsed_ms {
        println!("Took: {}ms", elapsed_ms);
    }
}
