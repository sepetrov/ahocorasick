# ahocorasick

[![Build Status](https://travis-ci.org/sepetrov/ahocorasick.svg?branch=master)](https://travis-ci.org/sepetrov/ahocorasick)

A [Golang][1] implementation of the [Aho-Corasick string-searching algorithm][2].

## Benchmark

The benchmarks are using dictionaries in different languages and sizes. The test
files are with size 1 MB, 10 MB and 100 MB.

Here is the result of the benchmarks on my laptop. You can download the `testdata`
from [here][3]. 

```bash
go test -bench=. -benchtime=20s -run=XXX -timeout=20m
goos: darwin
goarch: amd64
pkg: github.com/sepetrov/ahocorasick
Benchmark/bg/indexing-4         	                  10	2677145567 ns/op
Benchmark/bg/searching_in_1_MB_text-4         	      30	 681869022 ns/op
Benchmark/bg/searching_in_10_MB_text-4        	       3	7629313139 ns/op
Benchmark/bg/searching_in_100_MB_text-4       	       1	82646577887 ns/op
Benchmark/de/indexing-4                       	      14	1868422845 ns/op
Benchmark/de/searching_in_1_MB_text-4         	      90	 291651064 ns/op
Benchmark/de/searching_in_10_MB_text-4        	       3	6808002538 ns/op
Benchmark/de/searching_in_100_MB_text-4       	       1	23298900173 ns/op
Benchmark/en/indexing-4                       	      22	 973961808 ns/op
Benchmark/en/searching_in_1_MB_text-4         	     198	 133871837 ns/op
Benchmark/en/searching_in_10_MB_text-4        	      18	1297788853 ns/op
Benchmark/en/searching_in_100_MB_text-4       	       2	12225270857 ns/op
Benchmark/ru/indexing-4                       	       7	3514695014 ns/op
Benchmark/ru/searching_in_1_MB_text-4         	      48	 634960317 ns/op
Benchmark/ru/searching_in_10_MB_text-4        	       4	7170965736 ns/op
Benchmark/ru/searching_in_100_MB_text-4       	       1	48899360327 ns/op
Benchmark/sv/indexing-4                       	      32	 836909979 ns/op
Benchmark/sv/searching_in_1_MB_text-4         	     100	 221659739 ns/op
Benchmark/sv/searching_in_10_MB_text-4        	      10	2527256037 ns/op
Benchmark/sv/searching_in_100_MB_text-4       	       1	29016405323 ns/op
Benchmark/zh/indexing-4                       	     130	 170518657 ns/op
Benchmark/zh/searching_in_1_MB_text-4         	     195	 112330014 ns/op
Benchmark/zh/searching_in_10_MB_text-4        	      22	 986999921 ns/op
Benchmark/zh/searching_in_100_MB_text-4       	       3	9637587461 ns/op
PASS
ok  	github.com/sepetrov/ahocorasick	1001.639s
```

## License

See [LICENSE](LICENSE).

[1]: https://golang.org
[2]: https://en.wikipedia.org/wiki/Aho–Corasick_algorithm
[3]: https://github.com/sepetrov/ahocorasick/releases/tag/v0.1.0