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

use std::time::Instant;

mod cli;
mod fs;
mod lang;

#[cfg(not(target_os = "windows"))]
#[global_allocator]
static GLOBAL: jemallocator::Jemalloc = jemallocator::Jemalloc;

fn main() {
    let start = Instant::now();
    let opt = cli::get_options();

    let langs = fs::visit_path_parallel(&opt.path);

    let total_num_files = langs.iter().map(|l| l.num_files).sum();
    let total_num_lines = langs.iter().map(|l| l.num_lines).sum();

    let end = start.elapsed();
    let total_ms = (end.as_secs() * 1000) + end.subsec_millis() as u64;

    let out = cli::Output {
        languages: langs,
        total_num_files,
        total_num_lines,
        elapsed_ms: opt.timing.then_some(total_ms),
    };

    cli::write_output(&out, opt.format);
}
