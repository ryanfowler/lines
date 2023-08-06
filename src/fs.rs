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

use crossbeam::channel::{unbounded, Sender};
use ignore::{WalkBuilder, WalkState};
use rustc_hash::FxHashMap;
use std::fs::File;
use std::io::{self, Read};
use std::path::{Path, PathBuf};

use crate::cli;
use crate::lang;

pub fn visit_path_parallel(path: &PathBuf) -> Vec<cli::LangOut> {
    let (ch_s, ch_r) = unbounded();

    WalkBuilder::new(path)
        .threads(std::thread::available_parallelism().map_or(1, |v| v.get()) + 1)
        .build_parallel()
        .run(|| {
            let mut buf = [0u8; 1 << 14];
            let mut map = Map {
                s: ch_s.clone(),
                map: FxHashMap::default(),
            };
            Box::new(move |result| {
                let Ok(entry) = result else {
                    eprintln!("Error: {}", result.err().unwrap());
                    return WalkState::Continue;
                };
                let path = entry.path();
                if path.is_dir() {
                    return WalkState::Continue;
                }
                let Some(ext) = path.extension() else {
                    return WalkState::Continue;
                };
                let Some(ext_str) = ext.to_str() else {
                    return WalkState::Continue;
                };
                let Some(language) = lang::get_language(&ext_str.to_ascii_lowercase()) else {
                    return WalkState::Continue;
                };
                match lines_in_file(path, &mut buf) {
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

    let mut langs = FxHashMap::default();
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

fn map_to_vec(map: FxHashMap<lang::Language, LangResult>) -> Vec<cli::LangOut> {
    let mut langs: Vec<_> = map
        .into_iter()
        .map(|(key, val)| cli::LangOut {
            language: key,
            num_files: val.file_cnt,
            num_lines: val.line_cnt,
        })
        .collect();
    langs.sort_unstable_by(|a, b| a.num_lines.cmp(&b.num_lines).reverse());
    langs
}

struct Map {
    s: Sender<FxHashMap<lang::Language, LangResult>>,
    map: FxHashMap<lang::Language, LangResult>,
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
