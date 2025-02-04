# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.7] - 2025-01-30 Thu

Add: scientificNameString is the only added column

## [v0.3.6] - 2025-01-27 Mon

Add: removed all additional fields except `scientificNameString`.

## [v0.3.5] - 2025-01-23 Thu

Add: update modules.

## [v0.3.4] - 2024-10-08 Tue

Add: make tree hierarchy be preferred over flat hierarchy.

## [v0.3.3] - 2024-10-06 Sun

Add: 'domain' to hierarchy vocabulary.

## [v0.3.2] - 2024-09-27 Fri

Fix: CSV extensions processed by RFC CSV path.
Fix: the line count during csv read.

## [v0.3.1] - 2024-09-26 Thu

Add: move BadRow from ent to gnfmt external module

## [v0.3.0] - 2024-09-26 Thu

Add [#25]: add options for skipping or processing rows with wrong fields number.
Add [#24]: make CSV default output format instead of TSV.

## [v0.2.9] - 2024-09-18 Wed

Add [#23]: read/write CSV files with empty FieldsEnclosedBy differntly from
"proper" CSV.

## [v0.2.8] - 2024-09-13 Fri

Add [#22]: support '|' as a field separator.

## [v0.2.7] - 2024-09-13 Fri

Fix [#20]: metadata and filename for eml.xml now match in output DwCA.

## [v0.2.6] - 2024-03-18 Mon

Add: try ZIP file type for unknown types.
Add: no normalization for already normalized.

## [v0.2.5] - 2024-03-15 Fri

Add [#18]: add download from URL.

## [v0.2.4] - 2024-03-14 Thu

Add: stream methods return number of records.

## [v0.2.3] - 2024-03-14 Thu

Add: log for extensions.

## [v0.2.2] - 2024-03-14 Thu

Fix: allow absense of taxonomic status.

## [v0.2.1] - 2024-03-13 Wed

Fix [#17]: field separator for CSV.

## [v0.2.0] - 2024-03-08 Fri

Add: improve interfaces.
Fix: remove debug info.

## [v0.1.1] - 2024-03-07 Thu

Add [#15]: fix headers in csv files.

## [v0.1.0] - 2024-02-29 Thu

Add [#14]: configuration YAML file, first binary release.
Add [#13]: normalize cli command.
Add [#12]: cli command.
Add [#11]: normalize DwCA and save into a file.
Add [#10]: find how sci-name, hierarchy, synonyms are expressed in a DwCA file.
Add [#7]: normalize synonyms.
Add [#6]: normalize hierarhcies.
Add [#5]: normalize scientific name.
Add [#4]: convenience meta object.

## [v0.0.1] - 2024-02-12 Mon

Add [#3]: read core and extension files.
Add [#2]: load eml.xml and meta.xml
Add [#1]: extract dwca files into temporary dir.

## [v0.0.0] - 2024-02-02 Fri

Add: initial commit

## Footnotes

This document follows [changelog guidelines]

[v0.3.4]: https://github.com/gnames/dwca/compare/v0.3.3...v0.3.4
[v0.3.3]: https://github.com/gnames/dwca/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/gnames/dwca/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/gnames/dwca/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/gnames/dwca/compare/v0.2.9...v0.3.0
[v0.2.9]: https://github.com/gnames/dwca/compare/v0.2.8...v0.2.9
[v0.2.8]: https://github.com/gnames/dwca/compare/v0.2.7...v0.2.8
[v0.2.7]: https://github.com/gnames/dwca/compare/v0.2.6...v0.2.7
[v0.2.6]: https://github.com/gnames/dwca/compare/v0.2.5...v0.2.6
[v0.2.5]: https://github.com/gnames/dwca/compare/v0.2.4...v0.2.5
[v0.2.4]: https://github.com/gnames/dwca/compare/v0.2.3...v0.2.4
[v0.2.3]: https://github.com/gnames/dwca/compare/v0.2.2...v0.2.3
[v0.2.2]: https://github.com/gnames/dwca/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/gnames/dwca/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/gnames/dwca/compare/v0.1.1...v0.2.0
[v0.1.1]: https://github.com/gnames/dwca/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/gnames/dwca/compare/v0.0.1...v0.1.0
[v0.0.1]: https://github.com/gnames/dwca/compare/v0.0.0...v0.0.1
[v0.0.0]: https://github.com/gnames/dwca/tree/v0.0.0
[#30]: https://github.com/gnames/dwca/issues/30
[#29]: https://github.com/gnames/dwca/issues/29
[#28]: https://github.com/gnames/dwca/issues/28
[#27]: https://github.com/gnames/dwca/issues/27
[#26]: https://github.com/gnames/dwca/issues/26
[#25]: https://github.com/gnames/dwca/issues/25
[#24]: https://github.com/gnames/dwca/issues/24
[#23]: https://github.com/gnames/dwca/issues/23
[#22]: https://github.com/gnames/dwca/issues/22
[#21]: https://github.com/gnames/dwca/issues/21
[#20]: https://github.com/gnames/dwca/issues/20
[#19]: https://github.com/gnames/dwca/issues/19
[#18]: https://github.com/gnames/dwca/issues/18
[#17]: https://github.com/gnames/dwca/issues/17
[#16]: https://github.com/gnames/dwca/issues/16
[#15]: https://github.com/gnames/dwca/issues/15
[#14]: https://github.com/gnames/dwca/issues/14
[#13]: https://github.com/gnames/dwca/issues/13
[#12]: https://github.com/gnames/dwca/issues/12
[#11]: https://github.com/gnames/dwca/issues/11
[#10]: https://github.com/gnames/dwca/issues/10
[#9]: https://github.com/gnames/dwca/issues/9
[#8]: https://github.com/gnames/dwca/issues/8
[#7]: https://github.com/gnames/dwca/issues/7
[#6]: https://github.com/gnames/dwca/issues/6
[#5]: https://github.com/gnames/dwca/issues/5
[#4]: https://github.com/gnames/dwca/issues/4
[#3]: https://github.com/gnames/dwca/issues/3
[#2]: https://github.com/gnames/dwca/issues/2
[#1]: https://github.com/gnames/dwca/issues/1
[changelog guidelines]: https://keepachangelog.com/en/1.0.0/
