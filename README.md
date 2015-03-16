# golang + heroku で slackbot

## heroku

[セールスフォース・ドットコム](http://www.salesforce.com/)が提供するPaaS (Platform as a Service)。

### アカウントを作成

[https://www.heroku.com/](https://www.heroku.com/)にアクセスしてアカウントを作成する。

### コマンドラインツールをインストール

Mac OS X, Windows, Debian, Ubuntuは[heroku toolbelt](https://toolbelt.heroku.com/)のサイトから、[veca](http://gpo.zugaina.org/Overlays/vaca)をlaymanで追加すれば、emergeできる。

```
sudo emerge -Dav dev-util/heroku-client
```

取敢へず以下のコマンド邊りを憶える。

```
heroku app:create         # アプリケーションを作成
heroku auth:login         # ログイン
heroku config             # 環境変数の確認
heroku config:set FOO=bar # 環境変数FOOにbarを設定
heroku config:unset FOO   # 環境変数FOOを削除
heroku logs --tail        # ログの確認
```

## go

### 必須のツール

* godep (go get github.com/tools/godep)
  * 依存性を固定

### あると便利なツール

* gocode (go get github.com/nsf/gocode)
  * 補完
* goimports (go get golang.org/x/tools/cmd/goimports)
  * 自動import
* golint (go get github.com/golang/lint/golint)
  * lint

### botを書く

* [http://golang.org/doc](http://golang.org/doc)
* [http://golang.jp/](http://golang.jp/)

## デプロイ

基本的な手順。

* `godep` で依存性を固定
* *Procfile* を作成
* `heroku create` でアプリケーションを作成
  * [goのbuildpack](https://github.com/kr/heroku-buildpack-go.git)を使用
* Herokuにpush

### godep

`godep save` するとアプリケーションのディレクトリに `Godeps` と云ふディレクトリが作成される。
その下の *_workspace* に依存するパッケージが展開される。

```
godep save
git add Godeps
git commit -m "dependency"
```

*_workspace* 下に配置されたパッケージを用ゐてビルドする場合は、 `godep go build yachecker.go` とする。

### Procfile

webプロセスがどのコマンドを叩くかHerokuに教へる爲のファイル。

```
echo "web: appname" > Procfile
git add Procfile
git commit -m "procfile"
```

### アプリケーションを作成

Go用のbuildpackを使用して作成。

```
heroku create yachecker -b https://github.com/kr/heroku-buildpack-go.git
```

を實行すると以下のやうなメッセージが出力される。

```
Creating yachecker... done, stack is cedar-14
Buildpack set. Next release on yachecker will use https://github.com/kr/heroku-buildpack.git.
https://yachecker.herokuapp.com/ | https://git.heroku.com/yachecker.git
Git remote heroku added
```

このアプリケーションのURLとリポジトリは、

* https://yachecker.herokuapp.com/
* https://git.heroku.com/yachecker.git

で有る事が判る。

### push

```
git push heroku master
```

成功すると以下のやうに出力される。

```
Initializing repository, done.
Counting objects: 176, done.
Delta compression using up to 2 threads.
Compressing objects: 100% (163/163), done.
Writing objects: 100% (176/176), 252.24 KiB | 334.00 KiB/s, done.
Total 176 (delta 41), reused 0 (delta 0)

-----> Fetching custom git buildpack... done
-----> Go app detected
-----> Installing go1.4.2... done
-----> Running: godep go install -tags heroku ./...
-----> Discovering process types
       Procfile declares types -> web

-----> Compressing... done, 2.2MB
-----> Launching... done, v3
       https://yachecker.herokuapp.com/ deployed to Heroku

To git@heroku.com:yachecker.git
 * [new branch]      master -> master
```

### 完了

[https://yachecker.herokuacpp.com/](https://yachecker.herokuapp.com/)にアクセスして動作確認。
