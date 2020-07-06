use num_format::{Locale, ToFormattedString};
use prettytable::{format, Table};
use serde;
use serde::Serialize;
use serde_json;
use std::path::PathBuf;
use std::str::FromStr;
use std::string::ToString;
use structopt::StructOpt;

use crate::lang;

#[derive(Debug)]
pub enum Format {
    Json,
    JsonPretty,
    Table,
}

impl FromStr for Format {
    type Err = String;
    fn from_str(format: &str) -> Result<Self, Self::Err> {
        match format {
            "json" => Ok(Format::Json),
            "json-pretty" => Ok(Format::JsonPretty),
            "table" => Ok(Format::Table),
            _ => Err(format.to_string()),
        }
    }
}

#[derive(Debug, StructOpt)]
#[structopt(name = "lines", about = "Count lines of code.")]
pub struct Opt {
    /// Output format.
    #[structopt(short, long, default_value = "table")]
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
        Format::Json => write_json(&out),
        Format::JsonPretty => write_json_pretty(&out),
        Format::Table => write_table(&out),
    }
}

fn write_json(out: &Output) {
    println!("{}", serde_json::to_string(&out).unwrap());
}

fn write_json_pretty(out: &Output) {
    println!("{}", serde_json::to_string_pretty(&out).unwrap());
}

fn write_table(out: &Output) {
    let mut table = Table::new();
    table.set_format(*format::consts::FORMAT_NO_LINESEP_WITH_TITLE);
    table.set_titles(row!["Language", "Files", "Lines"]);
    for lang in out.languages.iter() {
        table.add_row(row![
            lang.language.as_str(),
            r->lang.num_files.to_formatted_string(&Locale::en),
            r->lang.num_lines.to_formatted_string(&Locale::en)]);
    }
    table.printstd();
    if let Some(elapsed_ms) = out.elapsed_ms {
        println!("Took: {}ms", elapsed_ms);
    }
}
