# Type Battle



## 必要な環境

- Go 1.25以降
- LinuxまたはWSLでウィンドウを表示できる環境
  - Windows 11のWSLではWSLgを利用できます

## Linux / WSLのセットアップ

Ebitengineのビルドに必要なX11、OpenGL、音声関連の開発パッケージをインストールします。

```bash
sudo apt update
sudo apt install -y gcc libc6-dev libgl1-mesa-dev libx11-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```

コマンドを複数行に分ける場合、行末のバックスラッシュ (`\`) の後ろに空白を入れないでください。貼り付け時のエラーを避けるには、上記の一行コマンドをそのまま使用してください。

## 実行

プロジェクトのディレクトリで次を実行します。

```bash
go run .
```

1000×600のウィンドウが開き、左側に200px幅の管理メニュー、右側に800×600のグリッドが表示されます。マスを左クリックすると、そのマスの状態が空と選択中の間で切り替わります。管理メニューには現在の`TIME`、グリッドの行・列数、選択中のマス数が表示されます。ウィンドウを閉じるとプログラムも終了します。

`TIME`は1から始まり、Spaceキーを押すたびに1つ進みます。

グリッドの列数と行数は、`main.go`の`gridColumns`と`gridRows`で変更できます。

## トラブルシューティング

### `X11/Xlib.h: No such file or directory`

Linux向けの開発パッケージが不足しています。「Linux / WSLのセットアップ」に記載した`apt install`を実行してください。

### ウィンドウが表示されない

WSLから実行している場合は、WSLgなどのGUI表示環境が利用できることを確認してください。
