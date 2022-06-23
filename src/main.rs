use std::time::Instant;

mod cli;
mod fs;
mod lang;

fn main() {
    let start = Instant::now();
    let opt = cli::get_options();

    let globs = vec!["!node_modules/", "!vendor/"];
    let langs = fs::visit_path_parallel(&opt.path, globs);

    let total_num_files = langs.iter().fold(0u64, |sum, l| sum + l.num_files);
    let total_num_lines = langs.iter().fold(0u64, |sum, l| sum + l.num_lines);

    let end = start.elapsed();
    let total_ms = (end.as_secs() * 1000) + end.subsec_millis() as u64;

    let out = cli::Output {
        languages: langs,
        total_num_files,
        total_num_lines,
        elapsed_ms: if opt.timing { Some(total_ms) } else { None },
    };

    cli::write_output(&out, opt.format);
}
