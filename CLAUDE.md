# youtube-summarize-scraper (YTSS) — 專案指令

> Go CLI,批次處理 YouTube 頻道:抓字幕 → 失敗則 whisper.cpp 轉錄 → 3-stage LLM 摘要(摘要 + 關鍵字 + Mermaid 流程圖)→ 寫成 markdown 進 Obsidian vault。每日 launchd cron 09:00 處理 9 個頻道 → `stephen_second_brain/learning/videos/`。這是 memory 裡的「YTSS 自動化」本體。

## Stack

- **Go 1.25**,Cobra CLI。yt-dlp / ffmpeg / whisper-cli **由原始碼編譯並 go:embed 進 binary**(無外部執行檔依賴)。
- 結構:`cmd/`(run/video/channel)｜`pipeline/`(編排)｜`fetcher/`(yt-dlp)｜`subtitle/`(4 階字幕 cascade)｜`transcriber/`(whisper)｜`summarizer/`(3-stage,5 種 LLM backend)｜`output/`(markdown + Obsidian MOC)｜`lang/`｜`config/`｜`embedded/`。

## 指令

```bash
make all        # 完整 build(下載 yt-dlp、編 ffmpeg、編 whisper、編 Go binary)
make build      # 只編 Go(假設 embedded/bin/ 已備)
make build-all  # cross-compile(darwin-arm64/amd64、linux-amd64,需各平台依賴)
go test ./... -count=1   # 60+ tests(fetcher/lang/output/subtitle/summarizer)
```

**verifier(loop gate)**:`go test ./... -count=1` → 🟢 GREEN(~0.7s)。

## 核心資訊流

```
頻道/影片 URL + config.yaml
  → yt-dlp 抓 metadata
  → 4 階字幕 cascade(manual→auto target → manual→auto any)
  → 全失敗 → ffmpeg 抽 WAV 16kHz → whisper-cli(語言別模型:kotoba-ja/belle-zh/medium-en)→ SRT
  → ⚠️3-stage LLM(summary 2000 tok / keywords / Mermaid 自動修語法)
  → 寫 output_dir/@handle/YYYY-MM-DD__id__title/(summary.md + frontmatter + mermaid;可 copy_to 跨 vault)
  → ledger skip 已處理
```

⚠️ 非決定性:LLM summary/keywords/mermaid。

## Gotchas

- **ffmpeg 從源碼編**(需 cmake/nasm/pkg-config);whisper 模型首次 on-demand 下載(~1.4-1.5GB)到 `~/.ytss/models/`。
- 受限影片支援 browser cookie 自動擷取,fallback `cookie.txt`(Netscape)。
- 語言要明寫 `zh-Hant,zh-Hans`(逗號=yt-dlp 優先序),不要裸 `zh`。
- 已知:config 中 `JayShettyPodcast` 頻道 404/改名;9 頻道常有 2 個 fail(下次 cron retry)。
- 相關但分離:`monkey-knowledge-youtube-skills`(平行的 YouTube skill set,補充用);**YTSS 是 production 主力**。memory 的「mk-youtube ffmpeg 缺口」是指那個 skill,不是本 CLI(本 CLI 自帶 embedded ffmpeg)。

## CI / Verifier

- 無 GitHub Actions。執行靠 launchd `com.stephen.ytss.plist` 每日 09:00。決定性 gate:`go test ./...`(綠)+ `go build`。

## Loop 啟動決策(2026-06-21)

- **這已是一條 launchd 每日 cron loop。** 決定性 verifier(`go test ./...`)綠;LLM 摘要品質非決定性、**無 eval**,品質靠 Stephen 讀 vault 筆記判斷(taste)。
- **閘門**:成熟運行中,**先不建品質 eval**(重複量足夠但品質判斷主觀,且現況可用)。待辦是修 config 中失效頻道,不是 loop。
- **loop 類型**:cron(launchd daily);dev 用 goal loop(`go test`)。
