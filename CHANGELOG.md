# Changelog

## [0.4.0](https://github.com/riadshalaby/agentinit/compare/v0.3.0...v0.4.0) (2026-04-11)


### Features

* **prompts:** inline critical workflow rules ([72ae946](https://github.com/riadshalaby/agentinit/commit/72ae946af72a9e7a43c0bc03d1d7a2277e16b0b4))
* **scaffold:** generate manifest for managed files ([416474c](https://github.com/riadshalaby/agentinit/commit/416474ca20d99a67e8cc85303b7312ffedd4a837))
* **update:** refresh managed scaffold files ([6d52ef6](https://github.com/riadshalaby/agentinit/commit/6d52ef697d2a46a38a68f8075b1be208afdefc58))
* **workflow:** merge agent rules into AGENTS.md ([030e8f1](https://github.com/riadshalaby/agentinit/commit/030e8f17b56954357595df5bc197f27ff2cca832))


### Bug Fixes

* **ai:** removed old agents config and commands ([7bebe7d](https://github.com/riadshalaby/agentinit/commit/7bebe7d818c30301d855e5e98eecbc3615f045ef))
* **workflow:** address scaffold merge regressions ([53e9301](https://github.com/riadshalaby/agentinit/commit/53e9301cdb68a9a44d2d99e592c192f3c984e9a3))


### Miscellaneous Chores

* **ai:** new claude settings ([24f9f5b](https://github.com/riadshalaby/agentinit/commit/24f9f5b592194333a0f7d7d96e8d87934817f244))
* **ai:** roadmap for v0.4.0 ([c16db5e](https://github.com/riadshalaby/agentinit/commit/c16db5e33aa9472c602d1938bd5dc08162ae30c4))
* **review:** T-001 PASS_WITH_NOTES — ready_for_test ([d073756](https://github.com/riadshalaby/agentinit/commit/d07375643af5133ee031dedb2659eb1a8db26c41))
* **review:** T-001 PASS_WITH_NOTES R2 — ready_for_test ([c6deb0c](https://github.com/riadshalaby/agentinit/commit/c6deb0cca1b491c4f3464745f7696121330bcf92))
* **review:** T-002 PASS_WITH_NOTES — ready_for_test ([4213480](https://github.com/riadshalaby/agentinit/commit/4213480978ee1dba3b622177563701cf242e697c))
* **review:** T-003 PASS — ready_for_test ([5fbc956](https://github.com/riadshalaby/agentinit/commit/5fbc956dc1da70587c57e657de32411a9693ea7b))
* **review:** T-004 PASS — ready_for_test ([4f740fd](https://github.com/riadshalaby/agentinit/commit/4f740fdeb7bc3fccb1a9410d5f6e61bb46692e74))
* **review:** T-005 PASS — ready_for_test ([441db2b](https://github.com/riadshalaby/agentinit/commit/441db2bf13b2547bbc6f03c7f06f63c14f682afc))
* start cycle mcp ([3c6ed8a](https://github.com/riadshalaby/agentinit/commit/3c6ed8ad39d47fa59a338c1a31c6c4a228d5efd0))
* **workflow:** validate restructured scaffold cycle ([8ce3ec1](https://github.com/riadshalaby/agentinit/commit/8ce3ec178e80515843b1e714c1bb6f075016aa07))

## [0.3.0](https://github.com/riadshalaby/agentinit/compare/v0.2.0...v0.3.0) (2026-04-09)


### Features

* **claude:** add tool preference guidance to CLAUDE files ([519fd40](https://github.com/riadshalaby/agentinit/commit/519fd40e3157c65157707f0346a3458f3621d363))
* **init:** add interactive setup wizard ([d887311](https://github.com/riadshalaby/agentinit/commit/d887311459cbc14fb750a381c6e271b325b883b5))
* **init:** add manual and auto workflow scaffolds ([c551b4f](https://github.com/riadshalaby/agentinit/commit/c551b4f6e1b2807823a37717e10b452df2b365a3))
* initial agentinit CLI scaffold generator ([e91f919](https://github.com/riadshalaby/agentinit/commit/e91f9193d1c4a6054b1b88cbd1d1347316e099ea))
* **mcp:** add MCP session management tools ([b345e81](https://github.com/riadshalaby/agentinit/commit/b345e811a553e7d09687b9b9810d29cf80df3a6e))
* **mcp:** add the stdio MCP server command ([88b7de2](https://github.com/riadshalaby/agentinit/commit/88b7de2f7142f291b2ee2d37a38c94b1e16c4547))
* **prereq:** add fd, bat, and jq to wizard setup checks ([665cb60](https://github.com/riadshalaby/agentinit/commit/665cb60e0d99e105ed9be2ccdea9c1196d945a6c))
* **prereq:** add optional advanced search tools ([29888b5](https://github.com/riadshalaby/agentinit/commit/29888b54b464a311d76bb4b5b9c85c7fd66f76d9))
* **prereq:** add platform-specific Claude and Codex installs ([b5c94d1](https://github.com/riadshalaby/agentinit/commit/b5c94d19839d58dc3a83754048e43c7425a15acb))
* **prompts:** add shared search-strategy guidance ([6cb7daf](https://github.com/riadshalaby/agentinit/commit/6cb7daf56a9b4954ea034ab781eb93f48688fc97))
* **scaffold:** always render PO artifacts ([50f5cf2](https://github.com/riadshalaby/agentinit/commit/50f5cf2629ac059347daa019cac50afccf71f724))
* **scaffold:** share init summaries across cli and wizard ([12d4b3f](https://github.com/riadshalaby/agentinit/commit/12d4b3f1b8da1eb841e57b3149bb3635644d0afe))
* **templates:** add the PO orchestration prompt and launcher ([3e30eb6](https://github.com/riadshalaby/agentinit/commit/3e30eb617aec29635e88928916e3bef9abbd5ed0))
* **workflow:** add cycle bootstrap pre-flight checks ([2164856](https://github.com/riadshalaby/agentinit/commit/21648563fc9b1cf41de40725fd776b2975091ac9))
* **workflow:** add shorthand commands and improvement roadmap ([1e04a49](https://github.com/riadshalaby/agentinit/commit/1e04a491be9828750324369e182a892ce8c55347))
* **workflow:** add the tester role and status flow ([dcd17e3](https://github.com/riadshalaby/agentinit/commit/dcd17e306979ab85bb68407408207e60c25307dd))
* **workflow:** categorize tools for the wizard and CLAUDE guidance ([156667e](https://github.com/riadshalaby/agentinit/commit/156667e1990a95dd58b601c924570db49163ebaa))
* **workflow:** split scaffolded agent rules into layered files ([0fa5780](https://github.com/riadshalaby/agentinit/commit/0fa5780c42e6f7d451750f6014bd6afd6d9cde45))


### Bug Fixes

* **agents:** address T-007 test failure ([25bfd1e](https://github.com/riadshalaby/agentinit/commit/25bfd1ebed77c51857266753cbc8f59b5e0e0759))
* **cli:** report build versions from Go module metadata ([9fc009c](https://github.com/riadshalaby/agentinit/commit/9fc009c6fcca68662a70bd97acde41bc41334bcb))
* **prereq:** address review findings ([cdf1bf6](https://github.com/riadshalaby/agentinit/commit/cdf1bf61d7b7638f4576e53e98930ac3ad5aa62d))
* **prereq:** address review findings ([4861a14](https://github.com/riadshalaby/agentinit/commit/4861a1438aa54e517c1c1a5e53e5e516bd77e2f2))
* **prereq:** install the tree-sitter CLI on macOS ([defe919](https://github.com/riadshalaby/agentinit/commit/defe91985dc35a2a45e1d27a0df45ab48a161467))
* **release:** add Windows binaries to GoReleaser artifacts ([e958c06](https://github.com/riadshalaby/agentinit/commit/e958c06c57cc3bb254a4846ad96ec31c86a766d4))
* **release:** use unprefixed release tags ([bd7e4cc](https://github.com/riadshalaby/agentinit/commit/bd7e4cc94cc73e0e18990c16f9cb0bd6263c2998))
* **scaffold:** address T-002 review findings ([cb3faab](https://github.com/riadshalaby/agentinit/commit/cb3faab9c710a34f4ec680028d9ad58657a90a13))
* **templates:** address review findings for the PO workflow ([1c09645](https://github.com/riadshalaby/agentinit/commit/1c096450031a2a1e4b4148f4e649ab9970e2f281))
* **workflow:** address review findings ([e2ca6fd](https://github.com/riadshalaby/agentinit/commit/e2ca6fdb81da188bd198ab931ebc38b4a378960c))
* **workflow:** address review findings ([76c17ac](https://github.com/riadshalaby/agentinit/commit/76c17ac5fa54d8b25155c4dcddc5a4097dffa67c))
* **workflow:** address review findings for cycle bootstrap ([af34bf7](https://github.com/riadshalaby/agentinit/commit/af34bf77e50c60d776ad49abff08f6c97cdf49d0))
* **workflow:** address tester review findings ([5584411](https://github.com/riadshalaby/agentinit/commit/558441158e9ff7597382462975df32cd53edfc69))
* **workflow:** correct the tester handoff status flow ([8d2f9f8](https://github.com/riadshalaby/agentinit/commit/8d2f9f839c59ce04bd78f191d47ac633556e0c8a))
* **workflow:** keep commit conventions in workflow-managed files ([405d588](https://github.com/riadshalaby/agentinit/commit/405d58875c0da32f5476f3364a9f90d8f7d17612))
* **workflow:** keep tester in manual scaffolds and ignore runtime reports ([89bc419](https://github.com/riadshalaby/agentinit/commit/89bc4191e224d3cd526d4d3272c4c075d6cdc5f1))


### Miscellaneous Chores

* added logo to README ([1473df3](https://github.com/riadshalaby/agentinit/commit/1473df385bbe68da33f8bf43281b2a995b327459))
* **ai:** close cycle ([288fe5e](https://github.com/riadshalaby/agentinit/commit/288fe5ebbd966ffed76a6923f3c5ef446a7014c6))
* **ai:** plan update ([9db7a94](https://github.com/riadshalaby/agentinit/commit/9db7a94523de86400e3ac9c28cfb2494a8074c11))
* **ai:** record review outcome ([64aa536](https://github.com/riadshalaby/agentinit/commit/64aa536248900d43ec8dadeb66e58c9e1d1098c3))
* **ai:** roadmap for v0.2.0 ([f047ada](https://github.com/riadshalaby/agentinit/commit/f047adae2f15b91c0f87c60f4d1463c265299502))
* **ai:** roadmap update ([e0fffbd](https://github.com/riadshalaby/agentinit/commit/e0fffbd829133b325ca439045fbdb3b58da9ad01))
* **hooks:** reject co-authored commit trailers ([e142dc0](https://github.com/riadshalaby/agentinit/commit/e142dc06eec5cc435a0d3198e623a4fb2032d074))
* **release:** add GitHub release builds with GoReleaser ([5d7714e](https://github.com/riadshalaby/agentinit/commit/5d7714e5b56ce3932ac0e469894c0b299c02ff92))
* **release:** add release-please automation ([76b1e4e](https://github.com/riadshalaby/agentinit/commit/76b1e4e90626fe4994c1d43608fe960197a6a8b6))
* remove references to non-existent gate-check scripts ([4e0418a](https://github.com/riadshalaby/agentinit/commit/4e0418a7398b841618bf2151bcbc9990fa148f55))
* remove unused scripts ([7bae5e1](https://github.com/riadshalaby/agentinit/commit/7bae5e1dd1c579851fe98dd450aede89ec54704d))
* **review:** advance T-001 to ready_for_test after passing review ([e1bf4a9](https://github.com/riadshalaby/agentinit/commit/e1bf4a95e2836da750a987346bcd8745f4eb7f23))
* **review:** advance T-001 to ready_for_test after passing review ([ecb1584](https://github.com/riadshalaby/agentinit/commit/ecb1584fda31be320759a6d85f40e3ee3871aa52))
* **review:** advance T-001 to ready_for_test after passing review ([593fb98](https://github.com/riadshalaby/agentinit/commit/593fb98319e037dbc296f36da477c46e0d711e08))
* **review:** advance T-002 to ready_for_test after passing re-review ([55a59ff](https://github.com/riadshalaby/agentinit/commit/55a59ff9a08e3b97c6c2faad350b8601e6c5b4cc))
* **review:** advance T-002 to ready_for_test after passing review ([373ad46](https://github.com/riadshalaby/agentinit/commit/373ad464227aba13534eb0a84f534c8f88d126a7))
* **review:** advance T-002 to ready_for_test after passing review ([6933610](https://github.com/riadshalaby/agentinit/commit/6933610e154e5e754c0d10d268cf7d063032b472))
* **review:** advance T-003 to ready_for_test after passing re-review ([2b02e44](https://github.com/riadshalaby/agentinit/commit/2b02e443b82eca6e1b1e7bd8e581eda00b6d6f65))
* **review:** advance T-003 to ready_for_test after passing review ([9bbd310](https://github.com/riadshalaby/agentinit/commit/9bbd3109f2cb18bd3f3eee0aaaf1f04a59aff2a7))
* **review:** advance T-003 to ready_for_test after passing review ([c4f02bb](https://github.com/riadshalaby/agentinit/commit/c4f02bb86d9d22e74b1aacaec7ae562361c61971))
* **review:** advance T-003 to ready_for_test after passing review ([834e7d2](https://github.com/riadshalaby/agentinit/commit/834e7d2fe7b391ffda50252278aa2442edaa73f4))
* **review:** advance T-004 to ready_for_test after passing review ([6cd9d31](https://github.com/riadshalaby/agentinit/commit/6cd9d319a6e567f9cfd9e0d5fe610de7fc9209a2))
* **review:** advance T-004 to ready_for_test after passing review ([467b851](https://github.com/riadshalaby/agentinit/commit/467b851494f29cf84080afe7f8358182af2982e5))
* **review:** advance T-005 to ready_for_test after passing review ([cc1f242](https://github.com/riadshalaby/agentinit/commit/cc1f2427d381ba9b3ab93cff2580418e4fecaa2f))
* **review:** advance T-006 to ready_for_test after passing review ([6b94b6d](https://github.com/riadshalaby/agentinit/commit/6b94b6d6f2a3e1ca0d74e12997f3913ef7ccc4b2))
* **review:** advance T-007 to ready_for_test after passing re-review ([6d418e1](https://github.com/riadshalaby/agentinit/commit/6d418e10eccd70f9ec29a19dcea8e553141960f7))
* **review:** advance T-007 to ready_for_test after passing review ([820f538](https://github.com/riadshalaby/agentinit/commit/820f53886da512a8e5871e53f5f426e37d9268d1))
* **review:** close cycle — all tasks done ([9304509](https://github.com/riadshalaby/agentinit/commit/9304509d30145f1057e88209cb78d0d741431c16))
* **review:** close cycle — all tasks done (T-001, T-002, T-003) ([4383751](https://github.com/riadshalaby/agentinit/commit/4383751f91cf7d3a8abbc4584bfad6fa69477012))
* **review:** log final review results ([6b504d2](https://github.com/riadshalaby/agentinit/commit/6b504d25f38f3ef657cefb4c57a74e7d97d832b8))
* **review:** record final review results ([1eea833](https://github.com/riadshalaby/agentinit/commit/1eea8335ff1948b911c9e1ad0400a6dead1a47ac))
* **review:** T-002 changes_requested — remove dead Workflow bridge in engine.go ([a6633f4](https://github.com/riadshalaby/agentinit/commit/a6633f4fa9aff39b5ba2a643d7aab428852cf7d1))
* **review:** T-003 changes_requested — README stale workflow flags and missing PO role ([dfbe0d9](https://github.com/riadshalaby/agentinit/commit/dfbe0d95e59a96b3afcc10b276cf9b5f3e0b7363))
* roadmap for next version ([c56dde0](https://github.com/riadshalaby/agentinit/commit/c56dde020023b7b302e776394f548665e3860d4d))
* roadmap for next version ([bd4a4b1](https://github.com/riadshalaby/agentinit/commit/bd4a4b113e91a33f8d36b55d8b2ffb4d3d9632e6))
* roadmap for next version - add PO workflow, ([4d90603](https://github.com/riadshalaby/agentinit/commit/4d90603bd19b36fb0f2eeea80548e134f37324a8))
* **scaffold:** remove redundant context guide ([bde40ab](https://github.com/riadshalaby/agentinit/commit/bde40abefdb2efced40fafd237fed9bbe74d6c36))
* start cycle automatic-agents-workflow ([1785d38](https://github.com/riadshalaby/agentinit/commit/1785d381e950fd0647823821d1c14f8fefb854d5))
* start cycle better-commands ([eb073f4](https://github.com/riadshalaby/agentinit/commit/eb073f498c027d7b15ab0ff0c395bdcd66a4c46e))
* start cycle cleanup ([350171d](https://github.com/riadshalaby/agentinit/commit/350171db943aa408fe6a6ec3879c633ecb2db934))
* start cycle hire-po-tester ([f25763a](https://github.com/riadshalaby/agentinit/commit/f25763a03b32c4ed3ed5c0039eead396af7ccefb))
* start cycle refine-auto ([5b9d195](https://github.com/riadshalaby/agentinit/commit/5b9d195fa0b2e0e0e252f4154f51333518538436))
* start cycle refine-workflow ([6beaee6](https://github.com/riadshalaby/agentinit/commit/6beaee6828e099c7e860489e1af42a62812fb1ef))
* start cycle single-mode ([2ecb9c6](https://github.com/riadshalaby/agentinit/commit/2ecb9c6ca79c8609cf2ccf4ee018386f87ffb3dd))
* **workflow:** disallow Co-Authored-By trailers in commit messages ([64a0e9d](https://github.com/riadshalaby/agentinit/commit/64a0e9dff43f6605a47ae806ae408d00564bcd56))
* **workflow:** finalize cycle task board ([af4323c](https://github.com/riadshalaby/agentinit/commit/af4323cd2ea5a734c672e2a4f84fef8b62f0ab64))
* **workflow:** record final review completion ([9a5c2c9](https://github.com/riadshalaby/agentinit/commit/9a5c2c9d2b7e53d6be751b133d140b220da4db7f))


### Documentation

* **agents:** promote hard rules to the top ([dd06775](https://github.com/riadshalaby/agentinit/commit/dd067753a702e62a07a23bdc02b18f447b94c986))
* **agents:** sync project files with scaffold templates ([1946dac](https://github.com/riadshalaby/agentinit/commit/1946dac06402f43c30857ee70af4a788ab47e37f))
* **handoff:** standardize handoff entry format ([eb01a1b](https://github.com/riadshalaby/agentinit/commit/eb01a1b136027813b69aff1a562feeff578f4424))
* **readme:** clarify manual, auto, and MCP workflow usage ([881d06b](https://github.com/riadshalaby/agentinit/commit/881d06b8dfb764a310cca9fd7981fabe313ae8b8))
* **readme:** describe interactive init wizard ([5e879e7](https://github.com/riadshalaby/agentinit/commit/5e879e74246f7631a25a3b68d824c50eea10a1a6))
* **readme:** improve onboarding clarity and extensibility for future workflows ([9bbd7f3](https://github.com/riadshalaby/agentinit/commit/9bbd7f3ee525ebb3040507a9fa1f60f18a3685c0))
* **scaffold:** address T-003 review findings ([6f1707d](https://github.com/riadshalaby/agentinit/commit/6f1707d7ca11b173f1cdfb1a4c17b075653ffe1a))
* **scaffold:** explain manual and auto runtime modes ([4e57566](https://github.com/riadshalaby/agentinit/commit/4e57566a265d48b65a9ebbe044cd78819ed05e85))
* **workflow:** add persistent session examples and coverage ([3164132](https://github.com/riadshalaby/agentinit/commit/31641322754b1b565a565e8435ecc15b2b52680b))
* **workflow:** align repo instruction files with layered layout ([1407cd3](https://github.com/riadshalaby/agentinit/commit/1407cd366ce1ad1c321adfc933606a606baa06ec))
* **workflow:** allow reviewer commits for review artifacts ([f136a0b](https://github.com/riadshalaby/agentinit/commit/f136a0b12d846e714163b9c4d7630ac915bbf626))
* **workflow:** define persistent session state transitions ([1920e81](https://github.com/riadshalaby/agentinit/commit/1920e812e52d7ebcbbf1e59e7b7f44be3890038e))
* **workflow:** document review rework flow ([e9b0748](https://github.com/riadshalaby/agentinit/commit/e9b07484b3b7391a32b3d8837e8069ed96b94a9f))
* **workflow:** record T-004 review and installer docs ([11585b0](https://github.com/riadshalaby/agentinit/commit/11585b0ea457ef4ec10a2d7b1e4d95eda94743ec))
* **workflow:** switch generated guidance to persistent agent sessions ([9dea910](https://github.com/riadshalaby/agentinit/commit/9dea91090083acd79570988068b1edc6da5caadf))
