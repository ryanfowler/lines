use crossbeam::channel::{unbounded, Sender};
use fnv::FnvHashMap;
use ignore::{overrides, WalkBuilder, WalkState};
use std::fs::File;
use std::io::{self, Read};
use std::path::{Path, PathBuf};

use crate::cli;
use crate::lang;

pub fn visit_path_parallel(path: &PathBuf, globs: Vec<&str>) -> Vec<cli::LangOut> {
    let (ch_s, ch_r) = unbounded();
    let overrides = parse_overrides(path, globs);

    WalkBuilder::new(path)
        .threads(num_cpus::get() + 1)
        .overrides(overrides)
        .build_parallel()
        .run(|| {
            let mut buf = [0u8; 1 << 14];
            let mut map = Map {
                s: ch_s.clone(),
                map: FnvHashMap::default(),
            };
            Box::new(move |result| {
                let entry = match result {
                    Err(err) => {
                        eprintln!("Error: {}", err);
                        return WalkState::Continue;
                    }
                    Ok(entry) => entry,
                };
                let path = entry.into_path();
                if path.is_dir() {
                    return WalkState::Continue;
                }
                let ext = match path.extension() {
                    None => return WalkState::Continue,
                    Some(ext) => ext,
                };
                let ext_str = match ext.to_str() {
                    None => return WalkState::Continue,
                    Some(ext_str) => ext_str,
                };
                let language = match lang::get_language(&ext_str.to_ascii_lowercase()) {
                    None => return WalkState::Continue,
                    Some(lang) => lang,
                };
                match lines_in_file(&path, &mut buf) {
                    Ok(lines) => {
                        let counter = map.map.entry(language).or_insert_with(LangResult::new);
                        counter.file_cnt += 1;
                        counter.line_cnt += lines;
                    }
                    Err(err) => {
                        eprintln!("Error: {:?}: {}", &path, err);
                    }
                }
                WalkState::Continue
            })
        });
    drop(ch_s);

    let mut langs = FnvHashMap::default();
    for m in ch_r.iter() {
        for (&key, val) in m.iter() {
            let l = langs.entry(key).or_insert_with(LangResult::new);
            l.file_cnt += val.file_cnt;
            l.line_cnt += val.line_cnt;
        }
    }
    map_to_vec(langs)
}

fn lines_in_file(path: &Path, buf: &mut [u8]) -> io::Result<u64> {
    let mut file = File::open(path)?;
    let mut cnt = 1;
    loop {
        let n = file.read(buf)?;
        if n == 0 {
            return Ok(cnt as u64);
        }
        cnt += bytecount::count(&buf[0..n], b'\n');
    }
}

fn parse_overrides(path: &PathBuf, globs: Vec<&str>) -> overrides::Override {
    let mut override_builder = overrides::OverrideBuilder::new(path);
    for glob in globs.iter() {
        override_builder.add(glob).unwrap();
    }
    override_builder.build().unwrap()
}

fn map_to_vec(map: FnvHashMap<lang::Language, LangResult>) -> Vec<cli::LangOut> {
    let mut langs = Vec::with_capacity(map.len());
    for (key, val) in map.iter() {
        langs.push(cli::LangOut {
            language: *key,
            num_files: val.file_cnt,
            num_lines: val.line_cnt,
        });
    }
    langs.sort_by(|a, b| a.num_lines.cmp(&b.num_lines).reverse());
    langs
}

struct Map {
    s: Sender<FnvHashMap<lang::Language, LangResult>>,
    map: FnvHashMap<lang::Language, LangResult>,
}

impl Drop for Map {
    fn drop(&mut self) {
        self.s.send(self.map.clone()).unwrap();
    }
}

#[derive(Clone, Copy, Debug)]
struct LangResult {
    file_cnt: u64,
    line_cnt: u64,
}

impl LangResult {
    fn new() -> Self {
        LangResult {
            file_cnt: 0,
            line_cnt: 0,
        }
    }
}
