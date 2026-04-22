# Changelog

## [0.8.3](https://github.com/riadshalaby/agentinit/compare/v0.8.2...v0.8.3) (2026-04-22)


### chore

* **ai:** close cycle ([b1977db](https://github.com/riadshalaby/agentinit/commit/b1977db912677bf34adaa9ceb65a0d79b464dedd))


### Bug Fixes

* **prompts:** align PO template with current MCP polling flow ([e569cae](https://github.com/riadshalaby/agentinit/commit/e569cae6d0fc0c248eb8e75ada101ed45fcc1b0f))
* **prompts:** align reviewer and handoff commit flow ([713968f](https://github.com/riadshalaby/agentinit/commit/713968fe1ab9b1247bf3c0ec632c7f90999ca43e))
* **prompts:** append cycle-close handoff entries ([2886bce](https://github.com/riadshalaby/agentinit/commit/2886bce5d76252a99e932ceb9d81e6ec18136eda))
* **prompts:** make implementer workflow test-first and adaptive ([58c4077](https://github.com/riadshalaby/agentinit/commit/58c4077262a4987b58ef3d31bb5129605c36ed3e))
* **prompts:** make reviewer verification mandatory ([a5a2c86](https://github.com/riadshalaby/agentinit/commit/a5a2c8654bb98a03a0a15096476b7a98164d1c22))
* **prompts:** preserve commit_task WIP commit messages ([b01abc5](https://github.com/riadshalaby/agentinit/commit/b01abc5e764cfecbe78effc8d591069ea011f9ac))
* **prompts:** remove WIP commits from implementer flow ([ef53683](https://github.com/riadshalaby/agentinit/commit/ef53683c6c9783d2c3f25ce1e43cd2244f0edc8b))
* **update:** make self-update checks catch managed drift ([0ba35d9](https://github.com/riadshalaby/agentinit/commit/0ba35d9496b9e00869759e928aa665c5c5f2c801))

## [0.8.2](https://github.com/riadshalaby/agentinit/compare/v0.8.1...v0.8.2) (2026-04-21)


### chore

* **ai:** close cycle ([ec1f3e9](https://github.com/riadshalaby/agentinit/commit/ec1f3e922c30ea280fca563af05cf922d0238c6b))


### Features

* **config:** default codex implementer effort to high ([289e8bf](https://github.com/riadshalaby/agentinit/commit/289e8bf8ffaece6a554ad44013e52aff46d6b34e))
* **update:** run tool checks after refreshing files ([b8c2c7e](https://github.com/riadshalaby/agentinit/commit/b8c2c7e0416d031f29ed3576b8cf286bddc17945))
* **wizard:** require git before scaffolding ([4db640e](https://github.com/riadshalaby/agentinit/commit/4db640ee7eaca75a6c2db8911fa49feb9b777c55))


### Bug Fixes

* **pr:** skip aide pr when no remote is configured ([86771da](https://github.com/riadshalaby/agentinit/commit/86771daffcbb59828b6312d107408226e8371e36))


### Documentation

* **readme:** add PATH setup after go install ([5833bd7](https://github.com/riadshalaby/agentinit/commit/5833bd762673322f8582692b0faf92ccd1d64591))

## [0.8.1](https://github.com/riadshalaby/agentinit/compare/v0.8.0...v0.8.1) (2026-04-18)


### Bug Fixes

* **mcp:** add structured session run results ([3eebbec](https://github.com/riadshalaby/agentinit/commit/3eebbec7413085d8f702470be6bd8d12e59ea72b))
* **mcp:** cap session output chunks ([94fb93e](https://github.com/riadshalaby/agentinit/commit/94fb93e155617f8ba3643dc739c8bfc86a4bb779))
* **mcp:** resume claude session runs ([9b910c3](https://github.com/riadshalaby/agentinit/commit/9b910c35896f3a76edf221713a47877866fdfe58))
* **po:** default the coordinator to cheaper models ([6861a20](https://github.com/riadshalaby/agentinit/commit/6861a20e0f49bc65b2d41d58e50c7e34033fdaa5))

## [0.8.0](https://github.com/riadshalaby/agentinit/compare/v0.7.2...v0.8.0) (2026-04-18)


### chore

* **ai:** close cycle ([6b73c6c](https://github.com/riadshalaby/agentinit/commit/6b73c6c0d27e8d9a01fadb01a612c22859d8704c))


### Features

* **cli:** add cross-platform cycle bootstrap command ([458ea82](https://github.com/riadshalaby/agentinit/commit/458ea822d65b05e4d1d3e5dbae165a292fed4cfd))
* **cli:** add cross-platform po launcher ([cc86b4a](https://github.com/riadshalaby/agentinit/commit/cc86b4a97902773b2950826c9a7b01ae53175f85))
* **cli:** add cross-platform role launch commands ([e83ec6b](https://github.com/riadshalaby/agentinit/commit/e83ec6b559b8513bf4d18f11cf1b73a931b89d4b))
* **cli:** add cycle close and pull request commands ([553560c](https://github.com/riadshalaby/agentinit/commit/553560c032eb46c066bb00d6061b5df2048d6743))
* **cli:** rename the agent binary to aide ([488ca13](https://github.com/riadshalaby/agentinit/commit/488ca1396dd2516dc73ef9ae6f93426ed9dfac29))
* **scaffold:** replace generated shell scripts with agentinit commands ([a2ea253](https://github.com/riadshalaby/agentinit/commit/a2ea2530a3044eed0c3b6e73204136d39ed66eeb))


### Bug Fixes

* **e2e:** restore tagged session manager test build ([84ac868](https://github.com/riadshalaby/agentinit/commit/84ac868f513310c90e36cb2b083e5343c0b0ec28))
* **mcp:** address review findings for session lifecycle context ([a232964](https://github.com/riadshalaby/agentinit/commit/a232964f78cf603d80751fa6f8146656bc314776))
* **mcp:** keep role model settings provider-aware ([0a7303e](https://github.com/riadshalaby/agentinit/commit/0a7303ebf651a9c1a6ac252182bd6a556716f353))
* **template:** broaden generated go and git tool permissions ([e380226](https://github.com/riadshalaby/agentinit/commit/e3802267b04721b34e42085c6d512b742b717244))
* **update:** reconcile desired-only managed files that already exist on disk ([3bca8f9](https://github.com/riadshalaby/agentinit/commit/3bca8f90c10fd1afca8ce275e301399776930164))

## [0.7.2](https://github.com/riadshalaby/agentinit/compare/v0.7.1...v0.7.2) (2026-04-16)


### Features

* **template:** add MCP permission and autoUpdatesChannel to scaffolded Claude settings ([98af946](https://github.com/riadshalaby/agentinit/commit/98af9460942c8dde349be92fa7d24ca616ff2966))


### Documentation

* **workflow:** finish_cycle amends HEAD when working tree is clean ([899fbfd](https://github.com/riadshalaby/agentinit/commit/899fbfd39c0c3c6cd073e02486b505d7ea65bff6))

## [0.7.1](https://github.com/riadshalaby/agentinit/compare/v0.7.0...v0.7.1) (2026-04-15)


### Features

* **mcp:** stream session output asynchronously ([553aab6](https://github.com/riadshalaby/agentinit/commit/553aab6bd1563cda76691cad464d8454f6135f78))
* **scaffold:** configure Claude MCP server ([c6bfd04](https://github.com/riadshalaby/agentinit/commit/c6bfd04764c5a4ab993e7f1902d9c8b7f2f92a46))
* **scaffold:** initialize projects on main by default ([c6ed145](https://github.com/riadshalaby/agentinit/commit/c6ed14531047b795d5b7ba78c3b2eefaeebe79f6))


### Miscellaneous Chores

* **ai:** close cycle ([c01a047](https://github.com/riadshalaby/agentinit/commit/c01a047ec94619407f5e515e5b5f8bb210baf439))
* **ai:** roadmap 0.7.1 ([f91e077](https://github.com/riadshalaby/agentinit/commit/f91e0779547c8d3c235d8d5454ed436c96c5ccfe))
* **ai:** update manifest ([e412b1d](https://github.com/riadshalaby/agentinit/commit/e412b1d8a09e2101a2a2996df50859341ca2f192))
* start cycle 0.7.1 ([5d2db83](https://github.com/riadshalaby/agentinit/commit/5d2db83710a2f3bd098d45510537011260967e73))

## [0.7.0](https://github.com/riadshalaby/agentinit/compare/v0.6.2...v0.7.0) (2026-04-14)


### Features

* **mcp:** add named session manager lifecycle ([b165a6c](https://github.com/riadshalaby/agentinit/commit/b165a6cb4fc1a3bc5a5c2640d7287e4182c90222))
* **mcp:** add persistent session store ([7c14bf1](https://github.com/riadshalaby/agentinit/commit/7c14bf19439b8ef69185c71d65b5045c164dc618))
* **mcp:** add provider adapters for codex and claude ([7dffe4d](https://github.com/riadshalaby/agentinit/commit/7dffe4d16532477eb2cf491285f0f05a535a595f))
* **mcp:** add typed config loading and provider validation ([806ff08](https://github.com/riadshalaby/agentinit/commit/806ff0870a7598b2c84143088a813f6a0328192e))
* **mcp:** wire real named-session MCP tools ([d370b61](https://github.com/riadshalaby/agentinit/commit/d370b61dcb63aa089ba162ab936d62e91c4983f0))


### Bug Fixes

* **ai:** manual scripts fixed ([f38e081](https://github.com/riadshalaby/agentinit/commit/f38e081826f6aabd314f7a2dcc9b0227260fe7ab))


### Miscellaneous Chores

* **ai:** close cycle ([ee434e9](https://github.com/riadshalaby/agentinit/commit/ee434e9b2bc0e5ac5dddb1be8983e3fa7b8ec495))
* **ai:** close cycle ([2ecce2d](https://github.com/riadshalaby/agentinit/commit/2ecce2d4a454dbab2ef4071c607f0238fc027c6e))
* **ai:** start cycle 0.7.0 ([53c2ac5](https://github.com/riadshalaby/agentinit/commit/53c2ac500887c7578d5534d8c0c9362bfd2157bd))
* start cycle 0.7.0 ([ab57e57](https://github.com/riadshalaby/agentinit/commit/ab57e5703d4165ea31fd4b3696ec4baa70248eaf))


### Documentation

* **mcp:** update scaffold prompts for named sessions ([7fc204f](https://github.com/riadshalaby/agentinit/commit/7fc204f07ae36aac4900ba43d11b8d145e40bfa9))

## [0.6.2](https://github.com/riadshalaby/agentinit/compare/v0.6.2...v0.6.2) (2026-04-14)


### Features

* **ai:** removed agents from tasks board ([668b150](https://github.com/riadshalaby/agentinit/commit/668b150d30497ecbdf5c759f12120128a6f8fa55))
* **claude:** add scaffolded tool-access permissions by project type ([293b8be](https://github.com/riadshalaby/agentinit/commit/293b8be6a98a0ce4f5603f7add518f168ab3b776))
* **claude:** add tool preference guidance to CLAUDE files ([519fd40](https://github.com/riadshalaby/agentinit/commit/519fd40e3157c65157707f0346a3458f3621d363))
* **claude:** scaffold Claude settings templates ([93580d3](https://github.com/riadshalaby/agentinit/commit/93580d38b2434bdd7e88b8b81df41b0cbfdecc9b))
* **config:** scaffold per-role agent, model, and effort defaults in .ai/config.json ([674143b](https://github.com/riadshalaby/agentinit/commit/674143b639c0c9fc2eb4b2354c7851a259619783))
* **init:** add interactive setup wizard ([d887311](https://github.com/riadshalaby/agentinit/commit/d887311459cbc14fb750a381c6e271b325b883b5))
* **init:** add manual and auto workflow scaffolds ([c551b4f](https://github.com/riadshalaby/agentinit/commit/c551b4f6e1b2807823a37717e10b452df2b365a3))
* initial agentinit CLI scaffold generator ([e91f919](https://github.com/riadshalaby/agentinit/commit/e91f9193d1c4a6054b1b88cbd1d1347316e099ea))
* **mcp:** add MCP session management tools ([b345e81](https://github.com/riadshalaby/agentinit/commit/b345e811a553e7d09687b9b9810d29cf80df3a6e))
* **mcp:** add the stdio MCP server command ([88b7de2](https://github.com/riadshalaby/agentinit/commit/88b7de2f7142f291b2ee2d37a38c94b1e16c4547))
* **mcp:** poll session output with get_output ([52d69ff](https://github.com/riadshalaby/agentinit/commit/52d69ffb0999f26248dc776f763c13e752562134))
* **mcp:** support codex role sessions across MCP commands ([f4e85e8](https://github.com/riadshalaby/agentinit/commit/f4e85e8ec44b2fd6b368cf46229e1b6fcf4fee4b))
* **mcp:** write MCP server debug logs to .ai/mcp-server.log ([bcc11f3](https://github.com/riadshalaby/agentinit/commit/bcc11f3ff06e2cddba92be8e086aadb719590c07))
* **planner:** add roadmap refinement guidance before start_plan ([360a0e7](https://github.com/riadshalaby/agentinit/commit/360a0e78320f4cebb0a61bbda0310a1f05251f44))
* **po:** support codex and validate ai-po agent selection ([63f97cf](https://github.com/riadshalaby/agentinit/commit/63f97cf5a329e40412dc772588cd944de3e3fadb))
* **prereq:** add fd, bat, and jq to wizard setup checks ([665cb60](https://github.com/riadshalaby/agentinit/commit/665cb60e0d99e105ed9be2ccdea9c1196d945a6c))
* **prereq:** add optional advanced search tools ([29888b5](https://github.com/riadshalaby/agentinit/commit/29888b54b464a311d76bb4b5b9c85c7fd66f76d9))
* **prereq:** add platform-specific Claude and Codex installs ([b5c94d1](https://github.com/riadshalaby/agentinit/commit/b5c94d19839d58dc3a83754048e43c7425a15acb))
* **prompts:** add shared search-strategy guidance ([6cb7daf](https://github.com/riadshalaby/agentinit/commit/6cb7daf56a9b4954ea034ab781eb93f48688fc97))
* **prompts:** inline critical workflow rules ([72ae946](https://github.com/riadshalaby/agentinit/commit/72ae946af72a9e7a43c0bc03d1d7a2277e16b0b4))
* **scaffold:** always render PO artifacts ([50f5cf2](https://github.com/riadshalaby/agentinit/commit/50f5cf2629ac059347daa019cac50afccf71f724))
* **scaffold:** generate manifest for managed files ([416474c](https://github.com/riadshalaby/agentinit/commit/416474ca20d99a67e8cc85303b7312ffedd4a837))
* **scaffold:** share init summaries across cli and wizard ([12d4b3f](https://github.com/riadshalaby/agentinit/commit/12d4b3f1b8da1eb841e57b3149bb3635644d0afe))
* **templates:** add the PO orchestration prompt and launcher ([3e30eb6](https://github.com/riadshalaby/agentinit/commit/3e30eb617aec29635e88928916e3bef9abbd5ed0))
* **update:** migrate legacy workflow files during scaffold refresh ([95aee2c](https://github.com/riadshalaby/agentinit/commit/95aee2c91f9bd16a6e27db0530e22559a7369ab8))
* **update:** refresh managed scaffold files ([6d52ef6](https://github.com/riadshalaby/agentinit/commit/6d52ef697d2a46a38a68f8075b1be208afdefc58))
* **workflow:** add a ready-to-commit stage to task flow ([a5b126e](https://github.com/riadshalaby/agentinit/commit/a5b126e1d7d1104e2ec6871668e0fccf2434fbfe))
* **workflow:** add cycle bootstrap pre-flight checks ([2164856](https://github.com/riadshalaby/agentinit/commit/21648563fc9b1cf41de40725fd776b2975091ac9))
* **workflow:** add shorthand commands and improvement roadmap ([1e04a49](https://github.com/riadshalaby/agentinit/commit/1e04a491be9828750324369e182a892ce8c55347))
* **workflow:** add the tester role and status flow ([dcd17e3](https://github.com/riadshalaby/agentinit/commit/dcd17e306979ab85bb68407408207e60c25307dd))
* **workflow:** categorize tools for the wizard and CLAUDE guidance ([156667e](https://github.com/riadshalaby/agentinit/commit/156667e1990a95dd58b601c924570db49163ebaa))
* **workflow:** merge agent rules into AGENTS.md ([030e8f1](https://github.com/riadshalaby/agentinit/commit/030e8f17b56954357595df5bc197f27ff2cca832))
* **workflow:** move finish_cycle to the implementer role ([89beec4](https://github.com/riadshalaby/agentinit/commit/89beec417a399bbfc91af67897ea35a1f620f694))
* **workflow:** remove the tester role from scaffolded projects ([5f30e3e](https://github.com/riadshalaby/agentinit/commit/5f30e3e506ecd778281f5ec8a4262a77b44fd869))
* **workflow:** require fresh file reads for every role command ([6756286](https://github.com/riadshalaby/agentinit/commit/6756286e22ca4c98883256684fbec75a040b321a))
* **workflow:** split scaffolded agent rules into layered files ([0fa5780](https://github.com/riadshalaby/agentinit/commit/0fa5780c42e6f7d451750f6014bd6afd6d9cde45))
* **workflow:** track cycle review and test logs in git ([3847efe](https://github.com/riadshalaby/agentinit/commit/3847efe7e4f87ac7868835859f5f8d8a08a9be9d))


### Bug Fixes

* **agents:** address T-007 test failure ([25bfd1e](https://github.com/riadshalaby/agentinit/commit/25bfd1ebed77c51857266753cbc8f59b5e0e0759))
* **ai:** removed old agents config and commands ([7bebe7d](https://github.com/riadshalaby/agentinit/commit/7bebe7d818c30301d855e5e98eecbc3615f045ef))
* **cli:** report build versions from Go module metadata ([9fc009c](https://github.com/riadshalaby/agentinit/commit/9fc009c6fcca68662a70bd97acde41bc41334bcb))
* **mcp:** force-stop hung sessions with SIGKILL ([02a7882](https://github.com/riadshalaby/agentinit/commit/02a7882b7aa6cfd564abc7f2df39597de08a72f2))
* **mcp:** preserve structured JSON tool results ([7bea686](https://github.com/riadshalaby/agentinit/commit/7bea68603cf15be79fb17a1b352c793ef36711f7))
* **mcp:** wait longer before cutting off MCP session output ([796bc62](https://github.com/riadshalaby/agentinit/commit/796bc62a90e46abc7581db4d90318f630e07fdff))
* **prereq:** address review findings ([cdf1bf6](https://github.com/riadshalaby/agentinit/commit/cdf1bf61d7b7638f4576e53e98930ac3ad5aa62d))
* **prereq:** address review findings ([4861a14](https://github.com/riadshalaby/agentinit/commit/4861a1438aa54e517c1c1a5e53e5e516bd77e2f2))
* **prereq:** install the tree-sitter CLI on macOS ([defe919](https://github.com/riadshalaby/agentinit/commit/defe91985dc35a2a45e1d27a0df45ab48a161467))
* **release:** add Windows binaries to GoReleaser artifacts ([e958c06](https://github.com/riadshalaby/agentinit/commit/e958c06c57cc3bb254a4846ad96ec31c86a766d4))
* **release:** build release assets from release-please tags ([dd9e50c](https://github.com/riadshalaby/agentinit/commit/dd9e50cacd0d818d9278de798abff404649282f7))
* **release:** use unprefixed release tags ([bd7e4cc](https://github.com/riadshalaby/agentinit/commit/bd7e4cc94cc73e0e18990c16f9cb0bd6263c2998))
* **scaffold:** address T-002 review findings ([cb3faab](https://github.com/riadshalaby/agentinit/commit/cb3faab9c710a34f4ec680028d9ad58657a90a13))
* **templates:** address review findings for the PO workflow ([1c09645](https://github.com/riadshalaby/agentinit/commit/1c096450031a2a1e4b4148f4e649ab9970e2f281))
* **workflow:** address review findings ([e2ca6fd](https://github.com/riadshalaby/agentinit/commit/e2ca6fdb81da188bd198ab931ebc38b4a378960c))
* **workflow:** address review findings ([76c17ac](https://github.com/riadshalaby/agentinit/commit/76c17ac5fa54d8b25155c4dcddc5a4097dffa67c))
* **workflow:** address review findings for cycle bootstrap ([af34bf7](https://github.com/riadshalaby/agentinit/commit/af34bf77e50c60d776ad49abff08f6c97cdf49d0))
* **workflow:** address scaffold merge regressions ([53e9301](https://github.com/riadshalaby/agentinit/commit/53e9301cdb68a9a44d2d99e592c192f3c984e9a3))
* **workflow:** address tester review findings ([5584411](https://github.com/riadshalaby/agentinit/commit/558441158e9ff7597382462975df32cd53edfc69))
* **workflow:** correct the tester handoff status flow ([8d2f9f8](https://github.com/riadshalaby/agentinit/commit/8d2f9f839c59ce04bd78f191d47ac633556e0c8a))
* **workflow:** keep commit conventions in workflow-managed files ([405d588](https://github.com/riadshalaby/agentinit/commit/405d58875c0da32f5476f3364a9f90d8f7d17612))
* **workflow:** keep tester in manual scaffolds and ignore runtime reports ([89bc419](https://github.com/riadshalaby/agentinit/commit/89bc4191e224d3cd526d4d3272c4c075d6cdc5f1))


### Miscellaneous Chores

* added logo to README ([1473df3](https://github.com/riadshalaby/agentinit/commit/1473df385bbe68da33f8bf43281b2a995b327459))
* **ai:** close cycle ([fe1d633](https://github.com/riadshalaby/agentinit/commit/fe1d6335db12440adbbcd52c6f2197341fb96be7))
* **ai:** close cycle ([5327cda](https://github.com/riadshalaby/agentinit/commit/5327cda465b42c1f05f415475630b412e6f3b549))
* **ai:** close cycle ([288fe5e](https://github.com/riadshalaby/agentinit/commit/288fe5ebbd966ffed76a6923f3c5ef446a7014c6))
* **ai:** new claude settings ([24f9f5b](https://github.com/riadshalaby/agentinit/commit/24f9f5b592194333a0f7d7d96e8d87934817f244))
* **ai:** plan update ([9db7a94](https://github.com/riadshalaby/agentinit/commit/9db7a94523de86400e3ac9c28cfb2494a8074c11))
* **ai:** record review outcome ([64aa536](https://github.com/riadshalaby/agentinit/commit/64aa536248900d43ec8dadeb66e58c9e1d1098c3))
* **ai:** roadmap for v0.2.0 ([f047ada](https://github.com/riadshalaby/agentinit/commit/f047adae2f15b91c0f87c60f4d1463c265299502))
* **ai:** roadmap for v0.4.0 ([c16db5e](https://github.com/riadshalaby/agentinit/commit/c16db5e33aa9472c602d1938bd5dc08162ae30c4))
* **ai:** roadmap for version 0.5.0 ([960ce0a](https://github.com/riadshalaby/agentinit/commit/960ce0aa7495fc54ab595246a47461bfe611bce8))
* **ai:** roadmap for version 0.5.1 ([349366a](https://github.com/riadshalaby/agentinit/commit/349366a7a793833393c3557d154fead3012ee098))
* **ai:** roadmap for version 0.7.0 ([a327cf1](https://github.com/riadshalaby/agentinit/commit/a327cf15b1133014fca585ca6355c9b269aa9301))
* **ai:** roadmap update ([e0fffbd](https://github.com/riadshalaby/agentinit/commit/e0fffbd829133b325ca439045fbdb3b58da9ad01))
* **claude:** allow more tools for claude ([a934a54](https://github.com/riadshalaby/agentinit/commit/a934a5452d988e027892c494e19cc3e1896b7af3))
* **docs:** new logo ([2a27b4b](https://github.com/riadshalaby/agentinit/commit/2a27b4bda9c3ab331ef88cc67d20cc7f7e1103a1))
* **hooks:** reject co-authored commit trailers ([e142dc0](https://github.com/riadshalaby/agentinit/commit/e142dc06eec5cc435a0d3198e623a4fb2032d074))
* **main:** release 0.3.0 ([11eef89](https://github.com/riadshalaby/agentinit/commit/11eef8970e2ad2eb3c6b5c9becba73b4c7c25390))
* **main:** release 0.4.0 ([e5e6b35](https://github.com/riadshalaby/agentinit/commit/e5e6b35b435a4b55ed741b1e32631afe6f9cde4c))
* **main:** release 0.5.0 ([d887207](https://github.com/riadshalaby/agentinit/commit/d88720778b47a538f07d693437031576c7fd9896))
* **main:** release 0.6.0 ([550ec45](https://github.com/riadshalaby/agentinit/commit/550ec4598c8d6ca107f2b0a7c471b32a6cf405a3))
* **main:** release 0.6.1 ([fa393f4](https://github.com/riadshalaby/agentinit/commit/fa393f46939272b71acd315d85ae6a6634ff3b39))
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
* **review:** T-001 PASS_WITH_NOTES — ready_for_test ([d073756](https://github.com/riadshalaby/agentinit/commit/d07375643af5133ee031dedb2659eb1a8db26c41))
* **review:** T-001 PASS_WITH_NOTES R2 — ready_for_test ([c6deb0c](https://github.com/riadshalaby/agentinit/commit/c6deb0cca1b491c4f3464745f7696121330bcf92))
* **review:** T-002 changes_requested — remove dead Workflow bridge in engine.go ([a6633f4](https://github.com/riadshalaby/agentinit/commit/a6633f4fa9aff39b5ba2a643d7aab428852cf7d1))
* **review:** T-002 PASS_WITH_NOTES — ready_for_test ([4213480](https://github.com/riadshalaby/agentinit/commit/4213480978ee1dba3b622177563701cf242e697c))
* **review:** T-003 changes_requested — README stale workflow flags and missing PO role ([dfbe0d9](https://github.com/riadshalaby/agentinit/commit/dfbe0d95e59a96b3afcc10b276cf9b5f3e0b7363))
* **review:** T-003 PASS — ready_for_test ([5fbc956](https://github.com/riadshalaby/agentinit/commit/5fbc956dc1da70587c57e657de32411a9693ea7b))
* **review:** T-004 PASS — ready_for_test ([4f740fd](https://github.com/riadshalaby/agentinit/commit/4f740fdeb7bc3fccb1a9410d5f6e61bb46692e74))
* **review:** T-005 PASS — ready_for_test ([441db2b](https://github.com/riadshalaby/agentinit/commit/441db2bf13b2547bbc6f03c7f06f63c14f682afc))
* roadmap for next version ([c56dde0](https://github.com/riadshalaby/agentinit/commit/c56dde020023b7b302e776394f548665e3860d4d))
* roadmap for next version ([bd4a4b1](https://github.com/riadshalaby/agentinit/commit/bd4a4b113e91a33f8d36b55d8b2ffb4d3d9632e6))
* roadmap for next version - add PO workflow, ([4d90603](https://github.com/riadshalaby/agentinit/commit/4d90603bd19b36fb0f2eeea80548e134f37324a8))
* **scaffold:** remove redundant context guide ([bde40ab](https://github.com/riadshalaby/agentinit/commit/bde40abefdb2efced40fafd237fed9bbe74d6c36))
* start cycle 0.5.0 ([f1eabc5](https://github.com/riadshalaby/agentinit/commit/f1eabc59c5859148833797d665d112684bcb0a22))
* start cycle 0.6.2 ([670dcfd](https://github.com/riadshalaby/agentinit/commit/670dcfda22ef4502f04c5cf1fc3ba70e75f04ac1))
* start cycle automatic-agents-workflow ([1785d38](https://github.com/riadshalaby/agentinit/commit/1785d381e950fd0647823821d1c14f8fefb854d5))
* start cycle better-commands ([eb073f4](https://github.com/riadshalaby/agentinit/commit/eb073f498c027d7b15ab0ff0c395bdcd66a4c46e))
* start cycle cleanup ([350171d](https://github.com/riadshalaby/agentinit/commit/350171db943aa408fe6a6ec3879c633ecb2db934))
* start cycle hire-po-tester ([f25763a](https://github.com/riadshalaby/agentinit/commit/f25763a03b32c4ed3ed5c0039eead396af7ccefb))
* start cycle mcp ([3c6ed8a](https://github.com/riadshalaby/agentinit/commit/3c6ed8ad39d47fa59a338c1a31c6c4a228d5efd0))
* start cycle refine-auto ([5b9d195](https://github.com/riadshalaby/agentinit/commit/5b9d195fa0b2e0e0e252f4154f51333518538436))
* start cycle refine-workflow ([6beaee6](https://github.com/riadshalaby/agentinit/commit/6beaee6828e099c7e860489e1af42a62812fb1ef))
* start cycle single-mode ([2ecb9c6](https://github.com/riadshalaby/agentinit/commit/2ecb9c6ca79c8609cf2ccf4ee018386f87ffb3dd))
* start cycle v0.5.1 ([3f161ec](https://github.com/riadshalaby/agentinit/commit/3f161ec9545089172991d1e7db990626d427e274))
* start cycle v0.7.0 ([51bad16](https://github.com/riadshalaby/agentinit/commit/51bad16f47d3c5a8c99703af21a06067de0d9e3c))
* **workflow:** close the v0.5.1 cycle ([c2506d3](https://github.com/riadshalaby/agentinit/commit/c2506d3f03bf7c25d32e10232f9344e7bdde3097))
* **workflow:** disallow Co-Authored-By trailers in commit messages ([64a0e9d](https://github.com/riadshalaby/agentinit/commit/64a0e9dff43f6605a47ae806ae408d00564bcd56))
* **workflow:** finalize cycle task board ([af4323c](https://github.com/riadshalaby/agentinit/commit/af4323cd2ea5a734c672e2a4f84fef8b62f0ab64))
* **workflow:** record final review completion ([9a5c2c9](https://github.com/riadshalaby/agentinit/commit/9a5c2c9d2b7e53d6be751b133d140b220da4db7f))
* **workflow:** validate restructured scaffold cycle ([8ce3ec1](https://github.com/riadshalaby/agentinit/commit/8ce3ec178e80515843b1e714c1bb6f075016aa07))


### Documentation

* **agents:** promote hard rules to the top ([dd06775](https://github.com/riadshalaby/agentinit/commit/dd067753a702e62a07a23bdc02b18f447b94c986))
* **agents:** sync project files with scaffold templates ([1946dac](https://github.com/riadshalaby/agentinit/commit/1946dac06402f43c30857ee70af4a788ab47e37f))
* **handoff:** standardize handoff entry format ([eb01a1b](https://github.com/riadshalaby/agentinit/commit/eb01a1b136027813b69aff1a562feeff578f4424))
* **po:** clarify post-planning auto-mode run control ([e837457](https://github.com/riadshalaby/agentinit/commit/e8374570d64ba5e27d68c33c0e052b6f4d5e9c19))
* **po:** define explicit work_task and work_all commands ([60096d1](https://github.com/riadshalaby/agentinit/commit/60096d1a3e10bb1cdb4ccf6291c044b4abf5f475))
* **readme:** clarify manual, auto, and MCP workflow usage ([881d06b](https://github.com/riadshalaby/agentinit/commit/881d06b8dfb764a310cca9fd7981fabe313ae8b8))
* **readme:** describe interactive init wizard ([5e879e7](https://github.com/riadshalaby/agentinit/commit/5e879e74246f7631a25a3b68d824c50eea10a1a6))
* **readme:** improve onboarding clarity and extensibility for future workflows ([9bbd7f3](https://github.com/riadshalaby/agentinit/commit/9bbd7f3ee525ebb3040507a9fa1f60f18a3685c0))
* **roadmap:** clarify required and optional roadmap sections ([e617092](https://github.com/riadshalaby/agentinit/commit/e61709213deee6dc10bcf40a11f8cc27ba1a8b03))
* **scaffold:** address T-003 review findings ([6f1707d](https://github.com/riadshalaby/agentinit/commit/6f1707d7ca11b173f1cdfb1a4c17b075653ffe1a))
* **scaffold:** explain manual and auto runtime modes ([4e57566](https://github.com/riadshalaby/agentinit/commit/4e57566a265d48b65a9ebbe044cd78819ed05e85))
* **workflow:** add persistent session examples and coverage ([3164132](https://github.com/riadshalaby/agentinit/commit/31641322754b1b565a565e8435ecc15b2b52680b))
* **workflow:** align repo instruction files with layered layout ([1407cd3](https://github.com/riadshalaby/agentinit/commit/1407cd366ce1ad1c321adfc933606a606baa06ec))
* **workflow:** allow reviewer commits for review artifacts ([f136a0b](https://github.com/riadshalaby/agentinit/commit/f136a0b12d846e714163b9c4d7630ac915bbf626))
* **workflow:** define persistent session state transitions ([1920e81](https://github.com/riadshalaby/agentinit/commit/1920e812e52d7ebcbbf1e59e7b7f44be3890038e))
* **workflow:** document review rework flow ([e9b0748](https://github.com/riadshalaby/agentinit/commit/e9b07484b3b7391a32b3d8837e8069ed96b94a9f))
* **workflow:** document task-scoped .ai artifact commits ([2502507](https://github.com/riadshalaby/agentinit/commit/25025077a27ade43ac5ad8062aa7127c1611063b))
* **workflow:** record T-004 review and installer docs ([11585b0](https://github.com/riadshalaby/agentinit/commit/11585b0ea457ef4ec10a2d7b1e4d95eda94743ec))
* **workflow:** switch generated guidance to persistent agent sessions ([9dea910](https://github.com/riadshalaby/agentinit/commit/9dea91090083acd79570988068b1edc6da5caadf))

## [0.6.1](https://github.com/riadshalaby/agentinit/compare/v0.6.0...v0.6.1) (2026-04-13)


### Features

* **ai:** removed agents from tasks board ([668b150](https://github.com/riadshalaby/agentinit/commit/668b150d30497ecbdf5c759f12120128a6f8fa55))
* **mcp:** poll session output with get_output ([52d69ff](https://github.com/riadshalaby/agentinit/commit/52d69ffb0999f26248dc776f763c13e752562134))
* **mcp:** write MCP server debug logs to .ai/mcp-server.log ([bcc11f3](https://github.com/riadshalaby/agentinit/commit/bcc11f3ff06e2cddba92be8e086aadb719590c07))


### Bug Fixes

* **mcp:** force-stop hung sessions with SIGKILL ([02a7882](https://github.com/riadshalaby/agentinit/commit/02a7882b7aa6cfd564abc7f2df39597de08a72f2))
* **mcp:** preserve structured JSON tool results ([7bea686](https://github.com/riadshalaby/agentinit/commit/7bea68603cf15be79fb17a1b352c793ef36711f7))


### Miscellaneous Chores

* **ai:** close cycle ([5327cda](https://github.com/riadshalaby/agentinit/commit/5327cda465b42c1f05f415475630b412e6f3b549))
* **ai:** roadmap for version 0.7.0 ([a327cf1](https://github.com/riadshalaby/agentinit/commit/a327cf15b1133014fca585ca6355c9b269aa9301))
* **docs:** new logo ([2a27b4b](https://github.com/riadshalaby/agentinit/commit/2a27b4bda9c3ab331ef88cc67d20cc7f7e1103a1))
* start cycle v0.7.0 ([51bad16](https://github.com/riadshalaby/agentinit/commit/51bad16f47d3c5a8c99703af21a06067de0d9e3c))


### Documentation

* **po:** clarify post-planning auto-mode run control ([e837457](https://github.com/riadshalaby/agentinit/commit/e8374570d64ba5e27d68c33c0e052b6f4d5e9c19))
* **workflow:** document task-scoped .ai artifact commits ([2502507](https://github.com/riadshalaby/agentinit/commit/25025077a27ade43ac5ad8062aa7127c1611063b))

## [0.6.0](https://github.com/riadshalaby/agentinit/compare/v0.5.1...v0.6.0) (2026-04-12)


### Features

* **claude:** add scaffolded tool-access permissions by project type ([293b8be](https://github.com/riadshalaby/agentinit/commit/293b8be6a98a0ce4f5603f7add518f168ab3b776))
* **claude:** add tool preference guidance to CLAUDE files ([519fd40](https://github.com/riadshalaby/agentinit/commit/519fd40e3157c65157707f0346a3458f3621d363))
* **claude:** scaffold Claude settings templates ([93580d3](https://github.com/riadshalaby/agentinit/commit/93580d38b2434bdd7e88b8b81df41b0cbfdecc9b))
* **config:** scaffold per-role agent, model, and effort defaults in .ai/config.json ([674143b](https://github.com/riadshalaby/agentinit/commit/674143b639c0c9fc2eb4b2354c7851a259619783))
* **init:** add interactive setup wizard ([d887311](https://github.com/riadshalaby/agentinit/commit/d887311459cbc14fb750a381c6e271b325b883b5))
* **init:** add manual and auto workflow scaffolds ([c551b4f](https://github.com/riadshalaby/agentinit/commit/c551b4f6e1b2807823a37717e10b452df2b365a3))
* initial agentinit CLI scaffold generator ([e91f919](https://github.com/riadshalaby/agentinit/commit/e91f9193d1c4a6054b1b88cbd1d1347316e099ea))
* **mcp:** add MCP session management tools ([b345e81](https://github.com/riadshalaby/agentinit/commit/b345e811a553e7d09687b9b9810d29cf80df3a6e))
* **mcp:** add the stdio MCP server command ([88b7de2](https://github.com/riadshalaby/agentinit/commit/88b7de2f7142f291b2ee2d37a38c94b1e16c4547))
* **planner:** add roadmap refinement guidance before start_plan ([360a0e7](https://github.com/riadshalaby/agentinit/commit/360a0e78320f4cebb0a61bbda0310a1f05251f44))
* **prereq:** add fd, bat, and jq to wizard setup checks ([665cb60](https://github.com/riadshalaby/agentinit/commit/665cb60e0d99e105ed9be2ccdea9c1196d945a6c))
* **prereq:** add optional advanced search tools ([29888b5](https://github.com/riadshalaby/agentinit/commit/29888b54b464a311d76bb4b5b9c85c7fd66f76d9))
* **prereq:** add platform-specific Claude and Codex installs ([b5c94d1](https://github.com/riadshalaby/agentinit/commit/b5c94d19839d58dc3a83754048e43c7425a15acb))
* **prompts:** add shared search-strategy guidance ([6cb7daf](https://github.com/riadshalaby/agentinit/commit/6cb7daf56a9b4954ea034ab781eb93f48688fc97))
* **prompts:** inline critical workflow rules ([72ae946](https://github.com/riadshalaby/agentinit/commit/72ae946af72a9e7a43c0bc03d1d7a2277e16b0b4))
* **scaffold:** always render PO artifacts ([50f5cf2](https://github.com/riadshalaby/agentinit/commit/50f5cf2629ac059347daa019cac50afccf71f724))
* **scaffold:** generate manifest for managed files ([416474c](https://github.com/riadshalaby/agentinit/commit/416474ca20d99a67e8cc85303b7312ffedd4a837))
* **scaffold:** share init summaries across cli and wizard ([12d4b3f](https://github.com/riadshalaby/agentinit/commit/12d4b3f1b8da1eb841e57b3149bb3635644d0afe))
* **templates:** add the PO orchestration prompt and launcher ([3e30eb6](https://github.com/riadshalaby/agentinit/commit/3e30eb617aec29635e88928916e3bef9abbd5ed0))
* **update:** migrate legacy workflow files during scaffold refresh ([95aee2c](https://github.com/riadshalaby/agentinit/commit/95aee2c91f9bd16a6e27db0530e22559a7369ab8))
* **update:** refresh managed scaffold files ([6d52ef6](https://github.com/riadshalaby/agentinit/commit/6d52ef697d2a46a38a68f8075b1be208afdefc58))
* **workflow:** add a ready-to-commit stage to task flow ([a5b126e](https://github.com/riadshalaby/agentinit/commit/a5b126e1d7d1104e2ec6871668e0fccf2434fbfe))
* **workflow:** add cycle bootstrap pre-flight checks ([2164856](https://github.com/riadshalaby/agentinit/commit/21648563fc9b1cf41de40725fd776b2975091ac9))
* **workflow:** add shorthand commands and improvement roadmap ([1e04a49](https://github.com/riadshalaby/agentinit/commit/1e04a491be9828750324369e182a892ce8c55347))
* **workflow:** add the tester role and status flow ([dcd17e3](https://github.com/riadshalaby/agentinit/commit/dcd17e306979ab85bb68407408207e60c25307dd))
* **workflow:** categorize tools for the wizard and CLAUDE guidance ([156667e](https://github.com/riadshalaby/agentinit/commit/156667e1990a95dd58b601c924570db49163ebaa))
* **workflow:** merge agent rules into AGENTS.md ([030e8f1](https://github.com/riadshalaby/agentinit/commit/030e8f17b56954357595df5bc197f27ff2cca832))
* **workflow:** move finish_cycle to the implementer role ([89beec4](https://github.com/riadshalaby/agentinit/commit/89beec417a399bbfc91af67897ea35a1f620f694))
* **workflow:** remove the tester role from scaffolded projects ([5f30e3e](https://github.com/riadshalaby/agentinit/commit/5f30e3e506ecd778281f5ec8a4262a77b44fd869))
* **workflow:** require fresh file reads for every role command ([6756286](https://github.com/riadshalaby/agentinit/commit/6756286e22ca4c98883256684fbec75a040b321a))
* **workflow:** split scaffolded agent rules into layered files ([0fa5780](https://github.com/riadshalaby/agentinit/commit/0fa5780c42e6f7d451750f6014bd6afd6d9cde45))
* **workflow:** track cycle review and test logs in git ([3847efe](https://github.com/riadshalaby/agentinit/commit/3847efe7e4f87ac7868835859f5f8d8a08a9be9d))


### Bug Fixes

* **agents:** address T-007 test failure ([25bfd1e](https://github.com/riadshalaby/agentinit/commit/25bfd1ebed77c51857266753cbc8f59b5e0e0759))
* **ai:** removed old agents config and commands ([7bebe7d](https://github.com/riadshalaby/agentinit/commit/7bebe7d818c30301d855e5e98eecbc3615f045ef))
* **cli:** report build versions from Go module metadata ([9fc009c](https://github.com/riadshalaby/agentinit/commit/9fc009c6fcca68662a70bd97acde41bc41334bcb))
* **prereq:** address review findings ([cdf1bf6](https://github.com/riadshalaby/agentinit/commit/cdf1bf61d7b7638f4576e53e98930ac3ad5aa62d))
* **prereq:** address review findings ([4861a14](https://github.com/riadshalaby/agentinit/commit/4861a1438aa54e517c1c1a5e53e5e516bd77e2f2))
* **prereq:** install the tree-sitter CLI on macOS ([defe919](https://github.com/riadshalaby/agentinit/commit/defe91985dc35a2a45e1d27a0df45ab48a161467))
* **release:** add Windows binaries to GoReleaser artifacts ([e958c06](https://github.com/riadshalaby/agentinit/commit/e958c06c57cc3bb254a4846ad96ec31c86a766d4))
* **release:** build release assets from release-please tags ([dd9e50c](https://github.com/riadshalaby/agentinit/commit/dd9e50cacd0d818d9278de798abff404649282f7))
* **release:** use unprefixed release tags ([bd7e4cc](https://github.com/riadshalaby/agentinit/commit/bd7e4cc94cc73e0e18990c16f9cb0bd6263c2998))
* **scaffold:** address T-002 review findings ([cb3faab](https://github.com/riadshalaby/agentinit/commit/cb3faab9c710a34f4ec680028d9ad58657a90a13))
* **templates:** address review findings for the PO workflow ([1c09645](https://github.com/riadshalaby/agentinit/commit/1c096450031a2a1e4b4148f4e649ab9970e2f281))
* **workflow:** address review findings ([e2ca6fd](https://github.com/riadshalaby/agentinit/commit/e2ca6fdb81da188bd198ab931ebc38b4a378960c))
* **workflow:** address review findings ([76c17ac](https://github.com/riadshalaby/agentinit/commit/76c17ac5fa54d8b25155c4dcddc5a4097dffa67c))
* **workflow:** address review findings for cycle bootstrap ([af34bf7](https://github.com/riadshalaby/agentinit/commit/af34bf77e50c60d776ad49abff08f6c97cdf49d0))
* **workflow:** address scaffold merge regressions ([53e9301](https://github.com/riadshalaby/agentinit/commit/53e9301cdb68a9a44d2d99e592c192f3c984e9a3))
* **workflow:** address tester review findings ([5584411](https://github.com/riadshalaby/agentinit/commit/558441158e9ff7597382462975df32cd53edfc69))
* **workflow:** correct the tester handoff status flow ([8d2f9f8](https://github.com/riadshalaby/agentinit/commit/8d2f9f839c59ce04bd78f191d47ac633556e0c8a))
* **workflow:** keep commit conventions in workflow-managed files ([405d588](https://github.com/riadshalaby/agentinit/commit/405d58875c0da32f5476f3364a9f90d8f7d17612))
* **workflow:** keep tester in manual scaffolds and ignore runtime reports ([89bc419](https://github.com/riadshalaby/agentinit/commit/89bc4191e224d3cd526d4d3272c4c075d6cdc5f1))


### Miscellaneous Chores

* added logo to README ([1473df3](https://github.com/riadshalaby/agentinit/commit/1473df385bbe68da33f8bf43281b2a995b327459))
* **ai:** close cycle ([288fe5e](https://github.com/riadshalaby/agentinit/commit/288fe5ebbd966ffed76a6923f3c5ef446a7014c6))
* **ai:** new claude settings ([24f9f5b](https://github.com/riadshalaby/agentinit/commit/24f9f5b592194333a0f7d7d96e8d87934817f244))
* **ai:** plan update ([9db7a94](https://github.com/riadshalaby/agentinit/commit/9db7a94523de86400e3ac9c28cfb2494a8074c11))
* **ai:** record review outcome ([64aa536](https://github.com/riadshalaby/agentinit/commit/64aa536248900d43ec8dadeb66e58c9e1d1098c3))
* **ai:** roadmap for v0.2.0 ([f047ada](https://github.com/riadshalaby/agentinit/commit/f047adae2f15b91c0f87c60f4d1463c265299502))
* **ai:** roadmap for v0.4.0 ([c16db5e](https://github.com/riadshalaby/agentinit/commit/c16db5e33aa9472c602d1938bd5dc08162ae30c4))
* **ai:** roadmap for version 0.5.0 ([960ce0a](https://github.com/riadshalaby/agentinit/commit/960ce0aa7495fc54ab595246a47461bfe611bce8))
* **ai:** roadmap for version 0.5.1 ([349366a](https://github.com/riadshalaby/agentinit/commit/349366a7a793833393c3557d154fead3012ee098))
* **ai:** roadmap update ([e0fffbd](https://github.com/riadshalaby/agentinit/commit/e0fffbd829133b325ca439045fbdb3b58da9ad01))
* **claude:** allow more tools for claude ([a934a54](https://github.com/riadshalaby/agentinit/commit/a934a5452d988e027892c494e19cc3e1896b7af3))
* **hooks:** reject co-authored commit trailers ([e142dc0](https://github.com/riadshalaby/agentinit/commit/e142dc06eec5cc435a0d3198e623a4fb2032d074))
* **main:** release 0.3.0 ([11eef89](https://github.com/riadshalaby/agentinit/commit/11eef8970e2ad2eb3c6b5c9becba73b4c7c25390))
* **main:** release 0.4.0 ([e5e6b35](https://github.com/riadshalaby/agentinit/commit/e5e6b35b435a4b55ed741b1e32631afe6f9cde4c))
* **main:** release 0.5.0 ([d887207](https://github.com/riadshalaby/agentinit/commit/d88720778b47a538f07d693437031576c7fd9896))
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
* **review:** T-001 PASS_WITH_NOTES — ready_for_test ([d073756](https://github.com/riadshalaby/agentinit/commit/d07375643af5133ee031dedb2659eb1a8db26c41))
* **review:** T-001 PASS_WITH_NOTES R2 — ready_for_test ([c6deb0c](https://github.com/riadshalaby/agentinit/commit/c6deb0cca1b491c4f3464745f7696121330bcf92))
* **review:** T-002 changes_requested — remove dead Workflow bridge in engine.go ([a6633f4](https://github.com/riadshalaby/agentinit/commit/a6633f4fa9aff39b5ba2a643d7aab428852cf7d1))
* **review:** T-002 PASS_WITH_NOTES — ready_for_test ([4213480](https://github.com/riadshalaby/agentinit/commit/4213480978ee1dba3b622177563701cf242e697c))
* **review:** T-003 changes_requested — README stale workflow flags and missing PO role ([dfbe0d9](https://github.com/riadshalaby/agentinit/commit/dfbe0d95e59a96b3afcc10b276cf9b5f3e0b7363))
* **review:** T-003 PASS — ready_for_test ([5fbc956](https://github.com/riadshalaby/agentinit/commit/5fbc956dc1da70587c57e657de32411a9693ea7b))
* **review:** T-004 PASS — ready_for_test ([4f740fd](https://github.com/riadshalaby/agentinit/commit/4f740fdeb7bc3fccb1a9410d5f6e61bb46692e74))
* **review:** T-005 PASS — ready_for_test ([441db2b](https://github.com/riadshalaby/agentinit/commit/441db2bf13b2547bbc6f03c7f06f63c14f682afc))
* roadmap for next version ([c56dde0](https://github.com/riadshalaby/agentinit/commit/c56dde020023b7b302e776394f548665e3860d4d))
* roadmap for next version ([bd4a4b1](https://github.com/riadshalaby/agentinit/commit/bd4a4b113e91a33f8d36b55d8b2ffb4d3d9632e6))
* roadmap for next version - add PO workflow, ([4d90603](https://github.com/riadshalaby/agentinit/commit/4d90603bd19b36fb0f2eeea80548e134f37324a8))
* **scaffold:** remove redundant context guide ([bde40ab](https://github.com/riadshalaby/agentinit/commit/bde40abefdb2efced40fafd237fed9bbe74d6c36))
* start cycle 0.5.0 ([f1eabc5](https://github.com/riadshalaby/agentinit/commit/f1eabc59c5859148833797d665d112684bcb0a22))
* start cycle automatic-agents-workflow ([1785d38](https://github.com/riadshalaby/agentinit/commit/1785d381e950fd0647823821d1c14f8fefb854d5))
* start cycle better-commands ([eb073f4](https://github.com/riadshalaby/agentinit/commit/eb073f498c027d7b15ab0ff0c395bdcd66a4c46e))
* start cycle cleanup ([350171d](https://github.com/riadshalaby/agentinit/commit/350171db943aa408fe6a6ec3879c633ecb2db934))
* start cycle hire-po-tester ([f25763a](https://github.com/riadshalaby/agentinit/commit/f25763a03b32c4ed3ed5c0039eead396af7ccefb))
* start cycle mcp ([3c6ed8a](https://github.com/riadshalaby/agentinit/commit/3c6ed8ad39d47fa59a338c1a31c6c4a228d5efd0))
* start cycle refine-auto ([5b9d195](https://github.com/riadshalaby/agentinit/commit/5b9d195fa0b2e0e0e252f4154f51333518538436))
* start cycle refine-workflow ([6beaee6](https://github.com/riadshalaby/agentinit/commit/6beaee6828e099c7e860489e1af42a62812fb1ef))
* start cycle single-mode ([2ecb9c6](https://github.com/riadshalaby/agentinit/commit/2ecb9c6ca79c8609cf2ccf4ee018386f87ffb3dd))
* start cycle v0.5.1 ([3f161ec](https://github.com/riadshalaby/agentinit/commit/3f161ec9545089172991d1e7db990626d427e274))
* **workflow:** close the v0.5.1 cycle ([c2506d3](https://github.com/riadshalaby/agentinit/commit/c2506d3f03bf7c25d32e10232f9344e7bdde3097))
* **workflow:** disallow Co-Authored-By trailers in commit messages ([64a0e9d](https://github.com/riadshalaby/agentinit/commit/64a0e9dff43f6605a47ae806ae408d00564bcd56))
* **workflow:** finalize cycle task board ([af4323c](https://github.com/riadshalaby/agentinit/commit/af4323cd2ea5a734c672e2a4f84fef8b62f0ab64))
* **workflow:** record final review completion ([9a5c2c9](https://github.com/riadshalaby/agentinit/commit/9a5c2c9d2b7e53d6be751b133d140b220da4db7f))
* **workflow:** validate restructured scaffold cycle ([8ce3ec1](https://github.com/riadshalaby/agentinit/commit/8ce3ec178e80515843b1e714c1bb6f075016aa07))


### Documentation

* **agents:** promote hard rules to the top ([dd06775](https://github.com/riadshalaby/agentinit/commit/dd067753a702e62a07a23bdc02b18f447b94c986))
* **agents:** sync project files with scaffold templates ([1946dac](https://github.com/riadshalaby/agentinit/commit/1946dac06402f43c30857ee70af4a788ab47e37f))
* **handoff:** standardize handoff entry format ([eb01a1b](https://github.com/riadshalaby/agentinit/commit/eb01a1b136027813b69aff1a562feeff578f4424))
* **readme:** clarify manual, auto, and MCP workflow usage ([881d06b](https://github.com/riadshalaby/agentinit/commit/881d06b8dfb764a310cca9fd7981fabe313ae8b8))
* **readme:** describe interactive init wizard ([5e879e7](https://github.com/riadshalaby/agentinit/commit/5e879e74246f7631a25a3b68d824c50eea10a1a6))
* **readme:** improve onboarding clarity and extensibility for future workflows ([9bbd7f3](https://github.com/riadshalaby/agentinit/commit/9bbd7f3ee525ebb3040507a9fa1f60f18a3685c0))
* **roadmap:** clarify required and optional roadmap sections ([e617092](https://github.com/riadshalaby/agentinit/commit/e61709213deee6dc10bcf40a11f8cc27ba1a8b03))
* **scaffold:** address T-003 review findings ([6f1707d](https://github.com/riadshalaby/agentinit/commit/6f1707d7ca11b173f1cdfb1a4c17b075653ffe1a))
* **scaffold:** explain manual and auto runtime modes ([4e57566](https://github.com/riadshalaby/agentinit/commit/4e57566a265d48b65a9ebbe044cd78819ed05e85))
* **workflow:** add persistent session examples and coverage ([3164132](https://github.com/riadshalaby/agentinit/commit/31641322754b1b565a565e8435ecc15b2b52680b))
* **workflow:** align repo instruction files with layered layout ([1407cd3](https://github.com/riadshalaby/agentinit/commit/1407cd366ce1ad1c321adfc933606a606baa06ec))
* **workflow:** allow reviewer commits for review artifacts ([f136a0b](https://github.com/riadshalaby/agentinit/commit/f136a0b12d846e714163b9c4d7630ac915bbf626))
* **workflow:** define persistent session state transitions ([1920e81](https://github.com/riadshalaby/agentinit/commit/1920e812e52d7ebcbbf1e59e7b7f44be3890038e))
* **workflow:** document review rework flow ([e9b0748](https://github.com/riadshalaby/agentinit/commit/e9b07484b3b7391a32b3d8837e8069ed96b94a9f))
* **workflow:** record T-004 review and installer docs ([11585b0](https://github.com/riadshalaby/agentinit/commit/11585b0ea457ef4ec10a2d7b1e4d95eda94743ec))
* **workflow:** switch generated guidance to persistent agent sessions ([9dea910](https://github.com/riadshalaby/agentinit/commit/9dea91090083acd79570988068b1edc6da5caadf))

## [0.5.0](https://github.com/riadshalaby/agentinit/compare/v0.4.0...v0.5.0) (2026-04-11)


### Features

* **claude:** add scaffolded tool-access permissions by project type ([293b8be](https://github.com/riadshalaby/agentinit/commit/293b8be6a98a0ce4f5603f7add518f168ab3b776))
* **claude:** scaffold Claude settings templates ([93580d3](https://github.com/riadshalaby/agentinit/commit/93580d38b2434bdd7e88b8b81df41b0cbfdecc9b))
* **config:** scaffold per-role agent, model, and effort defaults in .ai/config.json ([674143b](https://github.com/riadshalaby/agentinit/commit/674143b639c0c9fc2eb4b2354c7851a259619783))
* **workflow:** add a ready-to-commit stage to task flow ([a5b126e](https://github.com/riadshalaby/agentinit/commit/a5b126e1d7d1104e2ec6871668e0fccf2434fbfe))
* **workflow:** track cycle review and test logs in git ([3847efe](https://github.com/riadshalaby/agentinit/commit/3847efe7e4f87ac7868835859f5f8d8a08a9be9d))


### Bug Fixes

* **release:** build release assets from release-please tags ([dd9e50c](https://github.com/riadshalaby/agentinit/commit/dd9e50cacd0d818d9278de798abff404649282f7))


### Miscellaneous Chores

* **ai:** roadmap for version 0.5.0 ([960ce0a](https://github.com/riadshalaby/agentinit/commit/960ce0aa7495fc54ab595246a47461bfe611bce8))
* start cycle 0.5.0 ([f1eabc5](https://github.com/riadshalaby/agentinit/commit/f1eabc59c5859148833797d665d112684bcb0a22))

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
